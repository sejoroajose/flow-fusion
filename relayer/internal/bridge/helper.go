package bridge

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"math/big"
	"go.uber.org/zap"

	"golang.org/x/crypto/sha3"

	"flow-fusion/relayer/internal/oneinch"
)

// generateSecret generates a random secret for atomic swaps
func (s *FusionBridgeService) generateSecret() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return "0x" + hex.EncodeToString(bytes), nil
}

// hashSecret hashes a secret using Keccak256
func (s *FusionBridgeService) hashSecret(secret string) (string, error) {
	// Remove 0x prefix if present
	if len(secret) >= 2 && secret[:2] == "0x" {
		secret = secret[2:]
	}
	
	// Decode hex string to bytes
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret hex: %w", err)
	}
	
	// Hash using Keccak256
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(secretBytes)
	hash := hasher.Sum(nil)
	
	return "0x" + hex.EncodeToString(hash), nil
}

// generateOrderHash generates an order hash based on the quote and parameters
func (s *FusionBridgeService) generateOrderHash(quote *oneinch.FusionQuoteResponse, secretHashes []string, userAddress string) (string, error) {
	// This is a simplified implementation for the HTTP API version
	// In a production system, this would follow 1inch's order structure specification
	
	// Create a hash based on quote data, secret hashes, and user address
	hasher := sha3.NewLegacyKeccak256()
	
	// Write quote ID if available
	if quote.QuoteId != "" {
		hasher.Write([]byte(quote.QuoteId))
	}
	
	// Write user address
	hasher.Write([]byte(userAddress))
	
	// Write secret hashes
	for _, secretHash := range secretHashes {
		if len(secretHash) >= 2 && secretHash[:2] == "0x" {
			secretHash = secretHash[2:]
		}
		if secretBytes, err := hex.DecodeString(secretHash); err == nil {
			hasher.Write(secretBytes)
		}
	}
	
	// Write timestamp for uniqueness
	timestamp := time.Now().UnixNano()
	hasher.Write([]byte(fmt.Sprintf("%d", timestamp)))
	
	// Write amounts
	hasher.Write([]byte(quote.SrcTokenAmount))
	hasher.Write([]byte(quote.DstTokenAmount))
	
	hash := hasher.Sum(nil)
	return "0x" + hex.EncodeToString(hash), nil
}

// validateOrderHash validates that an order hash matches the expected format
func (s *FusionBridgeService) validateOrderHash(orderHash string) bool {
	// Check if it's a valid hex string with 0x prefix
	if len(orderHash) != 66 { // 0x + 64 hex characters
		return false
	}
	
	if orderHash[:2] != "0x" {
		return false
	}
	
	// Check if the rest are valid hex characters
	_, err := hex.DecodeString(orderHash[2:])
	return err == nil
}

// generateSecretPair generates a secret and its hash pair
func (s *FusionBridgeService) generateSecretPair() (secret, hash string, err error) {
	secret, err = s.generateSecret()
	if err != nil {
		return "", "", err
	}
	
	hash, err = s.hashSecret(secret)
	if err != nil {
		return "", "", err
	}
	
	return secret, hash, nil
}

// validateSecret validates that a secret matches its expected hash
func (s *FusionBridgeService) validateSecret(secret, expectedHash string) (bool, error) {
	calculatedHash, err := s.hashSecret(secret)
	if err != nil {
		return false, err
	}
	
	return calculatedHash == expectedHash, nil
}

// extractChainInfo extracts chain information from chain names
func (s *FusionBridgeService) extractChainInfo(chainName string) (chainID int, isSupported bool) {
	chainMap := map[string]int{
		"ethereum": 1,
		"polygon":  137,
		"arbitrum": 42161,
		"base":     8453,
		"optimism": 10,
		"avalanche": 43114,
		"bsc":      56,
		"cosmos":   1, // Custom chain ID for Cosmos
	}
	
	chainID, isSupported = chainMap[chainName]
	return chainID, isSupported
}

// formatTokenAmount formats token amounts with proper decimals
func (s *FusionBridgeService) formatTokenAmount(amount, decimals string) string {
	// This is a simplified implementation
	// In production, you'd want to handle decimal conversion properly
	return amount
}

