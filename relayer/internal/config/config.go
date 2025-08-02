package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"log"

	"github.com/joho/godotenv"
	vault "github.com/hashicorp/vault/api"
	"go.uber.org/zap"
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Ethereum   EthereumConfig   `json:"ethereum"`
	Cosmos     CosmosConfig     `json:"cosmos"`
	OneInch    OneInchConfig    `json:"oneinch"`
	TWAP       TWAPConfig       `json:"twap"`
	Redis      RedisConfig      `json:"redis"`
	Database   DatabaseConfig   `json:"database"`
	Security   SecurityConfig   `json:"security"`
	Monitoring MonitoringConfig `json:"monitoring"`
	Logging    LoggingConfig    `json:"logging"`
}

type ServerConfig struct {
	Port            int           `json:"port"`
	Host            string        `json:"host"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	MaxHeaderBytes  int           `json:"max_header_bytes"`
	TLSCertFile     string        `json:"tls_cert_file"`
	TLSKeyFile      string        `json:"tls_key_file"`
	EnableHTTPS     bool          `json:"enable_https"`
}

type EthereumConfig struct {
	RPCURL            string        `json:"rpc_url"`
	WSEndpoint        string        `json:"ws_endpoint"`
	PrivateKey        string        `json:"private_key"`
	ChainID           int64         `json:"chain_id"`
	ContractAddr      string        `json:"contract_addr"`
	GasLimit          uint64        `json:"gas_limit"`
	MaxGasPrice       string        `json:"max_gas_price"`
	GasPriceMultiplier float64      `json:"gas_price_multiplier"`
	ConfirmationBlocks uint64       `json:"confirmation_blocks"`
	RetryAttempts     int           `json:"retry_attempts"`
	RetryDelay        time.Duration `json:"retry_delay"`
	PoolInterval      time.Duration `json:"pool_interval"`
}

type CosmosConfig struct {
	RPCURL        string        `json:"rpc_url"`
	GRPCAddr      string        `json:"grpc_addr"`
	ChainID       string        `json:"chain_id"`
	Mnemonic      string        `json:"mnemonic"`
	Denom         string        `json:"denom"`
	ContractAddr  string        `json:"contract_addr"`
	GasAdjustment float64       `json:"gas_adjustment"`
	GasPrices     string        `json:"gas_prices"`
	KeyName       string        `json:"key_name"`
	RetryAttempts int           `json:"retry_attempts"`
	RetryDelay    time.Duration `json:"retry_delay"`
}

type OneInchConfig struct {
	APIKey        string        `json:"api_key"`
	BaseURL       string        `json:"base_url"`
	Timeout       time.Duration `json:"timeout"`
	RetryAttempts int           `json:"retry_attempts"`
	RetryDelay    time.Duration `json:"retry_delay"`
	RateLimit     int           `json:"rate_limit"`
}

type TWAPConfig struct {
	MaxIntervals     int           `json:"max_intervals"`
	MinIntervalTime  time.Duration `json:"min_interval_time"`
	MaxSlippage      float64       `json:"max_slippage"`
	MaxConcurrent    int           `json:"max_concurrent"`
	ExecutionTimeout time.Duration `json:"execution_timeout"`
	RetryAttempts    int           `json:"retry_attempts"`
	RetryDelay       time.Duration `json:"retry_delay"`
}

type RedisConfig struct {
	URL              string        `json:"url"`
	Password         string        `json:"password"`
	DB               int           `json:"db"`
	MaxRetries       int           `json:"max_retries"`
	MinRetryBackoff  time.Duration `json:"min_retry_backoff"`
	MaxRetryBackoff  time.Duration `json:"max_retry_backoff"`
	DialTimeout      time.Duration `json:"dial_timeout"`
	ReadTimeout      time.Duration `json:"read_timeout"`
	WriteTimeout     time.Duration `json:"write_timeout"`
	PoolSize         int           `json:"pool_size"`
	MinIdleConns     int           `json:"min_idle_conns"`
	MaxConnAge       time.Duration `json:"max_conn_age"`
	PoolTimeout      time.Duration `json:"pool_timeout"`
	IdleTimeout      time.Duration `json:"idle_timeout"`
	IdleCheckFreq    time.Duration `json:"idle_check_freq"`
}

type DatabaseConfig struct {
	URL             string        `json:"url"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	RetryAttempts   int           `json:"retry_attempts"`
	RetryDelay      time.Duration `json:"retry_delay"`
}

