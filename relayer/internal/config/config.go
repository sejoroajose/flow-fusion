package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig    `json:"server"`
	Ethereum  EthereumConfig  `json:"ethereum"`
	Cosmos    CosmosConfig    `json:"cosmos"`
	OneInch   OneInchConfig   `json:"oneinch"`
	TWAP      TWAPConfig      `json:"twap"`
	Redis     RedisConfig     `json:"redis"`
	Database  DatabaseConfig  `json:"database"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type EthereumConfig struct {
	RPCURL     string `json:"rpc_url"`
	PrivateKey string `json:"private_key"`
	ChainID    int64  `json:"chain_id"`
}

type CosmosConfig struct {
	RPCURL   string `json:"rpc_url"`
	GRPCAddr string `json:"grpc_addr"`
	ChainID  string `json:"chain_id"`
	Mnemonic string `json:"mnemonic"`
	Denom    string `json:"denom"`
}

type OneInchConfig struct {
	APIKey string `json:"api_key"`
}

type TWAPConfig struct {
	MaxIntervals    int           `json:"max_intervals"`
	MinIntervalTime time.Duration `json:"min_interval_time"`
	MaxSlippage     float64       `json:"max_slippage"`
}

type RedisConfig struct {
	URL string `json:"url"`
}

type DatabaseConfig struct {
	URL string `json:"url"`
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvInt("RELAYER_PORT", 8080),
		},
		Ethereum: EthereumConfig{
			RPCURL:     getEnvString("ETHEREUM_RPC", "http://localhost:8545"),
			PrivateKey: getEnvString("RELAYER_PRIVATE_KEY", ""),
			ChainID:    int64(getEnvInt("ETHEREUM_CHAIN_ID", 1)),
		},
		Cosmos: CosmosConfig{
			RPCURL:   getEnvString("COSMOS_RPC", "http://localhost:26657"),
			GRPCAddr: getEnvString("COSMOS_GRPC", "localhost:9090"),
			ChainID:  getEnvString("COSMOS_CHAIN_ID", "cosmos-1"),
			Mnemonic: getEnvString("COSMOS_MNEMONIC", ""),
			Denom:    getEnvString("COSMOS_DENOM", "uatom"),
		},
		OneInch: OneInchConfig{
			APIKey: getEnvString("ONEINCH_API_KEY", ""),
		},
		TWAP: TWAPConfig{
			MaxIntervals:    getEnvInt("TWAP_MAX_INTERVALS", 100),
			MinIntervalTime: time.Duration(getEnvInt("TWAP_MIN_INTERVAL_SECONDS", 60)) * time.Second,
			MaxSlippage:     getEnvFloat("TWAP_MAX_SLIPPAGE", 0.10), // 10%
		},
		Redis: RedisConfig{
			URL: getEnvString("REDIS_URL", "redis://localhost:6379"),
		},
		Database: DatabaseConfig{
			URL: getEnvString("DATABASE_URL", "postgresql://flow:fusion@localhost:5432/flowfusion"),
		},
	}

	// Validate required fields
	if config.Ethereum.PrivateKey == "" {
		return nil, fmt.Errorf("RELAYER_PRIVATE_KEY is required")
	}

	if config.OneInch.APIKey == "" {
		return nil, fmt.Errorf("ONEINCH_API_KEY is required")
	}

	return config, nil
}

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