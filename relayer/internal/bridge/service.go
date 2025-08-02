// relayer/internal/bridge/service.go
package bridge

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/ethereum"
	"flow-fusion/relayer/internal/oneinch"
)

// Production 1inch HTTP API Bridge Service
type FusionBridgeService struct {
	config           Config
	ethClient        *ethereum.Client
	cosmosClient     *cosmos.Client
	logger           *zap.Logger
	
	// 1inch HTTP client
	oneInchClient    *oneinch.HTTPClient
	
	// State management
	bridgeOrders     map[string]*BridgeOrder
	ordersMutex      sync.RWMutex
	
	// WebSocket connections
	wsUpgrader       websocket.Upgrader
	wsConnections    map[*websocket.Conn]bool
	wsConnectionsMutex sync.RWMutex
	
	// Metrics
	metrics          *BridgeMetrics
	
	// Private key for signing
	privateKey       *ecdsa.PrivateKey
	publicAddress    common.Address
}

type Config struct {
	EthereumClient *ethereum.Client
	CosmosClient   *cosmos.Client
	Logger         *zap.Logger
	OneInchAPIKey  string
	PrivateKey     string
	NodeURL        string
	ChainID        int
}

type BridgeOrder struct {
	ID                  string                 `json:"id"`
	UserAddress         string                 `json:"user_address"`
	SourceChain         string                 `json:"source_chain"`
	DestChain           string                 `json:"dest_chain"`
	SourceToken         string                 `json:"source_token"`
	DestToken           string                 `json:"dest_token"`
	Amount              *big.Int               `json:"amount"`
	Status              BridgeOrderStatus      `json:"status"`
	CreatedAt           time.Time              `json:"created_at"`
	
	// 1inch HTTP API specific fields
	OneInchOrderHash    string                 `json:"oneinch_order_hash"`
	Secrets             []string               `json:"-"` // Never expose secrets
	SecretHashes        []string               `json:"secret_hashes"`
	FusionPlusQuote     *oneinch.FusionQuoteResponse `json:"fusion_plus_quote,omitempty"`
	
	// Cross-chain tracking
	EthereumTxHash      string                 `json:"ethereum_tx_hash,omitempty"`
	CosmosTxHash        string                 `json:"cosmos_tx_hash,omitempty"`
	CompletedAt         *time.Time             `json:"completed_at,omitempty"`
	Error               string                 `json:"error,omitempty"`
	
	// TWAP support
	TWAPEnabled         bool                   `json:"twap_enabled"`
	TWAPConfig          *TWAPConfiguration     `json:"twap_config,omitempty"`
}

type TWAPConfiguration struct {
	TimeWindow    int     `json:"time_window"`
	IntervalCount int     `json:"interval_count"`
	MaxSlippage   float64 `json:"max_slippage"`
	StartTime     time.Time `json:"start_time"`
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

// Request/Response types for API endpoints
type CreateBridgeOrderRequest struct {
	UserAddress string             `json:"user_address" binding:"required"`
	SourceChain string             `json:"source_chain" binding:"required"`
	DestChain   string             `json:"dest_chain" binding:"required"`
	SourceToken string             `json:"source_token" binding:"required"`
	DestToken   string             `json:"dest_token" binding:"required"`
	Amount      string             `json:"amount" binding:"required"`
	TWAPEnabled bool               `json:"twap_enabled,omitempty"`
	TWAPConfig  *TWAPConfiguration `json:"twap_config,omitempty"`
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
	QuoteData       interface{}   `json:"quote_data"` // Raw 1inch quote data
}

func NewFusionBridgeService(config Config) (*FusionBridgeService, error) {
	// Parse private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	
	publicKey := privateKey.Public()
	publicAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	// Initialize 1inch HTTP client
	oneInchConfig := oneinch.Config{
		BaseURL: "https://api.1inch.dev",
		APIKey:  config.OneInchAPIKey,
		ChainID: uint64(config.ChainID),
		Timeout: 30 * time.Second,
	}

	oneInchClient := oneinch.NewHTTPClient(oneInchConfig, config.Logger)

	wsUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Configure proper CORS in production
		},
	}

	service := &FusionBridgeService{
		config:            config,
		ethClient:         config.EthereumClient,
		cosmosClient:      config.CosmosClient,
		logger:            config.Logger,
		oneInchClient:     oneInchClient,
		bridgeOrders:      make(map[string]*BridgeOrder),
		wsUpgrader:        wsUpgrader,
		wsConnections:     make(map[*websocket.Conn]bool),
		metrics:           NewBridgeMetrics(),
		privateKey:        privateKey,
		publicAddress:     publicAddress,
	}

	return service, nil
}

