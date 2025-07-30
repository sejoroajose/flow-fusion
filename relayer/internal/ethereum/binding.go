package ethereum

import (
	"errors"
	"math/big"
	"strings"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// FlowFusionEscrowFactoryTWAPConfig is an auto generated low-level Go binding around an user-defined struct.
type FlowFusionEscrowFactoryTWAPConfig struct {
	TotalAmount   *big.Int
	TimeWindow    *big.Int
	IntervalCount *big.Int
	MaxSlippage   *big.Int
	StartTime     *big.Int
}

// IBaseEscrowImmutables is an auto generated low-level Go binding around an user-defined struct.
type IBaseEscrowImmutables struct {
	Maker     common.Address
	Taker     common.Address
	Token     common.Address
	Amount    *big.Int
	HashLock  [32]byte
	Timelocks struct {
		Prepayment *big.Int
		Maturity   *big.Int
		Expiration *big.Int
	}
}

// FlowFusionEscrowFactoryMetaData contains all meta data concerning the FlowFusionEscrowFactory contract.
var FlowFusionEscrowFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIEscrowFactory\",\"name\":\"_escrowFactory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"IntervalAlreadyExecuted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInterval\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTWAPConfig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OrderAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OrderCancelled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OrderNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedResolver\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"structIBaseEscrow.Immutables\",\"name\":\"immutables\",\"type\":\"tuple\",\"components\":[{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hashLock\",\"type\":\"bytes32\"},{\"internalType\":\"structTimelocks\",\"name\":\"timelocks\",\"type\":\"tuple\",\"components\":[{\"internalType\":\"uint256\",\"name\":\"prepayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maturity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}]}]},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cosmosRecipient\",\"type\":\"string\"}],\"name\":\"CosmosEscrowCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"intervalIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"}],\"name\":\"TWAPIntervalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refundedAmount\",\"type\":\"uint256\"}],\"name\":\"TWAPOrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeWindow\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"intervalCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"cosmosRecipient\",\"type\":\"string\"}],\"name\":\"TWAPOrderCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"addResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"authorizedResolvers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"name\":\"cancelTWAPOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"intervalCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSlippage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"}],\"internalType\":\"structFlowFusionEscrowFactory.TWAPConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"cosmosRecipient\",\"type\":\"string\"}],\"name\":\"createTWAPOrder\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emergencyWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"escrowFactory\",\"outputs\":[{\"internalType\":\"contractIEscrowFactory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"intervalIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hashLock\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"prepayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maturity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expiration\",\"type\":\"uint256\"}],\"internalType\":\"structTimelocks\",\"name\":\"timelocks\",\"type\":\"tuple\"}],\"internalType\":\"structIBaseEscrow.Immutables\",\"name\":\"immutables\",\"type\":\"tuple\"}],\"name\":\"executeTWAPInterval\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"name\":\"getTWAPOrder\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"intervalCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSlippage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"}],\"internalType\":\"structFlowFusionEscrowFactory.TWAPConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"cosmosRecipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"executedIntervals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalExecuted\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"cancelled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"intervalIndex\",\"type\":\"uint256\"}],\"name\":\"isIntervalExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_INTERVALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_SLIPPAGE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_INTERVAL_DURATION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"resolver\",\"type\":\"address\"}],\"name\":\"removeResolver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"twapOrders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"intervalCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxSlippage\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"}],\"internalType\":\"structFlowFusionEscrowFactory.TWAPConfig\",\"name\":\"config\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"cosmosRecipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"executedIntervals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalExecuted\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"cancelled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FlowFusionEscrowFactoryABI is the input ABI used to generate the binding from.
var FlowFusionEscrowFactoryABI = FlowFusionEscrowFactoryMetaData.ABI

// FlowFusionEscrowFactory is an auto generated Go binding around an Ethereum contract.
type FlowFusionEscrowFactory struct {
	FlowFusionEscrowFactoryCaller     // Read-only binding to the contract
	FlowFusionEscrowFactoryTransactor // Write-only binding to the contract
	FlowFusionEscrowFactoryFilterer   // Log filterer for contract events
}

// FlowFusionEscrowFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlowFusionEscrowFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlowFusionEscrowFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlowFusionEscrowFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlowFusionEscrowFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlowFusionEscrowFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlowFusionEscrowFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlowFusionEscrowFactorySession struct {
	Contract     *FlowFusionEscrowFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth to use throughout this session
}

// FlowFusionEscrowFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlowFusionEscrowFactoryCallerSession struct {
	Contract *FlowFusionEscrowFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// FlowFusionEscrowFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlowFusionEscrowFactoryTransactorSession struct {
	Contract     *FlowFusionEscrowFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth to use throughout this session
}

