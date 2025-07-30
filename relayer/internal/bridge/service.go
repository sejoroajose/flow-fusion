package bridge

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/ethereum"
	"flow-fusion/relayer/internal/twap"
)

type Service struct {
	config         Config
	ethClient      *ethereum.Client
	cosmosClient   *cosmos.Client
	twapEngine     *twap.Engine
	logger         *zap.Logger
	
	// State management
	orders         map[string]*BridgeOrder
	ordersMutex    sync.RWMutex
	
	// WebSocket connections
	wsUpgrader     websocket.Upgrader
	wsConnections  map[*websocket.Conn]bool
	wsConnectionsMutex sync.RWMutex
	
	// Metrics
	metrics        *BridgeMetrics
}

type Config struct {
	EthereumClient *ethereum.Client
	CosmosClient   *cosmos.Client
	TWAPEngine     *twap.Engine
	Logger         *zap.Logger
}

type BridgeOrder struct {
	ID              string                 `json:"id"`
	UserAddress     string                 `json:"user_address"`
	SourceChain     string                 `json:"source_chain"`
	DestChain       string                 `json:"dest_chain"`
	SourceToken     string                 `json:"source_token"`
	DestToken       string                 `json:"dest_token"`
	Amount          *big.Int               `json:"amount"`
	Status          BridgeOrderStatus      `json:"status"`
	CreatedAt       time.Time              `json:"created_at"`
	EthereumTxHash  string                 `json:"ethereum_tx_hash,omitempty"`
	CosmosTxHash    string                 `json:"cosmos_tx_hash,omitempty"`
	SecretHash      string                 `json:"secret_hash"`
	Secret          []byte                 `json:"-"` // Never expose
	TimeLock        uint64                 `json:"time_lock"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
	Error           string                 `json:"error,omitempty"`
	TWAPOrderID     string                 `json:"twap_order_id,omitempty"`
}

type BridgeOrderStatus string
const (
	BridgeOrderStatusPending     BridgeOrderStatus = "pending"
	BridgeOrderStatusProcessing  BridgeOrderStatus = "processing"
	BridgeOrderStatusCompleted   BridgeOrderStatus = "completed"
	BridgeOrderStatusFailed      BridgeOrderStatus = "failed"
	BridgeOrderStatusCancelled   BridgeOrderStatus = "cancelled"
)

type BridgeMetrics struct {
	TotalOrders     int64     `json:"total_orders"`
	CompletedOrders int64     `json:"completed_orders"`
	FailedOrders    int64     `json:"failed_orders"`
	TotalVolume     *big.Int  `json:"total_volume"`
	StartTime       time.Time `json:"start_time"`
	mutex           sync.RWMutex
}

// Request/Response types
type CreateBridgeOrderRequest struct {
	UserAddress string `json:"user_address" binding:"required"`
	SourceChain string `json:"source_chain" binding:"required"`
	DestChain   string `json:"dest_chain" binding:"required"`
	SourceToken string `json:"source_token" binding:"required"`
	DestToken   string `json:"dest_token" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	UseTWAP     bool   `json:"use_twap,omitempty"`
	TWAPConfig  *TWAPOrderConfig `json:"twap_config,omitempty"`
}

type TWAPOrderConfig struct {
	TimeWindow    int     `json:"time_window"`
	IntervalCount int     `json:"interval_count"`
	MaxSlippage   float64 `json:"max_slippage"`
}

