
package twap

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"go.uber.org/zap"

	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/oneinch"
)

// TWAPScheduler manages the scheduling and execution of TWAP orders using 1inch HTTP API
type TWAPScheduler struct {
	engine              *ProductionTWAPEngine
	logger              *zap.Logger
	scheduledOrders     map[string]*TWAPOrder
	mutex               sync.RWMutex
	ticker              *time.Ticker
	workerPool          chan struct{}
	stopChan            chan struct{}
	maxConcurrentOrders int
	executionQueue      chan *ExecutionTask
}

type ExecutionTask struct {
	Order  *TWAPOrder
	Window *ExecutionWindow
	Retry  int
}

const (
	MaxConcurrentExecutions = 20
	MaxRetryAttempts       = 3
	RetryDelay            = 30 * time.Second
	ExecutionTimeout      = 10 * time.Minute
	HealthCheckInterval   = 30 * time.Second
)

func NewTWAPScheduler(logger *zap.Logger, engine *ProductionTWAPEngine) *TWAPScheduler {
	return &TWAPScheduler{
		engine:              engine,
		logger:              logger,
		scheduledOrders:     make(map[string]*TWAPOrder),
		ticker:              time.NewTicker(5 * time.Second),
		workerPool:          make(chan struct{}, MaxConcurrentExecutions),
		stopChan:            make(chan struct{}),
		maxConcurrentOrders: MaxConcurrentExecutions,
		executionQueue:      make(chan *ExecutionTask, 1000),
	}
}

func (s *TWAPScheduler) Start(ctx context.Context) {
	s.logger.Info("Starting Production TWAP Scheduler with HTTP API",
		zap.Int("max_concurrent", s.maxConcurrentOrders),
	)

	// Start worker pool
	for i := 0; i < s.maxConcurrentOrders; i++ {
		go s.executionWorker(ctx, i)
	}

	// Start main scheduler loop
	go s.schedulerLoop(ctx)

	// Start health monitoring
	go s.healthMonitor(ctx)

	defer s.ticker.Stop()
	defer close(s.workerPool)
	defer close(s.stopChan)
	defer close(s.executionQueue)

	<-ctx.Done()
	s.logger.Info("TWAP Scheduler stopped")
}

func (s *TWAPScheduler) ScheduleOrder(order *TWAPOrder) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.scheduledOrders[order.ID] = order
	s.logger.Info("TWAP order scheduled",
		zap.String("order_id", order.ID),
		zap.Time("start_time", order.StartTime),
		zap.Int("windows", len(order.ExecutionWindows)),
		zap.String("total_amount", order.TotalAmount.String()),
	)
}

func (s *TWAPScheduler) RemoveOrder(orderID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.scheduledOrders, orderID)
	s.logger.Info("TWAP order removed from scheduler", zap.String("order_id", orderID))
}

func (s *TWAPScheduler) schedulerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ticker.C:
			s.processScheduledOrders(ctx)
		}
	}
}

