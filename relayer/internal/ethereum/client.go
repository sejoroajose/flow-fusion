package ethereum

import (
	"context"
	"crypto/ecdsa"
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

type Client struct {
	config      Config
	ethClient   *ethclient.Client
	privateKey  *ecdsa.PrivateKey
	publicKey   common.Address
	chainID     *big.Int
	logger      *zap.Logger
}

type Config struct {
	RPCURL     string `json:"rpc_url"`
	PrivateKey string `json:"private_key"`
	ChainID    int64  `json:"chain_id"`
}

type EscrowCreatedEvent struct {
	OrderID       common.Hash
	EthereumTxHash common.Hash
	Amount        *big.Int
	Recipient     string
	SecretHash    string
	TimeStamp     uint64
}

func NewClient(config Config, logger *zap.Logger) (*Client, error) {
	// Connect to Ethereum client
	ethClient, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, err
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey[2:]) // Remove 0x prefix
	if err != nil {
		return nil, err
	}

	// Get public key
	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)
	
	chainID := big.NewInt(config.ChainID)

	client := &Client{
		config:     config,
		ethClient:  ethClient,
		privateKey: privateKey,
		publicKey:  publicKey,
		chainID:    chainID,
		logger:     logger,
	}

	logger.Info("Ethereum client initialized",
		zap.String("rpc_url", config.RPCURL),
		zap.String("address", publicKey.Hex()),
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

func (c *Client) CreateEscrow(ctx context.Context, params CreateEscrowParams) (common.Hash, error) {
	// Get nonce
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.publicKey)
	if err != nil {
		return common.Hash{}, err
	}

	// Get gas price
	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	// For demo purposes, create a mock transaction
	// In production, this would interact with the actual FlowFusionEscrowFactory contract
	
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return common.Hash{}, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = params.Amount
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// Mock transaction hash generation for demo
	txHash := common.BytesToHash(crypto.Keccak256([]byte(params.SecretHash + params.Recipient)))
	
	c.logger.Info("Created Ethereum escrow",
		zap.String("tx_hash", txHash.Hex()),
		zap.String("recipient", params.Recipient),
		zap.String("amount", params.Amount.String()),
	)

	return txHash, nil
}

func (c *Client) WaitForConfirmation(ctx context.Context, txHash common.Hash, confirmations int) (bool, error) {
	// For demo purposes, simulate waiting for confirmations
	time.Sleep(2 * time.Second)
	
	c.logger.Info("Transaction confirmed",
		zap.String("tx_hash", txHash.Hex()),
		zap.Int("confirmations", confirmations),
	)
	
	return true, nil
}

func (c *Client) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.ethClient.TransactionReceipt(ctx, txHash)
}

func (c *Client) WatchEscrowEvents(ctx context.Context, contractAddress common.Address, eventChan chan<- EscrowCreatedEvent) error {
	// Set up event filtering
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte("EscrowCreated(bytes32,address,uint256,string,string)"))},
		},
	}

	logs := make(chan types.Log)
	sub, err := c.ethClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				c.logger.Error("Event subscription error", zap.Error(err))
				return
			case vLog := <-logs:
				// Parse event (simplified for demo)
				event := EscrowCreatedEvent{
					OrderID:       vLog.Topics[1],
					EthereumTxHash: vLog.TxHash,
					Amount:        big.NewInt(0), // Would parse from log data
					Recipient:     "",            // Would parse from log data
					SecretHash:    "",            // Would parse from log data
					TimeStamp:     vLog.BlockNumber,
				}
				
				select {
				case eventChan <- event:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Client) RevealSecret(ctx context.Context, escrowAddress common.Address, secret string) (common.Hash, error) {
	// For demo purposes, simulate secret revelation
	txHash := common.BytesToHash(crypto.Keccak256([]byte("reveal_" + secret)))
	
	c.logger.Info("Revealed secret",
		zap.String("escrow", escrowAddress.Hex()),
		zap.String("tx_hash", txHash.Hex()),
	)
	
	return txHash, nil
}

func (c *Client) GetCurrentBlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

func (c *Client) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return c.ethClient.EstimateGas(ctx, msg)
}

// Helper types
type CreateEscrowParams struct {
	Recipient  string
	Amount     *big.Int
	SecretHash string
	TimeLimit  uint64
}

type EscrowStatus struct {
	Active    bool
	Completed bool
	Amount    *big.Int
	Recipient common.Address
}