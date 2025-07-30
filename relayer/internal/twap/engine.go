package twap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/ethereum"
)

const (
	OneInchAPIBaseURL   = "https://api.1inch.dev"
	OneInchSwapAPI      = "/swap/v6.0"
	OneInchQuoteAPI     = "/quote/v1.4"
	MaxRetries          = 3
	RetryDelay          = 2 * time.Second
	ConfirmationBlocks  = 3
	MaxConcurrentSwaps  = 10
	RateLimitPerSecond  = 5
	OrderTTL            = 24 * time.Hour
	ExecutionTimeout    = 30 * time.Minute
	CircuitBreakerDelay = 5 * time.Minute
)

// Engine represents the TWAP execution engine
type Engine struct {
	config            Config
	ethClient         *ethereum.Client
	cosmosClient      *cosmos.Client
	logger            *zap.Logger
	
	// State management
	orders            map[string]*TWAPOrder
	ordersMutex       sync.RWMutex
	redisClient       *redis.Client
	
	// API clients
	oneInchClient     *OneInchClient
	
	// Execution management
	scheduler         *ExecutionScheduler
	rateLimiter       *rate.Limiter
	metrics           *Metrics
	chainID           *big.Int
	
	// Circuit breaker
	circuitBreaker    *CircuitBreaker
}

// CircuitBreaker implements circuit breaker pattern for 1inch API
type CircuitBreaker struct {
	isTripped     bool
	lastTrip      time.Time
	failureCount  int
	threshold     int
	resetTimeout  time.Duration
	mutex         sync.RWMutex
}

// OneInchClient handles 1inch API interactions
type OneInchClient struct {
	apiKey      string
	httpClient  *http.Client
	logger      *zap.Logger
	rateLimiter *rate.Limiter
}

// Config holds engine configuration
type Config struct {
	RedisURL      string
	OneInchAPIKey string
	ChainID       string
	MaxGasPrice   *big.Int
	MaxIntervals     int           
	MinIntervalTime  time.Duration 
	MaxSlippage      float64       
	DatabaseURL      string  
}

// TWAPOrder represents a TWAP order
type TWAPOrder struct {
	ID              string             `json:"id"`
	UserAddress     string             `json:"user_address"`
	SourceToken     string             `json:"source_token"`
	DestToken       string             `json:"dest_token"`
	TotalAmount     *big.Int           `json:"total_amount"`
	IntervalCount   int                `json:"interval_count"`
	TimeWindow      int                `json:"time_window"`
	ExecutionPlan   []*ExecutionWindow `json:"execution_plan"`
	Status          OrderStatus        `json:"status"`
	CompletedWindows int               `json:"completed_windows"`
	TotalExecuted   *big.Int           `json:"total_executed"`
	CreatedAt       time.Time          `json:"created_at"`
	StartTime       time.Time          `json:"start_time"`
}