func (s *TWAPScheduler) processScheduledOrders(ctx context.Context) {
	now := time.Now()

	s.mutex.RLock()
	orders := make([]*TWAPOrder, 0, len(s.scheduledOrders))
	for _, order := range s.scheduledOrders {
		orders = append(orders, order)
	}
	s.mutex.RUnlock()

	for _, order := range orders {
		// Remove completed or failed orders
		if order.Status == TWAPOrderStatusCompleted || 
		   order.Status == TWAPOrderStatusFailed || 
		   order.Status == TWAPOrderStatusCancelled {
			s.RemoveOrder(order.ID)
			continue
		}

		// Start order execution if it's time
		if order.Status == TWAPOrderStatusPending && now.After(order.StartTime) {
			s.engine.ordersMutex.Lock()
			order.Status = TWAPOrderStatusExecuting
			s.engine.ordersMutex.Unlock()

			s.logger.Info("Starting TWAP order execution",
				zap.String("order_id", order.ID),
				zap.Time("scheduled_start", order.StartTime),
			)
		}

		// Process execution windows
		for _, window := range order.ExecutionWindows {
			if s.shouldExecuteWindow(window, now) {
				task := &ExecutionTask{
					Order:  order,
					Window: window,
					Retry:  0,
				}

				select {
				case s.executionQueue <- task:
					window.Status = WindowStatusExecuting
				default:
					s.logger.Warn("Execution queue full, delaying window",
						zap.String("order_id", order.ID),
						zap.Int("window", window.Index),
					)
				}
			} else if s.isWindowExpired(window, now) {
				window.Status = WindowStatusSkipped
				s.logger.Warn("Window execution expired",
					zap.String("order_id", order.ID),
					zap.Int("window", window.Index),
					zap.Time("end_time", window.EndTime),
				)
			}
		}

		// Check if order is complete
		if s.isOrderComplete(order) {
			s.engine.ordersMutex.Lock()
			order.Status = TWAPOrderStatusCompleted
			s.engine.ordersMutex.Unlock()

			s.logger.Info("TWAP order completed",
				zap.String("order_id", order.ID),
				zap.Int("completed_windows", order.CompletedWindows),
				zap.Int("total_windows", order.IntervalCount),
			)
		}
	}
}

func (s *TWAPScheduler) executionWorker(ctx context.Context, workerID int) {
	s.logger.Info("Starting TWAP execution worker", zap.Int("worker_id", workerID))

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-s.executionQueue:
			s.executeTask(ctx, task, workerID)
		}
	}
}

func (s *TWAPScheduler) executeTask(ctx context.Context, task *ExecutionTask, workerID int) {
	startTime := time.Now()

	s.logger.Info("Worker executing TWAP window",
		zap.Int("worker_id", workerID),
		zap.String("order_id", task.Order.ID),
		zap.Int("window", task.Window.Index),
		zap.Int("retry", task.Retry),
		zap.String("amount", task.Window.Amount.String()),
	)

	// Create execution context with timeout
	execCtx, cancel := context.WithTimeout(ctx, ExecutionTimeout)
	defer cancel()

	var err error
	if s.shouldUseFusionPlus(task.Order, task.Window) {
		err = s.executeFusionPlusWindow(execCtx, task.Order, task.Window)
	} else {
		err = s.executeOrderbookWindow(execCtx, task.Order, task.Window)
	}

	duration := time.Since(startTime)

	if err != nil {
		s.handleExecutionError(task, err, duration)
	} else {
		s.handleExecutionSuccess(task, duration, workerID)
	}
}

