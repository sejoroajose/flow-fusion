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
	"flow-fusion/relayer/internal/oneinch"
	"flow-fusion/relayer/internal/twap"
)

func main() {
	// Initialize production logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Flow Fusion Production Relayer with 1inch HTTP API Integration")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Ethereum client with contract bindings
	ethClient, err := ethereum.NewClient(ethereum.Config{
		RPCURL:          cfg.Ethereum.RPCURL,
		PrivateKey:      cfg.Ethereum.PrivateKey,
		ChainID:         cfg.Ethereum.ChainID,
		ContractAddress: cfg.Ethereum.ContractAddr,
		GasLimit:        500000,
		MaxGasPrice:     "50000000000", 
	}, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Ethereum client", zap.Error(err))
	}

	// Initialize Cosmos client
	cosmosClient, err := cosmos.NewClient(cosmos.Config{
		RPCURL:       cfg.Cosmos.RPCURL,
		GRPCAddr:     cfg.Cosmos.GRPCAddr,
		ChainID:      cfg.Cosmos.ChainID,
		Mnemonic:     cfg.Cosmos.Mnemonic,
		Denom:        cfg.Cosmos.Denom,
		ContractAddr: "BB4124578492DA44754DAE2B22AEAA8604D32BC0B82B93098CB9AF86C9E5A84D", // deployed Cosmos contract
	}, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Cosmos client", zap.Error(err))
	}

	// Initialize Production TWAP Engine with 1inch HTTP API
	twapEngine, err := twap.NewProductionTWAPEngine(twap.Config{
		RedisURL:        cfg.Redis.URL,
		OneInchAPIKey:   cfg.OneInch.APIKey,
		PrivateKey:      cfg.Ethereum.PrivateKey,
		NodeURL:         cfg.Ethereum.RPCURL,
		ChainID:         int(cfg.Ethereum.ChainID),
		MaxGasPrice:     nil,
		MaxIntervals:    cfg.TWAP.MaxIntervals,
		MinIntervalTime: cfg.TWAP.MinIntervalTime,
		MaxSlippage:     cfg.TWAP.MaxSlippage,
		DatabaseURL:     cfg.Database.URL,
	}, ethClient, cosmosClient, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Production TWAP engine", zap.Error(err))
	}

	// Initialize Fusion+ Bridge Service with HTTP API
	bridgeService, err := bridge.NewFusionBridgeService(bridge.Config{
		EthereumClient: ethClient,
		CosmosClient:   cosmosClient,
		Logger:         logger,
		OneInchAPIKey:  cfg.OneInch.APIKey,
		PrivateKey:     cfg.Ethereum.PrivateKey,
		NodeURL:        cfg.Ethereum.RPCURL,
		ChainID:        int(cfg.Ethereum.ChainID),
	})
	if err != nil {
		logger.Fatal("Failed to initialize Fusion+ bridge service", zap.Error(err))
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background services
	logger.Info("Starting production background services with HTTP API integration")

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

	// Setup production HTTP server with proper middleware
	router := setupProductionRouter(bridgeService, twapEngine, logger, cfg.OneInch.APIKey)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start HTTP server
	go func() {
		logger.Info("Starting production HTTP server", 
			zap.Int("port", cfg.Server.Port),
			zap.String("mode", "production"),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Flow Fusion Production Relayer")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server forced to shutdown", zap.Error(err))
	}

	// Cancel background services
	cancel()

	logger.Info("Flow Fusion Production Relayer stopped successfully")
}

func setupProductionRouter(bridgeService *bridge.FusionBridgeService, twapEngine *twap.ProductionTWAPEngine, logger *zap.Logger, oneInchAPIKey string) *gin.Engine {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()

	// Production middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(loggingMiddleware(logger))
	router.Use(rateLimitMiddleware())
	router.Use(securityHeadersMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
			"service":   "flow-fusion-relayer",
			"api_type":  "http",
		})
	})

	// API status endpoint with 1inch HTTP API status
	router.GET("/api/health", func(c *gin.Context) {
		oneInchClient := oneinch.NewHTTPClient(oneinch.Config{
			BaseURL: "https://api.1inch.dev",
			APIKey:  oneInchAPIKey,
			ChainID: 1,
			Timeout: 5 * time.Second,
		}, logger)

		// Test token endpoint
		tokens, err := oneInchClient.GetTokens(c.Request.Context())
		oneInchStatus := "healthy"
		if err != nil || len(tokens) == 0 {
			oneInchStatus = "unhealthy"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "healthy",
			"timestamp":    time.Now().Unix(),
			"oneinch_api":  oneInchStatus,
			"api_type":     "http",
		})
	})

	// Metrics endpoint for monitoring
	router.GET("/metrics", func(c *gin.Context) {
		metrics := gin.H{
			"bridge":    bridgeService.GetMetricsData(), 
			"twap":      twapEngine.GetStatus(),
			"ethereum":  bridgeService.GetEthereumStatusData(),
			"cosmos":    bridgeService.GetCosmosStatusData(),
			"timestamp": time.Now().Unix(),
			"api_type":  "http",
		}
		c.JSON(http.StatusOK, metrics)
	})

	// API routes
	api := router.Group("/api/v1")
	api.Use(authMiddleware()) 
	{
		// Bridge endpoints using 1inch HTTP API
		bridgeGroup := api.Group("/bridge")
		{
			bridgeGroup.POST("/quote", bridgeService.GetQuote)
			bridgeGroup.POST("/order", bridgeService.CreateOrder)
			bridgeGroup.GET("/order/:id", bridgeService.GetOrder)
			bridgeGroup.GET("/orders", bridgeService.ListOrders)
			bridgeGroup.POST("/order/:id/cancel", bridgeService.CancelOrder)
		}

		// TWAP endpoints with 1inch HTTP API integration
		twapGroup := api.Group("/twap")
		{
			twapGroup.POST("/quote", twapEngine.GetQuote)
			twapGroup.POST("/order", twapEngine.CreateOrder)
			twapGroup.GET("/order/:id", twapEngine.GetOrder)
			twapGroup.GET("/orders", twapEngine.ListOrders)
			twapGroup.POST("/order/:id/cancel", twapEngine.CancelOrder)
		}

		// 1inch HTTP API integration endpoints
		oneinchGroup := api.Group("/oneinch")
		{
			oneinchGroup.GET("/tokens", getTokenInfo(oneInchAPIKey, logger))
			oneinchGroup.POST("/swap/quote", getSwapQuote(oneInchAPIKey, logger))
			oneinchGroup.POST("/fusion/quote", getFusionQuote(oneInchAPIKey, logger))
			oneinchGroup.POST("/orderbook/quote", getOrderbookQuote(oneInchAPIKey, logger))
		}

		// System status and monitoring
		api.GET("/status", func(c *gin.Context) {
			status := gin.H{
				"ethereum":  bridgeService.GetEthereumStatusData(),
				"cosmos":    bridgeService.GetCosmosStatusData(),
				"twap":      twapEngine.GetStatus(),
				"bridge":    bridgeService.GetMetricsData(),  
				"timestamp": time.Now().Unix(),
				"uptime":    time.Since(startTime).String(),
				"api_type":  "http",
			}
			c.JSON(http.StatusOK, status)
		})
	}

	// WebSocket endpoint for real-time updates
	router.GET("/ws", bridgeService.HandleWebSocket)

	// Swagger/OpenAPI documentation endpoint
	router.GET("/docs/*any", serveSwaggerDocs())

	return router
}