func (s *FusionBridgeService) Start(ctx context.Context) error {
	s.logger.Info("Starting Production HTTP API Bridge Service",
		zap.String("public_address", s.publicAddress.Hex()),
		zap.Int("chain_id", s.config.ChainID),
	)

	// Start order monitoring and processing
	go s.monitorOrders(ctx)
	go s.processFusionPlusOrders(ctx)

	// Main processing loop
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("HTTP API Bridge Service stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := s.processOrders(ctx); err != nil {
				s.logger.Error("Error processing orders", zap.Error(err))
			}
		}
	}
}

func (s *FusionBridgeService) GetQuote(c *gin.Context) {
	var req BridgeQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote, err := s.calculateFusionPlusQuote(c.Request.Context(), req)
	if err != nil {
		s.logger.Error("Failed to calculate Fusion+ quote", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate quote"})
		return
	}

	c.JSON(http.StatusOK, quote)
}

func (s *FusionBridgeService) CreateOrder(c *gin.Context) {
	var req CreateBridgeOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	order, err := s.createFusionPlusOrder(c.Request.Context(), req, amount)
	if err != nil {
		s.logger.Error("Failed to create Fusion+ order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	s.broadcastOrderUpdate(order)
	c.JSON(http.StatusCreated, order)
}

func (s *FusionBridgeService) calculateFusionPlusQuote(ctx context.Context, req BridgeQuoteRequest) (*BridgeQuoteResponse, error) {
	// Determine chain IDs for Fusion+
	srcChain, dstChain, err := s.getChainIDs(req.SourceChain, req.DestChain)
	if err != nil {
		return nil, fmt.Errorf("invalid chain configuration: %w", err)
	}

	// Get Fusion+ quote using HTTP API
	quote, err := s.oneInchClient.GetFusionQuote(
		ctx,
		uint64(srcChain),
		uint64(dstChain),
		req.SourceToken,
		req.DestToken,
		req.Amount,
		s.publicAddress.Hex(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get Fusion+ quote: %w", err)
	}

	// Calculate estimated fees (gas + protocol fees)
	estimatedFees := s.calculateEstimatedFees(quote)

	return &BridgeQuoteResponse{
		EstimatedOutput: quote.DstTokenAmount,
		EstimatedFees:   estimatedFees,
		ExecutionTime:   5 * time.Minute, // Typical cross-chain execution time
		PriceImpact:     0.001, // Will be calculated from quote data
		Route:          []string{req.SourceChain, req.DestChain},
		QuoteData:      quote,
	}, nil
}

func (s *FusionBridgeService) createFusionPlusOrder(ctx context.Context, req CreateBridgeOrderRequest, amount *big.Int) (*BridgeOrder, error) {
	// Determine chain IDs
	srcChain, dstChain, err := s.getChainIDs(req.SourceChain, req.DestChain)
	if err != nil {
		return nil, fmt.Errorf("invalid chain configuration: %w", err)
	}

	// Generate order ID
	orderID := s.generateOrderID()

	// Get quote for the order using HTTP API
	quote, err := s.oneInchClient.GetFusionQuote(
		ctx,
		uint64(srcChain),
		uint64(dstChain),
		req.SourceToken,
		req.DestToken,
		amount.String(),
		req.UserAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote for order: %w", err)
	}

	// Find the recommended preset
	var selectedPreset *oneinch.FusionPreset
	for _, preset := range quote.Presets {
		if preset.Name == quote.RecommendedPreset {
			selectedPreset = &preset
			break
		}
	}
	if selectedPreset == nil {
		return nil, fmt.Errorf("recommended preset not found: %s", quote.RecommendedPreset)
	}

	// Generate secrets based on preset requirements
	secretsCount := selectedPreset.SecretsCount
	secrets := make([]string, secretsCount)
	secretHashes := make([]string, secretsCount)

	for i := 0; i < secretsCount; i++ {
		// Generate random secret - implement this method
		secret, err := s.generateSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate secret %d: %w", i, err)
		}
		secrets[i] = secret

		// Hash the secret - implement this method
		secretHash, err := s.hashSecret(secret)
		if err != nil {
			return nil, fmt.Errorf("failed to hash secret %d: %w", i, err)
		}
		secretHashes[i] = secretHash
	}

	// Create order hash - this would need to be implemented based on 1inch's order structure
	orderHash, err := s.generateOrderHash(quote, secretHashes, req.UserAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate order hash: %w", err)
	}

	// Place the Fusion+ order using HTTP API
	fusionOrderReq := oneinch.FusionOrderRequest{
		WalletAddress: req.UserAddress,
		OrderHash:     orderHash,
		SecretHashes:  secretHashes,
		Receiver:      req.UserAddress,
		Preset:        quote.RecommendedPreset,
	}

	fusionOrderResp, err := s.oneInchClient.PlaceFusionOrder(ctx, fusionOrderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to place Fusion+ order: %w", err)
	}

	if !fusionOrderResp.Success {
		return nil, fmt.Errorf("fusion+ order placement failed: %s", fusionOrderResp.Message)
	}

	// Create bridge order record
	bridgeOrder := &BridgeOrder{
		ID:               orderID,
		UserAddress:      req.UserAddress,
		SourceChain:      req.SourceChain,
		DestChain:        req.DestChain,
		SourceToken:      req.SourceToken,
		DestToken:        req.DestToken,
		Amount:           amount,
		Status:           BridgeOrderStatusPending,
		CreatedAt:        time.Now(),
		OneInchOrderHash: fusionOrderResp.OrderHash,
		Secrets:          secrets,
		SecretHashes:     secretHashes,
		FusionPlusQuote:  quote,
		TWAPEnabled:      req.TWAPEnabled,
		TWAPConfig:       req.TWAPConfig,
	}

	// Store order
	s.ordersMutex.Lock()
	s.bridgeOrders[orderID] = bridgeOrder
	s.ordersMutex.Unlock()

	// Update metrics
	s.metrics.mutex.Lock()
	s.metrics.TotalOrders++
	s.metrics.TotalVolume.Add(s.metrics.TotalVolume, amount)
	s.metrics.mutex.Unlock()

	s.logger.Info("Fusion+ bridge order created",
		zap.String("order_id", orderID),
		zap.String("oneinch_order_hash", fusionOrderResp.OrderHash),
		zap.String("user", req.UserAddress),
		zap.String("amount", amount.String()),
		zap.Bool("twap_enabled", req.TWAPEnabled),
	)

	return bridgeOrder, nil
}

func (s *FusionBridgeService) processFusionPlusOrders(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.ordersMutex.RLock()
			orders := make([]*BridgeOrder, 0)
			for _, order := range s.bridgeOrders {
				if order.Status == BridgeOrderStatusPending || order.Status == BridgeOrderStatusProcessing {
					orders = append(orders, order)
				}
			}
			s.ordersMutex.RUnlock()

			for _, order := range orders {
				go s.processFusionPlusOrder(ctx, order)
			}
		}
	}
}

func (s *FusionBridgeService) processFusionPlusOrder(ctx context.Context, order *BridgeOrder) {
	// Check order status using HTTP API
	fusionOrderStatus, err := s.oneInchClient.GetFusionOrderStatus(ctx, order.OneInchOrderHash)
	if err != nil {
		s.logger.Error("Failed to get Fusion+ order status",
			zap.String("order_id", order.ID),
			zap.String("oneinch_hash", order.OneInchOrderHash),
			zap.Error(err),
		)
		return
	}

	// Update order status based on Fusion+ status
	switch fusionOrderStatus.Status {
	case "pending":
		if order.Status != BridgeOrderStatusPending {
			order.Status = BridgeOrderStatusPending
			s.broadcastOrderUpdate(order)
		}
	case "executed":
		if order.Status != BridgeOrderStatusCompleted {
			order.Status = BridgeOrderStatusCompleted
			now := time.Now()
			order.CompletedAt = &now
			
			// Update metrics
			s.metrics.mutex.Lock()
			s.metrics.CompletedOrders++
			s.metrics.mutex.Unlock()
			
			s.broadcastOrderUpdate(order)
			s.logger.Info("Fusion+ order completed",
				zap.String("order_id", order.ID),
				zap.String("oneinch_hash", order.OneInchOrderHash),
				zap.String("fusion_status", fusionOrderStatus.Status),
			)
		}
	case "cancelled", "expired":
		if order.Status != BridgeOrderStatusFailed {
			order.Status = BridgeOrderStatusFailed
			order.Error = fmt.Sprintf("Order was %s", fusionOrderStatus.Status)
			
			// Update metrics
			s.metrics.mutex.Lock()
			s.metrics.FailedOrders++
			s.metrics.mutex.Unlock()
			
			s.broadcastOrderUpdate(order)
			s.logger.Warn("Fusion+ order failed",
				zap.String("order_id", order.ID),
				zap.String("oneinch_hash", order.OneInchOrderHash),
				zap.String("fusion_status", fusionOrderStatus.Status),
			)
		}
	default:
		s.logger.Debug("Fusion+ order status",
			zap.String("order_id", order.ID),
			zap.String("status", fusionOrderStatus.Status),
		)
	}

	// Check for ready-to-accept fills using HTTP API
	fills, err := s.oneInchClient.GetReadyToAcceptFills(ctx, order.OneInchOrderHash)
	if err != nil {
		s.logger.Error("Failed to get ready fills",
			zap.String("order_id", order.ID),
			zap.Error(err),
		)
		return
	}

	// Submit secrets for ready fills
	if len(fills.Fills) > 0 && len(order.Secrets) > 0 {
		order.Status = BridgeOrderStatusProcessing
		s.broadcastOrderUpdate(order)

		// Submit secret for the first fill (can be extended for multiple fills)
		secretResp, err := s.oneInchClient.SubmitSecret(ctx, order.OneInchOrderHash, order.Secrets[0])
		if err != nil {
			s.logger.Error("Failed to submit secret",
				zap.String("order_id", order.ID),
				zap.Error(err),
			)
			return
		}

		if !secretResp.Success {
			s.logger.Error("Secret submission failed",
				zap.String("order_id", order.ID),
				zap.String("message", secretResp.Message),
			)
			return
		}

		s.logger.Info("Secret submitted for Fusion+ order",
			zap.String("order_id", order.ID),
			zap.String("oneinch_hash", order.OneInchOrderHash),
		)

		// Create corresponding Cosmos escrow
		if err := s.createCosmosEscrow(ctx, order); err != nil {
			s.logger.Error("Failed to create Cosmos escrow",
				zap.String("order_id", order.ID),
				zap.Error(err),
			)
		}
	}
}

func (s *FusionBridgeService) createCosmosEscrow(ctx context.Context, order *BridgeOrder) error {
	// Create Cosmos escrow using the order details
	params := cosmos.CreateEscrowParams{
		EthereumTxHash: order.OneInchOrderHash, // Use 1inch order hash as reference
		SecretHash:     order.SecretHashes[0],  // Use first secret hash
		TimeLock:       uint64(time.Now().Add(24 * time.Hour).Unix()),
		Recipient:      order.UserAddress,
		EthereumSender: s.publicAddress.Hex(),
		Amount:         order.Amount,
	}

	txHash, err := s.cosmosClient.CreateEscrow(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create Cosmos escrow: %w", err)
	}

	order.CosmosTxHash = txHash
	s.logger.Info("Created Cosmos escrow",
		zap.String("order_id", order.ID),
		zap.String("cosmos_tx", txHash),
	)

	return nil
}

// Helper methods
func (s *FusionBridgeService) getChainIDs(sourceChain, destChain string) (int, int, error) {
	chainMap := map[string]int{
		"ethereum": 1,
		"polygon":  137,
		"arbitrum": 42161,
		"base":     8453,
		"cosmos":   1, // Custom chain ID for Cosmos
	}

	srcChainID, srcExists := chainMap[sourceChain]
	dstChainID, dstExists := chainMap[destChain]

	if !srcExists || !dstExists {
		return 0, 0, fmt.Errorf("unsupported chain: src=%s, dst=%s", sourceChain, destChain)
	}

	return srcChainID, dstChainID, nil
}

func (s *FusionBridgeService) calculateEstimatedFees(quote *oneinch.FusionQuoteResponse) string {
	// Implementation depends on the quote structure
	// This would calculate gas costs + protocol fees
	return "1000000" // Placeholder - implement based on actual quote data
}

func (s *FusionBridgeService) generateOrderID() string {
	return fmt.Sprintf("fusion_%d_%s", time.Now().UnixNano(), hex.EncodeToString([]byte{byte(len(s.bridgeOrders))}))
}


func (s *FusionBridgeService) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	s.ordersMutex.RLock()
	order, exists := s.bridgeOrders[orderID]
	s.ordersMutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (s *FusionBridgeService) ListOrders(c *gin.Context) {
	s.ordersMutex.RLock()
	orders := make([]*BridgeOrder, 0, len(s.bridgeOrders))
	for _, order := range s.bridgeOrders {
		orders = append(orders, order)
	}
	s.ordersMutex.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (s *FusionBridgeService) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	
	s.ordersMutex.Lock()
	order, exists := s.bridgeOrders[orderID]
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

func (s *FusionBridgeService) HandleWebSocket(c *gin.Context) {
	conn, err := s.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.logger.Error("Failed to upgrade WebSocket", zap.Error(err))
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

func (s *FusionBridgeService) calculateSuccessRate() float64 {
	if s.metrics.TotalOrders == 0 {
		return 0
	}
	return float64(s.metrics.CompletedOrders) / float64(s.metrics.TotalOrders) * 100
}

func (s *FusionBridgeService) broadcastOrderUpdate(order *BridgeOrder) {
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
			s.wsConnectionsMutex.Lock()
			delete(s.wsConnections, conn)
			s.wsConnectionsMutex.Unlock()
		}
	}
}

func (s *FusionBridgeService) monitorOrders(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Monitor order health and timeouts
			s.checkOrderHealth()
		}
	}
}

