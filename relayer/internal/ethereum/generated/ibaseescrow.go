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

// IBaseEscrowMetaData contains all meta data concerning the IBaseEscrow contract.
var IBaseEscrowMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"FACTORY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RESCUE_DELAY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancel\",\"inputs\":[{\"name\":\"immutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rescueFunds\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"immutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"immutables\",\"type\":\"tuple\",\"internalType\":\"structIBaseEscrow.Immutables\",\"components\":[{\"name\":\"orderHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hashlock\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"taker\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"token\",\"type\":\"uint256\",\"internalType\":\"Address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"safetyDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"timelocks\",\"type\":\"uint256\",\"internalType\":\"Timelocks\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"EscrowCancelled\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsRescued\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawal\",\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidCaller\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidImmutables\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecret\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NativeTokenSendingFailure\",\"inputs\":[]}]",
}

// IBaseEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use IBaseEscrowMetaData.ABI instead.
var IBaseEscrowABI = IBaseEscrowMetaData.ABI

// IBaseEscrow is an auto generated Go binding around an Ethereum contract.
type IBaseEscrow struct {
	IBaseEscrowCaller     // Read-only binding to the contract
	IBaseEscrowTransactor // Write-only binding to the contract
	IBaseEscrowFilterer   // Log filterer for contract events
}

// IBaseEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type IBaseEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBaseEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IBaseEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBaseEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IBaseEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBaseEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IBaseEscrowSession struct {
	Contract     *IBaseEscrow      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBaseEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IBaseEscrowCallerSession struct {
	Contract *IBaseEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// IBaseEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IBaseEscrowTransactorSession struct {
	Contract     *IBaseEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// IBaseEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type IBaseEscrowRaw struct {
	Contract *IBaseEscrow // Generic contract binding to access the raw methods on
}

// IBaseEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IBaseEscrowCallerRaw struct {
	Contract *IBaseEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// IBaseEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IBaseEscrowTransactorRaw struct {
	Contract *IBaseEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBaseEscrow creates a new instance of IBaseEscrow, bound to a specific deployed contract.
func NewIBaseEscrow(address common.Address, backend bind.ContractBackend) (*IBaseEscrow, error) {
	contract, err := bindIBaseEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBaseEscrow{IBaseEscrowCaller: IBaseEscrowCaller{contract: contract}, IBaseEscrowTransactor: IBaseEscrowTransactor{contract: contract}, IBaseEscrowFilterer: IBaseEscrowFilterer{contract: contract}}, nil
}

// NewIBaseEscrowCaller creates a new read-only instance of IBaseEscrow, bound to a specific deployed contract.
func NewIBaseEscrowCaller(address common.Address, caller bind.ContractCaller) (*IBaseEscrowCaller, error) {
	contract, err := bindIBaseEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBaseEscrowCaller{contract: contract}, nil
}

// NewIBaseEscrowTransactor creates a new write-only instance of IBaseEscrow, bound to a specific deployed contract.
func NewIBaseEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*IBaseEscrowTransactor, error) {
	contract, err := bindIBaseEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBaseEscrowTransactor{contract: contract}, nil
}

// NewIBaseEscrowFilterer creates a new log filterer instance of IBaseEscrow, bound to a specific deployed contract.
func NewIBaseEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*IBaseEscrowFilterer, error) {
	contract, err := bindIBaseEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBaseEscrowFilterer{contract: contract}, nil
}

