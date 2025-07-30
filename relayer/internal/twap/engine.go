package twap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/ethereum"
)

type Engine struct {
	config       Config
	ethClient    *ethereum.Client
	cosmosClient *cosmos.Client
	logger       *zap.Logger
	
	// State management
	orders       map[string]*TWAPOrder
	ordersMutex  sync.RWMutex
	redisClient  *redis.Client
	
	// Execution scheduler
	scheduler    *ExecutionScheduler
	
	// Metrics
	metrics      *Metrics
}

type Config struct {
	OneInchAPIKey   string        `json:"oneinch_api_key"`
	MaxIntervals    int           `json:"max_intervals"`
	MinIntervalTime time.Duration `json:"min_interval_time"`
	MaxSlippage     float64       `json:"max_slippage"`
	RedisURL        string        `json:"redis_url"`
	DatabaseURL     string        `json:"database_url"`
}

type TWAPOrder struct {
	ID               string              `json:"id"`
	UserAddress      string              `json:"user_address"`
	SourceChain      string              `json:"source_chain"`
	DestChain        string              `json:"dest_chain"`
	SourceToken      string              `json:"source_token"`
	DestToken        string              `json:"dest_token"`
	TotalAmount      *big.Int            `json:"total_amount"`
	TimeWindow       time.Duration       `json:"time_window"`
	IntervalCount    int                 `json:"interval_count"`
	MaxSlippage      float64             `json:"max_slippage"`
	Status           OrderStatus         `json:"status"`
	CreatedAt        time.Time           `json:"created_at"`
	StartTime        time.Time           `json:"start_time"`
	ExecutionPlan    []ExecutionWindow   `json:"execution_plan"`
	CompletedWindows int                 `json:"completed_windows"`
	TotalExecuted    *big.Int            `json:"total_executed"`
	Error            string              `json:"error,omitempty"`
}