type BridgeQuoteRequest struct {
	SourceChain string `json:"source_chain" binding:"required"`
	DestChain   string `json:"dest_chain" binding:"required"`
	SourceToken string `json:"source_token" binding:"required"`
	DestToken   string `json:"dest_token" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

type BridgeQuoteResponse struct {
	EstimatedOutput string        `json:"estimated_output"`
	EstimatedFees   string        `json:"estimated_fees"`
	ExecutionTime   time.Duration `json:"execution_time"`
	PriceImpact     float64       `json:"price_impact"`
	Route           []string      `json:"route"`
}

func NewService(config Config) (*Service, error) {
	wsUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for demo
		},
	}

	service := &Service{
		config:         config,
		ethClient:      config.EthereumClient,
		cosmosClient:   config.CosmosClient,
		twapEngine:     config.TWAPEngine,
		logger:         config.Logger,
		orders:         make(map[string]*BridgeOrder),
		wsUpgrader:     wsUpgrader,
		wsConnections:  make(map[*websocket.Conn]bool),
		metrics:        NewBridgeMetrics(),
	}

	return service, nil
}

func (s *Service) Start(ctx context.Context) error {
	s.logger.Info("Starting Bridge Service")

	// Start event watchers
	go s.watchEthereumEvents(ctx)
	go s.watchCosmosEvents(ctx)

	// Main processing loop
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Bridge Service stopping")
			return ctx.Err()
		case <-ticker.C:
			s.processOrders(ctx)
		}
	}
}

func (s *Service) GetQuote(c *gin.Context) {
	var req BridgeQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse amount
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	// Calculate quote
	quote, err := s.calculateQuote(req, amount)
	if err != nil {
		s.logger.Error("Failed to calculate quote", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate quote"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (s *Service) CreateOrder(c *gin.Context) {
	var req CreateBridgeOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse amount
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	// Create bridge order
	order, err := s.createBridgeOrder(req, amount)
	if err != nil {
		s.logger.Error("Failed to create bridge order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// Broadcast to WebSocket clients
	s.broadcastOrderUpdate(order)

	c.JSON(http.StatusCreated, order)
}

func (s *Service) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	s.ordersMutex.RLock()
	order, exists := s.orders[orderID]
	s.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (s *Service) ListOrders(c *gin.Context) {
	s.ordersMutex.RLock()
	orders := make([]*BridgeOrder, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	s.ordersMutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (s *Service) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	s.ordersMutex.Lock()
	order, exists := s.orders[orderID]
	if exists && order.Status == BridgeOrderStatusPending {
		order.Status = BridgeOrderStatusCancelled
		s.broadcastOrderUpdate(order)
	}
	s.ordersMutex.Unlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

func (s *Service) GetMetrics(c *gin.Context) {
	s.metrics.mutex.RLock()
	metrics := map[string]interface{}{
		"total_orders":     s.metrics.TotalOrders,
		"completed_orders": s.metrics.CompletedOrders,
		"failed_orders":    s.metrics.FailedOrders,
		"total_volume":     s.metrics.TotalVolume.String(),
		"uptime":          time.Since(s.metrics.StartTime).String(),
		"success_rate":    s.calculateSuccessRate(),
	}
	s.metrics.mutex.RUnlock()

	c.JSON(http.StatusOK, metrics)
}

func (s *Service) GetEthereumStatus(c *gin.Context) {
	balance, err := s.ethClient.GetBalance(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	blockNumber, err := s.ethClient.GetCurrentBlockNumber(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get block number"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "connected",
		"address":      s.ethClient.GetAddress().Hex(),
		"balance":      balance.String(),
		"block_number": blockNumber,
	})
}

func (s *Service) GetCosmosStatus(c *gin.Context) {
	balance, err := s.cosmosClient.GetBalance(context.Background(), "cosmos1demo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "connected",
		"chain_id":    s.cosmosClient.GetChainID(),
		"balance":     balance.String(),
		"contract":    s.cosmosClient.GetContractAddress(),
	})
}

func (s *Service) HandleWebSocket(c *gin.Context) {
	wsUpgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool {
            allowedOrigins := []string{"https://flow-fusion.com", "https://api.flow-fusion.com"}
            origin := r.Header.Get("Origin")
            for _, allowed := range allowedOrigins {
                if origin == allowed {
                    return true
                }
            }
            return false
        },
    }

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        s.logger.Error("Failed to upgrade WebSocket", zap.Error(err))
        return
    }

	token := c.Query("token")
    if !s.validateWebSocketToken(token) {
        conn.Close()
        return
    }

	defer conn.Close()

	// Add connection
	s.wsConnectionsMutex.Lock()
	s.wsConnections[conn] = true
	s.wsConnectionsMutex.Unlock()

	// Remove connection on exit
	defer func() {
		s.wsConnectionsMutex.Lock()
		delete(s.wsConnections, conn)
		s.wsConnectionsMutex.Unlock()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *Service) validateWebSocketToken(token string) bool {
    // Implement JWT or other token validation
    return true // Replace with actual validation logic
}

func (s *Service) calculateQuote(req BridgeQuoteRequest, amount *big.Int) (*BridgeQuoteResponse, error) {
	// For demo purposes, return mock quote
	estimatedOutput := new(big.Int).Mul(amount, big.NewInt(99))
	estimatedOutput.Div(estimatedOutput, big.NewInt(100)) // 99% of input

	return &BridgeQuoteResponse{
		EstimatedOutput: estimatedOutput.String(),
		EstimatedFees:   "10000", // 0.01 tokens
		ExecutionTime:   5 * time.Minute,
		PriceImpact:     0.001, // 0.1%
		Route:          []string{req.SourceChain, req.DestChain},
	}, nil
}

func (s *Service) createBridgeOrder(req CreateBridgeOrderRequest, amount *big.Int) (*BridgeOrder, error) {
	orderID := s.generateOrderID()
	
	order := &BridgeOrder{
		ID:          orderID,
		UserAddress: req.UserAddress,
		SourceChain: req.SourceChain,
		DestChain:   req.DestChain,
		SourceToken: req.SourceToken,
		DestToken:   req.DestToken,
		Amount:      amount,
		Status:      BridgeOrderStatusPending,
		CreatedAt:   time.Now(),
		SecretHash:  s.generateSecretHash(),
		Secret:      s.generateSecret(),
		TimeLock:    3600, // 1 hour
	}

	// If TWAP is requested, create TWAP order
	if req.UseTWAP && req.TWAPConfig != nil {
		// This would integrate with the TWAP engine
		order.TWAPOrderID = "twap_" + orderID
	}

	// Store order
	s.ordersMutex.Lock()
	s.orders[orderID] = order
	s.ordersMutex.Unlock()

	// Update metrics
	s.metrics.mutex.Lock()
	s.metrics.TotalOrders++
	s.metrics.TotalVolume.Add(s.metrics.TotalVolume, amount)
	s.metrics.mutex.Unlock()

	s.logger.Info("Bridge order created",
		zap.String("order_id", orderID),
		zap.String("user", req.UserAddress),
		zap.String("amount", amount.String()),
		zap.Bool("use_twap", req.UseTWAP),
	)

	return order, nil
}

func (s *Service) processOrders(ctx context.Context) {
	s.ordersMutex.RLock()
	orders := make([]*BridgeOrder, 0)
	for _, order := range s.orders {
		if order.Status == BridgeOrderStatusPending {
			orders = append(orders, order)
		}
	}
	s.ordersMutex.RUnlock()

	for _, order := range orders {
		go s.processOrder(ctx, order)
	}
}

func (s *Service) processOrder(ctx context.Context, order *BridgeOrder) {
	order.Status = BridgeOrderStatusProcessing
	s.broadcastOrderUpdate(order)

	// Simulate order processing
	time.Sleep(3 * time.Second)

	// For demo purposes, mark as completed
	order.Status = BridgeOrderStatusCompleted
	now := time.Now()
	order.CompletedAt = &now
	order.EthereumTxHash = "0xdemo_eth_tx"
	order.CosmosTxHash = "demo_cosmos_tx"

	// Update metrics
	s.metrics.mutex.Lock()
	s.metrics.CompletedOrders++
	s.metrics.mutex.Unlock()

	s.broadcastOrderUpdate(order)

	s.logger.Info("Bridge order completed",
		zap.String("order_id", order.ID),
		zap.Duration("duration", time.Since(order.CreatedAt)),
	)
}

func (s *Service) watchEthereumEvents(ctx context.Context) {
	// Implementation would watch for Ethereum escrow events
	s.logger.Info("Started Ethereum event watcher")
}

func (s *Service) watchCosmosEvents(ctx context.Context) {
	// Implementation would watch for Cosmos escrow events
	s.logger.Info("Started Cosmos event watcher")
}

func (s *Service) broadcastOrderUpdate(order *BridgeOrder) {
	s.wsConnectionsMutex.RLock()
	connections := make([]*websocket.Conn, 0, len(s.wsConnections))
	for conn := range s.wsConnections {
		connections = append(connections, conn)
	}
	s.wsConnectionsMutex.RUnlock()

	message := map[string]interface{}{
		"type":  "order_update",
		"order": order,
	}

	for _, conn := range connections {
		if err := conn.WriteJSON(message); err != nil {
			// Remove failed connection
			s.wsConnectionsMutex.Lock()
			delete(s.wsConnections, conn)
			s.wsConnectionsMutex.Unlock()
		}
	}
}

func (s *Service) generateOrderID() string {
	return fmt.Sprintf("bridge_%d", time.Now().UnixNano())
}

func (s *Service) generateSecretHash() string {
	return fmt.Sprintf("hash_%d", time.Now().UnixNano())
}

func (s *Service) generateSecret() []byte {
	return []byte(fmt.Sprintf("secret_%d", time.Now().UnixNano()))
}

func (s *Service) calculateSuccessRate() float64 {
	s.metrics.mutex.RLock()
	defer s.metrics.mutex.RUnlock()
	
	if s.metrics.TotalOrders == 0 {
		return 0
	}
	return float64(s.metrics.CompletedOrders) / float64(s.metrics.TotalOrders) * 100
}

func NewBridgeMetrics() *BridgeMetrics {
	return &BridgeMetrics{
		TotalVolume: big.NewInt(0),
		StartTime:   time.Now(),
	}
}

func (s *Service) GetEthereumStatusData() map[string]interface{} {
	balance, err := s.ethClient.GetBalance(context.Background())
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	blockNumber, err := s.ethClient.GetCurrentBlockNumber(context.Background())
	if err != nil {
		return map[string]interface{}{
			"status": "error", 
			"error":  err.Error(),
		}
	}

	return map[string]interface{}{
		"status":       "connected",
		"address":      s.ethClient.GetAddress().Hex(),
		"balance":      balance.String(),
		"block_number": blockNumber,
	}
}

// GetCosmosStatusData returns Cosmos status data without handling HTTP response
func (s *Service) GetCosmosStatusData() map[string]interface{} {
	balance, err := s.cosmosClient.GetBalance(context.Background(), "cosmos1demo")
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		}
	}

	return map[string]interface{}{
		"status":   "connected",
		"chain_id": s.cosmosClient.GetChainID(),
		"balance":  balance.String(),
		"contract": s.cosmosClient.GetContractAddress(),
	}
}