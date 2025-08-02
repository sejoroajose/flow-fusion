package oneinch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// HTTPClient implements 1inch API calls using direct HTTP requests
type HTTPClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
	chainID    uint64
}

// Config for the HTTP client
type Config struct {
	BaseURL string
	APIKey  string
	ChainID uint64
	Timeout time.Duration
}

// Response structures for 1inch API
type QuoteResponse struct {
	DstAmount          string `json:"dstAmount"`
	SrcAmount          string `json:"srcAmount"`
	EstimatedGas       string `json:"estimatedGas"`
	EstimatedPriceImpact string `json:"estimatedPriceImpact,omitempty"`
}

type SwapResponse struct {
	DstAmount    string      `json:"dstAmount"`
	SrcAmount    string      `json:"srcAmount"`
	Tx           Transaction `json:"tx"`
	EstimatedGas string      `json:"estimatedGas"`
}

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
	Gas      string `json:"gas"`
}

// Fusion+ specific structures
type FusionQuoteResponse struct {
	DstTokenAmount     string                 `json:"dstTokenAmount"`
	SrcTokenAmount     string                 `json:"srcTokenAmount"`
	Presets            []FusionPreset         `json:"presets"`
	RecommendedPreset  string                 `json:"recommendedPreset"`
	EstimatedGas       string                 `json:"estimatedGas"`
	Prices             []string               `json:"prices"`
	Volume             []string               `json:"volume"`
	SettlementAddress  string                 `json:"settlementAddress"`
	QuoteId            string                 `json:"quoteId"`
}

type FusionPreset struct {
	Name         string `json:"name"`
	SecretsCount int    `json:"secretsCount"`
	Timings      struct {
		AuctionStartDelay    int `json:"auctionStartDelay"`
		AuctionDuration      int `json:"auctionDuration"`
		InitialRateBump      int `json:"initialRateBump"`
		BankrollAddress      string `json:"bankrollAddress"`
	} `json:"timings"`
}

type FusionOrderRequest struct {
	WalletAddress    string            `json:"walletAddress"`
	OrderHash        string            `json:"orderHash"`
	SecretHashes     []string          `json:"secretHashes"`
	Receiver         string            `json:"receiver"`
	Preset           string            `json:"preset"`
	Extension        string            `json:"extension,omitempty"`
	Source           string            `json:"source,omitempty"`
}