// estimateExecutionTime estimates how long a bridge operation will take
func (s *FusionBridgeService) estimateExecutionTime(sourceChain, destChain string) time.Duration {
	// Base time for same-chain operations
	baseTime := 2 * time.Minute
	
	// Add time for cross-chain operations
	if sourceChain != destChain {
		baseTime += 3 * time.Minute
	}
	
	// Add extra time for Cosmos operations
	if sourceChain == "cosmos" || destChain == "cosmos" {
		baseTime += 2 * time.Minute
	}
	
	return baseTime
}

// calculateRouteOptimization calculates the optimal route for a bridge operation
func (s *FusionBridgeService) calculateRouteOptimization(sourceChain, destChain, sourceToken, destToken string) []string {
	// Simple direct route for now
	route := []string{sourceChain}
	
	if sourceChain != destChain {
		route = append(route, destChain)
	}
	
	return route
}

// validateTokenAddress validates that a token address is in the correct format
func (s *FusionBridgeService) validateTokenAddress(address string) bool {
	// Check if it's a valid Ethereum address format
	if len(address) != 42 {
		return false
	}
	
	if address[:2] != "0x" {
		return false
	}
	
	// Check if the rest are valid hex characters
	_, err := hex.DecodeString(address[2:])
	return err == nil
}

// parseAmount parses a string amount into a big.Int
func (s *FusionBridgeService) parseAmount(amountStr string) (*big.Int, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount format: %s", amountStr)
	}
	return amount, nil
}

// formatDuration formats a duration for API responses
func (s *FusionBridgeService) formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f seconds", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f minutes", d.Minutes())
	} else {
		return fmt.Sprintf("%.1f hours", d.Hours())
	}
}

// calculatePriceImpactFromQuote calculates price impact from a 1inch quote
func (s *FusionBridgeService) calculatePriceImpactFromQuote(quote *oneinch.FusionQuoteResponse) float64 {
	// This would be implemented based on the actual quote structure
	// For now, return a placeholder value
	return 0.001 // 0.1%
}

// getOptimalExecutionMethod determines the best execution method for an order
func (s *FusionBridgeService) getOptimalExecutionMethod(sourceChain, destChain, amount string) string {
	// Simple logic - use Fusion+ for cross-chain and large amounts
	if sourceChain != destChain {
		return "fusion_plus"
	}
	
	// Parse amount and check if it's large
	if amountBig, err := s.parseAmount(amount); err == nil {
		threshold := big.NewInt(1000000) // 1M units
		if amountBig.Cmp(threshold) > 0 {
			return "fusion_plus"
		}
	}
	
	return "orderbook"
}

// sanitizeInput sanitizes user input to prevent injection attacks
func (s *FusionBridgeService) sanitizeInput(input string) string {
	// Remove any potentially dangerous characters
	// This is a basic implementation - use a proper sanitization library in production
	return input
}

// logOrderCreation logs the creation of a new order with all relevant details
func (s *FusionBridgeService) logOrderCreation(order *BridgeOrder) {
	s.logger.Info("Bridge order created via HTTP API",
		zap.String("order_id", order.ID),
		zap.String("user_address", order.UserAddress),
		zap.String("source_chain", order.SourceChain),
		zap.String("dest_chain", order.DestChain),
		zap.String("source_token", order.SourceToken),
		zap.String("dest_token", order.DestToken),
		zap.String("amount", order.Amount.String()),
		zap.Bool("twap_enabled", order.TWAPEnabled),
		zap.String("status", string(order.Status)),
		zap.Time("created_at", order.CreatedAt),
	)
}

// logOrderUpdate logs updates to an order status
func (s *FusionBridgeService) logOrderUpdate(order *BridgeOrder, oldStatus BridgeOrderStatus) {
	s.logger.Info("Bridge order status updated",
		zap.String("order_id", order.ID),
		zap.String("old_status", string(oldStatus)),
		zap.String("new_status", string(order.Status)),
		zap.String("oneinch_hash", order.OneInchOrderHash),
	)
}

// getOrderExecutionProgress calculates the execution progress of an order
func (s *FusionBridgeService) getOrderExecutionProgress(order *BridgeOrder) float64 {
	switch order.Status {
	case BridgeOrderStatusPending:
		return 0.0
	case BridgeOrderStatusProcessing:
		return 0.5
	case BridgeOrderStatusCompleted:
		return 1.0
	case BridgeOrderStatusFailed, BridgeOrderStatusCancelled:
		return 0.0
	default:
		return 0.0
	}
}
