package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	DefaultGasLimit     = 500000
	ConfirmationBlocks  = 3
	MaxRetries         = 3
	RetryDelay         = 3 * time.Second
)

type Client struct {
	config           Config
	ethClient        *ethclient.Client
	privateKey       *ecdsa.PrivateKey
	publicKey        common.Address
	chainID          *big.Int
	logger           *zap.Logger
	contractAddress  common.Address
	contract         *FlowFusionEscrowFactory
	auth             *bind.TransactOpts
}

type Config struct {
	RPCURL          string `json:"rpc_url"`
	PrivateKey      string `json:"private_key"`
	ChainID         int64  `json:"chain_id"`
	ContractAddress string `json:"contract_address"`
	GasLimit        uint64 `json:"gas_limit"`
	MaxGasPrice     string `json:"max_gas_price"`
}

type TWAPOrderParams struct {
	OrderID         [32]byte
	Token           common.Address
	TotalAmount     *big.Int
	TimeWindow      *big.Int
	IntervalCount   *big.Int
	MaxSlippage     *big.Int
	StartTime       *big.Int
	CosmosRecipient string
}

type ExecuteIntervalParams struct {
	OrderID       [32]byte
	IntervalIndex *big.Int
	SecretHash    [32]byte
	Immutables    IBaseEscrowImmutables
}

type TWAPOrderInfo struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}

func NewClient(config Config, logger *zap.Logger) (*Client, error) {
	// Connect to Ethereum client
	ethClient, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	// Parse private key
	privateKeyHex := config.PrivateKey
	if privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)
	chainID := big.NewInt(config.ChainID)

	// Parse contract address
	contractAddress := common.HexToAddress(config.ContractAddress)
	if contractAddress == (common.Address{}) {
		return nil, fmt.Errorf("invalid contract address: %s", config.ContractAddress)
	}

	// Initialize contract binding
	contract, err := NewFlowFusionEscrowFactory(contractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to bind contract: %w", err)
	}

	// Create transact opts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Set gas limit
	if config.GasLimit > 0 {
		auth.GasLimit = config.GasLimit
	} else {
		auth.GasLimit = DefaultGasLimit
	}

	// Set max gas price if specified
	if config.MaxGasPrice != "" {
		maxGasPrice, ok := new(big.Int).SetString(config.MaxGasPrice, 10)
		if ok {
			auth.GasPrice = maxGasPrice
		}
	}

	client := &Client{
		config:          config,
		ethClient:       ethClient,
		privateKey:      privateKey,
		publicKey:       publicKey,
		chainID:         chainID,
		logger:          logger,
		contractAddress: contractAddress,
		contract:        contract,
		auth:            auth,
	}

	// Verify contract deployment
	code, err := ethClient.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract code: %w", err)
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract deployed at address: %s", config.ContractAddress)
	}

	logger.Info("Ethereum client initialized",
		zap.String("rpc_url", config.RPCURL),
		zap.String("address", publicKey.Hex()),
		zap.String("contract", contractAddress.Hex()),
		zap.Int64("chain_id", config.ChainID),
	)

	return client, nil
}

func (c *Client) GetAddress() common.Address {
	return c.publicKey
}

func (c *Client) GetBalance(ctx context.Context) (*big.Int, error) {
	return c.ethClient.BalanceAt(ctx, c.publicKey, nil)
}

func (c *Client) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