// ExecutionWindow represents a single TWAP execution window
type ExecutionWindow struct {
	Index           int           `json:"index"`
	Amount          *big.Int      `json:"amount"`
	Secret          []byte        `json:"-"`
	SecretHash      string        `json:"secret_hash"`
	EthereumTxHash  string        `json:"ethereum_tx_hash"`
	CosmosTxHash    string        `json:"cosmos_tx_hash"`
	ActualPrice     *big.Int      `json:"actual_price"`
	GasUsed         uint64        `json:"gas_used"`
	ExecutedAt      *time.Time    `json:"executed_at"`
	Status          WindowStatus  `json:"status"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
}

// Request/Response types
type CreateTWAPOrderRequest struct {
	UserAddress   string  `json:"user_address" binding:"required"`
	SourceToken   string  `json:"source_token" binding:"required"`
	DestToken     string  `json:"dest_token" binding:"required"`
	TotalAmount   string  `json:"total_amount" binding:"required"`
	TimeWindow    int     `json:"time_window" binding:"required,min=300"`
	IntervalCount int     `json:"interval_count" binding:"required,min=2"`
	MaxSlippage   float64 `json:"max_slippage" binding:"required,min=0.01,max=10"`
	StartDelay    int     `json:"start_delay,omitempty"`
}

type TWAPQuoteRequest struct {
	SourceToken   string  `json:"source_token" binding:"required"`
	DestToken     string  `json:"dest_token" binding:"required"`
	Amount        string  `json:"amount" binding:"required"`
	IntervalCount int     `json:"interval_count" binding:"required"`
	TimeWindow    int     `json:"time_window" binding:"required"`
	MaxSlippage   float64 `json:"max_slippage"`
}

type TWAPQuoteResponse struct {
	EstimatedOutput   string                `json:"estimated_output"`
	PriceImpact       float64               `json:"price_impact"`
	ImpactComparison  PriceImpactComparison `json:"impact_comparison"`
	ExecutionSchedule []SchedulePreview     `json:"execution_schedule"`
	Fees              FeeBreakdown          `json:"fees"`
	Recommendations   []string              `json:"recommendations"`
}

type PriceImpactComparison struct {
	InstantSwap        float64 `json:"instant_swap"`
	TWAPExecution      float64 `json:"twap_execution"`
	ImprovementPercent float64 `json:"improvement_percent"`
}

type SchedulePreview struct {
	IntervalIndex  int       `json:"interval_index"`
	ExecutionTime  time.Time `json:"execution_time"`
	Amount         string    `json:"amount"`
	EstimatedPrice string    `json:"estimated_price"`
}

type FeeBreakdown struct {
	EthereumGas string `json:"ethereum_gas"`
	CosmosGas   string `json:"cosmos_gas"`
	RelayerFee  string `json:"relayer_fee"`
	OneInchFee  string `json:"oneinch_fee"`
	Total       string `json:"total"`
}

// OrderStatus represents order states
type OrderStatus string
const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusExecuting OrderStatus = "executing"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// WindowStatus represents window states
type WindowStatus string
const (
	WindowStatusPending   WindowStatus = "pending"
	WindowStatusExecuting WindowStatus = "executing"
	WindowStatusCompleted WindowStatus = "completed"
	WindowStatusFailed    WindowStatus = "failed"
	WindowStatusSkipped   WindowStatus = "skipped"
)

// OneInch API response types
type OneInchQuoteResponse struct {
	DstAmount string `json:"dstAmount"`
	Gas       string `json:"gas"`
	GasPrice  string `json:"gasPrice"`
	Protocols [][]struct {
		Name             string  `json:"name"`
		Part             float64 `json:"part"`
		FromTokenAddress string  `json:"fromTokenAddress"`
		ToTokenAddress   string  `json:"toTokenAddress"`
	} `json:"protocols"`
}

type OneInchSwapResponse struct {
	DstAmount string `json:"dstAmount"`
	Tx        struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Data     string `json:"data"`
		Value    string `json:"value"`
		GasPrice string `json:"gasPrice"`
		Gas      int    `json:"gas"`
	} `json:"tx"`
}

// NewEngine creates a new TWAP engine
func NewEngine(config Config, ethClient *ethereum.Client, cosmosClient *cosmos.Client, logger *zap.Logger) (*Engine, error) {
	// Initialize Redis
	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	redisClient := redis.NewClient(opt)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Initialize 1inch client
	oneInchClient := &OneInchClient{
		apiKey: config.OneInchAPIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		logger:      logger,
		rateLimiter: rate.NewLimiter(rate.Limit(RateLimitPerSecond), 1),
	}

	// Parse chain ID
	chainID, ok := new(big.Int).SetString(config.ChainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain ID: %s", config.ChainID)
	}

	// Initialize circuit breaker
	circuitBreaker := &CircuitBreaker{
		threshold:    5,
		resetTimeout: CircuitBreakerDelay,
	}

	engine := &Engine{
		config:         config,
		ethClient:      ethClient,
		cosmosClient:   cosmosClient,
		logger:         logger,
		orders:         make(map[string]*TWAPOrder),
		redisClient:    redisClient,
		oneInchClient:  oneInchClient,
		rateLimiter:    rate.NewLimiter(rate.Limit(RateLimitPerSecond), MaxConcurrentSwaps),
		metrics:        NewMetrics(),
		chainID:        chainID,
		circuitBreaker: circuitBreaker,
	}

	// Initialize scheduler
	engine.scheduler = NewExecutionScheduler(logger, engine)

	return engine, nil
}

// Start runs the TWAP engine
func (e *Engine) Start(ctx context.Context) error {
	e.logger.Info("Starting Production TWAP Engine")

	// Restore persisted orders
	if err := e.restoreOrders(ctx); err != nil {
		e.logger.Error("Failed to restore orders", zap.Error(err))
	}

	// Start execution scheduler
	go e.scheduler.Start(ctx)

	// Start order monitoring
	go e.monitorOrders(ctx)

	// Main execution loop
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("TWAP Engine stopping")
			return e.shutdown(ctx)
		case <-ticker.C:
			if err := e.processExecutions(ctx); err != nil {
				e.logger.Error("Error processing executions", zap.Error(err))
			}
		}
	}
}

// GetQuote handles TWAP quote requests
func (e *Engine) GetQuote(c *gin.Context) {
	var req TWAPQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	totalAmount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || totalAmount.Sign() <= 0 {
		c.JSON(400, gin.H{"error": "invalid amount"})
		return
	}

	quote, err := e.calculateTWAPQuote(c.Request.Context(), req, totalAmount)
	if err != nil {
		e.logger.Error("Failed to calculate TWAP quote", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to calculate quote"})
		return
	}

	c.JSON(200, quote)
}

// CreateOrder handles TWAP order creation
func (e *Engine) CreateOrder(c *gin.Context) {
	var req CreateTWAPOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	totalAmount, ok := new(big.Int).SetString(req.TotalAmount, 10)
	if !ok || totalAmount.Sign() <= 0 {
		c.JSON(400, gin.H{"error": "invalid total_amount"})
		return
	}

	if err := e.validateTWAPRequest(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	order, err := e.createTWAPOrder(c.Request.Context(), req, totalAmount)
	if err != nil {
		e.logger.Error("Failed to create TWAP order", zap.Error(err))
		c.JSON(500, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(201, order)
}

// GetOrder retrieves order by ID
func (e *Engine) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	c.JSON(200, order)
}

// GetSchedule retrieves execution schedule for order
func (e *Engine) GetSchedule(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	c.JSON(200, gin.H{
		"order_id":          order.ID,
		"execution_plan":    order.ExecutionPlan,
		"completed_windows": order.CompletedWindows,
		"status":           order.Status,
	})
}

// ListOrders lists all orders
func (e *Engine) ListOrders(c *gin.Context) {
	e.ordersMutex.RLock()
	orders := make([]*TWAPOrder, 0, len(e.orders))
	for _, order := range e.orders {
		orders = append(orders, order)
	}
	e.ordersMutex.RUnlock()

	c.JSON(200, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

// CancelOrder cancels an order
func (e *Engine) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.Lock()
	order, exists := e.orders[orderID]
	if exists && order.Status == OrderStatusPending {
		order.Status = OrderStatusCancelled
	}
	e.ordersMutex.Unlock()
	
	if !exists {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	c.JSON(200, gin.H{"message": "order cancelled"})
}

// GetStatus returns engine status
func (e *Engine) GetStatus() map[string]interface{} {
	e.ordersMutex.RLock()
	orderCount := len(e.orders)
	e.ordersMutex.RUnlock()

	return map[string]interface{}{
		"status":           "running",
		"order_count":      orderCount,
		"metrics":          e.metrics.GetSummary(),
		"uptime":           time.Since(e.metrics.StartTime).String(),
		"circuit_breaker":  e.circuitBreaker.IsTripped(),
	}
}

// shutdown handles graceful shutdown
func (e *Engine) shutdown(ctx context.Context) error {
	e.logger.Info("Initiating graceful shutdown")
	
	if err := e.persistOrders(ctx); err != nil {
		e.logger.Error("Failed to persist orders during shutdown", zap.Error(err))
	}
	
	if err := e.redisClient.Close(); err != nil {
		e.logger.Error("Failed to close Redis connection", zap.Error(err))
	}
	
	return nil
}

// processExecutions processes ready executions
func (e *Engine) processExecutions(ctx context.Context) error {
	now := time.Now()
	
	e.ordersMutex.RLock()
	orders := make([]*TWAPOrder, 0, len(e.orders))
	for _, order := range e.orders {
		if order.Status == OrderStatusPending || order.Status == OrderStatusExecuting {
			orders = append(orders, order)
		}
	}
	e.ordersMutex.RUnlock()

	for _, order := range orders {
		if order.Status == OrderStatusPending && now.After(order.StartTime) {
			e.ordersMutex.Lock()
			order.Status = OrderStatusExecuting
			e.ordersMutex.Unlock()
		}

		for _, window := range order.ExecutionPlan {
			if window.Status == WindowStatusPending && now.After(window.StartTime) && now.Before(window.EndTime) {
				go func(o *TWAPOrder, w *ExecutionWindow) {
					if err := e.executeWindow(ctx, o, w); err != nil {
						e.logger.Error("Window execution failed",
							zap.String("order_id", o.ID),
							zap.Int("window", w.Index),
							zap.Error(err),
						)
					}
				}(order, window)
			}
		}

		if e.isOrderComplete(order) {
			e.ordersMutex.Lock()
			order.Status = OrderStatusCompleted
			e.ordersMutex.Unlock()
		}
	}

	return nil
}

// hashSecret generates a Keccak256 hash
func (e *Engine) hashSecret(secret []byte) string {
	hash := crypto.Keccak256(secret)
	return "0x" + hex.EncodeToString(hash)
}

// generateSecret creates a cryptographically secure random secret
func (e *Engine) generateSecret() ([]byte, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}
	return secret, nil
}

// calculateTWAPQuote generates a TWAP quote
func (e *Engine) calculateTWAPQuote(ctx context.Context, req TWAPQuoteRequest, totalAmount *big.Int) (*TWAPQuoteResponse, error) {
	if err := e.validateQuoteRequest(req); err != nil {
		return nil, err
	}

	// Check circuit breaker
	if e.circuitBreaker.IsTripped() {
		return nil, fmt.Errorf("1inch API circuit breaker is tripped, try again later")
	}

	// Get instant quote
	instantQuote, err := e.oneInchClient.GetQuote(ctx, e.chainID.String(), req.SourceToken, req.DestToken, totalAmount.String())
	if err != nil {
		e.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("failed to get instant quote: %w", err)
	}

	instantOutput, ok := new(big.Int).SetString(instantQuote.DstAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid instant quote amount: %s", instantQuote.DstAmount)
	}

	instantImpact := e.calculatePriceImpact(totalAmount, instantOutput)

	// Calculate interval quote
	intervalAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	intervalQuote, err := e.oneInchClient.GetQuote(ctx, e.chainID.String(), req.SourceToken, req.DestToken, intervalAmount.String())
	if err != nil {
		e.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("failed to get interval quote: %w", err)
	}

	intervalOutput, ok := new(big.Int).SetString(intervalQuote.DstAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid interval quote amount: %s", intervalQuote.DstAmount)
	}

	twapEstimatedOutput := new(big.Int).Mul(intervalOutput, big.NewInt(int64(req.IntervalCount)))
	twapImpact := e.calculatePriceImpact(totalAmount, twapEstimatedOutput)

	improvement := ((instantImpact - twapImpact) / instantImpact) * 100
	if improvement < 0 {
		improvement = 0
	}

	schedule := e.generateSchedulePreview(req, totalAmount)
	fees, err := e.calculateFees(req, instantQuote.Gas, instantQuote.GasPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fees: %w", err)
	}

	e.circuitBreaker.RecordSuccess()

	return &TWAPQuoteResponse{
		EstimatedOutput: twapEstimatedOutput.String(),
		PriceImpact:     twapImpact,
		ImpactComparison: PriceImpactComparison{
			InstantSwap:        instantImpact,
			TWAPExecution:      twapImpact,
			ImprovementPercent: improvement,
		},
		ExecutionSchedule: schedule,
		Fees:              fees,
		Recommendations:   e.generateRecommendations(instantImpact, twapImpact, req.IntervalCount, req.TimeWindow),
	}, nil
}

// validateQuoteRequest validates TWAP quote request
func (e *Engine) validateQuoteRequest(req TWAPQuoteRequest) error {
	if !common.IsHexAddress(req.SourceToken) || !common.IsHexAddress(req.DestToken) {
		return fmt.Errorf("invalid token addresses")
	}
	if req.IntervalCount <= 0 || req.IntervalCount > 1000 {
		return fmt.Errorf("invalid interval count: %d", req.IntervalCount)
	}
	if req.TimeWindow <= 0 || req.TimeWindow > 86400 {
		return fmt.Errorf("invalid time window: %d", req.TimeWindow)
	}
	return nil
}

// validateTWAPRequest validates TWAP order request
func (e *Engine) validateTWAPRequest(req CreateTWAPOrderRequest) error {
	if !common.IsHexAddress(req.SourceToken) || !common.IsHexAddress(req.DestToken) {
		return fmt.Errorf("invalid token addresses")
	}
	if req.IntervalCount <= 0 || req.IntervalCount > 100 {
		return fmt.Errorf("invalid interval count: %d", req.IntervalCount)
	}
	if req.TimeWindow <= 0 || req.TimeWindow > 86400 {
		return fmt.Errorf("invalid time window: %d", req.TimeWindow)
	}
	if req.MaxSlippage < 0 || req.MaxSlippage > 10 {
		return fmt.Errorf("invalid slippage: %f", req.MaxSlippage)
	}
	return nil
}

// createTWAPOrder creates a new TWAP order
func (e *Engine) createTWAPOrder(ctx context.Context, req CreateTWAPOrderRequest, totalAmount *big.Int) (*TWAPOrder, error) {
	orderID := e.generateOrderID()
	
	startTime := time.Now().Add(time.Duration(req.StartDelay) * time.Second)
	if req.StartDelay == 0 {
		startTime = startTime.Add(1 * time.Minute)
	}

	executionPlan, err := e.generateExecutionPlan(ctx, req, totalAmount, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate execution plan: %w", err)
	}

	order := &TWAPOrder{
		ID:              orderID,
		UserAddress:     req.UserAddress,
		SourceToken:     req.SourceToken,
		DestToken:       req.DestToken,
		TotalAmount:     totalAmount,
		IntervalCount:   req.IntervalCount,
		TimeWindow:      req.TimeWindow,
		ExecutionPlan:   executionPlan,
		Status:          OrderStatusPending,
		CompletedWindows: 0,
		TotalExecuted:   big.NewInt(0),
		CreatedAt:       time.Now(),
		StartTime:       startTime,
	}

	e.ordersMutex.Lock()
	e.orders[orderID] = order
	e.ordersMutex.Unlock()

	e.scheduler.ScheduleOrder(order)
	e.metrics.RecordOrderScheduled()

	if err := e.persistOrder(ctx, order); err != nil {
		e.logger.Warn("Failed to persist order", zap.String("order_id", orderID), zap.Error(err))
	}

	e.logger.Info("TWAP order created",
		zap.String("order_id", orderID),
		zap.String("user", req.UserAddress),
		zap.Int("intervals", req.IntervalCount),
		zap.String("total_amount", totalAmount.String()),
	)

	return order, nil
}

// generateOrderID creates a unique order ID
func (e *Engine) generateOrderID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateExecutionPlan creates the execution plan for a TWAP order
func (e *Engine) generateExecutionPlan(ctx context.Context, req CreateTWAPOrderRequest, totalAmount *big.Int, startTime time.Time) ([]*ExecutionWindow, error) {
	windows := make([]*ExecutionWindow, req.IntervalCount)
	intervalDuration := time.Duration(req.TimeWindow) * time.Second / time.Duration(req.IntervalCount)
	
	baseAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	remainder := new(big.Int).Mod(totalAmount, big.NewInt(int64(req.IntervalCount)))

	for i := 0; i < req.IntervalCount; i++ {
		windowStart := startTime.Add(time.Duration(i) * intervalDuration)
		windowEnd := windowStart.Add(intervalDuration)

		windowAmount := new(big.Int).Set(baseAmount)
		if i < int(remainder.Int64()) {
			windowAmount.Add(windowAmount, big.NewInt(1))
		}

		secret, err := e.generateSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate secret for window %d: %w", i, err)
		}
		secretHash := e.hashSecret(secret)

		windows[i] = &ExecutionWindow{
			Index:      i,
			Amount:     windowAmount,
			Secret:     secret,
			SecretHash: secretHash,
			Status:     WindowStatusPending,
			StartTime:  windowStart,
			EndTime:    windowEnd,
		}
	}

	return windows, nil
}

// executeWindow executes a single TWAP window
func (e *Engine) executeWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	if e.circuitBreaker.IsTripped() {
		return fmt.Errorf("circuit breaker tripped")
	}

	if err := e.rateLimiter.Wait(ctx); err != nil {
		return err
	}

	window.Status = WindowStatusExecuting
	executionStart := time.Now()
	window.ExecutedAt = &executionStart

	e.logger.Info("Executing TWAP window",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("amount", window.Amount.String()),
	)

	var ethTxHash common.Hash
	var cosmosTxHash string
	var actualPrice *big.Int
	var gasUsed uint64

	// Exponential backoff retry logic
	for attempt := 0; attempt < MaxRetries; attempt++ {
		// Create Ethereum escrow
		hash, err := e.createEthereumEscrow(ctx, order, window)
		if err != nil {
			e.logger.Warn("Failed to create Ethereum escrow",
				zap.Int("attempt", attempt+1),
				zap.Error(err))
			if attempt < MaxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * RetryDelay)
				continue
			}
			return e.handleWindowError(order, window, fmt.Errorf("ethereum escrow creation failed: %w", err))
		}
		ethTxHash = hash
		window.EthereumTxHash = hash.Hex()

		// Wait for Ethereum confirmation
		confirmed, err := e.waitForEthereumConfirmation(ctx, ethTxHash)
		if err != nil || !confirmed {
			e.logger.Warn("Ethereum confirmation failed",
				zap.Int("attempt", attempt+1),
				zap.Error(err))
			if attempt < MaxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * RetryDelay)
				continue
			}
			return e.handleWindowError(order, window, fmt.Errorf("ethereum confirmation failed: %w", err))
		}

		// Create Cosmos escrow
		cosmosTxHash, err = e.createCosmosEscrow(ctx, order, window)
		if err != nil {
			e.logger.Warn("Failed to create Cosmos escrow",
				zap.Int("attempt", attempt+1),
				zap.Error(err))
			if attempt < MaxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * RetryDelay)
				continue
			}
			return e.handleWindowError(order, window, fmt.Errorf("cosmos escrow creation failed: %w", err))
		}
		window.CosmosTxHash = cosmosTxHash

		// Wait for Cosmos confirmation
		confirmed, err = e.waitForCosmosConfirmation(ctx, cosmosTxHash)
		if err != nil || !confirmed {
			e.logger.Warn("Cosmos confirmation failed",
				zap.Int("attempt", attempt+1),
				zap.Error(err))
			if attempt < MaxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * RetryDelay)
				continue
			}
			return e.handleWindowError(order, window, fmt.Errorf("cosmos confirmation failed: %w", err))
		}

		// Execute swap
		actualPrice, gasUsed, err = e.executeSwap(ctx, order, window)
		if err != nil {
			e.circuitBreaker.RecordFailure()
			e.logger.Warn("Swap execution failed",
				zap.Int("attempt", attempt+1),
				zap.Error(err))
			if attempt < MaxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * RetryDelay)
				continue
			}
			return e.handleWindowError(order, window, fmt.Errorf("swap execution failed: %w", err))
		}

		// Success - break out of retry loop
		break
	}

	// Reveal secret
	if err := e.revealSecret(ctx, order, window); err != nil {
		e.logger.Warn("Failed to reveal secret, but swap completed",
			zap.String("order_id", order.ID),
			zap.Int("window", window.Index),
			zap.Error(err))
	}

	// Update window status
	window.Status = WindowStatusCompleted
	window.ActualPrice = actualPrice
	window.GasUsed = gasUsed

	// Update order
	e.ordersMutex.Lock()
	order.CompletedWindows++
	order.TotalExecuted.Add(order.TotalExecuted, window.Amount)
	e.ordersMutex.Unlock()

	// Persist execution data
	e.storeExecutionData(ctx, order.ID, window)

	e.logger.Info("TWAP window completed successfully",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("eth_tx", window.EthereumTxHash),
		zap.String("cosmos_tx", window.CosmosTxHash),
		zap.String("actual_price", actualPrice.String()),
		zap.Uint64("gas_used", gasUsed),
		zap.Duration("total_duration", time.Since(executionStart)),
	)

	e.circuitBreaker.RecordSuccess()
	e.metrics.RecordWindowExecution(order.ID, window.Index, time.Since(executionStart), window.GasUsed)
	return nil
}

// The implementation continues with more methods...
// [Methods for createEthereumEscrow, waitForEthereumConfirmation, createCosmosEscrow, etc.]
// [Circuit breaker implementation]
// [1inch client methods]
// [Helper methods]

// Circuit Breaker implementation
func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:    threshold,
		resetTimeout: resetTimeout,
	}
}

func (cb *CircuitBreaker) IsTripped() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	
	if cb.isTripped {
		if time.Since(cb.lastTrip) > cb.resetTimeout {
			cb.isTripped = false
			cb.failureCount = 0
			return false
		}
		return true
	}
	return false
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.failureCount++
	if cb.failureCount >= cb.threshold {
		cb.isTripped = true
		cb.lastTrip = time.Now()
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	
	cb.failureCount = 0
	cb.isTripped = false
}

// OneInch client methods
func (c *OneInchClient) GetQuote(ctx context.Context, chainID, src, dst, amount string) (*OneInchQuoteResponse, error) {
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }

    endpoint := fmt.Sprintf("%s%s/%s/quote", OneInchAPIBaseURL, OneInchQuoteAPI, chainID)
    
    params := url.Values{}
    params.Set("src", src)
    params.Set("dst", dst)
    params.Set("amount", amount)
    
    bodyBytes, err := c.makeRequest(ctx, "GET", endpoint+"?"+params.Encode(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get quote: %w", err)
    }

    var quote OneInchQuoteResponse
    if err := json.Unmarshal(bodyBytes, &quote); err != nil {
        return nil, fmt.Errorf("failed to decode quote response: %w", err)
    }

    return &quote, nil
}

func (c *OneInchClient) GetSwap(ctx context.Context, chainID, src, dst, amount, from string, slippage float64) (*OneInchSwapResponse, error) {
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, err
    }

    endpoint := fmt.Sprintf("%s%s/%s/swap", OneInchAPIBaseURL, OneInchSwapAPI, chainID)
    
    params := url.Values{}
    params.Set("src", src)
    params.Set("dst", dst)
    params.Set("amount", amount)
    params.Set("from", from)
    params.Set("slippage", strconv.FormatFloat(slippage, 'f', 2, 64))
    
    bodyBytes, err := c.makeRequest(ctx, "GET", endpoint+"?"+params.Encode(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get swap: %w", err)
    }

    var swap OneInchSwapResponse
    if err := json.Unmarshal(bodyBytes, &swap); err != nil {
        return nil, fmt.Errorf("failed to decode swap response: %w", err)
    }

    return &swap, nil
}


func (c *OneInchClient) makeRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
    var lastErr error
    for i := 0; i < MaxRetries; i++ {
        req, err := http.NewRequestWithContext(ctx, method, url, body)
        if err != nil {
            return nil, fmt.Errorf("failed to create request: %w", err)
        }

        req.Header.Set("Authorization", "Bearer "+c.apiKey)
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("User-Agent", "FlowFusion/1.0")

        resp, err := c.httpClient.Do(req)
        if err != nil {
            lastErr = fmt.Errorf("request failed: %w", err)
            if i < MaxRetries-1 {
                time.Sleep(RetryDelay)
            }
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            bodyBytes, _ := io.ReadAll(resp.Body)
            lastErr = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
            if i < MaxRetries-1 {
                time.Sleep(RetryDelay)
            }
            continue
        }

        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            lastErr = fmt.Errorf("failed to read response body: %w", err)
            if i < MaxRetries-1 {
                time.Sleep(RetryDelay)
            }
            continue
        }

        return bodyBytes, nil
    }
    return nil, fmt.Errorf("failed after %d retries: %w", MaxRetries, lastErr)
}

// createEthereumEscrow creates an Ethereum escrow transaction
func (e *Engine) createEthereumEscrow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) (common.Hash, error) {
	gasPrice, err := e.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %w", err)
	}

	if e.config.MaxGasPrice != nil && gasPrice.Cmp(e.config.MaxGasPrice) > 0 {
		gasPrice = e.config.MaxGasPrice
	}

	params := ethereum.CreateEscrowParams{
		Recipient:  order.UserAddress,
		Amount:     window.Amount,
		SecretHash: window.SecretHash,
		TimeLimit:  uint64(time.Now().Add(24 * time.Hour).Unix()),
	}

	txHash, err := e.ethClient.CreateEscrow(ctx, params)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create Ethereum escrow: %w", err)
	}

	return txHash, nil
}

// waitForEthereumConfirmation waits for Ethereum transaction confirmation
func (e *Engine) waitForEthereumConfirmation(ctx context.Context, txHash common.Hash) (bool, error) {
	timeout := time.After(10 * time.Minute)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for Ethereum confirmation")
		case <-ticker.C:
			receipt, err := e.ethClient.GetTransactionReceipt(ctx, txHash)
			if err != nil {
				continue // Transaction not yet mined
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				currentBlock, err := e.ethClient.GetCurrentBlockNumber(ctx)
				if err != nil {
					continue
				}
				
				if currentBlock-receipt.BlockNumber.Uint64() >= ConfirmationBlocks {
					return true, nil
				}
			} else {
				return false, fmt.Errorf("transaction failed with status: %d", receipt.Status)
			}
		}
	}
}

// createCosmosEscrow creates a Cosmos escrow transaction
func (e *Engine) createCosmosEscrow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) (string, error) {
	params := cosmos.CreateEscrowParams{
		EthereumTxHash: window.EthereumTxHash,
		SecretHash:     window.SecretHash,
		TimeLock:       uint64(time.Now().Add(24 * time.Hour).Unix()),
		Recipient:      order.UserAddress,
		EthereumSender: e.ethClient.GetAddress().Hex(),
		Amount:         window.Amount,
	}

	txHash, err := e.cosmosClient.CreateEscrow(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create Cosmos escrow: %w", err)
	}

	return txHash, nil
}

// waitForCosmosConfirmation waits for Cosmos transaction confirmation
func (e *Engine) waitForCosmosConfirmation(ctx context.Context, txHash string) (bool, error) {
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for Cosmos confirmation")
		case <-ticker.C:
			txResp, err := e.cosmosClient.GetTransactionStatus(ctx, txHash)
			if err != nil {
				continue // Transaction not yet processed
			}
			
			if txResp.Code == 0 {
				return true, nil
			} else {
				return false, fmt.Errorf("transaction failed with code %d: %s", txResp.Code, txResp.RawLog)
			}
		}
	}
}

// executeSwap executes the 1inch swap transaction
func (e *Engine) executeSwap(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) (*big.Int, uint64, error) {
	swapData, err := e.oneInchClient.GetSwap(ctx, e.chainID.String(), order.SourceToken, order.DestToken, 
		window.Amount.String(), e.ethClient.GetAddress().Hex(), 1.0) // 1% slippage
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get swap data: %w", err)
	}

	actualOutput, ok := new(big.Int).SetString(swapData.DstAmount, 10)
	if !ok {
		return nil, 0, fmt.Errorf("invalid swap output amount: %s", swapData.DstAmount)
	}

	// Execute the swap transaction through Ethereum client
	txData := common.FromHex(swapData.Tx.Data)
	value, _ := new(big.Int).SetString(swapData.Tx.Value, 10)
	gasPrice, _ := new(big.Int).SetString(swapData.Tx.GasPrice, 10)

	tx := types.NewTransaction(
		0, 
		common.HexToAddress(swapData.Tx.To),
		value,
		uint64(swapData.Tx.Gas),
		gasPrice,
		txData,
	)

	signedTx, err := e.ethClient.SignTransaction(ctx, tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to sign swap transaction: %w", err)
	}

	err = e.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send swap transaction: %w", err)
	}

	// Wait for transaction confirmation
	confirmed, err := e.waitForEthereumConfirmation(ctx, signedTx.Hash())
	if err != nil || !confirmed {
		return nil, 0, fmt.Errorf("swap transaction confirmation failed: %w", err)
	}

	return actualOutput, uint64(swapData.Tx.Gas), nil
}

// revealSecret reveals the secret to complete the atomic swap
func (e *Engine) revealSecret(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	secretHex := hex.EncodeToString(window.Secret)
	txHash, err := e.ethClient.RevealSecret(ctx, common.HexToAddress(order.UserAddress), secretHex)
	if err != nil {
		return fmt.Errorf("failed to reveal secret: %w", err)
	}

	confirmed, err := e.waitForEthereumConfirmation(ctx, txHash)
	if err != nil || !confirmed {
		return fmt.Errorf("secret revelation confirmation failed: %w", err)
	}

	return nil
}

// handleWindowError handles window execution errors
func (e *Engine) handleWindowError(order *TWAPOrder, window *ExecutionWindow, err error) error {
	window.Status = WindowStatusFailed
	
	e.logger.Error("TWAP window execution failed",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.Error(err),
	)

	errorData := map[string]interface{}{
		"order_id":    order.ID,
		"window":      window.Index,
		"error":       err.Error(),
		"timestamp":   time.Now(),
		"eth_tx":      window.EthereumTxHash,
		"cosmos_tx":   window.CosmosTxHash,
	}
	
	errorBytes, _ := json.Marshal(errorData)
	e.redisClient.LPush(context.Background(), "twap:errors", errorBytes)
	e.metrics.RecordWindowFailure(order.ID, window.Index)

	return err
}

// storeExecutionData persists execution data to Redis
func (e *Engine) storeExecutionData(ctx context.Context, orderID string, window *ExecutionWindow) {
	key := fmt.Sprintf("twap:execution:%s:%d", orderID, window.Index)
	data := map[string]interface{}{
		"order_id":      orderID,
		"window_index":  window.Index,
		"amount":        window.Amount.String(),
		"eth_tx":        window.EthereumTxHash,
		"cosmos_tx":     window.CosmosTxHash,
		"actual_price":  window.ActualPrice.String(),
		"gas_used":      window.GasUsed,
		"executed_at":   window.ExecutedAt,
		"secret_hash":   window.SecretHash,
		"status":        string(window.Status),
	}
	
	dataBytes, _ := json.Marshal(data)
	e.redisClient.Set(ctx, key, dataBytes, OrderTTL)
}

// persistOrder saves a single order to Redis
func (e *Engine) persistOrder(ctx context.Context, order *TWAPOrder) error {
	key := fmt.Sprintf("twap:order:%s", order.ID)
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	
	return e.redisClient.Set(ctx, key, data, OrderTTL).Err()
}

// persistOrders saves all orders to Redis
func (e *Engine) persistOrders(ctx context.Context) error {
	e.ordersMutex.RLock()
	defer e.ordersMutex.RUnlock()

	for _, order := range e.orders {
		if err := e.persistOrder(ctx, order); err != nil {
			e.logger.Error("Failed to persist order", 
				zap.String("order_id", order.ID), 
				zap.Error(err))
		}
	}
	return nil
}

// restoreOrders loads orders from Redis
func (e *Engine) restoreOrders(ctx context.Context) error {
	keys, err := e.redisClient.Keys(ctx, "twap:order:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get order keys: %w", err)
	}

	for _, key := range keys {
		data, err := e.redisClient.Get(ctx, key).Bytes()
		if err != nil {
			e.logger.Warn("Failed to restore order", zap.String("key", key), zap.Error(err))
			continue
		}

		var order TWAPOrder
		if err := json.Unmarshal(data, &order); err != nil {
			e.logger.Warn("Failed to unmarshal order", zap.String("key", key), zap.Error(err))
			continue
		}

		// Only restore active orders
		if order.Status == OrderStatusPending || order.Status == OrderStatusExecuting {
			e.ordersMutex.Lock()
			e.orders[order.ID] = &order
			e.ordersMutex.Unlock()
			
			// Re-schedule the order
			e.scheduler.ScheduleOrder(&order)
		}
	}
	
	e.logger.Info("Restored orders from persistence", zap.Int("count", len(keys)))
	return nil
}

// monitorOrders monitors order health and handles timeouts
func (e *Engine) monitorOrders(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := e.checkOrderHealth(ctx); err != nil {
				e.logger.Error("Error checking order health", zap.Error(err))
			}
		}
	}
}

// checkOrderHealth checks for stuck or timed out orders
func (e *Engine) checkOrderHealth(ctx context.Context) error {
	e.ordersMutex.RLock()
	orders := make([]*TWAPOrder, 0, len(e.orders))
	for _, order := range e.orders {
		if order.Status == OrderStatusExecuting {
			orders = append(orders, order)
		}
	}
	e.ordersMutex.RUnlock()

	for _, order := range orders {
		// Check for order timeout
		if time.Since(order.CreatedAt) > OrderTTL {
			e.handleOrderTimeout(order)
			continue
		}

		// Check for stuck windows
		for _, window := range order.ExecutionPlan {
			if window.Status == WindowStatusExecuting && window.ExecutedAt != nil {
				if time.Since(*window.ExecutedAt) > ExecutionTimeout {
					e.handleWindowError(order, window, fmt.Errorf("window execution timeout"))
				}
			}
		}
	}
	return nil
}

// handleOrderTimeout handles expired orders
func (e *Engine) handleOrderTimeout(order *TWAPOrder) {
	e.ordersMutex.Lock()
	order.Status = OrderStatusFailed
	e.ordersMutex.Unlock()

	e.logger.Warn("Order timed out",
		zap.String("order_id", order.ID),
		zap.Duration("age", time.Since(order.CreatedAt)),
	)

	e.metrics.RecordOrderCompletion(order.ID, false)
}

// isOrderComplete checks if an order has completed execution
func (e *Engine) isOrderComplete(order *TWAPOrder) bool {
	return order.CompletedWindows >= order.IntervalCount
}

// calculatePriceImpact calculates price impact percentage
func (e *Engine) calculatePriceImpact(inputAmount, outputAmount *big.Int) float64 {
	if inputAmount.Sign() == 0 {
		return 0
	}

	// Assuming 1:1 base rate for simplification
	expectedOutput := inputAmount
	if outputAmount.Cmp(expectedOutput) >= 0 {
		return 0 // No negative impact
	}

	impact := new(big.Int).Sub(expectedOutput, outputAmount)
	impact.Mul(impact, big.NewInt(10000))
	impact.Div(impact, expectedOutput)
	
	return float64(impact.Int64()) / 100.0
}

// generateSchedulePreview creates execution schedule preview
func (e *Engine) generateSchedulePreview(req TWAPQuoteRequest, totalAmount *big.Int) []SchedulePreview {
	schedule := make([]SchedulePreview, req.IntervalCount)
	intervalDuration := time.Duration(req.TimeWindow) * time.Second / time.Duration(req.IntervalCount)
	intervalAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	
	startTime := time.Now().Add(1 * time.Minute)
	for i := 0; i < req.IntervalCount; i++ {
		schedule[i] = SchedulePreview{
			IntervalIndex:  i,
			ExecutionTime:  startTime.Add(time.Duration(i) * intervalDuration),
			Amount:         intervalAmount.String(),
			EstimatedPrice: "0", // Can't predict future prices accurately
		}
	}
	return schedule
}

// calculateFees calculates total execution fees
func (e *Engine) calculateFees(req TWAPQuoteRequest, gasEstimate, gasPrice string) (FeeBreakdown, error) {
	gas, err := strconv.ParseUint(gasEstimate, 10, 64)
	if err != nil {
		return FeeBreakdown{}, fmt.Errorf("invalid gas estimate: %w", err)
	}
	
	price, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return FeeBreakdown{}, fmt.Errorf("invalid gas price: %s", gasPrice)
	}

	ethGasCost := new(big.Int).Mul(big.NewInt(int64(gas)), price)
	ethGasCost.Mul(ethGasCost, big.NewInt(int64(req.IntervalCount)))

	relayerFee := new(big.Int).Div(ethGasCost, big.NewInt(100)) // 1% relayer fee
	oneInchFee := new(big.Int).Div(ethGasCost, big.NewInt(333)) // ~0.3% 1inch fee
	cosmosGas := big.NewInt(1000) // Fixed Cosmos gas cost
	
	total := new(big.Int).Add(ethGasCost, relayerFee)
	total.Add(total, oneInchFee)
	total.Add(total, cosmosGas)

	return FeeBreakdown{
		EthereumGas: ethGasCost.String(),
		CosmosGas:   cosmosGas.String(),
		RelayerFee:  relayerFee.String(),
		OneInchFee:  oneInchFee.String(),
		Total:       total.String(),
	}, nil
}

// generateRecommendations provides execution recommendations
func (e *Engine) generateRecommendations(instantImpact, twapImpact float64, intervalCount, timeWindow int) []string {
	recommendations := []string{}
	
	if twapImpact < instantImpact {
		improvement := ((instantImpact - twapImpact) / instantImpact) * 100
		recommendations = append(recommendations, fmt.Sprintf("TWAP execution reduces price impact by %.1f%%", improvement))
	}
	
	if intervalCount < 10 {
		recommendations = append(recommendations, "Consider increasing interval count for better price improvement")
	}
	
	if timeWindow < 3600 {
		recommendations = append(recommendations, "Extending time window may provide better execution prices")
	}
	
	recommendations = append(recommendations, "Monitor market volatility during execution window")
	recommendations = append(recommendations, "Ensure sufficient gas balance for all transactions")
	
	return recommendations
}