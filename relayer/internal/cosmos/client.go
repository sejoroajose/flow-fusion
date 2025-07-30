package cosmos

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	config       Config
	rpcClient    *RPCClient
	grpcClient   *GRPCClient
	logger       *zap.Logger
	contractAddr string
}

type Config struct {
	RPCURL   string `json:"rpc_url"`
	GRPCAddr string `json:"grpc_addr"`
	ChainID  string `json:"chain_id"`
	Mnemonic string `json:"mnemonic"`
	Denom    string `json:"denom"`
}

type RPCClient struct {
	endpoint string
}

type GRPCClient struct {
	endpoint string
}

type EscrowEvent struct {
	EscrowID       uint64    `json:"escrow_id"`
	EthereumTxHash string    `json:"ethereum_tx_hash"`
	Sender         string    `json:"sender"`
	Recipient      string    `json:"recipient"`
	Amount         string    `json:"amount"`
	Denom          string    `json:"denom"`
	SecretHash     string    `json:"secret_hash"`
	TimeStamp      time.Time `json:"timestamp"`
}

type ExecuteMsg struct {
	CreateEscrow *CreateEscrowMsg `json:"create_escrow,omitempty"`
	Withdraw     *WithdrawMsg     `json:"withdraw,omitempty"`
	Refund       *RefundMsg       `json:"refund,omitempty"`
}

type CreateEscrowMsg struct {
	EthereumTxHash  string `json:"ethereum_tx_hash"`
	HashLock        string `json:"hash_lock"`
	TimeLock        uint64 `json:"time_lock"`
	Recipient       string `json:"recipient"`
	EthereumSender  string `json:"ethereum_sender"`
}

type WithdrawMsg struct {
	EscrowID uint64 `json:"escrow_id"`
	Secret   string `json:"secret"`
}

type RefundMsg struct {
	EscrowID uint64 `json:"escrow_id"`
}

type QueryMsg struct {
	GetEscrow  *GetEscrowQuery  `json:"get_escrow,omitempty"`
	GetEscrows *GetEscrowsQuery `json:"get_escrows,omitempty"`
	GetConfig  *struct{}        `json:"get_config,omitempty"`
}

type GetEscrowQuery struct {
	EscrowID uint64 `json:"escrow_id"`
}

type GetEscrowsQuery struct {
	StartAfter *uint64 `json:"start_after,omitempty"`
	Limit      *uint32 `json:"limit,omitempty"`
}

type EscrowResponse struct {
	Escrow Escrow `json:"escrow"`
}