type SecurityConfig struct {
	APIKeys           []string      `json:"api_keys"`
	JWTSecret         string        `json:"jwt_secret"`
	JWTExpiration     time.Duration `json:"jwt_expiration"`
	RateLimitPerIP    int           `json:"rate_limit_per_ip"`
	RateLimitWindow   time.Duration `json:"rate_limit_window"`
	MaxRequestSize    int64         `json:"max_request_size"`
	AllowedOrigins    []string      `json:"allowed_origins"`
	TrustedProxies    []string      `json:"trusted_proxies"`
	EnableCORS        bool          `json:"enable_cors"`
	EnableCSRF        bool          `json:"enable_csrf"`
}

type MonitoringConfig struct {
	EnableMetrics     bool          `json:"enable_metrics"`
	MetricsPort       int           `json:"metrics_port"`
	MetricsPath       string        `json:"metrics_path"`
	EnablePprof       bool          `json:"enable_pprof"`
	PprofPort         int           `json:"pprof_port"`
	HealthCheckPath   string        `json:"health_check_path"`
	ReadinessPath     string        `json:"readiness_path"`
	LivenessPath      string        `json:"liveness_path"`
	AlertWebhookURL   string        `json:"alert_webhook_url"`
	MetricsInterval   time.Duration `json:"metrics_interval"`
}

type LoggingConfig struct {
	Level            string `json:"level"`
	Format           string `json:"format"`
	Output           string `json:"output"`
	EnableStackTrace bool   `json:"enable_stack_trace"`
	EnableCaller     bool   `json:"enable_caller"`
	EnableSampling   bool   `json:"enable_sampling"`
	SamplingInitial  int    `json:"sampling_initial"`
	SamplingInterval int    `json:"sampling_interval"`
}

// Environment constants
const (
	EnvProduction  = "production"
	EnvStaging     = "staging"
	EnvDevelopment = "development"
	EnvTest        = "test"
)