// NewFlowFusionEscrowFactory creates a new instance of FlowFusionEscrowFactory, bound to a specific deployed contract.
func NewFlowFusionEscrowFactory(address common.Address, backend bind.ContractBackend) (*FlowFusionEscrowFactory, error) {
	contract, err := bindFlowFusionEscrowFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactory{
		FlowFusionEscrowFactoryCaller:     FlowFusionEscrowFactoryCaller{contract: contract},
		FlowFusionEscrowFactoryTransactor: FlowFusionEscrowFactoryTransactor{contract: contract},
		FlowFusionEscrowFactoryFilterer:   FlowFusionEscrowFactoryFilterer{contract: contract},
	}, nil
}

// bindFlowFusionEscrowFactory binds a generic wrapper to an already deployed contract.
func bindFlowFusionEscrowFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlowFusionEscrowFactoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// TWAP Order Creation
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) CreateTWAPOrder(opts *bind.TransactOpts, orderId [32]byte, token common.Address, config FlowFusionEscrowFactoryTWAPConfig, cosmosRecipient string) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "createTWAPOrder", orderId, token, config, cosmosRecipient)
}

// TWAP Interval Execution
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) ExecuteTWAPInterval(opts *bind.TransactOpts, orderId [32]byte, intervalIndex *big.Int, secretHash [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "executeTWAPInterval", orderId, intervalIndex, secretHash, immutables)
}

// Cancel TWAP Order
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) CancelTWAPOrder(opts *bind.TransactOpts, orderId [32]byte) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "cancelTWAPOrder", orderId)
}

// Get TWAP Order (view function)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) GetTWAPOrder(opts *bind.CallOpts, orderId [32]byte) (struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "getTWAPOrder", orderId)

	outstruct := new(struct {
		Maker             common.Address
		Token             common.Address
		Config            FlowFusionEscrowFactoryTWAPConfig
		CosmosRecipient   string
		ExecutedIntervals *big.Int
		TotalExecuted     *big.Int
		Cancelled         bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Maker = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Token = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Config = *abi.ConvertType(out[2], new(FlowFusionEscrowFactoryTWAPConfig)).(*FlowFusionEscrowFactoryTWAPConfig)
	outstruct.CosmosRecipient = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.ExecutedIntervals = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TotalExecuted = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Cancelled = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err
}

// Check if interval is executed (view function)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) IsIntervalExecuted(opts *bind.CallOpts, orderId [32]byte, intervalIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "isIntervalExecuted", orderId, intervalIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// Add authorized resolver
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) AddResolver(opts *bind.TransactOpts, resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "addResolver", resolver)
}

// Get escrow factory address
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) EscrowFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "escrowFactory")
	if err != nil {
		return common.Address{}, err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// Event filtering structures
type FlowFusionEscrowFactoryTWAPOrderCreated struct {
	OrderId         [32]byte
	Maker           common.Address
	Token           common.Address
	TotalAmount     *big.Int
	TimeWindow      *big.Int
	IntervalCount   *big.Int
	CosmosRecipient string
	Raw             types.Log
}

type FlowFusionEscrowFactoryTWAPIntervalExecuted struct {
	OrderId       [32]byte
	IntervalIndex *big.Int
	Amount        *big.Int
	SecretHash    [32]byte
	Raw           types.Log
}

type FlowFusionEscrowFactoryCosmosEscrowCreated struct {
	OrderId         [32]byte
	Immutables      IBaseEscrowImmutables
	CosmosRecipient string
	Raw             types.Log
}

// Event filtering functions
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) FilterTWAPOrderCreated(opts *bind.FilterOpts, orderId [][32]byte, maker []common.Address, token []common.Address) (*FlowFusionEscrowFactoryTWAPOrderCreatedIterator, error) {
	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.FilterLogs(opts, "TWAPOrderCreated", orderIdRule, makerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryTWAPOrderCreatedIterator{contract: _FlowFusionEscrowFactory.contract, event: "TWAPOrderCreated", logs: logs, sub: sub}, nil
}

func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) FilterCosmosEscrowCreated(opts *bind.FilterOpts, orderId [][32]byte) (*FlowFusionEscrowFactoryCosmosEscrowCreatedIterator, error) {
	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.FilterLogs(opts, "CosmosEscrowCreated", orderIdRule)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryCosmosEscrowCreatedIterator{contract: _FlowFusionEscrowFactory.contract, event: "CosmosEscrowCreated", logs: logs, sub: sub}, nil
}

type FlowFusionEscrowFactoryTWAPOrderCreatedIterator struct {
	Event    *FlowFusionEscrowFactoryTWAPOrderCreated
	contract *bind.BoundContract
	event    string
	logs     chan types.Log
	sub      ethereum.Subscription
	done     bool
	fail     error
}

