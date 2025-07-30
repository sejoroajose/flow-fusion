package cosmos

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"cosmossdk.io/math"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultGasAdjustment = 1.5
	DefaultGasLimit      = 300000
	MaxRetries          = 3
	RetryDelay          = 2 * time.Second
	ConfirmationTimeout = 60 * time.Second
)

type Client struct {
	config         Config
	clientCtx      client.Context
	txFactory      tx.Factory
	grpcConn       *grpc.ClientConn
	authClient     authtypes.QueryClient
	bankClient     banktypes.QueryClient
	wasmClient     types.QueryClient
	txClient       txtypes.ServiceClient
	logger         *zap.Logger
	contractAddr   string
	signingKey     cryptotypes.PrivKey
	address        sdk.AccAddress
	accountNumber  uint64
	sequence       uint64
}

type Config struct {
	RPCURL        string  `json:"rpc_url"`
	GRPCAddr      string  `json:"grpc_addr"`
	ChainID       string  `json:"chain_id"`
	Mnemonic      string  `json:"mnemonic"`
	Denom         string  `json:"denom"`
	ContractAddr  string  `json:"contract_addr"`
	GasAdjustment float64 `json:"gas_adjustment"`
	GasPrices     string  `json:"gas_prices"`
	KeyName       string  `json:"key_name"`
}

