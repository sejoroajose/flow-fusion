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

// IBaseEscrowImmutables is an auto generated low-level Go binding around an user-defined struct.
type IBaseEscrowImmutables struct {
	OrderHash     [32]byte
	Hashlock      [32]byte
	Maker         *big.Int
	Taker         *big.Int
	Token         *big.Int
	Amount        *big.Int
	SafetyDeposit *big.Int
	Timelocks     *big.Int
}

// IEscrowFactoryDstImmutablesComplement is an auto generated low-level Go binding around an user-defined struct.
type IEscrowFactoryDstImmutablesComplement struct {
	Maker         *big.Int
	Amount        *big.Int
	Token         *big.Int
	SafetyDeposit *big.Int
	ChainId       *big.Int
}

// IEscrowFactoryMetaData contains all meta data concerning the IEscrowFactory contract.
var IEscrowFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"ESCROW_DST_IMPLEMENTATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ESCROW_SRC_IMPLEMENTATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addressOfEscrowDst\",\"inputs\":[{\"name\":\"immutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addressOfEscrowSrc\",\"inputs\":[{\"name\":\"immutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createDstEscrow\",\"inputs\":[{\"name\":\"dstImmutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]},{\"name\":\"srcCancellationTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"createSrcEscrow\",\"inputs\":[{\"name\":\"srcImmutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]},{\"name\":\"dstCancellationTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"DstEscrowCreated\",\"inputs\":[{\"name\":\"escrow\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"taker\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"Address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SrcEscrowCreated\",\"inputs\":[{\"name\":\"srcImmutables\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]},{\"name\":\"dstImmutablesComplement\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIEscrowFactory.DstImmutablesComplement\",\"components\":[{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientEscrowBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCreationTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPartialFill\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecretsAmount\",\"inputs\":[]}]",
}

// IEscrowFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use IEscrowFactoryMetaData.ABI instead.
var IEscrowFactoryABI = IEscrowFactoryMetaData.ABI

// IEscrowFactory is an auto generated Go binding around an Ethereum contract.
type IEscrowFactory struct {
	IEscrowFactoryCaller     // Read-only binding to the contract
	IEscrowFactoryTransactor // Write-only binding to the contract
	IEscrowFactoryFilterer   // Log filterer for contract events
}

// IEscrowFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type IEscrowFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IEscrowFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IEscrowFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IEscrowFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IEscrowFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IEscrowFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IEscrowFactorySession struct {
	Contract     *IEscrowFactory   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IEscrowFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IEscrowFactoryCallerSession struct {
	Contract *IEscrowFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IEscrowFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IEscrowFactoryTransactorSession struct {
	Contract     *IEscrowFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IEscrowFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type IEscrowFactoryRaw struct {
	Contract *IEscrowFactory // Generic contract binding to access the raw methods on
}

// IEscrowFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IEscrowFactoryCallerRaw struct {
	Contract *IEscrowFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// IEscrowFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IEscrowFactoryTransactorRaw struct {
	Contract *IEscrowFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIEscrowFactory creates a new instance of IEscrowFactory, bound to a specific deployed contract.
func NewIEscrowFactory(address common.Address, backend bind.ContractBackend) (*IEscrowFactory, error) {
	contract, err := bindIEscrowFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IEscrowFactory{IEscrowFactoryCaller: IEscrowFactoryCaller{contract: contract}, IEscrowFactoryTransactor: IEscrowFactoryTransactor{contract: contract}, IEscrowFactoryFilterer: IEscrowFactoryFilterer{contract: contract}}, nil
}

// NewIEscrowFactoryCaller creates a new read-only instance of IEscrowFactory, bound to a specific deployed contract.
func NewIEscrowFactoryCaller(address common.Address, caller bind.ContractCaller) (*IEscrowFactoryCaller, error) {
	contract, err := bindIEscrowFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IEscrowFactoryCaller{contract: contract}, nil
}

// NewIEscrowFactoryTransactor creates a new write-only instance of IEscrowFactory, bound to a specific deployed contract.
func NewIEscrowFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*IEscrowFactoryTransactor, error) {
	contract, err := bindIEscrowFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IEscrowFactoryTransactor{contract: contract}, nil
}

// NewIEscrowFactoryFilterer creates a new log filterer instance of IEscrowFactory, bound to a specific deployed contract.
func NewIEscrowFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*IEscrowFactoryFilterer, error) {
	contract, err := bindIEscrowFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IEscrowFactoryFilterer{contract: contract}, nil
}

