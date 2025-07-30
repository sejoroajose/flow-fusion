package twap

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ExecutionScheduler manages the scheduling and execution of TWAP orders
type ExecutionScheduler struct {
	logger          *zap.Logger
	engine          *Engine
	scheduledOrders map[string]*TWAPOrder
	mutex           sync.RWMutex
	ticker          *time.Ticker
	workerPool      chan struct{}
	stopChan        chan struct{}
}

// Metrics tracks TWAP engine performance metrics
type Metrics struct {
	StartTime            time.Time
	TotalOrders          int64
	CompletedOrders      int64
	FailedOrders         int64
	TotalWindowsExecuted int64
	FailedWindows        int64
	AverageExecutionTime time.Duration
	TotalGasUsed         uint64
	LastExecutionTime    time.Time
	mutex                sync.RWMutex
}

// NewExecutionScheduler creates a new execution scheduler
func NewExecutionScheduler(logger *zap.Logger, engine *Engine) *ExecutionScheduler {
	return &ExecutionScheduler{
		logger:          logger,
		engine:          engine,
		scheduledOrders: make(map[string]*TWAPOrder),
		ticker:          time.NewTicker(5 * time.Second),
		workerPool:      make(chan struct{}, MaxConcurrentSwaps),
		stopChan:        make(chan struct{}),
	}
}

// Start begins the execution scheduler main loop
func (s *ExecutionScheduler) Start(ctx context.Context) {
	s.logger.Info("Starting TWAP execution scheduler")
	
	defer s.ticker.Stop()
	defer close(s.workerPool)
	defer close(s.stopChan)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("TWAP scheduler stopping")
			return
		case <-s.stopChan:
			s.logger.Info("TWAP scheduler received stop signal")
			return
		case <-s.ticker.C:
			if err := s.checkScheduledExecutions(ctx); err != nil {
				s.logger.Error("Error checking scheduled executions", zap.Error(err))
			}
		}
	}
}

// Stop gracefully stops the scheduler
func (s *ExecutionScheduler) Stop() {
	select {
	case s.stopChan <- struct{}{}:
	default:
		// Channel is full or closed, which is fine
	}
}

// ScheduleOrder adds an order to the execution schedule
func (s *ExecutionScheduler) ScheduleOrder(order *TWAPOrder) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.scheduledOrders[order.ID] = order
	s.logger.Info("Order scheduled for execution",
		zap.String("order_id", order.ID),
		zap.Time("start_time", order.StartTime),
		zap.Int("intervals", len(order.ExecutionPlan)),
	)
}

// RemoveOrder removes an order from the schedule
func (s *ExecutionScheduler) RemoveOrder(orderID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	delete(s.scheduledOrders, orderID)
	s.logger.Info("Order removed from scheduler",
		zap.String("order_id", orderID),
	)
}

// checkScheduledExecutions checks for ready-to-execute windows
func (s *ExecutionScheduler) checkScheduledExecutions(ctx context.Context) error {
	now := time.Now()
	
	s.mutex.RLock()
	orders := make([]*TWAPOrder, 0, len(s.scheduledOrders))
	for _, order := range s.scheduledOrders {
		orders = append(orders, order)
	}
	s.mutex.RUnlock()

	var wg sync.WaitGroup
	
	for _, order := range orders {
		// Remove completed or failed orders from schedule
		if order.Status == OrderStatusCompleted || order.Status == OrderStatusFailed || order.Status == OrderStatusCancelled {
			s.RemoveOrder(order.ID)
			continue
		}

		// Start order execution if it's time
		if order.Status == OrderStatusPending && now.After(order.StartTime) {
			s.engine.ordersMutex.Lock()
			order.Status = OrderStatusExecuting
			s.engine.ordersMutex.Unlock()
			
			s.logger.Info("Starting order execution",
				zap.String("order_id", order.ID),
				zap.Time("scheduled_start", order.StartTime),
			)
		}

		// Check each execution window
		for _, window := range order.ExecutionPlan {
			if window.Status != WindowStatusPending {
				continue
			}

			// Check if window is ready for execution
			if now.After(window.StartTime) && now.Before(window.EndTime) {
				// Try to acquire a worker from the pool
				select {
				case s.workerPool <- struct{}{}:
					wg.Add(1)
					go func(o *TWAPOrder, w *ExecutionWindow) {
						defer wg.Done()
						defer func() { <-s.workerPool }()
						
						s.executeWindowWithLogging(ctx, o, w)
					}(order, window)
				default:
					// Worker pool is full, log and skip this window for now
					s.logger.Warn("Worker pool full, delaying window execution",
						zap.String("order_id", order.ID),
						zap.Int("window", window.Index),
						zap.Time("window_start", window.StartTime),
					)
				}
			} else if now.After(window.EndTime) {
				// Window has expired, mark as skipped
				window.Status = WindowStatusSkipped
				s.logger.Warn("Window execution window expired",
					zap.String("order_id", order.ID),
					zap.Int("window", window.Index),
					zap.Time("window_end", window.EndTime),
				)
			}
		}

		// Check if order is complete
		if s.isOrderComplete(order) {
			s.engine.ordersMutex.Lock()
			order.Status = OrderStatusCompleted
			s.engine.ordersMutex.Unlock()
			
			s.engine.metrics.RecordOrderCompletion(order.ID, true)
			s.logger.Info("Order completed",
				zap.String("order_id", order.ID),
				zap.Int("completed_windows", order.CompletedWindows),
				zap.Int("total_windows", order.IntervalCount),
			)
		}
	}

	// Wait for all workers to complete before returning
	wg.Wait()
	return nil
}