func Load() (*Config, error) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	environment := getEnvString("ENVIRONMENT", EnvDevelopment)
	
	config := &Config{
		Server: ServerConfig{
			Port:           getEnvInt("SERVER_PORT", 8080),
			Host:           getEnvString("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:    getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:   getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:    getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
			MaxHeaderBytes: getEnvInt("SERVER_MAX_HEADER_BYTES", 1048576),
			TLSCertFile:    getEnvString("TLS_CERT_FILE", ""),
			TLSKeyFile:     getEnvString("TLS_KEY_FILE", ""),
			EnableHTTPS:    getEnvBool("ENABLE_HTTPS", false),
		},
		Ethereum: EthereumConfig{
			RPCURL:             getEnvString("ETHEREUM_RPC_URL", "http://localhost:8545"),
			WSEndpoint:         getEnvString("ETHEREUM_WS_ENDPOINT", "ws://localhost:8546"),
			ChainID:            int64(getEnvInt("ETHEREUM_CHAIN_ID", 1)),
			ContractAddr:       getEnvString("ETHEREUM_CONTRACT_ADDR", ""),
			GasLimit:           uint64(getEnvInt("ETHEREUM_GAS_LIMIT", 500000)),
			MaxGasPrice:        getEnvString("ETHEREUM_MAX_GAS_PRICE", "100000000000"), 
			GasPriceMultiplier: getEnvFloat("ETHEREUM_GAS_PRICE_MULTIPLIER", 1.1),
			ConfirmationBlocks: uint64(getEnvInt("ETHEREUM_CONFIRMATION_BLOCKS", 3)),
			RetryAttempts:      getEnvInt("ETHEREUM_RETRY_ATTEMPTS", 3),
			RetryDelay:         getEnvDuration("ETHEREUM_RETRY_DELAY", 3*time.Second),
			PoolInterval:       getEnvDuration("ETHEREUM_POOL_INTERVAL", 15*time.Second),
		},
		Cosmos: CosmosConfig{
			RPCURL:        getEnvString("COSMOS_RPC_URL", "http://localhost:26657"),
			GRPCAddr:      getEnvString("COSMOS_GRPC_ADDR", "localhost:9090"),
			ChainID:       getEnvString("COSMOS_CHAIN_ID", "cosmos-1"),
			Denom:         getEnvString("COSMOS_DENOM", "uatom"),
			ContractAddr:  getEnvString("COSMOS_CONTRACT_ADDR", ""),
			GasAdjustment: getEnvFloat("COSMOS_GAS_ADJUSTMENT", 1.5),
			GasPrices:     getEnvString("COSMOS_GAS_PRICES", "0.025uatom"),
			KeyName:       getEnvString("COSMOS_KEY_NAME", "relayer"),
			RetryAttempts: getEnvInt("COSMOS_RETRY_ATTEMPTS", 3),
			RetryDelay:    getEnvDuration("COSMOS_RETRY_DELAY", 2*time.Second),
		},
		OneInch: OneInchConfig{
			BaseURL:       getEnvString("ONEINCH_BASE_URL", "https://api.1inch.dev"),
			Timeout:       getEnvDuration("ONEINCH_TIMEOUT", 30*time.Second),
			RetryAttempts: getEnvInt("ONEINCH_RETRY_ATTEMPTS", 3),
			RetryDelay:    getEnvDuration("ONEINCH_RETRY_DELAY", 2*time.Second),
			RateLimit:     getEnvInt("ONEINCH_RATE_LIMIT", 5),
		},
		TWAP: TWAPConfig{
			MaxIntervals:     getEnvInt("TWAP_MAX_INTERVALS", 100),
			MinIntervalTime:  getEnvDuration("TWAP_MIN_INTERVAL_TIME", 60*time.Second),
			MaxSlippage:      getEnvFloat("TWAP_MAX_SLIPPAGE", 0.10), // 10%
			MaxConcurrent:    getEnvInt("TWAP_MAX_CONCURRENT", 20),
			ExecutionTimeout: getEnvDuration("TWAP_EXECUTION_TIMEOUT", 10*time.Minute),
			RetryAttempts:    getEnvInt("TWAP_RETRY_ATTEMPTS", 3),
			RetryDelay:       getEnvDuration("TWAP_RETRY_DELAY", 30*time.Second),
		},
		Redis: RedisConfig{
			URL:              getEnvString("REDIS_URL", "redis://localhost:6379"),
			Password:         getEnvString("REDIS_PASSWORD", ""),
			DB:               getEnvInt("REDIS_DB", 0),
			MaxRetries:       getEnvInt("REDIS_MAX_RETRIES", 3),
			MinRetryBackoff:  getEnvDuration("REDIS_MIN_RETRY_BACKOFF", 8*time.Millisecond),
			MaxRetryBackoff:  getEnvDuration("REDIS_MAX_RETRY_BACKOFF", 512*time.Millisecond),
			DialTimeout:      getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:      getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout:     getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
			PoolSize:         getEnvInt("REDIS_POOL_SIZE", 10),
			MinIdleConns:     getEnvInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxConnAge:       getEnvDuration("REDIS_MAX_CONN_AGE", 30*time.Minute),
			PoolTimeout:      getEnvDuration("REDIS_POOL_TIMEOUT", 4*time.Second),
			IdleTimeout:      getEnvDuration("REDIS_IDLE_TIMEOUT", 5*time.Minute),
			IdleCheckFreq:    getEnvDuration("REDIS_IDLE_CHECK_FREQ", 1*time.Minute),
		},
		Database: DatabaseConfig{
			URL:             getEnvString("DATABASE_URL", "postgresql://flow:fusion@localhost:5432/flowfusion"),
			MaxOpenConns:    getEnvInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DATABASE_CONN_MAX_LIFETIME", 1*time.Hour),
			ConnMaxIdleTime: getEnvDuration("DATABASE_CONN_MAX_IDLE_TIME", 10*time.Minute),
			RetryAttempts:   getEnvInt("DATABASE_RETRY_ATTEMPTS", 3),
			RetryDelay:      getEnvDuration("DATABASE_RETRY_DELAY", 1*time.Second),
		},
		Security: SecurityConfig{
			APIKeys:         getEnvStringSlice("API_KEYS", []string{}),
			JWTSecret:       getEnvString("JWT_SECRET", ""),
			JWTExpiration:   getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
			RateLimitPerIP:  getEnvInt("RATE_LIMIT_PER_IP", 1000),
			RateLimitWindow: getEnvDuration("RATE_LIMIT_WINDOW", 1*time.Hour),
			MaxRequestSize:  int64(getEnvInt("MAX_REQUEST_SIZE", 1048576)), // 1MB
			AllowedOrigins:  getEnvStringSlice("ALLOWED_ORIGINS", []string{"*"}),
			TrustedProxies:  getEnvStringSlice("TRUSTED_PROXIES", []string{}),
			EnableCORS:      getEnvBool("ENABLE_CORS", true),
			EnableCSRF:      getEnvBool("ENABLE_CSRF", false),
		},
		Monitoring: MonitoringConfig{
			EnableMetrics:   getEnvBool("ENABLE_METRICS", true),
			MetricsPort:     getEnvInt("METRICS_PORT", 9090),
			MetricsPath:     getEnvString("METRICS_PATH", "/metrics"),
			EnablePprof:     getEnvBool("ENABLE_PPROF", false),
			PprofPort:       getEnvInt("PPROF_PORT", 6060),
			HealthCheckPath: getEnvString("HEALTH_CHECK_PATH", "/health"),
			ReadinessPath:   getEnvString("READINESS_PATH", "/ready"),
			LivenessPath:    getEnvString("LIVENESS_PATH", "/live"),
			AlertWebhookURL: getEnvString("ALERT_WEBHOOK_URL", ""),
			MetricsInterval: getEnvDuration("METRICS_INTERVAL", 30*time.Second),
		},
		Logging: LoggingConfig{
			Level:            getEnvString("LOG_LEVEL", "info"),
			Format:           getEnvString("LOG_FORMAT", "json"),
			Output:           getEnvString("LOG_OUTPUT", "stdout"),
			EnableStackTrace: getEnvBool("LOG_ENABLE_STACK_TRACE", false),
			EnableCaller:     getEnvBool("LOG_ENABLE_CALLER", true),
			EnableSampling:   getEnvBool("LOG_ENABLE_SAMPLING", true),
			SamplingInitial:  getEnvInt("LOG_SAMPLING_INITIAL", 100),
			SamplingInterval: getEnvInt("LOG_SAMPLING_INTERVAL", 100),
		},
	}

	// Load secrets from secure sources
	if err := loadSecrets(config, environment); err != nil {
		return nil, fmt.Errorf("failed to load secrets: %w", err)
	}

	// Validate configuration
	if err := validate(config, environment); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadSecrets loads sensitive configuration from secure sources
func loadSecrets(config *Config, environment string) error {
	if environment == EnvProduction {
		if err := loadFromVault(config); err != nil {
			fmt.Printf("Warning: Failed to load from Vault, using environment variables: %v\n", err)
			loadSecretsFromEnv(config)
		}
	} else {
		loadSecretsFromEnv(config)
	}

	return nil
}

// loadFromVault loads secrets from HashiCorp Vault
func loadFromVault(config *Config) error {
	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		return fmt.Errorf("VAULT_ADDR not set")
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		return fmt.Errorf("VAULT_TOKEN not set")
	}

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = vaultAddr

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %w", err)
	}

	client.SetToken(vaultToken)

	// Read secrets from Vault
	secret, err := client.Logical().Read("secret/data/flow-fusion/relayer")
	if err != nil {
		return fmt.Errorf("failed to read secret from Vault: %w", err)
	}

	if secret == nil {
		return fmt.Errorf("no secret found at path secret/data/flow-fusion/relayer")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid secret format in Vault")
	}

	// Extract secrets
	if privateKey, ok := data["ethereum_private_key"].(string); ok {
		config.Ethereum.PrivateKey = privateKey
	}
	if mnemonic, ok := data["cosmos_mnemonic"].(string); ok {
		config.Cosmos.Mnemonic = mnemonic
	}
	if apiKey, ok := data["oneinch_api_key"].(string); ok {
		config.OneInch.APIKey = apiKey
	}
	if jwtSecret, ok := data["jwt_secret"].(string); ok {
		config.Security.JWTSecret = jwtSecret
	}
	if redisPassword, ok := data["redis_password"].(string); ok {
		config.Redis.Password = redisPassword
	}

	return nil
}

