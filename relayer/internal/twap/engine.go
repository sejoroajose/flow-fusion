// relayer/internal/twap/engine.go
package twap

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/ethereum"
	"flow-fusion/relayer/internal/oneinch"
)

// ProductionTWAPEngine implements TWAP execution using 1inch HTTP API
type ProductionTWAPEngine struct {
	config            Config
	ethClient         *ethereum.Client
	cosmosClient      *cosmos.Client
	logger            *zap.Logger
	
	// 1inch HTTP client
	oneInchClient     *oneinch.HTTPClient
	
	// State management
	twapOrders        map[string]*TWAPOrder
	ordersMutex       sync.RWMutex
	redisClient       *redis.Client
	
	// Execution management
	scheduler         *TWAPScheduler
	rateLimiter       *rate.Limiter
	metrics           *TWAPMetrics
	
	// Authentication
	privateKey        *ecdsa.PrivateKey
	publicAddress     common.Address
	chainID           int
}

type Config struct {
	RedisURL        string
	OneInchAPIKey   string
	PrivateKey      string
	NodeURL         string
	ChainID         int
	MaxGasPrice     *big.Int
	MaxIntervals    int
	MinIntervalTime time.Duration
	MaxSlippage     float64
	DatabaseURL     string
}

type TWAPOrder struct {
	ID              string                `json:"id"`
	UserAddress     string                `json:"user_address"`
	SourceToken     string                `json:"source_token"`
	DestToken       string                `json:"dest_token"`
	SourceChain     string                `json:"source_chain"`     
	DestChain       string                `json:"dest_chain"` 
	TotalAmount     *big.Int              `json:"total_amount"`
	IntervalCount   int                   `json:"interval_count"`
	TimeWindow      int                   `json:"time_window"`
	MaxSlippage     float64               `json:"max_slippage"`
	Status          TWAPOrderStatus       `json:"status"`
	CreatedAt       time.Time             `json:"created_at"`
	StartTime       time.Time             `json:"start_time"`
	
	// Execution windows
	ExecutionWindows []*ExecutionWindow   `json:"execution_windows"`
	CompletedWindows int                  `json:"completed_windows"`
	TotalExecuted    *big.Int             `json:"total_executed"`
	
	// 1inch integration using HTTP API
	FusionPlusOrders []*FusionPlusOrderRef `json:"fusion_plus_orders"`
	OrderbookOrders  []*OrderbookOrderRef  `json:"orderbook_orders"`
}

type ExecutionWindow struct {
	Index           int                   `json:"index"`
	Amount          *big.Int              `json:"amount"`
	StartTime       time.Time             `json:"start_time"`
	EndTime         time.Time             `json:"end_time"`
	Status          WindowStatus          `json:"status"`
	ExecutedAt      *time.Time            `json:"executed_at,omitempty"`
	
	// 1inch order details
	FusionPlusHash  string                `json:"fusion_plus_hash,omitempty"`
	OrderbookHash   string                `json:"orderbook_hash,omitempty"`
	ActualOutput    *big.Int              `json:"actual_output,omitempty"`
	GasUsed         uint64                `json:"gas_used,omitempty"`
	
	// Secrets for atomic swaps
	Secret          string                `json:"-"` // Never expose
	SecretHash      string                `json:"secret_hash"`
	
	// Cross-chain references
	EthereumTxHash  string                `json:"ethereum_tx_hash,omitempty"`
	CosmosTxHash    string                `json:"cosmos_tx_hash,omitempty"`
}

type FusionPlusOrderRef struct {
	OrderHash       string              `json:"order_hash"`
	WindowIndex     int                 `json:"window_index"`
	Amount          *big.Int            `json:"amount"`
	Status          string              `json:"status"`
	SecretSubmitted bool                `json:"secret_submitted"`
}

type OrderbookOrderRef struct {
	OrderHash       string              `json:"order_hash"`
	WindowIndex     int                 `json:"window_index"`
	Amount          *big.Int            `json:"amount"`
	Status          string              `json:"status"`
	Filled          *big.Int            `json:"filled"`
}