// bindIEscrowFactory binds a generic wrapper to an already deployed contract.
func bindIEscrowFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IEscrowFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IEscrowFactory *IEscrowFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IEscrowFactory.Contract.IEscrowFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IEscrowFactory *IEscrowFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.IEscrowFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IEscrowFactory *IEscrowFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.IEscrowFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IEscrowFactory *IEscrowFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IEscrowFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IEscrowFactory *IEscrowFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IEscrowFactory *IEscrowFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.contract.Transact(opts, method, params...)
}

// ESCROWDSTIMPLEMENTATION is a free data retrieval call binding the contract method 0xba551177.
//
// Solidity: function ESCROW_DST_IMPLEMENTATION() view returns(address)
func (_IEscrowFactory *IEscrowFactoryCaller) ESCROWDSTIMPLEMENTATION(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IEscrowFactory.contract.Call(opts, &out, "ESCROW_DST_IMPLEMENTATION")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ESCROWDSTIMPLEMENTATION is a free data retrieval call binding the contract method 0xba551177.
//
// Solidity: function ESCROW_DST_IMPLEMENTATION() view returns(address)
func (_IEscrowFactory *IEscrowFactorySession) ESCROWDSTIMPLEMENTATION() (common.Address, error) {
	return _IEscrowFactory.Contract.ESCROWDSTIMPLEMENTATION(&_IEscrowFactory.CallOpts)
}

// ESCROWDSTIMPLEMENTATION is a free data retrieval call binding the contract method 0xba551177.
//
// Solidity: function ESCROW_DST_IMPLEMENTATION() view returns(address)
func (_IEscrowFactory *IEscrowFactoryCallerSession) ESCROWDSTIMPLEMENTATION() (common.Address, error) {
	return _IEscrowFactory.Contract.ESCROWDSTIMPLEMENTATION(&_IEscrowFactory.CallOpts)
}

// ESCROWSRCIMPLEMENTATION is a free data retrieval call binding the contract method 0x7040f173.
//
// Solidity: function ESCROW_SRC_IMPLEMENTATION() view returns(address)
func (_IEscrowFactory *IEscrowFactoryCaller) ESCROWSRCIMPLEMENTATION(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IEscrowFactory.contract.Call(opts, &out, "ESCROW_SRC_IMPLEMENTATION")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ESCROWSRCIMPLEMENTATION is a free data retrieval call binding the contract method 0x7040f173.
//
// Solidity: function ESCROW_SRC_IMPLEMENTATION() view returns(address)
func (_IEscrowFactory *IEscrowFactorySession) ESCROWSRCIMPLEMENTATION() (common.Address, error) {
	return _IEscrowFactory.Contract.ESCROWSRCIMPLEMENTATION(&_IEscrowFactory.CallOpts)
}

// ESCROWSRCIMPLEMENTATION is a free data retrieval call binding the contract method 0x7040f173.
//
// Solidity: function ESCROW_SRC_IMPLEMENTATION() view returns(address)
func (_IEscrowFactory *IEscrowFactoryCallerSession) ESCROWSRCIMPLEMENTATION() (common.Address, error) {
	return _IEscrowFactory.Contract.ESCROWSRCIMPLEMENTATION(&_IEscrowFactory.CallOpts)
}

// AddressOfEscrowDst is a free data retrieval call binding the contract method 0xbe58e91c.
//
// Solidity: function addressOfEscrowDst((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) view returns(address)
func (_IEscrowFactory *IEscrowFactoryCaller) AddressOfEscrowDst(opts *bind.CallOpts, immutables IBaseEscrowImmutables) (common.Address, error) {
	var out []interface{}
	err := _IEscrowFactory.contract.Call(opts, &out, "addressOfEscrowDst", immutables)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AddressOfEscrowDst is a free data retrieval call binding the contract method 0xbe58e91c.
//
// Solidity: function addressOfEscrowDst((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) view returns(address)
func (_IEscrowFactory *IEscrowFactorySession) AddressOfEscrowDst(immutables IBaseEscrowImmutables) (common.Address, error) {
	return _IEscrowFactory.Contract.AddressOfEscrowDst(&_IEscrowFactory.CallOpts, immutables)
}

// AddressOfEscrowDst is a free data retrieval call binding the contract method 0xbe58e91c.
//
// Solidity: function addressOfEscrowDst((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) view returns(address)
func (_IEscrowFactory *IEscrowFactoryCallerSession) AddressOfEscrowDst(immutables IBaseEscrowImmutables) (common.Address, error) {
	return _IEscrowFactory.Contract.AddressOfEscrowDst(&_IEscrowFactory.CallOpts, immutables)
}

// AddressOfEscrowSrc is a free data retrieval call binding the contract method 0xfb6bd47e.
//
// Solidity: function addressOfEscrowSrc((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) view returns(address)
func (_IEscrowFactory *IEscrowFactoryCaller) AddressOfEscrowSrc(opts *bind.CallOpts, immutables IBaseEscrowImmutables) (common.Address, error) {
	var out []interface{}
	err := _IEscrowFactory.contract.Call(opts, &out, "addressOfEscrowSrc", immutables)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AddressOfEscrowSrc is a free data retrieval call binding the contract method 0xfb6bd47e.
//
// Solidity: function addressOfEscrowSrc((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) view returns(address)
func (_IEscrowFactory *IEscrowFactorySession) AddressOfEscrowSrc(immutables IBaseEscrowImmutables) (common.Address, error) {
	return _IEscrowFactory.Contract.AddressOfEscrowSrc(&_IEscrowFactory.CallOpts, immutables)
}

// AddressOfEscrowSrc is a free data retrieval call binding the contract method 0xfb6bd47e.
//
// Solidity: function addressOfEscrowSrc((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) view returns(address)
func (_IEscrowFactory *IEscrowFactoryCallerSession) AddressOfEscrowSrc(immutables IBaseEscrowImmutables) (common.Address, error) {
	return _IEscrowFactory.Contract.AddressOfEscrowSrc(&_IEscrowFactory.CallOpts, immutables)
}

// CreateDstEscrow is a paid mutator transaction binding the contract method 0xdea024e4.
//
// Solidity: function createDstEscrow((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) dstImmutables, uint256 srcCancellationTimestamp) payable returns()
func (_IEscrowFactory *IEscrowFactoryTransactor) CreateDstEscrow(opts *bind.TransactOpts, dstImmutables IBaseEscrowImmutables, srcCancellationTimestamp *big.Int) (*types.Transaction, error) {
	return _IEscrowFactory.contract.Transact(opts, "createDstEscrow", dstImmutables, srcCancellationTimestamp)
}

// CreateDstEscrow is a paid mutator transaction binding the contract method 0xdea024e4.
//
// Solidity: function createDstEscrow((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) dstImmutables, uint256 srcCancellationTimestamp) payable returns()
func (_IEscrowFactory *IEscrowFactorySession) CreateDstEscrow(dstImmutables IBaseEscrowImmutables, srcCancellationTimestamp *big.Int) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.CreateDstEscrow(&_IEscrowFactory.TransactOpts, dstImmutables, srcCancellationTimestamp)
}

// CreateDstEscrow is a paid mutator transaction binding the contract method 0xdea024e4.
//
// Solidity: function createDstEscrow((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) dstImmutables, uint256 srcCancellationTimestamp) payable returns()
func (_IEscrowFactory *IEscrowFactoryTransactorSession) CreateDstEscrow(dstImmutables IBaseEscrowImmutables, srcCancellationTimestamp *big.Int) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.CreateDstEscrow(&_IEscrowFactory.TransactOpts, dstImmutables, srcCancellationTimestamp)
}

// CreateSrcEscrow is a paid mutator transaction binding the contract method 0x82462834.
//
// Solidity: function createSrcEscrow((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) srcImmutables, uint256 dstCancellationTimestamp) payable returns()
func (_IEscrowFactory *IEscrowFactoryTransactor) CreateSrcEscrow(opts *bind.TransactOpts, srcImmutables IBaseEscrowImmutables, dstCancellationTimestamp *big.Int) (*types.Transaction, error) {
	return _IEscrowFactory.contract.Transact(opts, "createSrcEscrow", srcImmutables, dstCancellationTimestamp)
}

// CreateSrcEscrow is a paid mutator transaction binding the contract method 0x82462834.
//
// Solidity: function createSrcEscrow((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) srcImmutables, uint256 dstCancellationTimestamp) payable returns()
func (_IEscrowFactory *IEscrowFactorySession) CreateSrcEscrow(srcImmutables IBaseEscrowImmutables, dstCancellationTimestamp *big.Int) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.CreateSrcEscrow(&_IEscrowFactory.TransactOpts, srcImmutables, dstCancellationTimestamp)
}

// CreateSrcEscrow is a paid mutator transaction binding the contract method 0x82462834.
//
// Solidity: function createSrcEscrow((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) srcImmutables, uint256 dstCancellationTimestamp) payable returns()
func (_IEscrowFactory *IEscrowFactoryTransactorSession) CreateSrcEscrow(srcImmutables IBaseEscrowImmutables, dstCancellationTimestamp *big.Int) (*types.Transaction, error) {
	return _IEscrowFactory.Contract.CreateSrcEscrow(&_IEscrowFactory.TransactOpts, srcImmutables, dstCancellationTimestamp)
}

// IEscrowFactoryDstEscrowCreatedIterator is returned from FilterDstEscrowCreated and is used to iterate over the raw logs and unpacked data for DstEscrowCreated events raised by the IEscrowFactory contract.
type IEscrowFactoryDstEscrowCreatedIterator struct {
	Event *IEscrowFactoryDstEscrowCreated // Event containing the contract specifics and raw log

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
func (it *IEscrowFactoryDstEscrowCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IEscrowFactoryDstEscrowCreated)
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
		it.Event = new(IEscrowFactoryDstEscrowCreated)
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
func (it *IEscrowFactoryDstEscrowCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IEscrowFactoryDstEscrowCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IEscrowFactoryDstEscrowCreated represents a DstEscrowCreated event raised by the IEscrowFactory contract.
type IEscrowFactoryDstEscrowCreated struct {
	Escrow   common.Address
	Hashlock [32]byte
	Taker    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDstEscrowCreated is a free log retrieval operation binding the contract event 0xc30e111dcc74fddc2c3a4d98ffb97adec4485c0a687946bf5b22c2a99c7ff96d.
//
// Solidity: event DstEscrowCreated(address escrow, bytes32 hashlock, uint256 taker)
func (_IEscrowFactory *IEscrowFactoryFilterer) FilterDstEscrowCreated(opts *bind.FilterOpts) (*IEscrowFactoryDstEscrowCreatedIterator, error) {

	logs, sub, err := _IEscrowFactory.contract.FilterLogs(opts, "DstEscrowCreated")
	if err != nil {
		return nil, err
	}
	return &IEscrowFactoryDstEscrowCreatedIterator{contract: _IEscrowFactory.contract, event: "DstEscrowCreated", logs: logs, sub: sub}, nil
}

// WatchDstEscrowCreated is a free log subscription operation binding the contract event 0xc30e111dcc74fddc2c3a4d98ffb97adec4485c0a687946bf5b22c2a99c7ff96d.
//
// Solidity: event DstEscrowCreated(address escrow, bytes32 hashlock, uint256 taker)
func (_IEscrowFactory *IEscrowFactoryFilterer) WatchDstEscrowCreated(opts *bind.WatchOpts, sink chan<- *IEscrowFactoryDstEscrowCreated) (event.Subscription, error) {

	logs, sub, err := _IEscrowFactory.contract.WatchLogs(opts, "DstEscrowCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IEscrowFactoryDstEscrowCreated)
				if err := _IEscrowFactory.contract.UnpackLog(event, "DstEscrowCreated", log); err != nil {
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

// ParseDstEscrowCreated is a log parse operation binding the contract event 0xc30e111dcc74fddc2c3a4d98ffb97adec4485c0a687946bf5b22c2a99c7ff96d.
//
// Solidity: event DstEscrowCreated(address escrow, bytes32 hashlock, uint256 taker)
func (_IEscrowFactory *IEscrowFactoryFilterer) ParseDstEscrowCreated(log types.Log) (*IEscrowFactoryDstEscrowCreated, error) {
	event := new(IEscrowFactoryDstEscrowCreated)
	if err := _IEscrowFactory.contract.UnpackLog(event, "DstEscrowCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IEscrowFactorySrcEscrowCreatedIterator is returned from FilterSrcEscrowCreated and is used to iterate over the raw logs and unpacked data for SrcEscrowCreated events raised by the IEscrowFactory contract.
type IEscrowFactorySrcEscrowCreatedIterator struct {
	Event *IEscrowFactorySrcEscrowCreated // Event containing the contract specifics and raw log

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
func (it *IEscrowFactorySrcEscrowCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IEscrowFactorySrcEscrowCreated)
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
		it.Event = new(IEscrowFactorySrcEscrowCreated)
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
func (it *IEscrowFactorySrcEscrowCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IEscrowFactorySrcEscrowCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IEscrowFactorySrcEscrowCreated represents a SrcEscrowCreated event raised by the IEscrowFactory contract.
type IEscrowFactorySrcEscrowCreated struct {
	SrcImmutables           IBaseEscrowImmutables
	DstImmutablesComplement IEscrowFactoryDstImmutablesComplement
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterSrcEscrowCreated is a free log retrieval operation binding the contract event 0x0e534c62f0afd2fa0f0fa71198e8aa2d549f24daf2bb47de0d5486c7ce9288ca.
//
// Solidity: event SrcEscrowCreated((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) srcImmutables, (uint256,uint256,uint256,uint256,uint256) dstImmutablesComplement)
func (_IEscrowFactory *IEscrowFactoryFilterer) FilterSrcEscrowCreated(opts *bind.FilterOpts) (*IEscrowFactorySrcEscrowCreatedIterator, error) {

	logs, sub, err := _IEscrowFactory.contract.FilterLogs(opts, "SrcEscrowCreated")
	if err != nil {
		return nil, err
	}
	return &IEscrowFactorySrcEscrowCreatedIterator{contract: _IEscrowFactory.contract, event: "SrcEscrowCreated", logs: logs, sub: sub}, nil
}

// WatchSrcEscrowCreated is a free log subscription operation binding the contract event 0x0e534c62f0afd2fa0f0fa71198e8aa2d549f24daf2bb47de0d5486c7ce9288ca.
//
// Solidity: event SrcEscrowCreated((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) srcImmutables, (uint256,uint256,uint256,uint256,uint256) dstImmutablesComplement)
func (_IEscrowFactory *IEscrowFactoryFilterer) WatchSrcEscrowCreated(opts *bind.WatchOpts, sink chan<- *IEscrowFactorySrcEscrowCreated) (event.Subscription, error) {

	logs, sub, err := _IEscrowFactory.contract.WatchLogs(opts, "SrcEscrowCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IEscrowFactorySrcEscrowCreated)
				if err := _IEscrowFactory.contract.UnpackLog(event, "SrcEscrowCreated", log); err != nil {
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

// ParseSrcEscrowCreated is a log parse operation binding the contract event 0x0e534c62f0afd2fa0f0fa71198e8aa2d549f24daf2bb47de0d5486c7ce9288ca.
//
// Solidity: event SrcEscrowCreated((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) srcImmutables, (uint256,uint256,uint256,uint256,uint256) dstImmutablesComplement)
func (_IEscrowFactory *IEscrowFactoryFilterer) ParseSrcEscrowCreated(log types.Log) (*IEscrowFactorySrcEscrowCreated, error) {
	event := new(IEscrowFactorySrcEscrowCreated)
	if err := _IEscrowFactory.contract.UnpackLog(event, "SrcEscrowCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
