// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generated

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
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
	_ = abi.ConvertType
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
/* type IBaseEscrowImmutables struct {
	OrderHash     [32]byte
	Hashlock      [32]byte
	Maker         *big.Int
	Taker         *big.Int
	Token         *big.Int
	Amount        *big.Int
	SafetyDeposit *big.Int
	Timelocks     *big.Int
} */

// FlowFusionEscrowFactoryMetaData contains all meta data concerning the FlowFusionEscrowFactory contract.
var FlowFusionEscrowFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_escrowFactory\",\"type\":\"address\",\"internalType\":\"contractIEscrowFactory\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_INTERVALS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_SLIPPAGE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_INTERVAL_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addResolver\",\"inputs\":[{\"name\":\"resolver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedResolvers\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelTWAPOrder\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createTWAPOrder\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structFlowFusionEscrowFactory.TWAPConfig\",\"components\":[{\"name\":\"totalAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timeWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"intervalCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSlippage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"cosmosRecipient\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"emergencyWithdraw\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"escrowFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEscrowFactory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeTWAPInterval\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"intervalIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"secretHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"immutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getTWAPOrder\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structFlowFusionEscrowFactory.TWAPConfig\",\"components\":[{\"name\":\"totalAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timeWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"intervalCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSlippage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"cosmosRecipient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"executedIntervals\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalExecuted\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cancelled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isIntervalExecuted\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"intervalIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeResolver\",\"inputs\":[{\"name\":\"resolver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"twapOrders\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"config\",\"type\":\"tuple\",\"internalType\":\"structFlowFusionEscrowFactory.TWAPConfig\",\"components\":[{\"name\":\"totalAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timeWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"intervalCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxSlippage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"cosmosRecipient\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"executedIntervals\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalExecuted\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cancelled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"CosmosEscrowCreated\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"immutables\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]},{\"name\":\"cosmosRecipient\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TWAPIntervalExecuted\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"intervalIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"secretHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TWAPOrderCancelled\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"refundedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TWAPOrderCreated\",\"inputs\":[{\"name\":\"orderId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timeWindow\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"intervalCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cosmosRecipient\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IntervalAlreadyExecuted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInterval\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTWAPConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OrderAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OrderCancelled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OrderNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedResolver\",\"inputs\":[]}]",
}

// FlowFusionEscrowFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use FlowFusionEscrowFactoryMetaData.ABI instead.
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
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
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
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// FlowFusionEscrowFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlowFusionEscrowFactoryRaw struct {
	Contract *FlowFusionEscrowFactory // Generic contract binding to access the raw methods on
}

// FlowFusionEscrowFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlowFusionEscrowFactoryCallerRaw struct {
	Contract *FlowFusionEscrowFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// FlowFusionEscrowFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlowFusionEscrowFactoryTransactorRaw struct {
	Contract *FlowFusionEscrowFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlowFusionEscrowFactory creates a new instance of FlowFusionEscrowFactory, bound to a specific deployed contract.
func NewFlowFusionEscrowFactory(address common.Address, backend bind.ContractBackend) (*FlowFusionEscrowFactory, error) {
	contract, err := bindFlowFusionEscrowFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactory{FlowFusionEscrowFactoryCaller: FlowFusionEscrowFactoryCaller{contract: contract}, FlowFusionEscrowFactoryTransactor: FlowFusionEscrowFactoryTransactor{contract: contract}, FlowFusionEscrowFactoryFilterer: FlowFusionEscrowFactoryFilterer{contract: contract}}, nil
}

// NewFlowFusionEscrowFactoryCaller creates a new read-only instance of FlowFusionEscrowFactory, bound to a specific deployed contract.
func NewFlowFusionEscrowFactoryCaller(address common.Address, caller bind.ContractCaller) (*FlowFusionEscrowFactoryCaller, error) {
	contract, err := bindFlowFusionEscrowFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryCaller{contract: contract}, nil
}