type TWAPOrderStatus string
const (
	TWAPOrderStatusPending   TWAPOrderStatus = "pending"
	TWAPOrderStatusExecuting TWAPOrderStatus = "executing"
	TWAPOrderStatusCompleted TWAPOrderStatus = "completed"
	TWAPOrderStatusFailed    TWAPOrderStatus = "failed"
	TWAPOrderStatusCancelled TWAPOrderStatus = "cancelled"
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
	UserAddress   string  `json:"user_address" binding:"required"`
	SourceToken   string  `json:"source_token" binding:"required"`
	DestToken     string  `json:"dest_token" binding:"required"`
	SourceChain   string  `json:"source_chain,omitempty"`      
	DestChain     string  `json:"dest_chain,omitempty"`        
	TotalAmount   string  `json:"total_amount" binding:"required"`
	TimeWindow    int     `json:"time_window" binding:"required,min=300"`
	IntervalCount int     `json:"interval_count" binding:"required,min=2,max=100"`
	MaxSlippage   float64 `json:"max_slippage" binding:"required,min=0.01,max=10"`
	StartDelay    int     `json:"start_delay,omitempty"`
	UseFusionPlus bool    `json:"use_fusion_plus,omitempty"`
	UseOrderbook  bool    `json:"use_orderbook,omitempty"`
}

type TWAPQuoteResponse struct {
	EstimatedOutput   string                    `json:"estimated_output"`
	PriceImpact       float64                   `json:"price_impact"`
	ImpactComparison  PriceImpactComparison     `json:"impact_comparison"`
	ExecutionSchedule []SchedulePreview         `json:"execution_schedule"`
	Fees              FeeBreakdown              `json:"fees"`
	Recommendations   []string                  `json:"recommendations"`
	FusionPlusQuote   *oneinch.FusionQuoteResponse `json:"fusion_plus_quote,omitempty"`
	OrderbookQuote    *oneinch.SwapResponse        `json:"orderbook_quote,omitempty"`
}

type PriceImpactComparison struct {
	InstantSwap        float64 `json:"instant_swap"`
	TWAPExecution      float64 `json:"twap_execution"`
	ImprovementPercent float64 `json:"improvement_percent"`
}

type SchedulePreview struct {
	IntervalIndex    int       `json:"interval_index"`
	ExecutionTime    time.Time `json:"execution_time"`
	Amount           string    `json:"amount"`
	EstimatedPrice   string    `json:"estimated_price"`
	ExecutionMethod  string    `json:"execution_method"` // "fusion_plus" or "orderbook"
}

type FeeBreakdown struct {
	EthereumGas      string `json:"ethereum_gas"`
	CosmosGas        string `json:"cosmos_gas"`
	RelayerFee       string `json:"relayer_fee"`
	OneInchFee       string `json:"oneinch_fee"`
	Total            string `json:"total"`
}

func NewProductionTWAPEngine(config Config, ethClient *ethereum.Client, cosmosClient *cosmos.Client, logger *zap.Logger) (*ProductionTWAPEngine, error) {
	// Parse private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	
	publicKey := privateKey.Public()
	publicAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

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

	// Initialize 1inch HTTP client
	oneInchConfig := oneinch.Config{
		BaseURL: "https://api.1inch.dev",
		APIKey:  config.OneInchAPIKey,
		ChainID: uint64(config.ChainID),
		Timeout: 30 * time.Second,
	}

	oneInchClient := oneinch.NewHTTPClient(oneInchConfig, logger)

	engine := &ProductionTWAPEngine{
		config:            config,  
		ethClient:         ethClient,
		cosmosClient:      cosmosClient,
		logger:            logger,
		oneInchClient:     oneInchClient,
		twapOrders:        make(map[string]*TWAPOrder),
		redisClient:       redisClient,
		rateLimiter:       rate.NewLimiter(rate.Limit(5), 10), 
		metrics:           NewTWAPMetrics(),
		privateKey:        privateKey,
		publicAddress:     publicAddress,
		chainID:           config.ChainID,
	}

	// Initialize scheduler
	engine.scheduler = NewTWAPScheduler(logger, engine)

	return engine, nil
}

func (e *ProductionTWAPEngine) Start(ctx context.Context) error {
	e.logger.Info("Starting Production TWAP Engine with 1inch HTTP API Integration",
		zap.String("public_address", e.publicAddress.Hex()),
		zap.Int("chain_id", e.chainID),
	)

	// Restore persisted orders
	if err := e.restoreOrders(ctx); err != nil {
		e.logger.Error("Failed to restore TWAP orders", zap.Error(err))
	}

	// Start execution scheduler
	go e.scheduler.Start(ctx)

	// Start order monitoring
	go e.monitorTWAPOrders(ctx)

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
				e.logger.Error("Error processing TWAP executions", zap.Error(err))
			}
		}
	}
}