func (it *FlowFusionEscrowFactoryTWAPOrderCreatedIterator) Next() bool {
	if it.fail != nil {
		return false
	}
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowFusionEscrowFactoryTWAPOrderCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true
		default:
			return false
		}
	}
	select {
	case log := <-it.logs:
		it.Event = new(FlowFusionEscrowFactoryTWAPOrderCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true
	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *FlowFusionEscrowFactoryTWAPOrderCreatedIterator) Error() error {
	return it.fail
}

func (it *FlowFusionEscrowFactoryTWAPOrderCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FlowFusionEscrowFactoryCosmosEscrowCreatedIterator struct {
	Event    *FlowFusionEscrowFactoryCosmosEscrowCreated
	contract *bind.BoundContract
	event    string
	logs     chan types.Log
	sub      ethereum.Subscription
	done     bool
	fail     error
}

func (it *FlowFusionEscrowFactoryCosmosEscrowCreatedIterator) Next() bool {
	if it.fail != nil {
		return false
	}
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowFusionEscrowFactoryCosmosEscrowCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true
		default:
			return false
		}
	}
	select {
	case log := <-it.logs:
		it.Event = new(FlowFusionEscrowFactoryCosmosEscrowCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true
	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *FlowFusionEscrowFactoryCosmosEscrowCreatedIterator) Error() error {
	return it.fail
}

func (it *FlowFusionEscrowFactoryCosmosEscrowCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IEscrowFactoryABI is the ABI for the IEscrowFactory interface
var IEscrowFactoryABI = `[
	{
		"inputs": [
			{
				"components": [
					{"internalType": "bytes32", "name": "orderHash", "type": "bytes32"},
					{"internalType": "bytes32", "name": "hashLock", "type": "bytes32"},
					{"internalType": "address", "name": "maker", "type": "address"},
					{"internalType": "address", "name": "taker", "type": "address"},
					{"internalType": "address", "name": "token", "type": "address"},
					{"internalType": "uint256", "name": "amount", "type": "uint256"},
					{"internalType": "uint256", "name": "safetyDeposit", "type": "uint256"},
					{
						"components": [
							{"internalType": "uint256", "name": "prepayment", "type": "uint256"},
							{"internalType": "uint256", "name": "maturity", "type": "uint256"},
							{"internalType": "uint256", "name": "expiration", "type": "uint256"}
						],
						"name": "timelocks", "type": "tuple"
					}
				],
				"internalType": "struct IBaseEscrow.Immutables", "name": "immutables", "type": "tuple"
			},
			{"internalType": "uint256", "name": "timestamp", "type": "uint256"}
		],
		"name": "createSrcEscrow",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [
			{"internalType": "bytes32", "name": "orderHash", "type": "bytes32"}
		],
		"name": "escrows",
		"outputs": [
			{"internalType": "address", "name": "", "type": "address"}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

// IEscrowFactory is a Go binding for the IEscrowFactory contract
type IEscrowFactory struct {
	contract *bind.BoundContract
}

// NewIEscrowFactory creates a new instance of IEscrowFactory, bound to a specific deployed contract
func NewIEscrowFactory(address common.Address, backend bind.ContractBackend) (*IEscrowFactory, error) {
	parsed, err := abi.JSON(strings.NewReader(IEscrowFactoryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse IEscrowFactory ABI: %w", err)
	}
	contract := bind.NewBoundContract(address, parsed, backend, backend, backend)
	return &IEscrowFactory{contract: contract}, nil
}

// CreateSrcEscrow calls the createSrcEscrow function on the IEscrowFactory contract
func (f *IEscrowFactory) CreateSrcEscrow(opts *bind.TransactOpts, immutables IBaseEscrowImmutables, timestamp *big.Int) (*types.Transaction, error) {
	return f.contract.Transact(opts, "createSrcEscrow", immutables, timestamp)
}

// Escrows retrieves the escrow contract address for a given orderHash
func (f *IEscrowFactory) Escrows(opts *bind.CallOpts, orderHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := f.contract.Call(opts, &out, "escrows", orderHash)
	if err != nil {
		return common.Address{}, err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// IBaseEscrowABI is the assumed ABI for the IBaseEscrow interface
var IBaseEscrowABI = `[
	{
		"inputs": [
			{"internalType": "bytes32", "name": "secret", "type": "bytes32"}
		],
		"name": "revealSecret",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

// IBaseEscrow is a Go binding for the IBaseEscrow contract
type IBaseEscrow struct {
	contract *bind.BoundContract
}

// NewIBaseEscrow creates a new instance of IBaseEscrow, bound to a specific deployed contract
func NewIBaseEscrow(address common.Address, backend bind.ContractBackend) (*IBaseEscrow, error) {
	parsed, err := abi.JSON(strings.NewReader(IBaseEscrowABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse IBaseEscrow ABI: %w", err)
	}
	contract := bind.NewBoundContract(address, parsed, backend, backend, backend)
	return &IBaseEscrow{contract: contract}, nil
}

// RevealSecret calls the revealSecret function on the IBaseEscrow contract
func (e *IBaseEscrow) RevealSecret(opts *bind.TransactOpts, secret [32]byte) (*types.Transaction, error) {
	return e.contract.Transact(opts, "revealSecret", secret)
}
