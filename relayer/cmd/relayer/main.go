package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"flow-fusion/relayer/internal/bridge"
	"flow-fusion/relayer/internal/config"
	"flow-fusion/relayer/internal/cosmos"
	"flow-fusion/relayer/internal/ethereum"
	"flow-fusion/relayer/internal/twap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Flow Fusion Relayer")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Ethereum client
	ethClient, err := ethereum.NewClient(ethereum.Config{
		RPCURL:     cfg.Ethereum.RPCURL,
		PrivateKey: cfg.Ethereum.PrivateKey,
		ChainID:    cfg.Ethereum.ChainID,
	}, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Ethereum client", zap.Error(err))
	}

	// Initialize Cosmos client
	cosmosClient, err := cosmos.NewClient(cosmos.Config{
		RPCURL:     cfg.Cosmos.RPCURL,
		GRPCAddr:   cfg.Cosmos.GRPCAddr,
		ChainID:    cfg.Cosmos.ChainID,
		Mnemonic:   cfg.Cosmos.Mnemonic,
		Denom:      cfg.Cosmos.Denom,
	}, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Cosmos client", zap.Error(err))
	}

	// Initialize TWAP engine
	twapEngine, err := twap.NewEngine(twap.Config{
		OneInchAPIKey:    cfg.OneInch.APIKey,
		MaxIntervals:     cfg.TWAP.MaxIntervals,
		MinIntervalTime:  cfg.TWAP.MinIntervalTime,
		MaxSlippage:      cfg.TWAP.MaxSlippage,
		RedisURL:         cfg.Redis.URL,
		DatabaseURL:      cfg.Database.URL,
	}, ethClient, cosmosClient, logger)
	if err != nil {
		logger.Fatal("Failed to initialize TWAP engine", zap.Error(err))
	}

	// Initialize bridge service
	bridgeService, err := bridge.NewService(bridge.Config{
		EthereumClient: ethClient,
		CosmosClient:   cosmosClient,
		TWAPEngine:     twapEngine,
		Logger:         logger,
	})
	if err != nil {
		logger.Fatal("Failed to initialize bridge service", zap.Error(err))
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background services
	logger.Info("Starting background services")
	
	go func() {
		if err := twapEngine.Start(ctx); err != nil {
			logger.Error("TWAP engine error", zap.Error(err))
		}
	}()

	go func() {
		if err := bridgeService.Start(ctx); err != nil {
			logger.Error("Bridge service error", zap.Error(err))
		}
	}()

	// Setup HTTP server
	router := setupRouter(bridgeService, twapEngine, logger)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start HTTP server
	go func() {
		logger.Info("Starting HTTP server", zap.Int("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Flow Fusion Relayer")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	// Cancel background services
	cancel()

	logger.Info("Flow Fusion Relayer stopped")
}

func setupRouter(bridge *bridge.Service, twap *twap.Engine, logger *zap.Logger) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Logging middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("ip", param.ClientIP),
		)
		return ""
	}))

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().Unix(),
				"version":   "1.0.0",
			})
		})

		// Bridge endpoints
		bridgeGroup := v1.Group("/bridge")
		{
			bridgeGroup.POST("/quote", bridge.GetQuote)
			bridgeGroup.POST("/order", bridge.CreateOrder)
			bridgeGroup.GET("/order/:id", bridge.GetOrder)
			bridgeGroup.GET("/orders", bridge.ListOrders)
			bridgeGroup.POST("/order/:id/cancel", bridge.CancelOrder)
		}

		// TWAP endpoints
		twapGroup := v1.Group("/twap")
		{
			twapGroup.POST("/quote", twap.GetQuote)
			twapGroup.POST("/order", twap.CreateOrder)
			twapGroup.GET("/order/:id", twap.GetOrder)
			twapGroup.GET("/order/:id/schedule", twap.GetSchedule)
			twapGroup.GET("/orders", twap.ListOrders)
			twapGroup.POST("/order/:id/cancel", twap.CancelOrder)
		}

		// Metrics and monitoring
		v1.GET("/metrics", bridge.GetMetrics)
		v1.GET("/status", func(c *gin.Context) {
			status := map[string]interface{}{
				"ethereum": bridge.GetEthereumStatus(),
				"cosmos":   bridge.GetCosmosStatus(),
				"twap":     twap.GetStatus(),
			}
			c.JSON(http.StatusOK, status)
		})
	}

	// WebSocket endpoint for real-time updates
	router.GET("/ws", bridge.HandleWebSocket)

	return router
}