// NewFlowFusionEscrowFactoryTransactor creates a new write-only instance of FlowFusionEscrowFactory, bound to a specific deployed contract.
func NewFlowFusionEscrowFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*FlowFusionEscrowFactoryTransactor, error) {
	contract, err := bindFlowFusionEscrowFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryTransactor{contract: contract}, nil
}

// NewFlowFusionEscrowFactoryFilterer creates a new log filterer instance of FlowFusionEscrowFactory, bound to a specific deployed contract.
func NewFlowFusionEscrowFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*FlowFusionEscrowFactoryFilterer, error) {
	contract, err := bindFlowFusionEscrowFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryFilterer{contract: contract}, nil
}

// bindFlowFusionEscrowFactory binds a generic wrapper to an already deployed contract.
func bindFlowFusionEscrowFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FlowFusionEscrowFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlowFusionEscrowFactory.Contract.FlowFusionEscrowFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.FlowFusionEscrowFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.FlowFusionEscrowFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlowFusionEscrowFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.contract.Transact(opts, method, params...)
}

// MAXINTERVALS is a free data retrieval call binding the contract method 0xe6211922.
//
// Solidity: function MAX_INTERVALS() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) MAXINTERVALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "MAX_INTERVALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXINTERVALS is a free data retrieval call binding the contract method 0xe6211922.
//
// Solidity: function MAX_INTERVALS() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) MAXINTERVALS() (*big.Int, error) {
	return _FlowFusionEscrowFactory.Contract.MAXINTERVALS(&_FlowFusionEscrowFactory.CallOpts)
}

// MAXINTERVALS is a free data retrieval call binding the contract method 0xe6211922.
//
// Solidity: function MAX_INTERVALS() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) MAXINTERVALS() (*big.Int, error) {
	return _FlowFusionEscrowFactory.Contract.MAXINTERVALS(&_FlowFusionEscrowFactory.CallOpts)
}

// MAXSLIPPAGE is a free data retrieval call binding the contract method 0xf9759518.
//
// Solidity: function MAX_SLIPPAGE() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) MAXSLIPPAGE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "MAX_SLIPPAGE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSLIPPAGE is a free data retrieval call binding the contract method 0xf9759518.
//
// Solidity: function MAX_SLIPPAGE() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) MAXSLIPPAGE() (*big.Int, error) {
	return _FlowFusionEscrowFactory.Contract.MAXSLIPPAGE(&_FlowFusionEscrowFactory.CallOpts)
}

// MAXSLIPPAGE is a free data retrieval call binding the contract method 0xf9759518.
//
// Solidity: function MAX_SLIPPAGE() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) MAXSLIPPAGE() (*big.Int, error) {
	return _FlowFusionEscrowFactory.Contract.MAXSLIPPAGE(&_FlowFusionEscrowFactory.CallOpts)
}

// MININTERVALDURATION is a free data retrieval call binding the contract method 0xd92621f3.
//
// Solidity: function MIN_INTERVAL_DURATION() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) MININTERVALDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "MIN_INTERVAL_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MININTERVALDURATION is a free data retrieval call binding the contract method 0xd92621f3.
//
// Solidity: function MIN_INTERVAL_DURATION() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) MININTERVALDURATION() (*big.Int, error) {
	return _FlowFusionEscrowFactory.Contract.MININTERVALDURATION(&_FlowFusionEscrowFactory.CallOpts)
}

// MININTERVALDURATION is a free data retrieval call binding the contract method 0xd92621f3.
//
// Solidity: function MIN_INTERVAL_DURATION() view returns(uint256)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) MININTERVALDURATION() (*big.Int, error) {
	return _FlowFusionEscrowFactory.Contract.MININTERVALDURATION(&_FlowFusionEscrowFactory.CallOpts)
}

// AuthorizedResolvers is a free data retrieval call binding the contract method 0xd794d2d9.
//
// Solidity: function authorizedResolvers(address ) view returns(bool)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) AuthorizedResolvers(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "authorizedResolvers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedResolvers is a free data retrieval call binding the contract method 0xd794d2d9.
//
// Solidity: function authorizedResolvers(address ) view returns(bool)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) AuthorizedResolvers(arg0 common.Address) (bool, error) {
	return _FlowFusionEscrowFactory.Contract.AuthorizedResolvers(&_FlowFusionEscrowFactory.CallOpts, arg0)
}