type FusionOrderResponse struct {
	OrderHash string `json:"orderHash"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}

type FusionOrderStatus struct {
	OrderHash string `json:"orderHash"`
	Status    string `json:"status"` // "pending", "executed", "cancelled", "expired"
	Fills     []Fill `json:"fills"`
}

type Fill struct {
	FillHash   string `json:"fillHash"`
	TxHash     string `json:"txHash"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

type ReadyFillsResponse struct {
	Fills []ReadyFill `json:"fills"`
}

type ReadyFill struct {
	FillHash   string `json:"fillHash"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

type SecretSubmissionRequest struct {
	OrderHash string `json:"orderHash"`
	Secret    string `json:"secret"`
}

type SecretSubmissionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Orderbook structures
type OrderbookOrderRequest struct {
	Maker              string `json:"maker"`
	MakerAsset         string `json:"makerAsset"`
	TakerAsset         string `json:"takerAsset"`
	MakingAmount       string `json:"makingAmount"`
	TakingAmount       string `json:"takingAmount"`
	Salt               string `json:"salt"`
	MakerTraits        string `json:"makerTraits"`
	Signature          string `json:"signature"`
}

type OrderbookOrderResponse struct {
	Success   bool   `json:"success"`
	OrderHash string `json:"orderHash"`
	Message   string `json:"message,omitempty"`
}

type OrderbookOrder struct {
	OrderHash        string `json:"orderHash"`
	Maker            string `json:"maker"`
	MakerAsset       string `json:"makerAsset"`
	TakerAsset       string `json:"takerAsset"`
	MakingAmount     string `json:"makingAmount"`
	TakingAmount     string `json:"takingAmount"`
	RemainingAmount  string `json:"remainingMakerAmount"`
	Status           string `json:"status"`
	CreateDateTime   string `json:"createDateTime"`
}

type TokenInfo struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
	LogoURI  string `json:"logoURI"`
}

func NewHTTPClient(config Config, logger *zap.Logger) *HTTPClient {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.1inch.dev"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &HTTPClient{
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger:  logger,
		chainID: config.ChainID,
	}
}

func (c *HTTPClient) doRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	c.logger.Debug("Making 1inch API request",
		zap.String("method", method),
		zap.String("url", url),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.Error("1inch API request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(respBody)),
		)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Aggregation API methods
func (c *HTTPClient) GetQuote(ctx context.Context, src, dst, amount string, params map[string]string) (*QuoteResponse, error) {
	endpoint := fmt.Sprintf("/swap/v6.0/%d/quote", c.chainID)
	
	// Build query parameters
	values := url.Values{}
	values.Set("src", src)
	values.Set("dst", dst)
	values.Set("amount", amount)
	
	for key, value := range params {
		values.Set(key, value)
	}
	
	endpoint += "?" + values.Encode()

	var response QuoteResponse
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

func (c *HTTPClient) GetSwap(ctx context.Context, src, dst, amount, from string, slippage float64, params map[string]string) (*SwapResponse, error) {
	endpoint := fmt.Sprintf("/swap/v6.0/%d/swap", c.chainID)
	
	// Build query parameters
	values := url.Values{}
	values.Set("src", src)
	values.Set("dst", dst)
	values.Set("amount", amount)
	values.Set("from", from)
	values.Set("slippage", strconv.FormatFloat(slippage, 'f', 2, 64))
	
	for key, value := range params {
		values.Set(key, value)
	}
	
	endpoint += "?" + values.Encode()

	var response SwapResponse
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

// Fusion+ API methods
func (c *HTTPClient) GetFusionQuote(ctx context.Context, srcChain, dstChain uint64, srcToken, dstToken, amount, walletAddress string) (*FusionQuoteResponse, error) {
	endpoint := "/fusion-plus/quoter/v1.0/quote"
	
	values := url.Values{}
	values.Set("srcChain", strconv.FormatUint(srcChain, 10))
	values.Set("dstChain", strconv.FormatUint(dstChain, 10))
	values.Set("srcTokenAddress", srcToken)
	values.Set("dstTokenAddress", dstToken)
	values.Set("amount", amount)
	values.Set("walletAddress", walletAddress)
	values.Set("enableEstimate", "true")
	
	endpoint += "?" + values.Encode()

	var response FusionQuoteResponse
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

func (c *HTTPClient) PlaceFusionOrder(ctx context.Context, orderReq FusionOrderRequest) (*FusionOrderResponse, error) {
	endpoint := "/fusion-plus/relayer/v1.0/order"

	var response FusionOrderResponse
	err := c.doRequest(ctx, "POST", endpoint, orderReq, &response)
	return &response, err
}

func (c *HTTPClient) GetFusionOrderStatus(ctx context.Context, orderHash string) (*FusionOrderStatus, error) {
	endpoint := fmt.Sprintf("/fusion-plus/relayer/v1.0/order/%s", orderHash)

	var response FusionOrderStatus
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

func (c *HTTPClient) GetReadyToAcceptFills(ctx context.Context, orderHash string) (*ReadyFillsResponse, error) {
	endpoint := fmt.Sprintf("/fusion-plus/relayer/v1.0/order/%s/ready-to-accept-fills", orderHash)

	var response ReadyFillsResponse
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

func (c *HTTPClient) SubmitSecret(ctx context.Context, orderHash, secret string) (*SecretSubmissionResponse, error) {
	endpoint := "/fusion-plus/relayer/v1.0/secret"

	req := SecretSubmissionRequest{
		OrderHash: orderHash,
		Secret:    secret,
	}

	var response SecretSubmissionResponse
	err := c.doRequest(ctx, "POST", endpoint, req, &response)
	return &response, err
}

// Orderbook API methods
func (c *HTTPClient) CreateOrderbookOrder(ctx context.Context, orderReq OrderbookOrderRequest) (*OrderbookOrderResponse, error) {
	endpoint := fmt.Sprintf("/orderbook/v4.0/%d", c.chainID)

	var response OrderbookOrderResponse
	err := c.doRequest(ctx, "POST", endpoint, orderReq, &response)
	return &response, err
}

func (c *HTTPClient) GetOrderbookOrdersByMaker(ctx context.Context, makerAddress string, limit int) ([]OrderbookOrder, error) {
	endpoint := fmt.Sprintf("/orderbook/v4.0/%d/address/%s", c.chainID, makerAddress)
	
	if limit > 0 {
		endpoint += fmt.Sprintf("?limit=%d", limit)
	}

	var response []OrderbookOrder
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return response, err
}

func (c *HTTPClient) GetOrderbookOrder(ctx context.Context, orderHash string) (*OrderbookOrder, error) {
	endpoint := fmt.Sprintf("/orderbook/v4.0/%d/order/%s", c.chainID, orderHash)

	var response OrderbookOrder
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

// Token API methods
func (c *HTTPClient) GetTokens(ctx context.Context) (map[string]TokenInfo, error) {
	endpoint := fmt.Sprintf("/token/v1.2/%d", c.chainID)

	var response map[string]TokenInfo
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return response, err
}

func (c *HTTPClient) GetTokenInfo(ctx context.Context, tokenAddress string) (*TokenInfo, error) {
	endpoint := fmt.Sprintf("/token/v1.2/%d/%s", c.chainID, tokenAddress)

	var response TokenInfo
	err := c.doRequest(ctx, "GET", endpoint, nil, &response)
	return &response, err
}

// Utility methods
func (c *HTTPClient) SetChainID(chainID uint64) {
	c.chainID = chainID
}

func (c *HTTPClient) GetChainID() uint64 {
	return c.chainID
}