// executeWindowWithLogging executes a window with proper logging and metrics
func (s *ExecutionScheduler) executeWindowWithLogging(ctx context.Context, order *TWAPOrder, window *ExecutionWindow) {
	startTime := time.Now()
	
	s.logger.Info("Starting window execution",
		zap.String("order_id", order.ID),
		zap.Int("window", window.Index),
		zap.String("amount", window.Amount.String()),
	)

	err := s.engine.executeWindow(ctx, order, window)
	duration := time.Since(startTime)
	
	if err != nil {
		s.logger.Error("Window execution failed",
			zap.String("order_id", order.ID),
			zap.Int("window", window.Index),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		s.engine.metrics.RecordWindowFailure(order.ID, window.Index)
	} else {
		s.logger.Info("Window execution completed",
			zap.String("order_id", order.ID),
			zap.Int("window", window.Index),
			zap.Duration("duration", duration),
			zap.Uint64("gas_used", window.GasUsed),
		)
		s.engine.metrics.RecordWindowExecution(order.ID, window.Index, duration, window.GasUsed)
	}
}

// isOrderComplete checks if all windows in an order are completed
func (s *ExecutionScheduler) isOrderComplete(order *TWAPOrder) bool {
	completedCount := 0
	for _, window := range order.ExecutionPlan {
		if window.Status == WindowStatusCompleted {
			completedCount++
		}
	}
	return completedCount >= order.IntervalCount
}

// GetScheduledOrderCount returns the number of currently scheduled orders
func (s *ExecutionScheduler) GetScheduledOrderCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.scheduledOrders)
}

// GetScheduledOrders returns a copy of currently scheduled orders
func (s *ExecutionScheduler) GetScheduledOrders() map[string]*TWAPOrder {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	orders := make(map[string]*TWAPOrder, len(s.scheduledOrders))
	for id, order := range s.scheduledOrders {
		orders[id] = order
	}
	return orders
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

// RecordOrderScheduled increments the total orders counter
func (m *Metrics) RecordOrderScheduled() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.TotalOrders++
}

// RecordWindowExecution records a successful window execution
func (m *Metrics) RecordWindowExecution(orderID string, windowIndex int, duration time.Duration, gasUsed uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.TotalWindowsExecuted++
	m.TotalGasUsed += gasUsed
	m.LastExecutionTime = time.Now()
	
	// Update average execution time using a running average
	if m.TotalWindowsExecuted == 1 {
		m.AverageExecutionTime = duration
	} else {
		// Calculate weighted average
		totalTime := int64(m.AverageExecutionTime) * (m.TotalWindowsExecuted - 1)
		m.AverageExecutionTime = time.Duration((totalTime + int64(duration)) / m.TotalWindowsExecuted)
	}
}

// RecordWindowFailure records a failed window execution
func (m *Metrics) RecordWindowFailure(orderID string, windowIndex int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.FailedWindows++
}

// RecordOrderCompletion records order completion or failure
func (m *Metrics) RecordOrderCompletion(orderID string, success bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if success {
		m.CompletedOrders++
	} else {
		m.FailedOrders++
	}
}

// GetSummary returns a summary of current metrics
func (m *Metrics) GetSummary() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	successRate := float64(0)
	if m.TotalOrders > 0 {
		successRate = float64(m.CompletedOrders) / float64(m.TotalOrders) * 100
	}
	
	windowSuccessRate := float64(0)
	if m.TotalWindowsExecuted+m.FailedWindows > 0 {
		windowSuccessRate = float64(m.TotalWindowsExecuted) / float64(m.TotalWindowsExecuted+m.FailedWindows) * 100
	}
	
	return map[string]interface{}{
		"total_orders":           m.TotalOrders,
		"completed_orders":       m.CompletedOrders,
		"failed_orders":          m.FailedOrders,
		"order_success_rate":     successRate,
		"total_windows_executed": m.TotalWindowsExecuted,
		"failed_windows":         m.FailedWindows,
		"window_success_rate":    windowSuccessRate,
		"avg_execution_time":     m.AverageExecutionTime.String(),
		"total_gas_used":         m.TotalGasUsed,
		"last_execution":         m.LastExecutionTime.Format(time.RFC3339),
		"uptime":                 time.Since(m.StartTime).String(),
	}
}

// GetDetailedMetrics returns detailed metrics for monitoring
func (m *Metrics) GetDetailedMetrics() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	summary := m.GetSummary()
	
	// Add more detailed information
	summary["start_time"] = m.StartTime.Format(time.RFC3339)
	summary["avg_execution_seconds"] = m.AverageExecutionTime.Seconds()
	summary["uptime_seconds"] = time.Since(m.StartTime).Seconds()
	
	// Calculate rates
	uptimeHours := time.Since(m.StartTime).Hours()
	if uptimeHours > 0 {
		summary["orders_per_hour"] = float64(m.TotalOrders) / uptimeHours
		summary["windows_per_hour"] = float64(m.TotalWindowsExecuted) / uptimeHours
		summary["gas_per_hour"] = float64(m.TotalGasUsed) / uptimeHours
	}
	
	return summary
}

// Reset resets all metrics (useful for testing)
func (m *Metrics) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.StartTime = time.Now()
	m.TotalOrders = 0
	m.CompletedOrders = 0
	m.FailedOrders = 0
	m.TotalWindowsExecuted = 0
	m.FailedWindows = 0
	m.AverageExecutionTime = 0
	m.TotalGasUsed = 0
	m.LastExecutionTime = time.Time{}
}