func (e *ProductionTWAPEngine) GetQuote(c *gin.Context) {
	var req CreateTWAPOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	totalAmount, ok := new(big.Int).SetString(req.TotalAmount, 10)
	if !ok || totalAmount.Sign() <= 0 {
		c.JSON(400, gin.H{"error": "invalid total amount"})
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

func (e *ProductionTWAPEngine) CreateOrder(c *gin.Context) {
	var req CreateTWAPOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := e.validateTWAPRequest(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	totalAmount, ok := new(big.Int).SetString(req.TotalAmount, 10)
	if !ok || totalAmount.Sign() <= 0 {
		c.JSON(400, gin.H{"error": "invalid total amount"})
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

func (e *ProductionTWAPEngine) calculateTWAPQuote(ctx context.Context, req CreateTWAPOrderRequest, totalAmount *big.Int) (*TWAPQuoteResponse, error) {
	// Get instant swap quote for comparison using HTTP API
	instantQuote, err := e.oneInchClient.GetSwap(
		ctx,
		req.SourceToken,
		req.DestToken,
		totalAmount.String(),
		req.UserAddress,
		req.MaxSlippage,
		map[string]string{
			"disableEstimate": "false",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get instant quote: %w", err)
	}

	// Calculate interval amount and get interval quote
	intervalAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	intervalQuote, err := e.oneInchClient.GetSwap(
		ctx,
		req.SourceToken,
		req.DestToken,
		intervalAmount.String(),
		req.UserAddress,
		req.MaxSlippage,
		map[string]string{
			"disableEstimate": "false",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get interval quote: %w", err)
	}

	// Calculate estimated TWAP output
	intervalOutput, ok := new(big.Int).SetString(intervalQuote.DstAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid interval output amount: %s", intervalQuote.DstAmount)
	}
	twapEstimatedOutput := new(big.Int).Mul(intervalOutput, big.NewInt(int64(req.IntervalCount)))

	// Compare price impacts
	instantOutput, ok := new(big.Int).SetString(instantQuote.DstAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid instant output amount: %s", instantQuote.DstAmount)
	}

	instantImpact := e.calculatePriceImpact(totalAmount, instantOutput)
	twapImpact := e.calculatePriceImpact(totalAmount, twapEstimatedOutput)

	improvement := ((instantImpact - twapImpact) / instantImpact) * 100
	if improvement < 0 {
		improvement = 0
	}

	// Get Fusion+ quote if requested
	var fusionPlusQuote *oneinch.FusionQuoteResponse
	if req.UseFusionPlus {
		fusionQuote, err := e.getFusionPlusQuote(ctx, req, intervalAmount)
		if err != nil {
			e.logger.Warn("Failed to get Fusion+ quote", zap.Error(err))
		} else {
			fusionPlusQuote = fusionQuote
		}
	}

	// Generate execution schedule
	schedule := e.generateExecutionSchedule(req, totalAmount)

	// Calculate fees
	var gasEstimate int = 300000 // Default gas estimate
	if instantQuote.EstimatedGas != "" {
		if gasStr := instantQuote.EstimatedGas; gasStr != "" {
			if parsed, parseErr := new(big.Int).SetString(gasStr, 10); parseErr {
				gasEstimate = int(parsed.Int64())
			}
		}
	}
	fees := e.calculateFees(req, gasEstimate)

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
		Recommendations:   e.generateRecommendations(instantImpact, twapImpact, req.IntervalCount),
		FusionPlusQuote:   fusionPlusQuote,
		OrderbookQuote:    intervalQuote,
	}, nil
}

func (e *ProductionTWAPEngine) getFusionPlusQuote(ctx context.Context, req CreateTWAPOrderRequest, intervalAmount *big.Int) (*oneinch.FusionQuoteResponse, error) {
	// Get source and destination chain IDs
	srcChain := uint64(e.chainID)
	dstChain := uint64(1) // Cosmos or target chain

	return e.oneInchClient.GetFusionQuote(
		ctx,
		srcChain,
		dstChain,
		req.SourceToken,
		req.DestToken,
		intervalAmount.String(),
		req.UserAddress,
	)
}

func (e *ProductionTWAPEngine) createTWAPOrder(ctx context.Context, req CreateTWAPOrderRequest, totalAmount *big.Int) (*TWAPOrder, error) {
	orderID := e.generateOrderID()
	
	startTime := time.Now().Add(time.Duration(req.StartDelay) * time.Second)
	if req.StartDelay == 0 {
		startTime = startTime.Add(1 * time.Minute) 
	}

	sourceChain := req.SourceChain
	destChain := req.DestChain
	if sourceChain == "" {
		sourceChain = "ethereum"
	}
	if destChain == "" {
		destChain = "cosmos"
	}

	windows, err := e.generateExecutionWindows(ctx, req, totalAmount, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate execution windows: %w", err)
	}

	order := &TWAPOrder{
		ID:               orderID,
		UserAddress:      req.UserAddress,
		SourceToken:      req.SourceToken,
		DestToken:        req.DestToken,
		SourceChain:      sourceChain,    
		DestChain:        destChain,      
		TotalAmount:      totalAmount,
		IntervalCount:    req.IntervalCount,
		TimeWindow:       req.TimeWindow,
		MaxSlippage:      req.MaxSlippage,
		Status:           TWAPOrderStatusPending,
		CreatedAt:        time.Now(),
		StartTime:        startTime,
		ExecutionWindows: windows,
		CompletedWindows: 0,
		TotalExecuted:    big.NewInt(0),
		FusionPlusOrders: make([]*FusionPlusOrderRef, 0),
		OrderbookOrders:  make([]*OrderbookOrderRef, 0),
	}

	// Store order
	e.ordersMutex.Lock()
	e.twapOrders[orderID] = order
	e.ordersMutex.Unlock()

	// Schedule execution
	e.scheduler.ScheduleOrder(order)

	// Persist order
	if err := e.persistOrder(ctx, order); err != nil {
		e.logger.Warn("Failed to persist TWAP order", zap.String("order_id", orderID), zap.Error(err))
	}

	e.logger.Info("TWAP order created",
		zap.String("order_id", orderID),
		zap.String("user", req.UserAddress),
		zap.String("source_chain", sourceChain),
		zap.String("dest_chain", destChain),
		zap.Int("intervals", req.IntervalCount),
		zap.String("total_amount", totalAmount.String()),
		zap.Bool("fusion_plus", req.UseFusionPlus),
		zap.Bool("orderbook", req.UseOrderbook),
	)

	return order, nil
}

func (e *ProductionTWAPEngine) generateExecutionWindows(ctx context.Context, req CreateTWAPOrderRequest, totalAmount *big.Int, startTime time.Time) ([]*ExecutionWindow, error) {
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

		// Generate secret for atomic swap
		secret, err := e.generateSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate secret for window %d: %w", i, err)
		}
		secretHash := e.hashSecret(secret)

		windows[i] = &ExecutionWindow{
			Index:      i,
			Amount:     windowAmount,
			StartTime:  windowStart,
			EndTime:    windowEnd,
			Status:     WindowStatusPending,
			Secret:     secret,
			SecretHash: secretHash,
		}
	}

	return windows, nil
}

func (e *ProductionTWAPEngine) processExecutions(ctx context.Context) error {
	now := time.Now()
	
	e.ordersMutex.RLock()
	orders := make([]*TWAPOrder, 0, len(e.twapOrders))
	for _, order := range e.twapOrders {
		if order.Status == TWAPOrderStatusPending || order.Status == TWAPOrderStatusExecuting {
			orders = append(orders, order)
		}
	}
	e.ordersMutex.RUnlock()

	for _, order := range orders {
		// Start order execution if it's time
		if order.Status == TWAPOrderStatusPending && now.After(order.StartTime) {
			e.ordersMutex.Lock()
			order.Status = TWAPOrderStatusExecuting
			e.ordersMutex.Unlock()
		}

		// Process execution windows
		for _, window := range order.ExecutionWindows {
			if window.Status == WindowStatusPending && now.After(window.StartTime) && now.Before(window.EndTime) {
				go e.executeWindow(ctx, order, window)
			}
		}

		// Check if order is complete
		if e.isOrderComplete(order) {
			e.ordersMutex.Lock()
			order.Status = TWAPOrderStatusCompleted
			e.ordersMutex.Unlock()
		}
	}

	return nil
}

func (e *ProductionTWAPEngine) executeWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) {
	if err := e.rateLimiter.Wait(ctx); err != nil {
		e.logger.Error("Rate limit exceeded", zap.Error(err))
		return
	}

	window.Status = WindowStatusExecuting
	executionStart := time.Now()
	window.ExecutedAt = &executionStart

	e.logger.Info("Executing TWAP window",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("amount", window.Amount.String()),
	)

	// Choose execution method (Fusion+ or Orderbook)
	var err error
	if e.shouldUseFusionPlus(order, window) {
		err = e.executeFusionPlusWindow(ctx, order, window)
	} else {
		err = e.executeOrderbookWindow(ctx, order, window)
	}

	if err != nil {
		window.Status = WindowStatusFailed
		e.logger.Error("TWAP window execution failed",
			zap.String("order_id", order.ID),
			zap.Int("window", window.Index),
			zap.Error(err),
		)
		return
	}

	window.Status = WindowStatusCompleted
	e.ordersMutex.Lock()
	order.CompletedWindows++
	order.TotalExecuted.Add(order.TotalExecuted, window.Amount)
	e.ordersMutex.Unlock()

	e.logger.Info("TWAP window completed",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.Duration("execution_time", time.Since(executionStart)),
	)
}

func (e *ProductionTWAPEngine) executeFusionPlusWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	// Implementation for Fusion+ execution using HTTP API
	e.logger.Info("Executing via Fusion+ HTTP API",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
	)
	
	// Get Fusion+ quote
	srcChain := uint64(e.chainID)
	dstChain := uint64(1) // Target chain
	
	quote, err := e.oneInchClient.GetFusionQuote(
		ctx,
		srcChain,
		dstChain,
		order.SourceToken,
		order.DestToken,
		window.Amount.String(),
		order.UserAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to get Fusion+ quote: %w", err)
	}

	// Create and place Fusion+ order
	orderHash, err := e.generateOrderHash()
	if err != nil {
		return fmt.Errorf("failed to generate order hash: %w", err)
	}

	fusionOrderReq := oneinch.FusionOrderRequest{
		WalletAddress: order.UserAddress,
		OrderHash:     orderHash,
		SecretHashes:  []string{window.SecretHash},
		Receiver:      order.UserAddress,
		Preset:        quote.RecommendedPreset,
	}

	fusionOrderResp, err := e.oneInchClient.PlaceFusionOrder(ctx, fusionOrderReq)
	if err != nil {
		return fmt.Errorf("failed to place Fusion+ order: %w", err)
	}

	if !fusionOrderResp.Success {
		return fmt.Errorf("fusion+ order placement failed: %s", fusionOrderResp.Message)
	}

	window.FusionPlusHash = fusionOrderResp.OrderHash

	// Track the order
	fusionRef := &FusionPlusOrderRef{
		OrderHash:       fusionOrderResp.OrderHash,
		WindowIndex:     window.Index,
		Amount:          window.Amount,
		Status:          "pending",
		SecretSubmitted: false,
	}
	order.FusionPlusOrders = append(order.FusionPlusOrders, fusionRef)

	return nil
}

func (e *ProductionTWAPEngine) executeOrderbookWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	// Implementation for Orderbook execution using HTTP API
	e.logger.Info("Executing via Orderbook HTTP API",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
	)
	
	// This would require implementing orderbook order creation
	// For now, return success as placeholder
	return nil
}

// Helper methods
func (e *ProductionTWAPEngine) shouldUseFusionPlus(order *TWAPOrder, window *ExecutionWindow) bool {
	// Decision logic for choosing execution method
	// Could be based on amount, market conditions, etc.
	return window.Amount.Cmp(big.NewInt(1000000)) > 0 // Use Fusion+ for larger amounts
}

func (e *ProductionTWAPEngine) generateSecret() (string, error) {
	// Implement random secret generation
	return e.oneInchClient.GenerateRandomBytes32()
}

func (e *ProductionTWAPEngine) hashSecret(secret string) string {
	hash, _ := e.oneInchClient.HashSecret(secret)
	return hash
}

func (e *ProductionTWAPEngine) generateOrderHash() (string, error) {
	// Generate a unique order hash
	return fmt.Sprintf("0x%x", time.Now().UnixNano()), nil
}

func (e *ProductionTWAPEngine) calculatePriceImpact(inputAmount, outputAmount *big.Int) float64 {
	if inputAmount.Sign() == 0 {
		return 0
	}
	
	// Simple price impact calculation
	// In production, this should use proper price feeds
	impact := new(big.Int).Sub(inputAmount, outputAmount)
	impact.Mul(impact, big.NewInt(10000))
	impact.Div(impact, inputAmount)
	
	return float64(impact.Int64()) / 100.0
}

func (e *ProductionTWAPEngine) generateExecutionSchedule(req CreateTWAPOrderRequest, totalAmount *big.Int) []SchedulePreview {
	schedule := make([]SchedulePreview, req.IntervalCount)
	intervalDuration := time.Duration(req.TimeWindow) * time.Second / time.Duration(req.IntervalCount)
	intervalAmount := new(big.Int).Div(totalAmount, big.NewInt(int64(req.IntervalCount)))
	
	startTime := time.Now().Add(1 * time.Minute)
	for i := 0; i < req.IntervalCount; i++ {
		method := "orderbook"
		if intervalAmount.Cmp(big.NewInt(1000000)) > 0 {
			method = "fusion_plus"
		}
		
		schedule[i] = SchedulePreview{
			IntervalIndex:   i,
			ExecutionTime:   startTime.Add(time.Duration(i) * intervalDuration),
			Amount:          intervalAmount.String(),
			EstimatedPrice:  "0", // Would be calculated from quotes
			ExecutionMethod: method,
		}
	}
	return schedule
}

func (e *ProductionTWAPEngine) calculateFees(req CreateTWAPOrderRequest, gasEstimate int) FeeBreakdown {
	// Simple fee calculation - implement proper fee estimation
	ethGas := big.NewInt(int64(gasEstimate) * int64(req.IntervalCount))
	relayerFee := new(big.Int).Div(ethGas, big.NewInt(100)) // 1%
	oneInchFee := new(big.Int).Div(ethGas, big.NewInt(333)) // ~0.3%
	cosmosGas := big.NewInt(1000)
	
	total := new(big.Int).Add(ethGas, relayerFee)
	total.Add(total, oneInchFee)
	total.Add(total, cosmosGas)

	return FeeBreakdown{
		EthereumGas: ethGas.String(),
		CosmosGas:   cosmosGas.String(),
		RelayerFee:  relayerFee.String(),
		OneInchFee:  oneInchFee.String(),
		Total:       total.String(),
	}
}

func (e *ProductionTWAPEngine) generateRecommendations(instantImpact, twapImpact float64, intervalCount int) []string {
	recommendations := []string{}
	
	if twapImpact < instantImpact {
		improvement := ((instantImpact - twapImpact) / instantImpact) * 100
		recommendations = append(recommendations, fmt.Sprintf("TWAP execution reduces price impact by %.1f%%", improvement))
	}
	
	if intervalCount < 10 {
		recommendations = append(recommendations, "Consider increasing interval count for better price improvement")
	}
	
	recommendations = append(recommendations, "Monitor market volatility during execution")
	recommendations = append(recommendations, "Use Fusion+ for larger amounts to minimize MEV")
	
	return recommendations
}

func (e *ProductionTWAPEngine) validateTWAPRequest(req CreateTWAPOrderRequest) error {
	if !common.IsHexAddress(req.SourceToken) || !common.IsHexAddress(req.DestToken) {
		return fmt.Errorf("invalid token addresses")
	}
	if req.IntervalCount <= 0 || req.IntervalCount > e.config.MaxIntervals {
		return fmt.Errorf("invalid interval count: %d", req.IntervalCount)
	}
	if req.TimeWindow <= 0 || req.TimeWindow > 86400 {
		return fmt.Errorf("invalid time window: %d", req.TimeWindow)
	}
	if req.MaxSlippage < 0 || req.MaxSlippage > e.config.MaxSlippage {
		return fmt.Errorf("invalid slippage: %f", req.MaxSlippage)
	}
	return nil
}

func (e *ProductionTWAPEngine) generateOrderID() string {
	return fmt.Sprintf("twap_%d_%s", time.Now().UnixNano(), hex.EncodeToString([]byte{byte(len(e.twapOrders))}))
}

func (e *ProductionTWAPEngine) isOrderComplete(order *TWAPOrder) bool {
	return order.CompletedWindows >= order.IntervalCount
}

// Additional methods for GetOrder, ListOrders, CancelOrder, GetStatus, etc.
func (e *ProductionTWAPEngine) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.RLock()
	order, exists := e.twapOrders[orderID]
	e.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	c.JSON(200, order)
}