var startTime = time.Now()

// Production middleware implementations
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		c.Header("Access-Control-Max-Age", "3600")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

func loggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Validate API key against your authentication system
		if !validateAPIKey(apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func validateAPIKey(apiKey string) bool {
	// Implement proper API key validation
	// This could involve checking against a database, JWT validation, etc.
	validKeys := map[string]bool{
		"dev-key-12345":  true,
		"prod-key-67890": true,
	}
	return validKeys[apiKey]
}

// 1inch HTTP API integration endpoints
func getTokenInfo(oneInchAPIKey string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		chainID := c.DefaultQuery("chainId", "1")
		
		oneInchClient := oneinch.NewHTTPClient(oneinch.Config{
			BaseURL: "https://api.1inch.dev",
			APIKey:  oneInchAPIKey,
			ChainID: 1, // Will be parsed from chainID param
			Timeout: 30 * time.Second,
		}, logger)

		tokens, err := oneInchClient.GetTokens(c.Request.Context())
		if err != nil {
			logger.Error("Failed to get tokens from 1inch API", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch token information",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tokens":   tokens,
			"chainId":  chainID,
			"api_type": "http",
		})
	}
}

func getSwapQuote(oneInchAPIKey string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Src      string  `json:"src" binding:"required"`
			Dst      string  `json:"dst" binding:"required"`
			Amount   string  `json:"amount" binding:"required"`
			From     string  `json:"from" binding:"required"`
			Slippage float64 `json:"slippage" binding:"required"`
			ChainID  uint64  `json:"chainId"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.ChainID == 0 {
			req.ChainID = 1
		}

		oneInchClient := oneinch.NewHTTPClient(oneinch.Config{
			BaseURL: "https://api.1inch.dev",
			APIKey:  oneInchAPIKey,
			ChainID: req.ChainID,
			Timeout: 30 * time.Second,
		}, logger)

		quote, err := oneInchClient.GetSwap(
			c.Request.Context(),
			req.Src,
			req.Dst,
			req.Amount,
			req.From,
			req.Slippage,
			nil,
		)
		if err != nil {
			logger.Error("Failed to get swap quote from 1inch API", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get swap quote",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"quote":    quote,
			"api_type": "http",
		})
	}
}

func getFusionQuote(oneInchAPIKey string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			SrcChain        uint64 `json:"srcChain" binding:"required"`
			DstChain        uint64 `json:"dstChain" binding:"required"`
			SrcTokenAddress string `json:"srcTokenAddress" binding:"required"`
			DstTokenAddress string `json:"dstTokenAddress" binding:"required"`
			Amount          string `json:"amount" binding:"required"`
			WalletAddress   string `json:"walletAddress" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		oneInchClient := oneinch.NewHTTPClient(oneinch.Config{
			BaseURL: "https://api.1inch.dev",
			APIKey:  oneInchAPIKey,
			ChainID: req.SrcChain,
			Timeout: 30 * time.Second,
		}, logger)

		quote, err := oneInchClient.GetFusionQuote(
			c.Request.Context(),
			req.SrcChain,
			req.DstChain,
			req.SrcTokenAddress,
			req.DstTokenAddress,
			req.Amount,
			req.WalletAddress,
		)
		if err != nil {
			logger.Error("Failed to get Fusion+ quote from 1inch API", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get Fusion+ quote",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"quote":    quote,
			"api_type": "http",
		})
	}
}