func (s *TWAPScheduler) executeFusionPlusWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	s.logger.Info("Executing window via Fusion+ HTTP API",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
	)

	srcChain := uint64(s.engine.chainID)
	dstChain := uint64(1) // Target chain

	// Get Fusion+ quote using HTTP API
	quote, err := s.engine.oneInchClient.GetFusionQuote(
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

	// Find the recommended preset
	var selectedPreset *oneinch.FusionPreset
	for _, preset := range quote.Presets {
		if preset.Name == quote.RecommendedPreset {
			selectedPreset = &preset
			break
		}
	}
	if selectedPreset == nil {
		return fmt.Errorf("recommended preset not found: %s", quote.RecommendedPreset)
	}

	// Generate secrets based on preset requirements
	secretsCount := selectedPreset.SecretsCount
	secrets := make([]string, secretsCount)
	secretHashes := make([]string, secretsCount)

	for i := 0; i < secretsCount; i++ {
		secret, err := s.generateSecret()
		if err != nil {
			return fmt.Errorf("failed to generate secret %d: %w", i, err)
		}
		secrets[i] = secret

		secretHash, err := s.hashSecret(secret)
		if err != nil {
			return fmt.Errorf("failed to hash secret %d: %w", i, err)
		}
		secretHashes[i] = secretHash
	}

	// Generate order hash
	orderHash, err := s.generateOrderHash(quote, secretHashes, order.UserAddress)
	if err != nil {
		return fmt.Errorf("failed to generate order hash: %w", err)
	}

	// Place Fusion+ order using HTTP API
	fusionOrderReq := oneinch.FusionOrderRequest{
		WalletAddress: order.UserAddress,
		OrderHash:     orderHash,
		SecretHashes:  secretHashes,
		Receiver:      order.UserAddress,
		Preset:        quote.RecommendedPreset,
	}

	fusionOrderResp, err := s.engine.oneInchClient.PlaceFusionOrder(ctx, fusionOrderReq)
	if err != nil {
		return fmt.Errorf("failed to place Fusion+ order: %w", err)
	}

	if !fusionOrderResp.Success {
		return fmt.Errorf("fusion+ order placement failed: %s", fusionOrderResp.Message)
	}

	window.FusionPlusHash = fusionOrderResp.OrderHash
	window.Secret = secrets[0] 
	window.SecretHash = secretHashes[0]

	fusionOrderRef := &FusionPlusOrderRef{
		OrderHash:       fusionOrderResp.OrderHash,
		WindowIndex:     window.Index,
		Amount:          window.Amount,
		Status:          "pending",
		SecretSubmitted: false,
	}
	order.FusionPlusOrders = append(order.FusionPlusOrders, fusionOrderRef)

	// Monitor and submit secret when ready
	go s.monitorFusionPlusOrder(ctx, order, window, fusionOrderRef, secrets[0])

	s.logger.Info("Fusion+ order placed via HTTP API",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("fusion_hash", fusionOrderResp.OrderHash),
	)

	return nil
}

func (s *TWAPScheduler) executeOrderbookWindow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	s.logger.Info("Executing window via Orderbook HTTP API",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
	)

	// Calculate taking amount based on expected output
	takingAmount := s.calculateExpectedOutput(window.Amount, order.SourceToken, order.DestToken)

	// Create orderbook order using HTTP API
	orderbookOrderReq := oneinch.OrderbookOrderRequest{
		Maker:        order.UserAddress,
		MakerAsset:   order.SourceToken,
		TakerAsset:   order.DestToken,
		MakingAmount: window.Amount.String(),
		TakingAmount: takingAmount.String(),
		Salt:         fmt.Sprintf("%d", time.Now().UnixNano()),
		MakerTraits:  "0", // Default traits
		Signature:    "0x", // Would need proper signature
	}

	orderbookOrderResp, err := s.engine.oneInchClient.CreateOrderbookOrder(ctx, orderbookOrderReq)
	if err != nil {
		return fmt.Errorf("failed to create orderbook order: %w", err)
	}

	if !orderbookOrderResp.Success {
		return fmt.Errorf("orderbook order creation failed: %s", orderbookOrderResp.Message)
	}

	// Update window with orderbook order details
	window.OrderbookHash = orderbookOrderResp.OrderHash

	// Add to order's orderbook orders tracking
	orderbookOrderRef := &OrderbookOrderRef{
		OrderHash:   orderbookOrderResp.OrderHash,
		WindowIndex: window.Index,
		Amount:      window.Amount,
		Status:      "pending",
		Filled:      big.NewInt(0),
	}
	order.OrderbookOrders = append(order.OrderbookOrders, orderbookOrderRef)

	// Monitor orderbook order status
	go s.monitorOrderbookOrder(ctx, order, window, orderbookOrderRef)

	s.logger.Info("Orderbook order created via HTTP API",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("orderbook_hash", orderbookOrderResp.OrderHash),
	)

	return nil
}

