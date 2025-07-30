package twap

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ExecutionScheduler struct {
	logger      *zap.Logger
	scheduledOrders map[string]*TWAPOrder
	mutex       sync.RWMutex
	ticker      *time.Ticker
}

type Metrics struct {
	StartTime           time.Time
	TotalOrders         int64
	CompletedOrders     int64
	FailedOrders        int64
	TotalWindowsExecuted int64
	AverageExecutionTime time.Duration
	mutex               sync.RWMutex
}

func NewExecutionScheduler(logger *zap.Logger) *ExecutionScheduler {
	return &ExecutionScheduler{
		logger:          logger,
		scheduledOrders: make(map[string]*TWAPOrder),
		ticker:          time.NewTicker(5 * time.Second),
	}
}

func (s *ExecutionScheduler) Start(ctx context.Context) {
	s.logger.Info("Starting TWAP execution scheduler")
	
	defer s.ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("TWAP scheduler stopping")
			return
		case <-s.ticker.C:
			s.checkScheduledExecutions()
		}
	}
}

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

func (s *ExecutionScheduler) RemoveOrder(orderID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	delete(s.scheduledOrders, orderID)
}

func (s *ExecutionScheduler) checkScheduledExecutions() {
	now := time.Now()
	
	s.mutex.RLock()
	orders := make([]*TWAPOrder, 0, len(s.scheduledOrders))
	for _, order := range s.scheduledOrders {
		orders = append(orders, order)
	}
	s.mutex.RUnlock()
	
	for _, order := range orders {
		if order.Status == OrderStatusCompleted || order.Status == OrderStatusFailed || order.Status == OrderStatusCancelled {
			s.RemoveOrder(order.ID)
			continue
		}
		
		// Check if any windows are ready for execution
		for _, window := range order.ExecutionPlan {
			if window.Status == WindowStatusPending && 
			   now.After(window.StartTime) && 
			   now.Before(window.EndTime.Add(30*time.Second)) { // 30s grace period
				s.logger.Debug("Window ready for execution",
					zap.String("order_id", order.ID),
					zap.Int("window", window.Index),
					zap.Time("start_time", window.StartTime),
				)
			}
		}
	}
}

func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

func (m *Metrics) RecordWindowExecution(orderID string, windowIndex int, duration time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.TotalWindowsExecuted++
	
	// Update average execution time
	if m.TotalWindowsExecuted == 1 {
		m.AverageExecutionTime = duration
	} else {
		m.AverageExecutionTime = time.Duration(
			(int64(m.AverageExecutionTime)*m.TotalWindowsExecuted + int64(duration)) / (m.TotalWindowsExecuted + 1),
		)
	}
}

func (m *Metrics) RecordOrderCompletion(orderID string, success bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.TotalOrders++
	if success {
		m.CompletedOrders++
	} else {
		m.FailedOrders++
	}
}

func (m *Metrics) GetSummary() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	successRate := float64(0)
	if m.TotalOrders > 0 {
		successRate = float64(m.CompletedOrders) / float64(m.TotalOrders) * 100
	}
	
	return map[string]interface{}{
		"total_orders":           m.TotalOrders,
		"completed_orders":       m.CompletedOrders,
		"failed_orders":          m.FailedOrders,
		"success_rate":          successRate,
		"total_windows_executed": m.TotalWindowsExecuted,
		"avg_execution_time":     m.AverageExecutionTime.String(),
		"uptime":                time.Since(m.StartTime).String(),
	}
}