// AuthorizedResolvers is a free data retrieval call binding the contract method 0xd794d2d9.
//
// Solidity: function authorizedResolvers(address ) view returns(bool)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) AuthorizedResolvers(arg0 common.Address) (bool, error) {
	return _FlowFusionEscrowFactory.Contract.AuthorizedResolvers(&_FlowFusionEscrowFactory.CallOpts, arg0)
}

// EscrowFactory is a free data retrieval call binding the contract method 0xbdd1daaa.
//
// Solidity: function escrowFactory() view returns(address)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) EscrowFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "escrowFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EscrowFactory is a free data retrieval call binding the contract method 0xbdd1daaa.
//
// Solidity: function escrowFactory() view returns(address)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) EscrowFactory() (common.Address, error) {
	return _FlowFusionEscrowFactory.Contract.EscrowFactory(&_FlowFusionEscrowFactory.CallOpts)
}

// EscrowFactory is a free data retrieval call binding the contract method 0xbdd1daaa.
//
// Solidity: function escrowFactory() view returns(address)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) EscrowFactory() (common.Address, error) {
	return _FlowFusionEscrowFactory.Contract.EscrowFactory(&_FlowFusionEscrowFactory.CallOpts)
}

// GetTWAPOrder is a free data retrieval call binding the contract method 0xbe22fc9b.
//
// Solidity: function getTWAPOrder(bytes32 orderId) view returns(address maker, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient, uint256 executedIntervals, uint256 totalExecuted, bool cancelled)
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

// GetTWAPOrder is a free data retrieval call binding the contract method 0xbe22fc9b.
//
// Solidity: function getTWAPOrder(bytes32 orderId) view returns(address maker, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient, uint256 executedIntervals, uint256 totalExecuted, bool cancelled)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) GetTWAPOrder(orderId [32]byte) (struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}, error) {
	return _FlowFusionEscrowFactory.Contract.GetTWAPOrder(&_FlowFusionEscrowFactory.CallOpts, orderId)
}

// GetTWAPOrder is a free data retrieval call binding the contract method 0xbe22fc9b.
//
// Solidity: function getTWAPOrder(bytes32 orderId) view returns(address maker, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient, uint256 executedIntervals, uint256 totalExecuted, bool cancelled)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) GetTWAPOrder(orderId [32]byte) (struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}, error) {
	return _FlowFusionEscrowFactory.Contract.GetTWAPOrder(&_FlowFusionEscrowFactory.CallOpts, orderId)
}

// IsIntervalExecuted is a free data retrieval call binding the contract method 0x50c30967.
//
// Solidity: function isIntervalExecuted(bytes32 orderId, uint256 intervalIndex) view returns(bool)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) IsIntervalExecuted(opts *bind.CallOpts, orderId [32]byte, intervalIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "isIntervalExecuted", orderId, intervalIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsIntervalExecuted is a free data retrieval call binding the contract method 0x50c30967.
//
// Solidity: function isIntervalExecuted(bytes32 orderId, uint256 intervalIndex) view returns(bool)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) IsIntervalExecuted(orderId [32]byte, intervalIndex *big.Int) (bool, error) {
	return _FlowFusionEscrowFactory.Contract.IsIntervalExecuted(&_FlowFusionEscrowFactory.CallOpts, orderId, intervalIndex)
}