func (s *FusionBridgeService) checkOrderHealth() {
	s.ordersMutex.RLock()
	orders := make([]*BridgeOrder, 0, len(s.bridgeOrders))
	for _, order := range s.bridgeOrders {
		if order.Status == BridgeOrderStatusProcessing {
			orders = append(orders, order)
		}
	}
	s.ordersMutex.RUnlock()

	for _, order := range orders {
		// Check for order timeout (24 hours)
		if time.Since(order.CreatedAt) > 24*time.Hour {
			s.ordersMutex.Lock()
			order.Status = BridgeOrderStatusFailed
			order.Error = "Order timeout"
			s.ordersMutex.Unlock()

			s.metrics.mutex.Lock()
			s.metrics.FailedOrders++
			s.metrics.mutex.Unlock()

			s.broadcastOrderUpdate(order)
			s.logger.Warn("Order timed out",
				zap.String("order_id", order.ID),
				zap.Duration("age", time.Since(order.CreatedAt)),
			)
		}
	}
}

func (s *FusionBridgeService) processOrders(ctx context.Context) error {
	// Additional order processing logic can be added here
	return nil
}

func (s *FusionBridgeService) GetEthereumStatusData() map[string]interface{} {
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
		"status":         "connected",
		"address":        s.ethClient.GetAddress().Hex(),
		"balance":        balance.String(),
		"block_number":   blockNumber,
		"fusion_address": s.publicAddress.Hex(),
	}
}

func (s *FusionBridgeService) GetCosmosStatusData() map[string]interface{} {
	balance, err := s.cosmosClient.GetBalance(context.Background(), "")
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

func NewBridgeMetrics() *BridgeMetrics {
	return &BridgeMetrics{
		TotalVolume: big.NewInt(0),
		StartTime:   time.Now(),
	}
}

func (s *FusionBridgeService) GetMetricsData() map[string]interface{} {
	s.metrics.mutex.RLock()
	defer s.metrics.mutex.RUnlock()
	
	return map[string]interface{}{
		"total_orders":     s.metrics.TotalOrders,
		"completed_orders": s.metrics.CompletedOrders,
		"failed_orders":    s.metrics.FailedOrders,
		"total_volume":     s.metrics.TotalVolume.String(),
		"uptime":          time.Since(s.metrics.StartTime).String(),
		"success_rate":    s.calculateSuccessRate(),
		"public_address":  s.publicAddress.Hex(),
		"chain_id":        s.config.ChainID,
	}
}