// loadSecretsFromEnv loads secrets from environment variables
func loadSecretsFromEnv(config *Config) {
	config.Ethereum.PrivateKey = getEnvString("ETHEREUM_PRIVATE_KEY", "")
	config.Cosmos.Mnemonic = getEnvString("COSMOS_MNEMONIC", "")
	config.OneInch.APIKey = getEnvString("ONEINCH_API_KEY", "")
	config.Security.JWTSecret = getEnvString("JWT_SECRET", "")
	config.Redis.Password = getEnvString("REDIS_PASSWORD", "")
}

// validate validates the configuration
func validate(config *Config, environment string) error {
	// Validate required fields
	if config.Ethereum.PrivateKey == "" {
		return fmt.Errorf("ETHEREUM_PRIVATE_KEY is required")
	}

	if config.OneInch.APIKey == "" {
		return fmt.Errorf("ONEINCH_API_KEY is required")
	}

	if config.Ethereum.ContractAddr == "" {
		return fmt.Errorf("ETHEREUM_CONTRACT_ADDR is required")
	}

	if config.Cosmos.ContractAddr == "" {
		return fmt.Errorf("COSMOS_CONTRACT_ADDR is required")
	}

	// Production-specific validations
	if environment == EnvProduction {
		if config.Security.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET is required in production")
		}

		if len(config.Security.APIKeys) == 0 {
			return fmt.Errorf("API_KEYS must be configured in production")
		}

		if !config.Server.EnableHTTPS {
			fmt.Println("Warning: HTTPS is not enabled in production")
		}
	}

	// Validate chain ID
	if config.Ethereum.ChainID <= 0 {
		return fmt.Errorf("invalid Ethereum chain ID: %d", config.Ethereum.ChainID)
	}

	// Validate TWAP configuration
	if config.TWAP.MaxIntervals <= 0 || config.TWAP.MaxIntervals > 1000 {
		return fmt.Errorf("invalid TWAP max intervals: %d", config.TWAP.MaxIntervals)
	}

	if config.TWAP.MaxSlippage <= 0 || config.TWAP.MaxSlippage > 1.0 {
		return fmt.Errorf("invalid TWAP max slippage: %f", config.TWAP.MaxSlippage)
	}

	return nil
}