// IsIntervalExecuted is a free data retrieval call binding the contract method 0x50c30967.
//
// Solidity: function isIntervalExecuted(bytes32 orderId, uint256 intervalIndex) view returns(bool)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) IsIntervalExecuted(orderId [32]byte, intervalIndex *big.Int) (bool, error) {
	return _FlowFusionEscrowFactory.Contract.IsIntervalExecuted(&_FlowFusionEscrowFactory.CallOpts, orderId, intervalIndex)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) Owner() (common.Address, error) {
	return _FlowFusionEscrowFactory.Contract.Owner(&_FlowFusionEscrowFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) Owner() (common.Address, error) {
	return _FlowFusionEscrowFactory.Contract.Owner(&_FlowFusionEscrowFactory.CallOpts)
}

// TwapOrders is a free data retrieval call binding the contract method 0x4149bc28.
//
// Solidity: function twapOrders(bytes32 ) view returns(address maker, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient, uint256 executedIntervals, uint256 totalExecuted, bool cancelled)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCaller) TwapOrders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}, error) {
	var out []interface{}
	err := _FlowFusionEscrowFactory.contract.Call(opts, &out, "twapOrders", arg0)

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

// TwapOrders is a free data retrieval call binding the contract method 0x4149bc28.
//
// Solidity: function twapOrders(bytes32 ) view returns(address maker, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient, uint256 executedIntervals, uint256 totalExecuted, bool cancelled)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) TwapOrders(arg0 [32]byte) (struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}, error) {
	return _FlowFusionEscrowFactory.Contract.TwapOrders(&_FlowFusionEscrowFactory.CallOpts, arg0)
}

// TwapOrders is a free data retrieval call binding the contract method 0x4149bc28.
//
// Solidity: function twapOrders(bytes32 ) view returns(address maker, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient, uint256 executedIntervals, uint256 totalExecuted, bool cancelled)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryCallerSession) TwapOrders(arg0 [32]byte) (struct {
	Maker             common.Address
	Token             common.Address
	Config            FlowFusionEscrowFactoryTWAPConfig
	CosmosRecipient   string
	ExecutedIntervals *big.Int
	TotalExecuted     *big.Int
	Cancelled         bool
}, error) {
	return _FlowFusionEscrowFactory.Contract.TwapOrders(&_FlowFusionEscrowFactory.CallOpts, arg0)
}

// AddResolver is a paid mutator transaction binding the contract method 0x390a0b6f.
//
// Solidity: function addResolver(address resolver) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) AddResolver(opts *bind.TransactOpts, resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "addResolver", resolver)
}

// AddResolver is a paid mutator transaction binding the contract method 0x390a0b6f.
//
// Solidity: function addResolver(address resolver) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) AddResolver(resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.AddResolver(&_FlowFusionEscrowFactory.TransactOpts, resolver)
}

// AddResolver is a paid mutator transaction binding the contract method 0x390a0b6f.
//
// Solidity: function addResolver(address resolver) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) AddResolver(resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.AddResolver(&_FlowFusionEscrowFactory.TransactOpts, resolver)
}

// CancelTWAPOrder is a paid mutator transaction binding the contract method 0x333e7c29.
//
// Solidity: function cancelTWAPOrder(bytes32 orderId) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) CancelTWAPOrder(opts *bind.TransactOpts, orderId [32]byte) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "cancelTWAPOrder", orderId)
}

// CancelTWAPOrder is a paid mutator transaction binding the contract method 0x333e7c29.
//
// Solidity: function cancelTWAPOrder(bytes32 orderId) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) CancelTWAPOrder(orderId [32]byte) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.CancelTWAPOrder(&_FlowFusionEscrowFactory.TransactOpts, orderId)
}

// CancelTWAPOrder is a paid mutator transaction binding the contract method 0x333e7c29.
//
// Solidity: function cancelTWAPOrder(bytes32 orderId) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) CancelTWAPOrder(orderId [32]byte) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.CancelTWAPOrder(&_FlowFusionEscrowFactory.TransactOpts, orderId)
}

// CreateTWAPOrder is a paid mutator transaction binding the contract method 0x12521df8.
//
// Solidity: function createTWAPOrder(bytes32 orderId, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient) payable returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) CreateTWAPOrder(opts *bind.TransactOpts, orderId [32]byte, token common.Address, config FlowFusionEscrowFactoryTWAPConfig, cosmosRecipient string) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "createTWAPOrder", orderId, token, config, cosmosRecipient)
}