type ExecutionWindow struct {
	Index            int           `json:"index"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
	Amount           *big.Int      `json:"amount"`
	MaxPrice         *big.Int      `json:"max_price"`
	Status           WindowStatus  `json:"status"`
	EthereumTxHash   string        `json:"ethereum_tx_hash,omitempty"`
	CosmosTxHash     string        `json:"cosmos_tx_hash,omitempty"`
	Secret           []byte        `json:"-"` // Never expose in JSON
	SecretHash       string        `json:"secret_hash"`
	ExecutedAt       *time.Time    `json:"executed_at,omitempty"`
	ActualPrice      *big.Int      `json:"actual_price,omitempty"`
	GasUsed          uint64        `json:"gas_used,omitempty"`
}

type OrderStatus string
const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusExecuting OrderStatus = "executing"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusFailed    OrderStatus = "failed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type WindowStatus string
const (
	WindowStatusPending   WindowStatus = "pending"
	WindowStatusExecuting WindowStatus = "executing"
	WindowStatusCompleted WindowStatus = "completed"
	WindowStatusFailed    WindowStatus = "failed"
	WindowStatusSkipped   WindowStatus = "skipped"
)

// Request/Response types
type CreateTWAPOrderRequest struct {
	UserAddress      string        `json:"user_address" binding:"required"`
	SourceChain      string        `json:"source_chain" binding:"required"`
	DestChain        string        `json:"dest_chain" binding:"required"`
	SourceToken      string        `json:"source_token" binding:"required"`
	DestToken        string        `json:"dest_token" binding:"required"`
	TotalAmount      string        `json:"total_amount" binding:"required"`
	TimeWindow       int           `json:"time_window" binding:"required,min=300"`     // Min 5 minutes
	IntervalCount    int           `json:"interval_count" binding:"required,min=2"`
	MaxSlippage      float64       `json:"max_slippage" binding:"required,min=0.01,max=10"`
	StartDelay       int           `json:"start_delay,omitempty"` // Seconds to wait before starting
}

type TWAPQuoteRequest struct {
	SourceChain      string        `json:"source_chain" binding:"required"`
	DestChain        string        `json:"dest_chain" binding:"required"`
	SourceToken      string        `json:"source_token" binding:"required"`
	DestToken        string        `json:"dest_token" binding:"required"`
	TotalAmount      string        `json:"total_amount" binding:"required"`
	TimeWindow       int           `json:"time_window" binding:"required"`
	IntervalCount    int           `json:"interval_count" binding:"required"`
}

type TWAPQuoteResponse struct {
	EstimatedOutput     string                 `json:"estimated_output"`
	PriceImpact         float64                `json:"price_impact"`
	ImpactComparison    PriceImpactComparison  `json:"impact_comparison"`
	ExecutionSchedule   []SchedulePreview      `json:"execution_schedule"`
	Fees               FeeBreakdown           `json:"fees"`
	Recommendations    []string               `json:"recommendations"`
}

type PriceImpactComparison struct {
	InstantSwap        float64 `json:"instant_swap"`
	TWAPExecution      float64 `json:"twap_execution"`
	ImprovementPercent float64 `json:"improvement_percent"`
}

type SchedulePreview struct {
	IntervalIndex     int       `json:"interval_index"`
	ExecutionTime     time.Time `json:"execution_time"`
	Amount            string    `json:"amount"`
	EstimatedPrice    string    `json:"estimated_price"`
}

type FeeBreakdown struct {
	EthereumGas       string `json:"ethereum_gas"`
	CosmosGas         string `json:"cosmos_gas"`
	RelayerFee        string `json:"relayer_fee"`
	OneInchFee        string `json:"oneinch_fee"`
	Total             string `json:"total"`
}

func NewEngine(config Config, ethClient *ethereum.Client, cosmosClient *cosmos.Client, logger *zap.Logger) (*Engine, error) {
	// Initialize Redis client
	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	
	redisClient := redis.NewClient(opt)
	
	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	scheduler := NewExecutionScheduler(logger)
	metrics := NewMetrics()

	engine := &Engine{
		config:       config,
		ethClient:    ethClient,
		cosmosClient: cosmosClient,
		logger:       logger,
		orders:       make(map[string]*TWAPOrder),
		redisClient:  redisClient,
		scheduler:    scheduler,
		metrics:      metrics,
	}

	return engine, nil
}

func (e *Engine) Start(ctx context.Context) error {
	e.logger.Info("Starting TWAP Engine")

	// Start execution scheduler
	go e.scheduler.Start(ctx)

	// Main execution loop
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			e.logger.Info("TWAP Engine stopping")
			return ctx.Err()
		case <-ticker.C:
			e.processExecutions(ctx)
		}
	}
}

func (e *Engine) GetQuote(c *gin.Context) {
	var req TWAPQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse amount
	totalAmount, ok := new(big.Int).SetString(req.TotalAmount, 10)
	if !ok || totalAmount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid total_amount"})
		return
	}

	// Calculate TWAP quote
	quote, err := e.calculateTWAPQuote(req, totalAmount)
	if err != nil {
		e.logger.Error("Failed to calculate TWAP quote", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate quote"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (e *Engine) CreateOrder(c *gin.Context) {
	var req CreateTWAPOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse total amount
	totalAmount, ok := new(big.Int).SetString(req.TotalAmount, 10)
	if !ok || totalAmount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid total_amount"})
		return
	}

	// Validate configuration
	if err := e.validateTWAPRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create TWAP order
	order, err := e.createTWAPOrder(req, totalAmount)
	if err != nil {
		e.logger.Error("Failed to create TWAP order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (e *Engine) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (e *Engine) GetSchedule(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.RLock()
	order, exists := e.orders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":        order.ID,
		"execution_plan":  order.ExecutionPlan,
		"completed_windows": order.CompletedWindows,
		"status":          order.Status,
	})
}

func (e *Engine) ListOrders(c *gin.Context) {
	e.ordersMutex.RLock()
	orders := make([]*TWAPOrder, 0, len(e.orders))
	for _, order := range e.orders {
		orders = append(orders, order)
	}
	e.ordersMutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (e *Engine) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.Lock()
	order, exists := e.orders[orderID]
	if exists && order.Status == OrderStatusPending {
		order.Status = OrderStatusCancelled
	}
	e.ordersMutex.Unlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

func (e *Engine) GetStatus(c *gin.Context) {
	e.ordersMutex.RLock()
	orderCount := len(e.orders)
	e.ordersMutex.RUnlock()

	status := gin.H{
		"status":      "running",
		"order_count": orderCount,
		"metrics":     e.metrics.GetSummary(),
		"uptime":      time.Since(e.metrics.StartTime).String(),
	}

	c.JSON(http.StatusOK, status)
}

// Private methods

func (e *Engine) validateTWAPRequest(req CreateTWAPOrderRequest) error {
	if req.IntervalCount > e.config.MaxIntervals {
		return fmt.Errorf("interval_count exceeds maximum of %d", e.config.MaxIntervals)
	}

	intervalDuration := time.Duration(req.TimeWindow) * time.Second / time.Duration(req.IntervalCount)
	if intervalDuration < e.config.MinIntervalTime {
		return fmt.Errorf("interval duration too short, minimum is %v", e.config.MinIntervalTime)
	}

	if req.MaxSlippage > e.config.MaxSlippage {
		return fmt.Errorf("max_slippage exceeds maximum of %.2f%%", e.config.MaxSlippage*100)
	}

	return nil
}

func (e *Engine) createTWAPOrder(req CreateTWAPOrderRequest, totalAmount *big.Int) (*TWAPOrder, error) {
	orderID := e.generateOrderID()
	
	// Calculate start time
	startTime := time.Now().Add(time.Duration(req.StartDelay) * time.Second)
	if req.StartDelay == 0 {
		startTime = startTime.Add(1 * time.Minute) // Default 1 minute delay
	}

	// Generate execution plan
	executionPlan, err := e.generateExecutionPlan(req, totalAmount, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate execution plan: %w", err)
	}

	order := &TWAPOrder{
		ID:               orderID,
		UserAddress:      req.UserAddress,
		SourceChain:      req.SourceChain,
		DestChain:        req.DestChain,
		SourceToken:      req.SourceToken,
		DestToken:        req.DestToken,
		TotalAmount:      totalAmount,
		TimeWindow:       time.Duration(req.TimeWindow) * time.Second,
		IntervalCount:    req.IntervalCount,
		MaxSlippage:      req.MaxSlippage,
		Status:           OrderStatusPending,
		CreatedAt:        time.Now(),
		StartTime:        startTime,
		ExecutionPlan:    executionPlan,
		CompletedWindows: 0,
		TotalExecuted:    big.NewInt(0),
	}

	// Store order
	e.ordersMutex.Lock()
	e.orders[orderID] = order
	e.ordersMutex.Unlock()

	// Schedule execution
	e.scheduler.ScheduleOrder(order)

	e.logger.Info("TWAP order created",
		zap.String("order_id", orderID),
		zap.String("user", req.UserAddress),
		zap.Int("intervals", req.IntervalCount),
		zap.String("total_amount", totalAmount.String()),
	)

	return order, nil
}

func (e *Engine) generateOrderID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (e *Engine) generateExecutionPlan(req CreateTWAPOrderRequest, totalAmount *big.Int, startTime time.Time) ([]ExecutionWindow, error) {
	windows := make([]ExecutionWindow, req.IntervalCount)
	intervalDuration := time.Duration(req.TimeWindow) * time.Second / time.Duration(req.IntervalCount)
	
	// Calculate base amount per window
	baseAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	remainder := new(big.Int).Mod(totalAmount, big.NewInt(int64(req.IntervalCount)))

	for i := 0; i < req.IntervalCount; i++ {
		windowStart := startTime.Add(time.Duration(i) * intervalDuration)
		windowEnd := windowStart.Add(intervalDuration)

		// Distribute remainder across first few windows
		windowAmount := new(big.Int).Set(baseAmount)
		if i < int(remainder.Int64()) {
			windowAmount.Add(windowAmount, big.NewInt(1))
		}

		// Generate secret for this window
		secret := make([]byte, 32)
		rand.Read(secret)
		secretHash := e.hashSecret(secret)

		// Calculate max price (this would use 1inch API in production)
		maxPrice, err := e.calculateMaxPrice(req.SourceToken, req.DestToken, windowAmount, req.MaxSlippage)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate max price for window %d: %w", i, err)
		}

		windows[i] = ExecutionWindow{
			Index:      i,
			StartTime:  windowStart,
			EndTime:    windowEnd,
			Amount:     windowAmount,
			MaxPrice:   maxPrice,
			Status:     WindowStatusPending,
			Secret:     secret,
			SecretHash: secretHash,
		}
	}

	return windows, nil
}

func (e *Engine) hashSecret(secret []byte) string {
	// Simple SHA256 hash - in production use the same hash as smart contracts
	hash := make([]byte, 32)
	copy(hash, secret) // Simplified for demo
	return hex.EncodeToString(hash)
}

func (e *Engine) calculateMaxPrice(sourceToken, destToken string, amount *big.Int, maxSlippage float64) (*big.Int, error) {
	// This would integrate with 1inch API to get real market prices
	// For demo purposes, return a placeholder value
	basePrice := big.NewInt(1000000) // 1:1 ratio with 6 decimals
	slippageMultiplier := big.NewInt(int64((1 + maxSlippage) * 1000000))
	maxPrice := new(big.Int).Mul(basePrice, slippageMultiplier)
	maxPrice.Div(maxPrice, big.NewInt(1000000))
	return maxPrice, nil
}

func (e *Engine) calculateTWAPQuote(req TWAPQuoteRequest, totalAmount *big.Int) (*TWAPQuoteResponse, error) {
	// This would integrate with 1inch API for real price calculations
	// For demo purposes, return placeholder values
	
	// Simulate better price impact with TWAP
	instantImpact := 0.15  // 0.15% impact for instant swap
	twapImpact := 0.08     // 0.08% impact with TWAP
	improvement := ((instantImpact - twapImpact) / instantImpact) * 100

	// Generate execution schedule preview
	schedule := make([]SchedulePreview, req.IntervalCount)
	intervalDuration := time.Duration(req.TimeWindow) * time.Second / time.Duration(req.IntervalCount)
	intervalAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	
	startTime := time.Now().Add(1 * time.Minute)
	for i := 0; i < req.IntervalCount; i++ {
		schedule[i] = SchedulePreview{
			IntervalIndex:  i,
			ExecutionTime:  startTime.Add(time.Duration(i) * intervalDuration),
			Amount:         intervalAmount.String(),
			EstimatedPrice: "1000000", // Placeholder
		}
	}

	return &TWAPQuoteResponse{
		EstimatedOutput: new(big.Int).Mul(totalAmount, big.NewInt(990000)).String(), // 99% of input
		PriceImpact:     twapImpact,
		ImpactComparison: PriceImpactComparison{
			InstantSwap:        instantImpact,
			TWAPExecution:      twapImpact,
			ImprovementPercent: improvement,
		},
		ExecutionSchedule: schedule,
		Fees: FeeBreakdown{
			EthereumGas: "0.01",
			CosmosGas:   "0.001",
			RelayerFee:  "0.005",
			OneInchFee:  "0.001",
			Total:       "0.0165",
		},
		Recommendations: []string{
			"TWAP execution reduces price impact by 46.7%",
			"Optimal interval count for this amount is 12",
			"Consider extending time window for better rates",
		},
	}, nil
}

func (e *Engine) processExecutions(ctx context.Context) {
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
		e.processOrderExecution(ctx, order, now)
	}
}

func (e *Engine) processOrderExecution(ctx context.Context, order *TWAPOrder, now time.Time) {
	// Check if we need to start execution
	if order.Status == OrderStatusPending && now.After(order.StartTime) {
		order.Status = OrderStatusExecuting
		e.logger.Info("Starting TWAP order execution", zap.String("order_id", order.ID))
	}

	if order.Status != OrderStatusExecuting {
		return
	}

	// Process pending windows
	for i := range order.ExecutionPlan {
		window := &order.ExecutionPlan[i]
		
		if window.Status != WindowStatusPending {
			continue
		}

		// Check if it's time to execute this window
		if now.After(window.StartTime) && now.Before(window.EndTime) {
			go e.executeWindow(ctx, order, window)
		}
	}

	// Check if order is complete
	if e.isOrderComplete(order) {
		order.Status = OrderStatusCompleted
		e.logger.Info("TWAP order completed", zap.String("order_id", order.ID))
	}
}

func (e *Engine) executeWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) {
	window.Status = WindowStatusExecuting
	executionStart := time.Now()
	
	e.logger.Info("Executing TWAP window",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("amount", window.Amount.String()),
	)

	// For demo purposes, simulate successful execution
	// In production, this would:
	// 1. Create Fusion+ order on Ethereum
	// 2. Wait for confirmation
	// 3. Execute Cosmos transaction
	// 4. Reveal secret
	
	time.Sleep(2 * time.Second) // Simulate execution time

	// Update window status
	window.Status = WindowStatusCompleted
	window.ExecutedAt = &executionStart
	window.EthereumTxHash = "0x" + hex.EncodeToString([]byte("demo_eth_tx"))
	window.CosmosTxHash = "demo_cosmos_tx"
	window.ActualPrice = window.MaxPrice // In production, get actual execution price
	window.GasUsed = 150000 // Demo value

	// Update order totals
	order.CompletedWindows++
	order.TotalExecuted.Add(order.TotalExecuted, window.Amount)

	e.logger.Info("TWAP window completed",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("eth_tx", window.EthereumTxHash),
		zap.String("cosmos_tx", window.CosmosTxHash),
	)

	// Update metrics
	e.metrics.RecordWindowExecution(order.ID, window.Index, time.Since(executionStart))
}

func (e *Engine) isOrderComplete(order *TWAPOrder) bool {
	return order.CompletedWindows >= order.IntervalCount
}