func (e *ProductionTWAPEngine) ListOrders(c *gin.Context) {
	e.ordersMutex.RLock()
	orders := make([]*TWAPOrder, 0, len(e.twapOrders))
	for _, order := range e.twapOrders {
		orders = append(orders, order)
	}
	e.ordersMutex.RUnlock()

	c.JSON(200, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (e *ProductionTWAPEngine) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	e.ordersMutex.Lock()
	order, exists := e.twapOrders[orderID]
	if exists && order.Status == TWAPOrderStatusPending {
		order.Status = TWAPOrderStatusCancelled
	}
	e.ordersMutex.Unlock()
	
	if !exists {
		c.JSON(404, gin.H{"error": "order not found"})
		return
	}

	c.JSON(200, gin.H{"message": "order cancelled"})
}

func (e *ProductionTWAPEngine) GetStatus() map[string]interface{} {
	e.ordersMutex.RLock()
	orderCount := len(e.twapOrders)
	e.ordersMutex.RUnlock()

	return map[string]interface{}{
		"status":         "running",
		"order_count":    orderCount,
		"public_address": e.publicAddress.Hex(),
		"chain_id":       e.chainID,
		"metrics":        e.metrics.GetSummary(),
	}
}

// Persistence and monitoring methods
func (e *ProductionTWAPEngine) persistOrder(ctx context.Context, order *TWAPOrder) error {
	key := fmt.Sprintf("twap:order:%s", order.ID)
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	
	return e.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

func (e *ProductionTWAPEngine) restoreOrders(ctx context.Context) error {
	keys, err := e.redisClient.Keys(ctx, "twap:order:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get order keys: %w", err)
	}

	for _, key := range keys {
		data, err := e.redisClient.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var order TWAPOrder
		if err := json.Unmarshal(data, &order); err != nil {
			continue
		}

		if order.Status == TWAPOrderStatusPending || order.Status == TWAPOrderStatusExecuting {
			e.ordersMutex.Lock()
			e.twapOrders[order.ID] = &order
			e.ordersMutex.Unlock()
			
			e.scheduler.ScheduleOrder(&order)
		}
	}
	
	e.logger.Info("Restored TWAP orders from persistence", zap.Int("count", len(keys)))
	return nil
}

func (e *ProductionTWAPEngine) monitorTWAPOrders(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Monitor order health and handle timeouts
			e.checkOrderHealth()
		}
	}
}