// CreateTWAPOrder is a paid mutator transaction binding the contract method 0x12521df8.
//
// Solidity: function createTWAPOrder(bytes32 orderId, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient) payable returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) CreateTWAPOrder(orderId [32]byte, token common.Address, config FlowFusionEscrowFactoryTWAPConfig, cosmosRecipient string) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.CreateTWAPOrder(&_FlowFusionEscrowFactory.TransactOpts, orderId, token, config, cosmosRecipient)
}

// CreateTWAPOrder is a paid mutator transaction binding the contract method 0x12521df8.
//
// Solidity: function createTWAPOrder(bytes32 orderId, address token, (uint256,uint256,uint256,uint256,uint256) config, string cosmosRecipient) payable returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) CreateTWAPOrder(orderId [32]byte, token common.Address, config FlowFusionEscrowFactoryTWAPConfig, cosmosRecipient string) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.CreateTWAPOrder(&_FlowFusionEscrowFactory.TransactOpts, orderId, token, config, cosmosRecipient)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) EmergencyWithdraw(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "emergencyWithdraw", token, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) EmergencyWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.EmergencyWithdraw(&_FlowFusionEscrowFactory.TransactOpts, token, amount)
}

// EmergencyWithdraw is a paid mutator transaction binding the contract method 0x95ccea67.
//
// Solidity: function emergencyWithdraw(address token, uint256 amount) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) EmergencyWithdraw(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.EmergencyWithdraw(&_FlowFusionEscrowFactory.TransactOpts, token, amount)
}

// ExecuteTWAPInterval is a paid mutator transaction binding the contract method 0x289a1279.
//
// Solidity: function executeTWAPInterval(bytes32 orderId, uint256 intervalIndex, bytes32 secretHash, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) payable returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) ExecuteTWAPInterval(opts *bind.TransactOpts, orderId [32]byte, intervalIndex *big.Int, secretHash [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "executeTWAPInterval", orderId, intervalIndex, secretHash, immutables)
}

// ExecuteTWAPInterval is a paid mutator transaction binding the contract method 0x289a1279.
//
// Solidity: function executeTWAPInterval(bytes32 orderId, uint256 intervalIndex, bytes32 secretHash, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) payable returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) ExecuteTWAPInterval(orderId [32]byte, intervalIndex *big.Int, secretHash [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.ExecuteTWAPInterval(&_FlowFusionEscrowFactory.TransactOpts, orderId, intervalIndex, secretHash, immutables)
}

// ExecuteTWAPInterval is a paid mutator transaction binding the contract method 0x289a1279.
//
// Solidity: function executeTWAPInterval(bytes32 orderId, uint256 intervalIndex, bytes32 secretHash, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) payable returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) ExecuteTWAPInterval(orderId [32]byte, intervalIndex *big.Int, secretHash [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.ExecuteTWAPInterval(&_FlowFusionEscrowFactory.TransactOpts, orderId, intervalIndex, secretHash, immutables)
}

// RemoveResolver is a paid mutator transaction binding the contract method 0xa1333536.
//
// Solidity: function removeResolver(address resolver) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) RemoveResolver(opts *bind.TransactOpts, resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "removeResolver", resolver)
}

// RemoveResolver is a paid mutator transaction binding the contract method 0xa1333536.
//
// Solidity: function removeResolver(address resolver) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) RemoveResolver(resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.RemoveResolver(&_FlowFusionEscrowFactory.TransactOpts, resolver)
}

// RemoveResolver is a paid mutator transaction binding the contract method 0xa1333536.
//
// Solidity: function removeResolver(address resolver) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) RemoveResolver(resolver common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.RemoveResolver(&_FlowFusionEscrowFactory.TransactOpts, resolver)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.RenounceOwnership(&_FlowFusionEscrowFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.RenounceOwnership(&_FlowFusionEscrowFactory.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.TransferOwnership(&_FlowFusionEscrowFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FlowFusionEscrowFactory.Contract.TransferOwnership(&_FlowFusionEscrowFactory.TransactOpts, newOwner)
}