func (s *TWAPScheduler) monitorFusionPlusOrder(ctx context.Context, order *TWAPOrder, window *ExecutionWindow, fusionRef *FusionPlusOrderRef, secret string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	timeout := time.After(ExecutionTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timeout:
			s.logger.Warn("Fusion+ order monitoring timeout",
				zap.String("order_id", order.ID),
				zap.String("fusion_hash", fusionRef.OrderHash),
			)
			return
		case <-ticker.C:
			// Check order status using HTTP API
			fusionOrder, err := s.engine.oneInchClient.GetFusionOrderStatus(ctx, fusionRef.OrderHash)
			if err != nil {
				s.logger.Error("Failed to get Fusion+ order status",
					zap.String("fusion_hash", fusionRef.OrderHash),
					zap.Error(err),
				)
				continue
			}

			fusionRef.Status = fusionOrder.Status

			// Check for ready fills using HTTP API
			fills, err := s.engine.oneInchClient.GetReadyToAcceptFills(ctx, fusionRef.OrderHash)
			if err != nil {
				s.logger.Error("Failed to get ready fills",
					zap.String("fusion_hash", fusionRef.OrderHash),
					zap.Error(err),
				)
				continue
			}

			// Submit secret if fills are ready and not already submitted
			if len(fills.Fills) > 0 && !fusionRef.SecretSubmitted {
				secretResp, err := s.engine.oneInchClient.SubmitSecret(ctx, fusionRef.OrderHash, secret)
				if err != nil {
					s.logger.Error("Failed to submit secret",
						zap.String("fusion_hash", fusionRef.OrderHash),
						zap.Error(err),
					)
					continue
				}

				if !secretResp.Success {
					s.logger.Error("Secret submission failed",
						zap.String("fusion_hash", fusionRef.OrderHash),
						zap.String("message", secretResp.Message),
					)
					continue
				}

				fusionRef.SecretSubmitted = true
				s.logger.Info("Secret submitted for Fusion+ order via HTTP API",
					zap.String("order_id", order.ID),
					zap.String("fusion_hash", fusionRef.OrderHash),
				)

				// Create corresponding Cosmos escrow
				if err := s.createCosmosEscrow(ctx, order, window); err != nil {
					s.logger.Error("Failed to create Cosmos escrow",
						zap.String("order_id", order.ID),
						zap.Error(err),
					)
				}
			}

			// Check if order is executed
			if fusionOrder.Status == "executed" {
				window.Status = WindowStatusCompleted
				s.engine.ordersMutex.Lock()
				order.CompletedWindows++
				order.TotalExecuted.Add(order.TotalExecuted, window.Amount)
				s.engine.ordersMutex.Unlock()

				s.logger.Info("Fusion+ window execution completed",
					zap.String("order_id", order.ID),
					zap.Int("window", window.Index),
				)
				return
			}
		}
	}
}

func (s *TWAPScheduler) monitorOrderbookOrder(ctx context.Context, order *TWAPOrder, window *ExecutionWindow, orderbookRef *OrderbookOrderRef) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	timeout := time.After(ExecutionTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timeout:
			s.logger.Warn("Orderbook order monitoring timeout",
				zap.String("order_id", order.ID),
				zap.String("orderbook_hash", orderbookRef.OrderHash),
			)
			return
		case <-ticker.C:
			// Check order status by fetching orders for the user using HTTP API
			orders, err := s.engine.oneInchClient.GetOrderbookOrdersByMaker(ctx, order.UserAddress, 100)
			if err != nil {
				s.logger.Error("Failed to get orderbook orders",
					zap.String("orderbook_hash", orderbookRef.OrderHash),
					zap.Error(err),
				)
				continue
			}

			// Find our specific order
			var foundOrder *oneinch.OrderbookOrder
			for i := range orders {
				if orders[i].OrderHash == orderbookRef.OrderHash {
					foundOrder = &orders[i]
					break
				}
			}

			if foundOrder != nil {
				orderbookRef.Status = foundOrder.Status
				
				if foundOrder.RemainingAmount != "" {
					if remaining, ok := new(big.Int).SetString(foundOrder.RemainingAmount, 10); ok {
						filled := new(big.Int).Sub(window.Amount, remaining)
						orderbookRef.Filled = filled
					}
				}

				// Check if order is fully filled
				isCompleted := foundOrder.Status == "filled" || 
							   foundOrder.Status == "completed" ||
							   (orderbookRef.Filled != nil && orderbookRef.Filled.Cmp(window.Amount) >= 0)

				if isCompleted {
					window.Status = WindowStatusCompleted
					s.engine.ordersMutex.Lock()
					order.CompletedWindows++
					order.TotalExecuted.Add(order.TotalExecuted, window.Amount)
					s.engine.ordersMutex.Unlock()

					s.logger.Info("Orderbook window execution completed",
						zap.String("order_id", order.ID),
						zap.Int("window", window.Index),
						zap.String("filled_amount", orderbookRef.Filled.String()),
					)
					return
				}
			} else {
				s.logger.Debug("Order not found in user's orders",
					zap.String("order_id", order.ID),
					zap.String("orderbook_hash", orderbookRef.OrderHash),
				)
			}
		}
	}
}