func (e *ProductionTWAPEngine) checkOrderHealth() {
	// Implementation for order health monitoring
}

func (e *ProductionTWAPEngine) shutdown(ctx context.Context) error {
	e.logger.Info("Shutting down TWAP Engine")
	
	if err := e.persistAllOrders(ctx); err != nil {
		e.logger.Error("Failed to persist orders during shutdown", zap.Error(err))
	}
	
	if err := e.redisClient.Close(); err != nil {
		e.logger.Error("Failed to close Redis connection", zap.Error(err))
	}
	
	return nil
}

func (e *ProductionTWAPEngine) persistAllOrders(ctx context.Context) error {
	e.ordersMutex.RLock()
	defer e.ordersMutex.RUnlock()

	for _, order := range e.twapOrders {
		if err := e.persistOrder(ctx, order); err != nil {
			e.logger.Error("Failed to persist order", 
				zap.String("order_id", order.ID), 
				zap.Error(err))
		}
	}
	return nil
}

type TWAPMetrics struct {
	StartTime time.Time
	mutex     sync.RWMutex
}

func NewTWAPMetrics() *TWAPMetrics {
	return &TWAPMetrics{
		StartTime: time.Now(),
	}
}

func (m *TWAPMetrics) GetSummary() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	return map[string]interface{}{
		"uptime": time.Since(m.StartTime).String(),
	}
}