type EscrowEvent struct {
	EscrowID       uint64    `json:"escrow_id"`
	EthereumTxHash string    `json:"ethereum_tx_hash"`
	Sender         string    `json:"sender"`
	Recipient      string    `json:"recipient"`
	Amount         sdk.Coin  `json:"amount"`
	SecretHash     string    `json:"secret_hash"`
	TimeStamp      time.Time `json:"timestamp"`
	BlockHeight    int64     `json:"block_height"`
	TxHash         string    `json:"tx_hash"`
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
	Amount           sdk.Coin  `json:"amount"`
	HashLock         string    `json:"hash_lock"`
	TimeLock         uint64    `json:"time_lock"`
	UnlockTime       time.Time `json:"unlock_time"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	Secret           string    `json:"secret,omitempty"`
}

type TransactionResponse struct {
	TxHash     string             `json:"txhash"`
	Height     string             `json:"height"`
	Code       uint32             `json:"code"`
	Codespace  string             `json:"codespace"`
	RawLog     string             `json:"raw_log"`
	GasWanted  string             `json:"gas_wanted"`
	GasUsed    string             `json:"gas_used"`
	Events     []abcitypes.Event  `json:"events"`
	Data       string             `json:"data"`
}

func NewClient(config Config, logger *zap.Logger) (*Client, error) {
	if config.GasAdjustment == 0 {
		config.GasAdjustment = DefaultGasAdjustment
	}
	if config.GasPrices == "" {
		config.GasPrices = fmt.Sprintf("0.025%s", config.Denom)
	}
	if config.KeyName == "" {
		config.KeyName = "relayer"
	}

	grpcConn, err := grpc.Dial(config.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC: %w", err)
	}

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	
	cdc := codec.NewProtoCodec(interfaceRegistry)

	kr, err := keyring.New("flow-fusion", keyring.BackendMemory, "", nil, cdc)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyring: %w", err)
	}

	var keyInfo *keyring.Record
	if config.Mnemonic != "" {
		keyInfo, err = kr.NewAccount(config.KeyName, config.Mnemonic, "", "", hd.Secp256k1)
		if err != nil {
			return nil, fmt.Errorf("failed to import mnemonic: %w", err)
		}
	} else {
		keyInfo, err = kr.Key(config.KeyName)
		if err != nil {
			return nil, fmt.Errorf("failed to get key from keyring: %w", err)
		}
	}

	addr, err := keyInfo.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	clientCtx := client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(interfaceRegistry).
		WithTxConfig(authtx.NewTxConfig(cdc, authtx.DefaultSignModes)).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(kr).
		WithChainID(config.ChainID).
		WithGRPCClient(grpcConn)

	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	wasmClient := types.NewQueryClient(grpcConn)
	txClient := txtypes.NewServiceClient(grpcConn)

	txFactory := tx.Factory{}.
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithChainID(config.ChainID).
		WithTxConfig(clientCtx.TxConfig).
		WithGasAdjustment(config.GasAdjustment).
		WithGasPrices(config.GasPrices).
		WithKeybase(kr).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	client := &Client{
		config:       config,
		clientCtx:    clientCtx,
		txFactory:    txFactory,
		grpcConn:     grpcConn,
		authClient:   authClient,
		bankClient:   bankClient,
		wasmClient:   wasmClient,
		txClient:     txClient,
		logger:       logger,
		contractAddr: config.ContractAddr,
		address:      addr,
	}

	if err := client.updateAccountInfo(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to get initial account info: %w", err)
	}

	logger.Info("Cosmos client initialized",
		zap.String("chain_id", config.ChainID),
		zap.String("address", addr.String()),
		zap.String("contract_addr", config.ContractAddr),
	)

	return client, nil
}

func (c *Client) updateAccountInfo(ctx context.Context) error {
	req := &authtypes.QueryAccountRequest{Address: c.address.String()}
	resp, err := c.authClient.Account(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	var account authtypes.AccountI
	if err := c.clientCtx.InterfaceRegistry.UnpackAny(resp.Account, &account); err != nil {
		return fmt.Errorf("failed to unpack account: %w", err)
	}

	c.accountNumber = account.GetAccountNumber()
	c.sequence = account.GetSequence()

	return nil
}

func (c *Client) CreateEscrow(ctx context.Context, params CreateEscrowParams) (string, error) {
	msg := &ExecuteMsg{
		CreateEscrow: &CreateEscrowMsg{
			EthereumTxHash: params.EthereumTxHash,
			HashLock:       params.SecretHash,
			TimeLock:       params.TimeLock,
			Recipient:      params.Recipient,
			EthereumSender: params.EthereumSender,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal execute message: %w", err)
	}

	funds := sdk.NewCoins(sdk.NewCoin(c.config.Denom, math.NewIntFromBigInt(params.Amount)))

	executeMsg := &types.MsgExecuteContract{
		Sender:   c.address.String(),
		Contract: c.contractAddr,
		Msg:      msgBytes,
		Funds:    funds,
	}

	txHash, err := c.broadcastTx(ctx, executeMsg)
	if err != nil {
		return "", fmt.Errorf("failed to create escrow: %w", err)
	}

	c.logger.Info("Created Cosmos escrow",
		zap.String("tx_hash", txHash),
		zap.String("ethereum_tx_hash", params.EthereumTxHash),
		zap.String("recipient", params.Recipient),
		zap.String("amount", params.Amount.String()),
	)

	return txHash, nil
}

func (c *Client) WithdrawEscrow(ctx context.Context, escrowID uint64, secret string) (string, error) {
	msg := &ExecuteMsg{
		Withdraw: &WithdrawMsg{
			EscrowID: escrowID,
			Secret:   secret,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal withdraw message: %w", err)
	}

	executeMsg := &types.MsgExecuteContract{
		Sender:   c.address.String(),
		Contract: c.contractAddr,
		Msg:      msgBytes,
		Funds:    nil,
	}

	txHash, err := c.broadcastTx(ctx, executeMsg)
	if err != nil {
		return "", fmt.Errorf("failed to withdraw escrow: %w", err)
	}

	c.logger.Info("Withdrew from Cosmos escrow",
		zap.String("tx_hash", txHash),
		zap.Uint64("escrow_id", escrowID),
	)

	return txHash, nil
}

func (c *Client) RefundEscrow(ctx context.Context, escrowID uint64) (string, error) {
	msg := &ExecuteMsg{
		Refund: &RefundMsg{
			EscrowID: escrowID,
		},
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal refund message: %w", err)
	}

	executeMsg := &types.MsgExecuteContract{
		Sender:   c.address.String(),
		Contract: c.contractAddr,
		Msg:      msgBytes,
		Funds:    nil,
	}

	txHash, err := c.broadcastTx(ctx, executeMsg)
	if err != nil {
		return "", fmt.Errorf("failed to refund escrow: %w", err)
	}

	c.logger.Info("Refunded Cosmos escrow",
		zap.String("tx_hash", txHash),
		zap.Uint64("escrow_id", escrowID),
	)

	return txHash, nil
}

func (c *Client) GetEscrow(ctx context.Context, escrowID uint64) (*Escrow, error) {
	query := &QueryMsg{
		GetEscrow: &GetEscrowQuery{
			EscrowID: escrowID,
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	req := &types.QuerySmartContractStateRequest{
		Address:   c.contractAddr,
		QueryData: queryBytes,
	}

	resp, err := c.wasmClient.SmartContractState(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query escrow: %w", err)
	}

	var escrowResp EscrowResponse
	if err := json.Unmarshal(resp.Data, &escrowResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal escrow response: %w", err)
	}

	return &escrowResp.Escrow, nil
}

func (c *Client) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	if address == "" {
		address = c.address.String()
	}

	req := &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   c.config.Denom,
	}

	resp, err := c.bankClient.Balance(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	if resp.Balance == nil {
		return big.NewInt(0), nil
	}

	balance, ok := new(big.Int).SetString(resp.Balance.Amount.String(), 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse balance amount: %s", resp.Balance.Amount.String())
	}

	return balance, nil
}

func (c *Client) GetTransactionStatus(ctx context.Context, txHash string) (*TransactionResponse, error) {
	req := &txtypes.GetTxRequest{Hash: txHash}
	resp, err := c.txClient.GetTx(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &TransactionResponse{
		TxHash:    resp.TxResponse.TxHash,
		Height:    strconv.FormatInt(resp.TxResponse.Height, 10),
		Code:      resp.TxResponse.Code,
		Codespace: resp.TxResponse.Codespace,
		RawLog:    resp.TxResponse.RawLog,
		GasWanted: strconv.FormatInt(resp.TxResponse.GasWanted, 10),
		GasUsed:   strconv.FormatInt(resp.TxResponse.GasUsed, 10),
		Events:    resp.TxResponse.Events,
		Data:      resp.TxResponse.Data,
	}, nil
}

func (c *Client) WaitForConfirmation(ctx context.Context, txHash string) (bool, error) {
	timeout := time.After(ConfirmationTimeout)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for confirmation")
		case <-ticker.C:
			txResp, err := c.GetTransactionStatus(ctx, txHash)
			if err != nil {
				continue // Transaction not found yet
			}
			
			if txResp.Code == 0 {
				return true, nil
			} else {
				return false, fmt.Errorf("transaction failed with code %d: %s", txResp.Code, txResp.RawLog)
			}
		}
	}
}

func (c *Client) broadcastTx(ctx context.Context, msgs ...sdk.Msg) (string, error) {
	if err := c.updateAccountInfo(ctx); err != nil {
		return "", fmt.Errorf("failed to update account info: %w", err)
	}

	txFactory := c.txFactory.
		WithAccountNumber(c.accountNumber).
		WithSequence(c.sequence)

	txBuilder, err := txFactory.BuildUnsignedTx(msgs...)
	if err != nil {
		return "", fmt.Errorf("failed to build unsigned tx: %w", err)
	}

	if txFactory.Gas() == 0 {
		_, adjusted, err := tx.CalculateGas(c.clientCtx, txFactory, msgs...)
		if err != nil {
			adjusted = DefaultGasLimit
		}
		txFactory = txFactory.WithGas(adjusted)
		txBuilder.SetGasLimit(adjusted)
	}

	err = tx.Sign(ctx, txFactory, c.config.KeyName, txBuilder, true)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	txBytes, err := c.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", fmt.Errorf("failed to encode transaction: %w", err)
	}

	var txResp *sdk.TxResponse
	for i := 0; i < MaxRetries; i++ {
		txResp, err = c.clientCtx.BroadcastTx(txBytes)
		if err == nil && txResp.Code == 0 {
			break
		}
		
		if i < MaxRetries-1 {
			var code uint32 = 999
			if txResp != nil {
				code = txResp.Code
			}
			c.logger.Warn("Transaction broadcast failed, retrying",
				zap.Error(err),
				zap.Int("attempt", i+1),
				zap.Uint32("code", code),
			)
			time.Sleep(RetryDelay)
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction after %d retries: %w", MaxRetries, err)
	}

	if txResp.Code != 0 {
		return "", fmt.Errorf("transaction failed with code %d: %s", txResp.Code, txResp.RawLog)
	}

	c.sequence++

	return txResp.TxHash, nil
}

func (c *Client) GetChainID() string {
	return c.config.ChainID
}

func (c *Client) GetContractAddress() string {
	return c.contractAddr
}

func (c *Client) GetAddress() sdk.AccAddress {
	return c.address
}

func (c *Client) Close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

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