func (s *TWAPScheduler) createCosmosEscrow(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) error {
	// Create Cosmos escrow for cross-chain completion
	params := cosmos.CreateEscrowParams{
		EthereumTxHash: window.FusionPlusHash,
		SecretHash:     window.SecretHash,
		TimeLock:       uint64(time.Now().Add(24 * time.Hour).Unix()),
		Recipient:      order.UserAddress,
		EthereumSender: s.engine.publicAddress.Hex(),
		Amount:         window.Amount,
	}

	txHash, err := s.engine.cosmosClient.CreateEscrow(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create Cosmos escrow: %w", err)
	}

	window.CosmosTxHash = txHash
	s.logger.Info("Created Cosmos escrow for TWAP window",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("cosmos_tx", txHash),
	)

	return nil
}

func (s *TWAPScheduler) calculateExpectedOutput(inputAmount *big.Int, srcToken, dstToken string) *big.Int {
	// Simple calculation - in production, use proper price feeds
	// This is a placeholder that assumes 1:1 ratio
	return new(big.Int).Set(inputAmount)
}

func (s *TWAPScheduler) shouldUseFusionPlus(order *TWAPOrder, window *ExecutionWindow) bool {
	// Decision logic for execution method
	// Use Fusion+ for larger amounts or when cross-chain is involved
	threshold := big.NewInt(1000000) // 1M units threshold
	
	// Set defaults if chains are not specified
	sourceChain := order.SourceChain
	destChain := order.DestChain
	if sourceChain == "" {
		sourceChain = "ethereum"
	}
	if destChain == "" {
		destChain = "cosmos"
	}
	
	return window.Amount.Cmp(threshold) > 0 || sourceChain != destChain
}

// Helper methods for crypto operations
func (s *TWAPScheduler) generateSecret() (string, error) {
	// Generate random secret - this needs to be implemented
	return s.engine.oneInchClient.GenerateRandomBytes32()
}

func (s *TWAPScheduler) hashSecret(secret string) (string, error) {
	// Hash the secret - this needs to be implemented
	return s.engine.oneInchClient.HashSecret(secret)
}

func (s *TWAPScheduler) generateOrderHash(quote *oneinch.FusionQuoteResponse, secretHashes []string, userAddress string) (string, error) {
	// Generate order hash based on order parameters
	// This is a simplified implementation - real implementation would follow 1inch order structure
	return fmt.Sprintf("0x%x", time.Now().UnixNano()), nil
}

func (s *TWAPScheduler) generateSimpleOrderHash() string {
	// Generate a simple order hash
	return fmt.Sprintf("0x%x", time.Now().UnixNano())
}

func (s *TWAPScheduler) shouldExecuteWindow(window *ExecutionWindow, now time.Time) bool {
	return window.Status == WindowStatusPending && 
		   now.After(window.StartTime) && 
		   now.Before(window.EndTime)
}