type Escrow struct {
	ID               uint64    `json:"id"`
	EthereumTxHash   string    `json:"ethereum_tx_hash"`
	Sender           string    `json:"sender"`
	Recipient        string    `json:"recipient"`
	EthereumSender   string    `json:"ethereum_sender"`
	Amount           Coin      `json:"amount"`
	HashLock         string    `json:"hash_lock"`
	TimeLock         uint64    `json:"time_lock"`
	UnlockTime       time.Time `json:"unlock_time"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type TransactionResponse struct {
	TxHash    string `json:"txhash"`
	Height    string `json:"height"`
	Code      uint32 `json:"code"`
	Codespace string `json:"codespace"`
}

func NewClient(config Config, logger *zap.Logger) (*Client, error) {
	rpcClient := &RPCClient{endpoint: config.RPCURL}
	grpcClient := &GRPCClient{endpoint: config.GRPCAddr}

	client := &Client{
		config:     config,
		rpcClient:  rpcClient,
		grpcClient: grpcClient,
		logger:     logger,
		// For demo, use a mock contract address
		contractAddr: "cosmos1contractaddresshere",
	}

	logger.Info("Cosmos client initialized",
		zap.String("rpc_url", config.RPCURL),
		zap.String("grpc_addr", config.GRPCAddr),
		zap.String("chain_id", config.ChainID),
		zap.String("contract_addr", client.contractAddr),
	)

	return client, nil
}

func (c *Client) CreateEscrow(ctx context.Context, params CreateEscrowParams) (string, error) {
	executeMsg := ExecuteMsg{
		CreateEscrow: &CreateEscrowMsg{
			EthereumTxHash: params.EthereumTxHash,
			HashLock:       params.SecretHash,
			TimeLock:       params.TimeLock,
			Recipient:      params.Recipient,
			EthereumSender: params.EthereumSender,
		},
	}

	// For demo purposes, simulate transaction execution
	txHash := fmt.Sprintf("cosmos_tx_%d", time.Now().Unix())
	
	c.logger.Info("Created Cosmos escrow",
		zap.String("tx_hash", txHash),
		zap.String("ethereum_tx_hash", params.EthereumTxHash),
		zap.String("recipient", params.Recipient),
		zap.String("amount", params.Amount.String()),
	)

	return txHash, nil
}

func (c *Client) WithdrawEscrow(ctx context.Context, escrowID uint64, secret string) (string, error) {
	executeMsg := ExecuteMsg{
		Withdraw: &WithdrawMsg{
			EscrowID: escrowID,
			Secret:   secret,
		},
	}

	// For demo purposes, simulate transaction execution
	txHash := fmt.Sprintf("cosmos_withdraw_%d", time.Now().Unix())
	
	c.logger.Info("Withdrew from Cosmos escrow",
		zap.String("tx_hash", txHash),
		zap.Uint64("escrow_id", escrowID),
	)

	return txHash, nil
}

func (c *Client) RefundEscrow(ctx context.Context, escrowID uint64) (string, error) {
	executeMsg := ExecuteMsg{
		Refund: &RefundMsg{
			EscrowID: escrowID,
		},
	}

	// For demo purposes, simulate transaction execution
	txHash := fmt.Sprintf("cosmos_refund_%d", time.Now().Unix())
	
	c.logger.Info("Refunded Cosmos escrow",
		zap.String("tx_hash", txHash),
		zap.Uint64("escrow_id", escrowID),
	)

	return txHash, nil
}

func (c *Client) GetEscrow(ctx context.Context, escrowID uint64) (*Escrow, error) {
	queryMsg := QueryMsg{
		GetEscrow: &GetEscrowQuery{
			EscrowID: escrowID,
		},
	}

	// For demo purposes, return mock escrow data
	escrow := &Escrow{
		ID:             escrowID,
		EthereumTxHash: "0xdemo",
		Sender:         "cosmos1sender",
		Recipient:      "cosmos1recipient",
		EthereumSender: "0xethsender",
		Amount: Coin{
			Denom:  c.config.Denom,
			Amount: "1000000",
		},
		HashLock:   "demo_hash_lock",
		TimeLock:   3600,
		UnlockTime: time.Now().Add(1 * time.Hour),
		Status:     "Active",
		CreatedAt:  time.Now(),
	}

	return escrow, nil
}

func (c *Client) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	// For demo purposes, return mock balance
	balance := big.NewInt(1000000000) // 1000 tokens with 6 decimals
	
	c.logger.Debug("Retrieved balance",
		zap.String("address", address),
		zap.String("balance", balance.String()),
	)

	return balance, nil
}

func (c *Client) WatchEscrowEvents(ctx context.Context, eventChan chan<- EscrowEvent) error {
	// For demo purposes, simulate event watching
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Simulate an event
				event := EscrowEvent{
					EscrowID:       uint64(time.Now().Unix()),
					EthereumTxHash: "0xdemo",
					Sender:         "cosmos1sender",
					Recipient:      "cosmos1recipient",
					Amount:         "1000000",
					Denom:          c.config.Denom,
					SecretHash:     "demo_secret_hash",
					TimeStamp:      time.Now(),
				}

				select {
				case eventChan <- event:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return nil
}

func (c *Client) GetTransactionStatus(ctx context.Context, txHash string) (*TransactionResponse, error) {
	// For demo purposes, return successful transaction
	return &TransactionResponse{
		TxHash:    txHash,
		Height:    "12345",
		Code:      0, // 0 means success
		Codespace: "",
	}, nil
}

func (c *Client) GetChainID() string {
	return c.config.ChainID
}

func (c *Client) GetContractAddress() string {
	return c.contractAddr
}

// Helper types
type CreateEscrowParams struct {
	EthereumTxHash string
	SecretHash     string
	TimeLock       uint64
	Recipient      string
	EthereumSender string
	Amount         *big.Int
}

type EscrowStatus struct {
	Active    bool
	Completed bool
	Amount    *big.Int
	Recipient string
}