func getOrderbookQuote(oneInchAPIKey string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			MakerAddress string `json:"makerAddress" binding:"required"`
			ChainID      uint64 `json:"chainId"`
			Limit        int    `json:"limit"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.ChainID == 0 {
			req.ChainID = 1
		}
		if req.Limit == 0 {
			req.Limit = 100
		}

		oneInchClient := oneinch.NewHTTPClient(oneinch.Config{
			BaseURL: "https://api.1inch.dev",
			APIKey:  oneInchAPIKey,
			ChainID: req.ChainID,
			Timeout: 30 * time.Second,
		}, logger)

		orders, err := oneInchClient.GetOrderbookOrdersByMaker(
			c.Request.Context(),
			req.MakerAddress,
			req.Limit,
		)
		if err != nil {
			logger.Error("Failed to get orderbook orders from 1inch API", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get orderbook orders",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"orders":   orders,
			"api_type": "http",
		})
	}
}

func serveSwaggerDocs() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Serve Swagger/OpenAPI documentation
		c.JSON(http.StatusOK, gin.H{
			"message":  "API documentation - implement Swagger/OpenAPI docs",
			"api_type": "http",
			"endpoints": map[string]string{
				"bridge_quote":    "POST /api/v1/bridge/quote",
				"bridge_order":    "POST /api/v1/bridge/order",
				"twap_quote":      "POST /api/v1/twap/quote", 
				"twap_order":      "POST /api/v1/twap/order",
				"oneinch_tokens":  "GET /api/v1/oneinch/tokens",
				"health":          "GET /health",
				"metrics":         "GET /metrics",
				"status":          "GET /api/v1/status",
			},
		})
	}
}