func (s *TWAPScheduler) isWindowExpired(window *ExecutionWindow, now time.Time) bool {
	return window.Status == WindowStatusPending && now.After(window.EndTime)
}

func (s *TWAPScheduler) isOrderComplete(order *TWAPOrder) bool {
	completedCount := 0
	for _, window := range order.ExecutionWindows {
		if window.Status == WindowStatusCompleted {
			completedCount++
		}
	}
	return completedCount >= order.IntervalCount
}

func (s *TWAPScheduler) handleExecutionError(task *ExecutionTask, err error, duration time.Duration) {
	s.logger.Error("TWAP window execution failed",
		zap.String("order_id", task.Order.ID),
		zap.Int("window", task.Window.Index),
		zap.Int("retry", task.Retry),
		zap.Duration("duration", duration),
		zap.Error(err),
	)

	// Retry logic
	if task.Retry < MaxRetryAttempts {
		task.Retry++
		task.Window.Status = WindowStatusPending // Reset to pending for retry

		// Schedule retry with delay
		go func() {
			time.Sleep(RetryDelay)
			select {
			case s.executionQueue <- task:
			default:
				s.logger.Error("Failed to schedule retry, queue full",
					zap.String("order_id", task.Order.ID),
					zap.Int("window", task.Window.Index),
				)
			}
		}()

		s.logger.Info("Scheduled retry for failed window",
			zap.String("order_id", task.Order.ID),
			zap.Int("window", task.Window.Index),
			zap.Int("retry", task.Retry),
		)
	} else {
		// Mark window as failed after max retries
		task.Window.Status = WindowStatusFailed
		s.logger.Error("Window permanently failed after max retries",
			zap.String("order_id", task.Order.ID),
			zap.Int("window", task.Window.Index),
			zap.Int("max_retries", MaxRetryAttempts),
		)
	}
}

func (s *TWAPScheduler) handleExecutionSuccess(task *ExecutionTask, duration time.Duration, workerID int) {
	executedAt := time.Now()
	task.Window.ExecutedAt = &executedAt

	s.logger.Info("TWAP window execution successful",
		zap.Int("worker_id", workerID),
		zap.String("order_id", task.Order.ID),
		zap.Int("window", task.Window.Index),
		zap.Duration("duration", duration),
		zap.String("amount", task.Window.Amount.String()),
	)
}

func (s *TWAPScheduler) healthMonitor(ctx context.Context) {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.performHealthCheck()
		}
	}
}

func (s *TWAPScheduler) performHealthCheck() {
	s.mutex.RLock()
	orderCount := len(s.scheduledOrders)
	queueSize := len(s.executionQueue)
	s.mutex.RUnlock()

	s.logger.Debug("TWAP Scheduler health check",
		zap.Int("scheduled_orders", orderCount),
		zap.Int("queue_size", queueSize),
		zap.Int("max_queue", cap(s.executionQueue)),
	)

	// Alert if queue is nearly full
	if queueSize > cap(s.executionQueue)*80/100 {
		s.logger.Warn("Execution queue is nearly full",
			zap.Int("current_size", queueSize),
			zap.Int("max_size", cap(s.executionQueue)),
		)
	}
}

func (s *TWAPScheduler) GetScheduledOrderCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.scheduledOrders)
}

func (s *TWAPScheduler) GetQueueSize() int {
	return len(s.executionQueue)
}

func (s *TWAPScheduler) GetSchedulerMetrics() map[string]interface{} {
	s.mutex.RLock()
	orderCount := len(s.scheduledOrders)
	s.mutex.RUnlock()

	queueSize := len(s.executionQueue)
	queueCapacity := cap(s.executionQueue)

	return map[string]interface{}{
		"scheduled_orders":    orderCount,
		"queue_size":         queueSize,
		"queue_capacity":     queueCapacity,
		"queue_utilization":  float64(queueSize) / float64(queueCapacity) * 100,
		"max_concurrent":     s.maxConcurrentOrders,
	}
}