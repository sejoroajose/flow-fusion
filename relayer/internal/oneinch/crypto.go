package oneinch

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// GenerateRandomBytes32 generates 32 random bytes and returns as hex string
func (c *HTTPClient) GenerateRandomBytes32() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return "0x" + hex.EncodeToString(bytes), nil
}

// HashSecret hashes a secret using Keccak256 (Ethereum's hash function)
func (c *HTTPClient) HashSecret(secret string) (string, error) {
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

// HashSecretSHA256 hashes a secret using SHA256 (alternative)
func (c *HTTPClient) HashSecretSHA256(secret string) (string, error) {
	// Remove 0x prefix if present
	if len(secret) >= 2 && secret[:2] == "0x" {
		secret = secret[2:]
	}
	
	// Decode hex string to bytes
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret hex: %w", err)
	}
	
	// Hash using SHA256
	hash := sha256.Sum256(secretBytes)
	
	return "0x" + hex.EncodeToString(hash[:]), nil
}

// ValidateSecret validates that a secret matches its hash
func (c *HTTPClient) ValidateSecret(secret, expectedHash string) (bool, error) {
	calculatedHash, err := c.HashSecret(secret)
	if err != nil {
		return false, err
	}
	
	return calculatedHash == expectedHash, nil
}

// GenerateSecretPair generates a secret and its hash
func (c *HTTPClient) GenerateSecretPair() (secret, hash string, err error) {
	secret, err = c.GenerateRandomBytes32()
	if err != nil {
		return "", "", err
	}
	
	hash, err = c.HashSecret(secret)
	if err != nil {
		return "", "", err
	}
	
	return secret, hash, nil
}

// Standalone crypto functions that don't require the client

// GenerateRandomBytes32Standalone generates 32 random bytes and returns as hex string
func GenerateRandomBytes32Standalone() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return "0x" + hex.EncodeToString(bytes), nil
}

// HashSecretStandalone hashes a secret using Keccak256
func HashSecretStandalone(secret string) (string, error) {
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

// GenerateSecretPairStandalone generates a secret and its hash without client
func GenerateSecretPairStandalone() (secret, hash string, err error) {
	secret, err = GenerateRandomBytes32Standalone()
	if err != nil {
		return "", "", err
	}
	
	hash, err = HashSecretStandalone(secret)
	if err != nil {
		return "", "", err
	}
	
	return secret, hash, nil
}