// Helper functions for environment variable parsing
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if str := os.Getenv(key); str != "" {
		if value, err := strconv.Atoi(str); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if str := os.Getenv(key); str != "" {
		if value, err := strconv.ParseFloat(str, 64); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if str := os.Getenv(key); str != "" {
		if value, err := strconv.ParseBool(str); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if str := os.Getenv(key); str != "" {
		if value, err := time.ParseDuration(str); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if str := os.Getenv(key); str != "" {
		// Simple comma-separated parsing
		// In production, consider using JSON or more sophisticated parsing
		var result []string
		for _, item := range strings.Split(str, ",") {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

// GetEnvironment returns the current environment
func GetEnvironment() string {
	return getEnvString("ENVIRONMENT", EnvDevelopment)
}

// IsProduction returns true if running in production
func IsProduction() bool {
	return GetEnvironment() == EnvProduction
}

// IsDevelopment returns true if running in development
func IsDevelopment() bool {
	return GetEnvironment() == EnvDevelopment
}

// LogConfiguration logs the non-sensitive configuration (for debugging)
func LogConfiguration(config *Config, logger *zap.Logger) {
	logger.Info("Configuration loaded",
		zap.String("environment", GetEnvironment()),
		zap.Int("server_port", config.Server.Port),
		zap.String("server_host", config.Server.Host),
		zap.Int64("ethereum_chain_id", config.Ethereum.ChainID),
		zap.String("cosmos_chain_id", config.Cosmos.ChainID),
		zap.String("ethereum_contract", config.Ethereum.ContractAddr),
		zap.String("cosmos_contract", config.Cosmos.ContractAddr),
		zap.Int("twap_max_intervals", config.TWAP.MaxIntervals),
		zap.Float64("twap_max_slippage", config.TWAP.MaxSlippage),
		zap.Bool("enable_metrics", config.Monitoring.EnableMetrics),
		zap.String("log_level", config.Logging.Level),
		zap.Bool("enable_https", config.Server.EnableHTTPS),
	)
}