// bindIBaseEscrow binds a generic wrapper to an already deployed contract.
func bindIBaseEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBaseEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBaseEscrow *IBaseEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBaseEscrow.Contract.IBaseEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBaseEscrow *IBaseEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.IBaseEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBaseEscrow *IBaseEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.IBaseEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBaseEscrow *IBaseEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBaseEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBaseEscrow *IBaseEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBaseEscrow *IBaseEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.contract.Transact(opts, method, params...)
}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_IBaseEscrow *IBaseEscrowCaller) FACTORY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IBaseEscrow.contract.Call(opts, &out, "FACTORY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_IBaseEscrow *IBaseEscrowSession) FACTORY() (common.Address, error) {
	return _IBaseEscrow.Contract.FACTORY(&_IBaseEscrow.CallOpts)
}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_IBaseEscrow *IBaseEscrowCallerSession) FACTORY() (common.Address, error) {
	return _IBaseEscrow.Contract.FACTORY(&_IBaseEscrow.CallOpts)
}

// RESCUEDELAY is a free data retrieval call binding the contract method 0xf56cd69c.
//
// Solidity: function RESCUE_DELAY() view returns(uint256)
func (_IBaseEscrow *IBaseEscrowCaller) RESCUEDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IBaseEscrow.contract.Call(opts, &out, "RESCUE_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RESCUEDELAY is a free data retrieval call binding the contract method 0xf56cd69c.
//
// Solidity: function RESCUE_DELAY() view returns(uint256)
func (_IBaseEscrow *IBaseEscrowSession) RESCUEDELAY() (*big.Int, error) {
	return _IBaseEscrow.Contract.RESCUEDELAY(&_IBaseEscrow.CallOpts)
}

// RESCUEDELAY is a free data retrieval call binding the contract method 0xf56cd69c.
//
// Solidity: function RESCUE_DELAY() view returns(uint256)
func (_IBaseEscrow *IBaseEscrowCallerSession) RESCUEDELAY() (*big.Int, error) {
	return _IBaseEscrow.Contract.RESCUEDELAY(&_IBaseEscrow.CallOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0x90d3252f.
//
// Solidity: function cancel((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowTransactor) Cancel(opts *bind.TransactOpts, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.contract.Transact(opts, "cancel", immutables)
}

// Cancel is a paid mutator transaction binding the contract method 0x90d3252f.
//
// Solidity: function cancel((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowSession) Cancel(immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.Cancel(&_IBaseEscrow.TransactOpts, immutables)
}

// Cancel is a paid mutator transaction binding the contract method 0x90d3252f.
//
// Solidity: function cancel((bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowTransactorSession) Cancel(immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.Cancel(&_IBaseEscrow.TransactOpts, immutables)
}

// RescueFunds is a paid mutator transaction binding the contract method 0x4649088b.
//
// Solidity: function rescueFunds(address token, uint256 amount, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowTransactor) RescueFunds(opts *bind.TransactOpts, token common.Address, amount *big.Int, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.contract.Transact(opts, "rescueFunds", token, amount, immutables)
}

// RescueFunds is a paid mutator transaction binding the contract method 0x4649088b.
//
// Solidity: function rescueFunds(address token, uint256 amount, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowSession) RescueFunds(token common.Address, amount *big.Int, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.RescueFunds(&_IBaseEscrow.TransactOpts, token, amount, immutables)
}

// RescueFunds is a paid mutator transaction binding the contract method 0x4649088b.
//
// Solidity: function rescueFunds(address token, uint256 amount, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowTransactorSession) RescueFunds(token common.Address, amount *big.Int, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.RescueFunds(&_IBaseEscrow.TransactOpts, token, amount, immutables)
}

// Withdraw is a paid mutator transaction binding the contract method 0x23305703.
//
// Solidity: function withdraw(bytes32 secret, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowTransactor) Withdraw(opts *bind.TransactOpts, secret [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.contract.Transact(opts, "withdraw", secret, immutables)
}

// Withdraw is a paid mutator transaction binding the contract method 0x23305703.
//
// Solidity: function withdraw(bytes32 secret, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowSession) Withdraw(secret [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.Withdraw(&_IBaseEscrow.TransactOpts, secret, immutables)
}