func (c *Client) CreateTWAPOrder(ctx context.Context, params TWAPOrderParams) (*types.Transaction, error) {
	c.logger.Info("Creating TWAP order",
		zap.String("order_id", fmt.Sprintf("0x%x", params.OrderID)),
		zap.String("token", params.Token.Hex()),
		zap.String("amount", params.TotalAmount.String()),
		zap.String("cosmos_recipient", params.CosmosRecipient),
	)

	// Update nonce and gas price
	if err := c.updateTransactOpts(ctx); err != nil {
		return nil, fmt.Errorf("failed to update transaction options: %w", err)
	}

	config := FlowFusionEscrowFactoryTWAPConfig{
		TotalAmount:   params.TotalAmount,
		TimeWindow:    params.TimeWindow,
		IntervalCount: params.IntervalCount,
		MaxSlippage:   params.MaxSlippage,
		StartTime:     params.StartTime,
	}

	var tx *types.Transaction
	var err error

	for i := 0; i < MaxRetries; i++ {
		tx, err = c.contract.CreateTWAPOrder(c.auth, params.OrderID, params.Token, config, params.CosmosRecipient)
		if err == nil {
			break
		}

		c.logger.Warn("Failed to create TWAP order, retrying",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)

		if i < MaxRetries-1 {
			time.Sleep(RetryDelay)
			// Update nonce for retry
			if err := c.updateTransactOpts(ctx); err != nil {
				return nil, fmt.Errorf("failed to update nonce for retry: %w", err)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create TWAP order after %d retries: %w", MaxRetries, err)
	}

	c.logger.Info("TWAP order transaction sent",
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.String("order_id", fmt.Sprintf("0x%x", params.OrderID)),
	)

	return tx, nil
}

func (c *Client) ExecuteTWAPInterval(ctx context.Context, params ExecuteIntervalParams) (*types.Transaction, error) {
	c.logger.Info("Executing TWAP interval",
		zap.String("order_id", fmt.Sprintf("0x%x", params.OrderID)),
		zap.String("interval_index", params.IntervalIndex.String()),
		zap.String("secret_hash", fmt.Sprintf("0x%x", params.SecretHash)),
	)

	// Update nonce and gas price
	if err := c.updateTransactOpts(ctx); err != nil {
		return nil, fmt.Errorf("failed to update transaction options: %w", err)
	}

	var tx *types.Transaction
	var err error

	for i := 0; i < MaxRetries; i++ {
		tx, err = c.contract.ExecuteTWAPInterval(c.auth, params.OrderID, params.IntervalIndex, params.SecretHash, params.Immutables)
		if err == nil {
			break
		}

		c.logger.Warn("Failed to execute TWAP interval, retrying",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)

		if i < MaxRetries-1 {
			time.Sleep(RetryDelay)
			if err := c.updateTransactOpts(ctx); err != nil {
				return nil, fmt.Errorf("failed to update nonce for retry: %w", err)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute TWAP interval after %d retries: %w", MaxRetries, err)
	}

	c.logger.Info("TWAP interval execution transaction sent",
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.String("order_id", fmt.Sprintf("0x%x", params.OrderID)),
		zap.String("interval_index", params.IntervalIndex.String()),
	)

	return tx, nil
}

func (c *Client) CancelTWAPOrder(ctx context.Context, orderID [32]byte) (*types.Transaction, error) {
	c.logger.Info("Cancelling TWAP order",
		zap.String("order_id", fmt.Sprintf("0x%x", orderID)),
	)

	if err := c.updateTransactOpts(ctx); err != nil {
		return nil, fmt.Errorf("failed to update transaction options: %w", err)
	}

	tx, err := c.contract.CancelTWAPOrder(c.auth, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel TWAP order: %w", err)
	}

	c.logger.Info("TWAP order cancellation transaction sent",
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.String("order_id", fmt.Sprintf("0x%x", orderID)),
	)

	return tx, nil
}

func (c *Client) GetTWAPOrder(ctx context.Context, orderID [32]byte) (*TWAPOrderInfo, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	
	result, err := c.contract.GetTWAPOrder(callOpts, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get TWAP order: %w", err)
	}

	return &TWAPOrderInfo{
		Maker:             result.Maker,
		Token:             result.Token,
		Config:            result.Config,
		CosmosRecipient:   result.CosmosRecipient,
		ExecutedIntervals: result.ExecutedIntervals,
		TotalExecuted:     result.TotalExecuted,
		Cancelled:         result.Cancelled,
	}, nil
}

func (c *Client) IsIntervalExecuted(ctx context.Context, orderID [32]byte, intervalIndex *big.Int) (bool, error) {
	callOpts := &bind.CallOpts{Context: ctx}
	
	executed, err := c.contract.IsIntervalExecuted(callOpts, orderID, intervalIndex)
	if err != nil {
		return false, fmt.Errorf("failed to check interval execution: %w", err)
	}

	return executed, nil
}

func (c *Client) AddAuthorizedResolver(ctx context.Context, resolver common.Address) (*types.Transaction, error) {
	c.logger.Info("Adding authorized resolver",
		zap.String("resolver", resolver.Hex()),
	)

	if err := c.updateTransactOpts(ctx); err != nil {
		return nil, fmt.Errorf("failed to update transaction options: %w", err)
	}

	tx, err := c.contract.AddResolver(c.auth, resolver)
	if err != nil {
		return nil, fmt.Errorf("failed to add resolver: %w", err)
	}

	c.logger.Info("Add resolver transaction sent",
		zap.String("tx_hash", tx.Hash().Hex()),
		zap.String("resolver", resolver.Hex()),
	)

	return tx, nil
}

func (c *Client) WaitForConfirmation(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	c.logger.Info("Waiting for transaction confirmation",
		zap.String("tx_hash", txHash.Hex()),
		zap.Uint64("required_confirmations", ConfirmationBlocks),
	)

	// First wait for the transaction to be mined
	receipt, err := c.waitForReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	// Wait for confirmations
	for {
		select {
		case <-ctx.Done():
			return receipt, ctx.Err()
		case <-time.After(15 * time.Second):
			currentBlock, err := c.ethClient.BlockNumber(ctx)
			if err != nil {
				c.logger.Warn("Failed to get current block number", zap.Error(err))
				continue
			}

			confirmations := currentBlock - receipt.BlockNumber.Uint64()
			if confirmations >= ConfirmationBlocks {
				c.logger.Info("Transaction confirmed",
					zap.String("tx_hash", txHash.Hex()),
					zap.Uint64("confirmations", confirmations),
					zap.Uint64("block_number", receipt.BlockNumber.Uint64()),
				)
				return receipt, nil
			}

			c.logger.Debug("Waiting for more confirmations",
				zap.String("tx_hash", txHash.Hex()),
				zap.Uint64("current_confirmations", confirmations),
				zap.Uint64("required_confirmations", ConfirmationBlocks),
			)
		}
	}
}

func (c *Client) waitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
			
		}
	}
}

func (c *Client) WatchTWAPOrderEvents(ctx context.Context, eventChan chan<- FlowFusionEscrowFactoryTWAPOrderCreated) error {
	c.logger.Info("Starting to watch TWAP order events")

	// Set up event filtering for TWAP order creation
	filterOpts := &bind.FilterOpts{
		Start:   nil, // Start from latest block
		Context: ctx,
	}

	iter, err := c.contract.FilterTWAPOrderCreated(filterOpts, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create event filter: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		if iter.Error() != nil {
			c.logger.Error("Event iterator error", zap.Error(iter.Error()))
			continue
		}

		event := *iter.Event
		c.logger.Info("TWAP order created event",
			zap.String("order_id", fmt.Sprintf("0x%x", event.OrderId)),
			zap.String("maker", event.Maker.Hex()),
			zap.String("token", event.Token.Hex()),
			zap.String("total_amount", event.TotalAmount.String()),
		)

		select {
		case eventChan <- event:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return iter.Error()
}

func (c *Client) WatchIntervalExecutionEvents(ctx context.Context, eventChan chan<- FlowFusionEscrowFactoryTWAPIntervalExecuted) error {
	c.logger.Info("Starting to watch TWAP interval execution events")

	filterOpts := &bind.FilterOpts{
		Start:   nil,
		Context: ctx,
	}

	iter, err := c.contract.FilterTWAPIntervalExecuted(filterOpts, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create interval execution filter: %w", err)
	}
	defer iter.Close()

	for iter.Next() {
		if iter.Error() != nil {
			c.logger.Error("Interval execution event iterator error", zap.Error(iter.Error()))
			continue
		}

		event := *iter.Event
		c.logger.Info("TWAP interval executed event",
			zap.String("order_id", fmt.Sprintf("0x%x", event.OrderId)),
			zap.String("interval_index", event.IntervalIndex.String()),
			zap.String("amount", event.Amount.String()),
			zap.String("secret_hash", fmt.Sprintf("0x%x", event.SecretHash)),
		)

		select {
		case eventChan <- event:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return iter.Error()
}

func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.ethClient.EstimateGas(ctx, msg)
}

func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.ethClient.SuggestGasPrice(ctx)
}

func (c *Client) updateTransactOpts(ctx context.Context) error {
	// Update nonce
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.publicKey)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	c.auth.Nonce = big.NewInt(int64(nonce))

	// Update gas price if not set
	if c.auth.GasPrice == nil {
		gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to suggest gas price: %w", err)
		}
		c.auth.GasPrice = gasPrice
	}

	c.auth.Context = ctx
	return nil
}

// Helper function to convert string to bytes32
func StringToBytes32(s string) [32]byte {
	var bytes32 [32]byte
	copy(bytes32[:], []byte(s))
	return bytes32
}

// Helper function to convert hex string to bytes32
func HexToBytes32(hex string) ([32]byte, error) {
	var bytes32 [32]byte
	if hex[:2] == "0x" {
		hex = hex[2:]
	}
	bytes, err := common.FromHex("0x" + hex)
	if err != nil {
		return bytes32, err
	}
	copy(bytes32[:], bytes)
	return bytes32, nil
}