// FlowFusionEscrowFactoryCosmosEscrowCreatedIterator is returned from FilterCosmosEscrowCreated and is used to iterate over the raw logs and unpacked data for CosmosEscrowCreated events raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryCosmosEscrowCreatedIterator struct {
	Event *FlowFusionEscrowFactoryCosmosEscrowCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowFusionEscrowFactoryCosmosEscrowCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowFusionEscrowFactoryCosmosEscrowCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowFusionEscrowFactoryCosmosEscrowCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowFusionEscrowFactoryCosmosEscrowCreated represents a CosmosEscrowCreated event raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryCosmosEscrowCreated struct {
	OrderId         [32]byte
	Immutables      IBaseEscrowImmutables
	CosmosRecipient string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCosmosEscrowCreated is a free log retrieval operation binding the contract event 0xc67a6eed83f9112ae5c6f418b154a38a812b99b7946903b95f48067259c8388c.
//
// Solidity: event CosmosEscrowCreated(bytes32 indexed orderId, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables, string cosmosRecipient)
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

// WatchCosmosEscrowCreated is a free log subscription operation binding the contract event 0xc67a6eed83f9112ae5c6f418b154a38a812b99b7946903b95f48067259c8388c.
//
// Solidity: event CosmosEscrowCreated(bytes32 indexed orderId, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables, string cosmosRecipient)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) WatchCosmosEscrowCreated(opts *bind.WatchOpts, sink chan<- *FlowFusionEscrowFactoryCosmosEscrowCreated, orderId [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.WatchLogs(opts, "CosmosEscrowCreated", orderIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowFusionEscrowFactoryCosmosEscrowCreated)
				if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "CosmosEscrowCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCosmosEscrowCreated is a log parse operation binding the contract event 0xc67a6eed83f9112ae5c6f418b154a38a812b99b7946903b95f48067259c8388c.
//
// Solidity: event CosmosEscrowCreated(bytes32 indexed orderId, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables, string cosmosRecipient)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) ParseCosmosEscrowCreated(log types.Log) (*FlowFusionEscrowFactoryCosmosEscrowCreated, error) {
	event := new(FlowFusionEscrowFactoryCosmosEscrowCreated)
	if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "CosmosEscrowCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowFusionEscrowFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryOwnershipTransferredIterator struct {
	Event *FlowFusionEscrowFactoryOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowFusionEscrowFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowFusionEscrowFactoryOwnershipTransferred)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowFusionEscrowFactoryOwnershipTransferred)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowFusionEscrowFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowFusionEscrowFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowFusionEscrowFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FlowFusionEscrowFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryOwnershipTransferredIterator{contract: _FlowFusionEscrowFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FlowFusionEscrowFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowFusionEscrowFactoryOwnershipTransferred)
				if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*FlowFusionEscrowFactoryOwnershipTransferred, error) {
	event := new(FlowFusionEscrowFactoryOwnershipTransferred)
	if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowFusionEscrowFactoryTWAPIntervalExecutedIterator is returned from FilterTWAPIntervalExecuted and is used to iterate over the raw logs and unpacked data for TWAPIntervalExecuted events raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryTWAPIntervalExecutedIterator struct {
	Event *FlowFusionEscrowFactoryTWAPIntervalExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowFusionEscrowFactoryTWAPIntervalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowFusionEscrowFactoryTWAPIntervalExecuted)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowFusionEscrowFactoryTWAPIntervalExecuted)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowFusionEscrowFactoryTWAPIntervalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowFusionEscrowFactoryTWAPIntervalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowFusionEscrowFactoryTWAPIntervalExecuted represents a TWAPIntervalExecuted event raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryTWAPIntervalExecuted struct {
	OrderId       [32]byte
	IntervalIndex *big.Int
	Amount        *big.Int
	SecretHash    [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTWAPIntervalExecuted is a free log retrieval operation binding the contract event 0xf95c8cd56b169f7c73d5c50943f1850a7137a3def1e51cf36fb4db0b832356af.
//
// Solidity: event TWAPIntervalExecuted(bytes32 indexed orderId, uint256 indexed intervalIndex, uint256 amount, bytes32 secretHash)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) FilterTWAPIntervalExecuted(opts *bind.FilterOpts, orderId [][32]byte, intervalIndex []*big.Int) (*FlowFusionEscrowFactoryTWAPIntervalExecutedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var intervalIndexRule []interface{}
	for _, intervalIndexItem := range intervalIndex {
		intervalIndexRule = append(intervalIndexRule, intervalIndexItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.FilterLogs(opts, "TWAPIntervalExecuted", orderIdRule, intervalIndexRule)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryTWAPIntervalExecutedIterator{contract: _FlowFusionEscrowFactory.contract, event: "TWAPIntervalExecuted", logs: logs, sub: sub}, nil
}

// WatchTWAPIntervalExecuted is a free log subscription operation binding the contract event 0xf95c8cd56b169f7c73d5c50943f1850a7137a3def1e51cf36fb4db0b832356af.
//
// Solidity: event TWAPIntervalExecuted(bytes32 indexed orderId, uint256 indexed intervalIndex, uint256 amount, bytes32 secretHash)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) WatchTWAPIntervalExecuted(opts *bind.WatchOpts, sink chan<- *FlowFusionEscrowFactoryTWAPIntervalExecuted, orderId [][32]byte, intervalIndex []*big.Int) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var intervalIndexRule []interface{}
	for _, intervalIndexItem := range intervalIndex {
		intervalIndexRule = append(intervalIndexRule, intervalIndexItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.WatchLogs(opts, "TWAPIntervalExecuted", orderIdRule, intervalIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowFusionEscrowFactoryTWAPIntervalExecuted)
				if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "TWAPIntervalExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTWAPIntervalExecuted is a log parse operation binding the contract event 0xf95c8cd56b169f7c73d5c50943f1850a7137a3def1e51cf36fb4db0b832356af.
//
// Solidity: event TWAPIntervalExecuted(bytes32 indexed orderId, uint256 indexed intervalIndex, uint256 amount, bytes32 secretHash)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) ParseTWAPIntervalExecuted(log types.Log) (*FlowFusionEscrowFactoryTWAPIntervalExecuted, error) {
	event := new(FlowFusionEscrowFactoryTWAPIntervalExecuted)
	if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "TWAPIntervalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowFusionEscrowFactoryTWAPOrderCancelledIterator is returned from FilterTWAPOrderCancelled and is used to iterate over the raw logs and unpacked data for TWAPOrderCancelled events raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryTWAPOrderCancelledIterator struct {
	Event *FlowFusionEscrowFactoryTWAPOrderCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowFusionEscrowFactoryTWAPOrderCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowFusionEscrowFactoryTWAPOrderCancelled)
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
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowFusionEscrowFactoryTWAPOrderCancelled)
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowFusionEscrowFactoryTWAPOrderCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowFusionEscrowFactoryTWAPOrderCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowFusionEscrowFactoryTWAPOrderCancelled represents a TWAPOrderCancelled event raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryTWAPOrderCancelled struct {
	OrderId        [32]byte
	Maker          common.Address
	RefundedAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterTWAPOrderCancelled is a free log retrieval operation binding the contract event 0xda08c4e9e894e13c604fc091b32436f34a9dcb0d2f79f093b1f2ac6898c101e7.
//
// Solidity: event TWAPOrderCancelled(bytes32 indexed orderId, address indexed maker, uint256 refundedAmount)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) FilterTWAPOrderCancelled(opts *bind.FilterOpts, orderId [][32]byte, maker []common.Address) (*FlowFusionEscrowFactoryTWAPOrderCancelledIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.FilterLogs(opts, "TWAPOrderCancelled", orderIdRule, makerRule)
	if err != nil {
		return nil, err
	}
	return &FlowFusionEscrowFactoryTWAPOrderCancelledIterator{contract: _FlowFusionEscrowFactory.contract, event: "TWAPOrderCancelled", logs: logs, sub: sub}, nil
}

// WatchTWAPOrderCancelled is a free log subscription operation binding the contract event 0xda08c4e9e894e13c604fc091b32436f34a9dcb0d2f79f093b1f2ac6898c101e7.
//
// Solidity: event TWAPOrderCancelled(bytes32 indexed orderId, address indexed maker, uint256 refundedAmount)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) WatchTWAPOrderCancelled(opts *bind.WatchOpts, sink chan<- *FlowFusionEscrowFactoryTWAPOrderCancelled, orderId [][32]byte, maker []common.Address) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}

	logs, sub, err := _FlowFusionEscrowFactory.contract.WatchLogs(opts, "TWAPOrderCancelled", orderIdRule, makerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowFusionEscrowFactoryTWAPOrderCancelled)
				if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "TWAPOrderCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTWAPOrderCancelled is a log parse operation binding the contract event 0xda08c4e9e894e13c604fc091b32436f34a9dcb0d2f79f093b1f2ac6898c101e7.
//
// Solidity: event TWAPOrderCancelled(bytes32 indexed orderId, address indexed maker, uint256 refundedAmount)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) ParseTWAPOrderCancelled(log types.Log) (*FlowFusionEscrowFactoryTWAPOrderCancelled, error) {
	event := new(FlowFusionEscrowFactoryTWAPOrderCancelled)
	if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "TWAPOrderCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowFusionEscrowFactoryTWAPOrderCreatedIterator is returned from FilterTWAPOrderCreated and is used to iterate over the raw logs and unpacked data for TWAPOrderCreated events raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryTWAPOrderCreatedIterator struct {
	Event *FlowFusionEscrowFactoryTWAPOrderCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowFusionEscrowFactoryTWAPOrderCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
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
	// Iterator still in progress, wait for either a data or an error event
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

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowFusionEscrowFactoryTWAPOrderCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowFusionEscrowFactoryTWAPOrderCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowFusionEscrowFactoryTWAPOrderCreated represents a TWAPOrderCreated event raised by the FlowFusionEscrowFactory contract.
type FlowFusionEscrowFactoryTWAPOrderCreated struct {
	OrderId         [32]byte
	Maker           common.Address
	Token           common.Address
	TotalAmount     *big.Int
	TimeWindow      *big.Int
	IntervalCount   *big.Int
	CosmosRecipient string
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTWAPOrderCreated is a free log retrieval operation binding the contract event 0x9d84b5f427b1488867bce8ddae5560fb91496d2964ae15ae199a46fdf19ebf33.
//
// Solidity: event TWAPOrderCreated(bytes32 indexed orderId, address indexed maker, address indexed token, uint256 totalAmount, uint256 timeWindow, uint256 intervalCount, string cosmosRecipient)
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

// WatchTWAPOrderCreated is a free log subscription operation binding the contract event 0x9d84b5f427b1488867bce8ddae5560fb91496d2964ae15ae199a46fdf19ebf33.
//
// Solidity: event TWAPOrderCreated(bytes32 indexed orderId, address indexed maker, address indexed token, uint256 totalAmount, uint256 timeWindow, uint256 intervalCount, string cosmosRecipient)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) WatchTWAPOrderCreated(opts *bind.WatchOpts, sink chan<- *FlowFusionEscrowFactoryTWAPOrderCreated, orderId [][32]byte, maker []common.Address, token []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _FlowFusionEscrowFactory.contract.WatchLogs(opts, "TWAPOrderCreated", orderIdRule, makerRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowFusionEscrowFactoryTWAPOrderCreated)
				if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "TWAPOrderCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTWAPOrderCreated is a log parse operation binding the contract event 0x9d84b5f427b1488867bce8ddae5560fb91496d2964ae15ae199a46fdf19ebf33.
//
// Solidity: event TWAPOrderCreated(bytes32 indexed orderId, address indexed maker, address indexed token, uint256 totalAmount, uint256 timeWindow, uint256 intervalCount, string cosmosRecipient)
func (_FlowFusionEscrowFactory *FlowFusionEscrowFactoryFilterer) ParseTWAPOrderCreated(log types.Log) (*FlowFusionEscrowFactoryTWAPOrderCreated, error) {
	event := new(FlowFusionEscrowFactoryTWAPOrderCreated)
	if err := _FlowFusionEscrowFactory.contract.UnpackLog(event, "TWAPOrderCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