// Withdraw is a paid mutator transaction binding the contract method 0x23305703.
//
// Solidity: function withdraw(bytes32 secret, (bytes32,bytes32,uint256,uint256,uint256,uint256,uint256,uint256) immutables) returns()
func (_IBaseEscrow *IBaseEscrowTransactorSession) Withdraw(secret [32]byte, immutables IBaseEscrowImmutables) (*types.Transaction, error) {
	return _IBaseEscrow.Contract.Withdraw(&_IBaseEscrow.TransactOpts, secret, immutables)
}

// IBaseEscrowEscrowCancelledIterator is returned from FilterEscrowCancelled and is used to iterate over the raw logs and unpacked data for EscrowCancelled events raised by the IBaseEscrow contract.
type IBaseEscrowEscrowCancelledIterator struct {
	Event *IBaseEscrowEscrowCancelled // Event containing the contract specifics and raw log

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
func (it *IBaseEscrowEscrowCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBaseEscrowEscrowCancelled)
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
		it.Event = new(IBaseEscrowEscrowCancelled)
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
func (it *IBaseEscrowEscrowCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBaseEscrowEscrowCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBaseEscrowEscrowCancelled represents a EscrowCancelled event raised by the IBaseEscrow contract.
type IBaseEscrowEscrowCancelled struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEscrowCancelled is a free log retrieval operation binding the contract event 0x6e3be9294e58d10b9c8053cfd5e09871b67e442fe394d6b0870d336b9df984a9.
//
// Solidity: event EscrowCancelled()
func (_IBaseEscrow *IBaseEscrowFilterer) FilterEscrowCancelled(opts *bind.FilterOpts) (*IBaseEscrowEscrowCancelledIterator, error) {

	logs, sub, err := _IBaseEscrow.contract.FilterLogs(opts, "EscrowCancelled")
	if err != nil {
		return nil, err
	}
	return &IBaseEscrowEscrowCancelledIterator{contract: _IBaseEscrow.contract, event: "EscrowCancelled", logs: logs, sub: sub}, nil
}

// WatchEscrowCancelled is a free log subscription operation binding the contract event 0x6e3be9294e58d10b9c8053cfd5e09871b67e442fe394d6b0870d336b9df984a9.
//
// Solidity: event EscrowCancelled()
func (_IBaseEscrow *IBaseEscrowFilterer) WatchEscrowCancelled(opts *bind.WatchOpts, sink chan<- *IBaseEscrowEscrowCancelled) (event.Subscription, error) {

	logs, sub, err := _IBaseEscrow.contract.WatchLogs(opts, "EscrowCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBaseEscrowEscrowCancelled)
				if err := _IBaseEscrow.contract.UnpackLog(event, "EscrowCancelled", log); err != nil {
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

// ParseEscrowCancelled is a log parse operation binding the contract event 0x6e3be9294e58d10b9c8053cfd5e09871b67e442fe394d6b0870d336b9df984a9.
//
// Solidity: event EscrowCancelled()
func (_IBaseEscrow *IBaseEscrowFilterer) ParseEscrowCancelled(log types.Log) (*IBaseEscrowEscrowCancelled, error) {
	event := new(IBaseEscrowEscrowCancelled)
	if err := _IBaseEscrow.contract.UnpackLog(event, "EscrowCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBaseEscrowFundsRescuedIterator is returned from FilterFundsRescued and is used to iterate over the raw logs and unpacked data for FundsRescued events raised by the IBaseEscrow contract.
type IBaseEscrowFundsRescuedIterator struct {
	Event *IBaseEscrowFundsRescued // Event containing the contract specifics and raw log

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
func (it *IBaseEscrowFundsRescuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBaseEscrowFundsRescued)
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
		it.Event = new(IBaseEscrowFundsRescued)
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
func (it *IBaseEscrowFundsRescuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBaseEscrowFundsRescuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBaseEscrowFundsRescued represents a FundsRescued event raised by the IBaseEscrow contract.
type IBaseEscrowFundsRescued struct {
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFundsRescued is a free log retrieval operation binding the contract event 0xc4474c2790e13695f6d2b6f1d8e164290b55370f87a542fd7711abe0a1bf40ac.
//
// Solidity: event FundsRescued(address token, uint256 amount)
func (_IBaseEscrow *IBaseEscrowFilterer) FilterFundsRescued(opts *bind.FilterOpts) (*IBaseEscrowFundsRescuedIterator, error) {

	logs, sub, err := _IBaseEscrow.contract.FilterLogs(opts, "FundsRescued")
	if err != nil {
		return nil, err
	}
	return &IBaseEscrowFundsRescuedIterator{contract: _IBaseEscrow.contract, event: "FundsRescued", logs: logs, sub: sub}, nil
}

// WatchFundsRescued is a free log subscription operation binding the contract event 0xc4474c2790e13695f6d2b6f1d8e164290b55370f87a542fd7711abe0a1bf40ac.
//
// Solidity: event FundsRescued(address token, uint256 amount)
func (_IBaseEscrow *IBaseEscrowFilterer) WatchFundsRescued(opts *bind.WatchOpts, sink chan<- *IBaseEscrowFundsRescued) (event.Subscription, error) {

	logs, sub, err := _IBaseEscrow.contract.WatchLogs(opts, "FundsRescued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBaseEscrowFundsRescued)
				if err := _IBaseEscrow.contract.UnpackLog(event, "FundsRescued", log); err != nil {
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

// ParseFundsRescued is a log parse operation binding the contract event 0xc4474c2790e13695f6d2b6f1d8e164290b55370f87a542fd7711abe0a1bf40ac.
//
// Solidity: event FundsRescued(address token, uint256 amount)
func (_IBaseEscrow *IBaseEscrowFilterer) ParseFundsRescued(log types.Log) (*IBaseEscrowFundsRescued, error) {
	event := new(IBaseEscrowFundsRescued)
	if err := _IBaseEscrow.contract.UnpackLog(event, "FundsRescued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IBaseEscrowWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the IBaseEscrow contract.
type IBaseEscrowWithdrawalIterator struct {
	Event *IBaseEscrowWithdrawal // Event containing the contract specifics and raw log

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
func (it *IBaseEscrowWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IBaseEscrowWithdrawal)
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
		it.Event = new(IBaseEscrowWithdrawal)
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
func (it *IBaseEscrowWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IBaseEscrowWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IBaseEscrowWithdrawal represents a Withdrawal event raised by the IBaseEscrow contract.
type IBaseEscrowWithdrawal struct {
	Secret [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x0ce781a18c10c8289803c7c4cfd532d797113c4b41c9701ffad7d0a632ac555b.
//
// Solidity: event Withdrawal(bytes32 secret)
func (_IBaseEscrow *IBaseEscrowFilterer) FilterWithdrawal(opts *bind.FilterOpts) (*IBaseEscrowWithdrawalIterator, error) {

	logs, sub, err := _IBaseEscrow.contract.FilterLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return &IBaseEscrowWithdrawalIterator{contract: _IBaseEscrow.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x0ce781a18c10c8289803c7c4cfd532d797113c4b41c9701ffad7d0a632ac555b.
//
// Solidity: event Withdrawal(bytes32 secret)
func (_IBaseEscrow *IBaseEscrowFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *IBaseEscrowWithdrawal) (event.Subscription, error) {

	logs, sub, err := _IBaseEscrow.contract.WatchLogs(opts, "Withdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IBaseEscrowWithdrawal)
				if err := _IBaseEscrow.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x0ce781a18c10c8289803c7c4cfd532d797113c4b41c9701ffad7d0a632ac555b.
//
// Solidity: event Withdrawal(bytes32 secret)
func (_IBaseEscrow *IBaseEscrowFilterer) ParseWithdrawal(log types.Log) (*IBaseEscrowWithdrawal, error) {
	event := new(IBaseEscrowWithdrawal)
	if err := _IBaseEscrow.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
