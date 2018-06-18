// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethcontract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// RootChainABI is the input ABI used to generate the binding from.
const RootChainABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"lastParentBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"childBlockInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"defaultHashes\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"leaf\",\"type\":\"bytes32\"},{\"name\":\"index\",\"type\":\"uint64\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"getRoot\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"NUM_COINS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentChildBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentDepositBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"coins\",\"outputs\":[{\"name\":\"uid\",\"type\":\"uint64\"},{\"name\":\"denomination\",\"type\":\"uint32\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"exitSlots\",\"outputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"authority\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"exits\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"created_at\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"leaf\",\"type\":\"bytes32\"},{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"tokenID\",\"type\":\"uint64\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"checkMembership\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"childChain\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"created_at\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"slot\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"depositBlockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"denomination\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"slot\",\"type\":\"uint64\"},{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"created_at\",\"type\":\"uint256\"}],\"name\":\"StartedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"slot\",\"type\":\"uint64\"}],\"name\":\"ChallengedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"slot\",\"type\":\"uint64\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"FinalizedExit\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_cryptoCards\",\"type\":\"address\"}],\"name\":\"setCryptoCards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"submitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"slot\",\"type\":\"uint64\"},{\"name\":\"prevTxBytes\",\"type\":\"bytes\"},{\"name\":\"exitingTxBytes\",\"type\":\"bytes\"},{\"name\":\"prevTxInclusionProof\",\"type\":\"bytes\"},{\"name\":\"exitingTxInclusionProof\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"prevTxIncBlock\",\"type\":\"uint256\"},{\"name\":\"exitingTxIncBlock\",\"type\":\"uint256\"}],\"name\":\"startExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"getSig\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"slot\",\"type\":\"uint64\"},{\"name\":\"challengingTransaction\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"challengeBefore\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"slot\",\"type\":\"uint64\"},{\"name\":\"challengingTransaction\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"challengeBetween\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"slot\",\"type\":\"uint64\"},{\"name\":\"challengingTransaction\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"challengeAfter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"slot\",\"type\":\"uint64\"},{\"name\":\"challengingTransaction\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"respondChallengeBefore\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"slot\",\"type\":\"uint64\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getDepositBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_uid\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// RootChain is an auto generated Go binding around an Ethereum contract.
type RootChain struct {
	RootChainCaller     // Read-only binding to the contract
	RootChainTransactor // Write-only binding to the contract
	RootChainFilterer   // Log filterer for contract events
}

// RootChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootChainSession struct {
	Contract     *RootChain        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootChainCallerSession struct {
	Contract *RootChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RootChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootChainTransactorSession struct {
	Contract     *RootChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RootChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootChainRaw struct {
	Contract *RootChain // Generic contract binding to access the raw methods on
}

// RootChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootChainCallerRaw struct {
	Contract *RootChainCaller // Generic read-only contract binding to access the raw methods on
}

// RootChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootChainTransactorRaw struct {
	Contract *RootChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootChain creates a new instance of RootChain, bound to a specific deployed contract.
func NewRootChain(address common.Address, backend bind.ContractBackend) (*RootChain, error) {
	contract, err := bindRootChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RootChain{RootChainCaller: RootChainCaller{contract: contract}, RootChainTransactor: RootChainTransactor{contract: contract}, RootChainFilterer: RootChainFilterer{contract: contract}}, nil
}

// NewRootChainCaller creates a new read-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainCaller(address common.Address, caller bind.ContractCaller) (*RootChainCaller, error) {
	contract, err := bindRootChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainCaller{contract: contract}, nil
}

// NewRootChainTransactor creates a new write-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainTransactor(address common.Address, transactor bind.ContractTransactor) (*RootChainTransactor, error) {
	contract, err := bindRootChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainTransactor{contract: contract}, nil
}

// NewRootChainFilterer creates a new log filterer instance of RootChain, bound to a specific deployed contract.
func NewRootChainFilterer(address common.Address, filterer bind.ContractFilterer) (*RootChainFilterer, error) {
	contract, err := bindRootChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootChainFilterer{contract: contract}, nil
}

// bindRootChain binds a generic wrapper to an already deployed contract.
func bindRootChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RootChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.RootChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transact(opts, method, params...)
}

// NUMCOINS is a free data retrieval call binding the contract method 0x66e30d97.
//
// Solidity: function NUM_COINS() constant returns(uint64)
func (_RootChain *RootChainCaller) NUMCOINS(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "NUM_COINS")
	return *ret0, err
}

// NUMCOINS is a free data retrieval call binding the contract method 0x66e30d97.
//
// Solidity: function NUM_COINS() constant returns(uint64)
func (_RootChain *RootChainSession) NUMCOINS() (uint64, error) {
	return _RootChain.Contract.NUMCOINS(&_RootChain.CallOpts)
}

// NUMCOINS is a free data retrieval call binding the contract method 0x66e30d97.
//
// Solidity: function NUM_COINS() constant returns(uint64)
func (_RootChain *RootChainCallerSession) NUMCOINS() (uint64, error) {
	return _RootChain.Contract.NUMCOINS(&_RootChain.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_RootChain *RootChainCaller) Authority(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "authority")
	return *ret0, err
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_RootChain *RootChainSession) Authority() (common.Address, error) {
	return _RootChain.Contract.Authority(&_RootChain.CallOpts)
}

// Authority is a free data retrieval call binding the contract method 0xbf7e214f.
//
// Solidity: function authority() constant returns(address)
func (_RootChain *RootChainCallerSession) Authority() (common.Address, error) {
	return _RootChain.Contract.Authority(&_RootChain.CallOpts)
}

// CheckMembership is a free data retrieval call binding the contract method 0xf586df65.
//
// Solidity: function checkMembership(leaf bytes32, root bytes32, tokenID uint64, proof bytes) constant returns(bool)
func (_RootChain *RootChainCaller) CheckMembership(opts *bind.CallOpts, leaf [32]byte, root [32]byte, tokenID uint64, proof []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "checkMembership", leaf, root, tokenID, proof)
	return *ret0, err
}

// CheckMembership is a free data retrieval call binding the contract method 0xf586df65.
//
// Solidity: function checkMembership(leaf bytes32, root bytes32, tokenID uint64, proof bytes) constant returns(bool)
func (_RootChain *RootChainSession) CheckMembership(leaf [32]byte, root [32]byte, tokenID uint64, proof []byte) (bool, error) {
	return _RootChain.Contract.CheckMembership(&_RootChain.CallOpts, leaf, root, tokenID, proof)
}

// CheckMembership is a free data retrieval call binding the contract method 0xf586df65.
//
// Solidity: function checkMembership(leaf bytes32, root bytes32, tokenID uint64, proof bytes) constant returns(bool)
func (_RootChain *RootChainCallerSession) CheckMembership(leaf [32]byte, root [32]byte, tokenID uint64, proof []byte) (bool, error) {
	return _RootChain.Contract.CheckMembership(&_RootChain.CallOpts, leaf, root, tokenID, proof)
}

// ChildBlockInterval is a free data retrieval call binding the contract method 0x38a9e0bc.
//
// Solidity: function childBlockInterval() constant returns(uint256)
func (_RootChain *RootChainCaller) ChildBlockInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "childBlockInterval")
	return *ret0, err
}

// ChildBlockInterval is a free data retrieval call binding the contract method 0x38a9e0bc.
//
// Solidity: function childBlockInterval() constant returns(uint256)
func (_RootChain *RootChainSession) ChildBlockInterval() (*big.Int, error) {
	return _RootChain.Contract.ChildBlockInterval(&_RootChain.CallOpts)
}

// ChildBlockInterval is a free data retrieval call binding the contract method 0x38a9e0bc.
//
// Solidity: function childBlockInterval() constant returns(uint256)
func (_RootChain *RootChainCallerSession) ChildBlockInterval() (*big.Int, error) {
	return _RootChain.Contract.ChildBlockInterval(&_RootChain.CallOpts)
}

// ChildChain is a free data retrieval call binding the contract method 0xf95643b1.
//
// Solidity: function childChain( uint256) constant returns(root bytes32, created_at uint256)
func (_RootChain *RootChainCaller) ChildChain(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Root      [32]byte
	CreatedAt *big.Int
}, error) {
	ret := new(struct {
		Root      [32]byte
		CreatedAt *big.Int
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "childChain", arg0)
	return *ret, err
}

// ChildChain is a free data retrieval call binding the contract method 0xf95643b1.
//
// Solidity: function childChain( uint256) constant returns(root bytes32, created_at uint256)
func (_RootChain *RootChainSession) ChildChain(arg0 *big.Int) (struct {
	Root      [32]byte
	CreatedAt *big.Int
}, error) {
	return _RootChain.Contract.ChildChain(&_RootChain.CallOpts, arg0)
}

// ChildChain is a free data retrieval call binding the contract method 0xf95643b1.
//
// Solidity: function childChain( uint256) constant returns(root bytes32, created_at uint256)
func (_RootChain *RootChainCallerSession) ChildChain(arg0 *big.Int) (struct {
	Root      [32]byte
	CreatedAt *big.Int
}, error) {
	return _RootChain.Contract.ChildChain(&_RootChain.CallOpts, arg0)
}

// Coins is a free data retrieval call binding the contract method 0xb8bdaf62.
//
// Solidity: function coins( uint64) constant returns(uid uint64, denomination uint32, owner address, state uint8)
func (_RootChain *RootChainCaller) Coins(opts *bind.CallOpts, arg0 uint64) (struct {
	Uid          uint64
	Denomination uint32
	Owner        common.Address
	State        uint8
}, error) {
	ret := new(struct {
		Uid          uint64
		Denomination uint32
		Owner        common.Address
		State        uint8
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "coins", arg0)
	return *ret, err
}

// Coins is a free data retrieval call binding the contract method 0xb8bdaf62.
//
// Solidity: function coins( uint64) constant returns(uid uint64, denomination uint32, owner address, state uint8)
func (_RootChain *RootChainSession) Coins(arg0 uint64) (struct {
	Uid          uint64
	Denomination uint32
	Owner        common.Address
	State        uint8
}, error) {
	return _RootChain.Contract.Coins(&_RootChain.CallOpts, arg0)
}

// Coins is a free data retrieval call binding the contract method 0xb8bdaf62.
//
// Solidity: function coins( uint64) constant returns(uid uint64, denomination uint32, owner address, state uint8)
func (_RootChain *RootChainCallerSession) Coins(arg0 uint64) (struct {
	Uid          uint64
	Denomination uint32
	Owner        common.Address
	State        uint8
}, error) {
	return _RootChain.Contract.Coins(&_RootChain.CallOpts, arg0)
}

// CurrentChildBlock is a free data retrieval call binding the contract method 0x7a95f1e8.
//
// Solidity: function currentChildBlock() constant returns(uint256)
func (_RootChain *RootChainCaller) CurrentChildBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "currentChildBlock")
	return *ret0, err
}

// CurrentChildBlock is a free data retrieval call binding the contract method 0x7a95f1e8.
//
// Solidity: function currentChildBlock() constant returns(uint256)
func (_RootChain *RootChainSession) CurrentChildBlock() (*big.Int, error) {
	return _RootChain.Contract.CurrentChildBlock(&_RootChain.CallOpts)
}

// CurrentChildBlock is a free data retrieval call binding the contract method 0x7a95f1e8.
//
// Solidity: function currentChildBlock() constant returns(uint256)
func (_RootChain *RootChainCallerSession) CurrentChildBlock() (*big.Int, error) {
	return _RootChain.Contract.CurrentChildBlock(&_RootChain.CallOpts)
}

// CurrentDepositBlock is a free data retrieval call binding the contract method 0xa98c7f2c.
//
// Solidity: function currentDepositBlock() constant returns(uint256)
func (_RootChain *RootChainCaller) CurrentDepositBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "currentDepositBlock")
	return *ret0, err
}

// CurrentDepositBlock is a free data retrieval call binding the contract method 0xa98c7f2c.
//
// Solidity: function currentDepositBlock() constant returns(uint256)
func (_RootChain *RootChainSession) CurrentDepositBlock() (*big.Int, error) {
	return _RootChain.Contract.CurrentDepositBlock(&_RootChain.CallOpts)
}

// CurrentDepositBlock is a free data retrieval call binding the contract method 0xa98c7f2c.
//
// Solidity: function currentDepositBlock() constant returns(uint256)
func (_RootChain *RootChainCallerSession) CurrentDepositBlock() (*big.Int, error) {
	return _RootChain.Contract.CurrentDepositBlock(&_RootChain.CallOpts)
}

// DefaultHashes is a free data retrieval call binding the contract method 0x48419ad8.
//
// Solidity: function defaultHashes( uint256) constant returns(bytes32)
func (_RootChain *RootChainCaller) DefaultHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "defaultHashes", arg0)
	return *ret0, err
}

// DefaultHashes is a free data retrieval call binding the contract method 0x48419ad8.
//
// Solidity: function defaultHashes( uint256) constant returns(bytes32)
func (_RootChain *RootChainSession) DefaultHashes(arg0 *big.Int) ([32]byte, error) {
	return _RootChain.Contract.DefaultHashes(&_RootChain.CallOpts, arg0)
}

// DefaultHashes is a free data retrieval call binding the contract method 0x48419ad8.
//
// Solidity: function defaultHashes( uint256) constant returns(bytes32)
func (_RootChain *RootChainCallerSession) DefaultHashes(arg0 *big.Int) ([32]byte, error) {
	return _RootChain.Contract.DefaultHashes(&_RootChain.CallOpts, arg0)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_RootChain *RootChainCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "depositCount")
	return *ret0, err
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_RootChain *RootChainSession) DepositCount() (*big.Int, error) {
	return _RootChain.Contract.DepositCount(&_RootChain.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_RootChain *RootChainCallerSession) DepositCount() (*big.Int, error) {
	return _RootChain.Contract.DepositCount(&_RootChain.CallOpts)
}

// ExitSlots is a free data retrieval call binding the contract method 0xbcd5df39.
//
// Solidity: function exitSlots( uint256) constant returns(uint64)
func (_RootChain *RootChainCaller) ExitSlots(opts *bind.CallOpts, arg0 *big.Int) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "exitSlots", arg0)
	return *ret0, err
}

// ExitSlots is a free data retrieval call binding the contract method 0xbcd5df39.
//
// Solidity: function exitSlots( uint256) constant returns(uint64)
func (_RootChain *RootChainSession) ExitSlots(arg0 *big.Int) (uint64, error) {
	return _RootChain.Contract.ExitSlots(&_RootChain.CallOpts, arg0)
}

// ExitSlots is a free data retrieval call binding the contract method 0xbcd5df39.
//
// Solidity: function exitSlots( uint256) constant returns(uint64)
func (_RootChain *RootChainCallerSession) ExitSlots(arg0 *big.Int) (uint64, error) {
	return _RootChain.Contract.ExitSlots(&_RootChain.CallOpts, arg0)
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(owner address, created_at uint256)
func (_RootChain *RootChainCaller) Exits(opts *bind.CallOpts, arg0 uint64) (struct {
	Owner     common.Address
	CreatedAt *big.Int
}, error) {
	ret := new(struct {
		Owner     common.Address
		CreatedAt *big.Int
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "exits", arg0)
	return *ret, err
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(owner address, created_at uint256)
func (_RootChain *RootChainSession) Exits(arg0 uint64) (struct {
	Owner     common.Address
	CreatedAt *big.Int
}, error) {
	return _RootChain.Contract.Exits(&_RootChain.CallOpts, arg0)
}

// Exits is a free data retrieval call binding the contract method 0xd6463d40.
//
// Solidity: function exits( uint64) constant returns(owner address, created_at uint256)
func (_RootChain *RootChainCallerSession) Exits(arg0 uint64) (struct {
	Owner     common.Address
	CreatedAt *big.Int
}, error) {
	return _RootChain.Contract.Exits(&_RootChain.CallOpts, arg0)
}

// GetDepositBlock is a free data retrieval call binding the contract method 0xbcd59261.
//
// Solidity: function getDepositBlock() constant returns(uint256)
func (_RootChain *RootChainCaller) GetDepositBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "getDepositBlock")
	return *ret0, err
}

// GetDepositBlock is a free data retrieval call binding the contract method 0xbcd59261.
//
// Solidity: function getDepositBlock() constant returns(uint256)
func (_RootChain *RootChainSession) GetDepositBlock() (*big.Int, error) {
	return _RootChain.Contract.GetDepositBlock(&_RootChain.CallOpts)
}

// GetDepositBlock is a free data retrieval call binding the contract method 0xbcd59261.
//
// Solidity: function getDepositBlock() constant returns(uint256)
func (_RootChain *RootChainCallerSession) GetDepositBlock() (*big.Int, error) {
	return _RootChain.Contract.GetDepositBlock(&_RootChain.CallOpts)
}

// GetRoot is a free data retrieval call binding the contract method 0x509a7e54.
//
// Solidity: function getRoot(leaf bytes32, index uint64, proof bytes) constant returns(bytes32)
func (_RootChain *RootChainCaller) GetRoot(opts *bind.CallOpts, leaf [32]byte, index uint64, proof []byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "getRoot", leaf, index, proof)
	return *ret0, err
}

// GetRoot is a free data retrieval call binding the contract method 0x509a7e54.
//
// Solidity: function getRoot(leaf bytes32, index uint64, proof bytes) constant returns(bytes32)
func (_RootChain *RootChainSession) GetRoot(leaf [32]byte, index uint64, proof []byte) ([32]byte, error) {
	return _RootChain.Contract.GetRoot(&_RootChain.CallOpts, leaf, index, proof)
}

// GetRoot is a free data retrieval call binding the contract method 0x509a7e54.
//
// Solidity: function getRoot(leaf bytes32, index uint64, proof bytes) constant returns(bytes32)
func (_RootChain *RootChainCallerSession) GetRoot(leaf [32]byte, index uint64, proof []byte) ([32]byte, error) {
	return _RootChain.Contract.GetRoot(&_RootChain.CallOpts, leaf, index, proof)
}

// GetSig is a free data retrieval call binding the contract method 0x245cab33.
//
// Solidity: function getSig(sigs bytes, i uint256) constant returns(bytes)
func (_RootChain *RootChainCaller) GetSig(opts *bind.CallOpts, sigs []byte, i *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "getSig", sigs, i)
	return *ret0, err
}

// GetSig is a free data retrieval call binding the contract method 0x245cab33.
//
// Solidity: function getSig(sigs bytes, i uint256) constant returns(bytes)
func (_RootChain *RootChainSession) GetSig(sigs []byte, i *big.Int) ([]byte, error) {
	return _RootChain.Contract.GetSig(&_RootChain.CallOpts, sigs, i)
}

// GetSig is a free data retrieval call binding the contract method 0x245cab33.
//
// Solidity: function getSig(sigs bytes, i uint256) constant returns(bytes)
func (_RootChain *RootChainCallerSession) GetSig(sigs []byte, i *big.Int) ([]byte, error) {
	return _RootChain.Contract.GetSig(&_RootChain.CallOpts, sigs, i)
}

// LastParentBlock is a free data retrieval call binding the contract method 0x117546c5.
//
// Solidity: function lastParentBlock() constant returns(uint256)
func (_RootChain *RootChainCaller) LastParentBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "lastParentBlock")
	return *ret0, err
}

// LastParentBlock is a free data retrieval call binding the contract method 0x117546c5.
//
// Solidity: function lastParentBlock() constant returns(uint256)
func (_RootChain *RootChainSession) LastParentBlock() (*big.Int, error) {
	return _RootChain.Contract.LastParentBlock(&_RootChain.CallOpts)
}

// LastParentBlock is a free data retrieval call binding the contract method 0x117546c5.
//
// Solidity: function lastParentBlock() constant returns(uint256)
func (_RootChain *RootChainCallerSession) LastParentBlock() (*big.Int, error) {
	return _RootChain.Contract.LastParentBlock(&_RootChain.CallOpts)
}

// ChallengeAfter is a paid mutator transaction binding the contract method 0xa888683f.
//
// Solidity: function challengeAfter(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactor) ChallengeAfter(opts *bind.TransactOpts, slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeAfter", slot, challengingTransaction, proof)
}

// ChallengeAfter is a paid mutator transaction binding the contract method 0xa888683f.
//
// Solidity: function challengeAfter(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainSession) ChallengeAfter(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeAfter(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// ChallengeAfter is a paid mutator transaction binding the contract method 0xa888683f.
//
// Solidity: function challengeAfter(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactorSession) ChallengeAfter(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeAfter(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// ChallengeBefore is a paid mutator transaction binding the contract method 0x2c9ecf0a.
//
// Solidity: function challengeBefore(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactor) ChallengeBefore(opts *bind.TransactOpts, slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeBefore", slot, challengingTransaction, proof)
}

// ChallengeBefore is a paid mutator transaction binding the contract method 0x2c9ecf0a.
//
// Solidity: function challengeBefore(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainSession) ChallengeBefore(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeBefore(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// ChallengeBefore is a paid mutator transaction binding the contract method 0x2c9ecf0a.
//
// Solidity: function challengeBefore(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactorSession) ChallengeBefore(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeBefore(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// ChallengeBetween is a paid mutator transaction binding the contract method 0x759b4929.
//
// Solidity: function challengeBetween(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactor) ChallengeBetween(opts *bind.TransactOpts, slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeBetween", slot, challengingTransaction, proof)
}

// ChallengeBetween is a paid mutator transaction binding the contract method 0x759b4929.
//
// Solidity: function challengeBetween(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainSession) ChallengeBetween(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeBetween(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// ChallengeBetween is a paid mutator transaction binding the contract method 0x759b4929.
//
// Solidity: function challengeBetween(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactorSession) ChallengeBetween(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeBetween(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_RootChain *RootChainTransactor) FinalizeExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "finalizeExits")
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_RootChain *RootChainSession) FinalizeExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeExits(&_RootChain.TransactOpts)
}

// FinalizeExits is a paid mutator transaction binding the contract method 0xc6ab44cd.
//
// Solidity: function finalizeExits() returns()
func (_RootChain *RootChainTransactorSession) FinalizeExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeExits(&_RootChain.TransactOpts)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xf0b9e5ba.
//
// Solidity: function onERC721Received(_from address, _uid uint256, _data bytes) returns(bytes4)
func (_RootChain *RootChainTransactor) OnERC721Received(opts *bind.TransactOpts, _from common.Address, _uid *big.Int, _data []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "onERC721Received", _from, _uid, _data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xf0b9e5ba.
//
// Solidity: function onERC721Received(_from address, _uid uint256, _data bytes) returns(bytes4)
func (_RootChain *RootChainSession) OnERC721Received(_from common.Address, _uid *big.Int, _data []byte) (*types.Transaction, error) {
	return _RootChain.Contract.OnERC721Received(&_RootChain.TransactOpts, _from, _uid, _data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0xf0b9e5ba.
//
// Solidity: function onERC721Received(_from address, _uid uint256, _data bytes) returns(bytes4)
func (_RootChain *RootChainTransactorSession) OnERC721Received(_from common.Address, _uid *big.Int, _data []byte) (*types.Transaction, error) {
	return _RootChain.Contract.OnERC721Received(&_RootChain.TransactOpts, _from, _uid, _data)
}

// RespondChallengeBefore is a paid mutator transaction binding the contract method 0x813b0b86.
//
// Solidity: function respondChallengeBefore(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactor) RespondChallengeBefore(opts *bind.TransactOpts, slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "respondChallengeBefore", slot, challengingTransaction, proof)
}

// RespondChallengeBefore is a paid mutator transaction binding the contract method 0x813b0b86.
//
// Solidity: function respondChallengeBefore(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainSession) RespondChallengeBefore(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.RespondChallengeBefore(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// RespondChallengeBefore is a paid mutator transaction binding the contract method 0x813b0b86.
//
// Solidity: function respondChallengeBefore(slot uint64, challengingTransaction bytes, proof bytes) returns()
func (_RootChain *RootChainTransactorSession) RespondChallengeBefore(slot uint64, challengingTransaction []byte, proof []byte) (*types.Transaction, error) {
	return _RootChain.Contract.RespondChallengeBefore(&_RootChain.TransactOpts, slot, challengingTransaction, proof)
}

// SetCryptoCards is a paid mutator transaction binding the contract method 0x329be5d5.
//
// Solidity: function setCryptoCards(_cryptoCards address) returns()
func (_RootChain *RootChainTransactor) SetCryptoCards(opts *bind.TransactOpts, _cryptoCards common.Address) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "setCryptoCards", _cryptoCards)
}

// SetCryptoCards is a paid mutator transaction binding the contract method 0x329be5d5.
//
// Solidity: function setCryptoCards(_cryptoCards address) returns()
func (_RootChain *RootChainSession) SetCryptoCards(_cryptoCards common.Address) (*types.Transaction, error) {
	return _RootChain.Contract.SetCryptoCards(&_RootChain.TransactOpts, _cryptoCards)
}

// SetCryptoCards is a paid mutator transaction binding the contract method 0x329be5d5.
//
// Solidity: function setCryptoCards(_cryptoCards address) returns()
func (_RootChain *RootChainTransactorSession) SetCryptoCards(_cryptoCards common.Address) (*types.Transaction, error) {
	return _RootChain.Contract.SetCryptoCards(&_RootChain.TransactOpts, _cryptoCards)
}

// StartExit is a paid mutator transaction binding the contract method 0x5b6a6efb.
//
// Solidity: function startExit(slot uint64, prevTxBytes bytes, exitingTxBytes bytes, prevTxInclusionProof bytes, exitingTxInclusionProof bytes, sigs bytes, prevTxIncBlock uint256, exitingTxIncBlock uint256) returns()
func (_RootChain *RootChainTransactor) StartExit(opts *bind.TransactOpts, slot uint64, prevTxBytes []byte, exitingTxBytes []byte, prevTxInclusionProof []byte, exitingTxInclusionProof []byte, sigs []byte, prevTxIncBlock *big.Int, exitingTxIncBlock *big.Int) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "startExit", slot, prevTxBytes, exitingTxBytes, prevTxInclusionProof, exitingTxInclusionProof, sigs, prevTxIncBlock, exitingTxIncBlock)
}

// StartExit is a paid mutator transaction binding the contract method 0x5b6a6efb.
//
// Solidity: function startExit(slot uint64, prevTxBytes bytes, exitingTxBytes bytes, prevTxInclusionProof bytes, exitingTxInclusionProof bytes, sigs bytes, prevTxIncBlock uint256, exitingTxIncBlock uint256) returns()
func (_RootChain *RootChainSession) StartExit(slot uint64, prevTxBytes []byte, exitingTxBytes []byte, prevTxInclusionProof []byte, exitingTxInclusionProof []byte, sigs []byte, prevTxIncBlock *big.Int, exitingTxIncBlock *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.StartExit(&_RootChain.TransactOpts, slot, prevTxBytes, exitingTxBytes, prevTxInclusionProof, exitingTxInclusionProof, sigs, prevTxIncBlock, exitingTxIncBlock)
}

// StartExit is a paid mutator transaction binding the contract method 0x5b6a6efb.
//
// Solidity: function startExit(slot uint64, prevTxBytes bytes, exitingTxBytes bytes, prevTxInclusionProof bytes, exitingTxInclusionProof bytes, sigs bytes, prevTxIncBlock uint256, exitingTxIncBlock uint256) returns()
func (_RootChain *RootChainTransactorSession) StartExit(slot uint64, prevTxBytes []byte, exitingTxBytes []byte, prevTxInclusionProof []byte, exitingTxInclusionProof []byte, sigs []byte, prevTxIncBlock *big.Int, exitingTxIncBlock *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.StartExit(&_RootChain.TransactOpts, slot, prevTxBytes, exitingTxBytes, prevTxInclusionProof, exitingTxInclusionProof, sigs, prevTxIncBlock, exitingTxIncBlock)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xbaa47694.
//
// Solidity: function submitBlock(root bytes32) returns()
func (_RootChain *RootChainTransactor) SubmitBlock(opts *bind.TransactOpts, root [32]byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "submitBlock", root)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xbaa47694.
//
// Solidity: function submitBlock(root bytes32) returns()
func (_RootChain *RootChainSession) SubmitBlock(root [32]byte) (*types.Transaction, error) {
	return _RootChain.Contract.SubmitBlock(&_RootChain.TransactOpts, root)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xbaa47694.
//
// Solidity: function submitBlock(root bytes32) returns()
func (_RootChain *RootChainTransactorSession) SubmitBlock(root [32]byte) (*types.Transaction, error) {
	return _RootChain.Contract.SubmitBlock(&_RootChain.TransactOpts, root)
}

// Withdraw is a paid mutator transaction binding the contract method 0x750f0acc.
//
// Solidity: function withdraw(slot uint64) returns()
func (_RootChain *RootChainTransactor) Withdraw(opts *bind.TransactOpts, slot uint64) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "withdraw", slot)
}

// Withdraw is a paid mutator transaction binding the contract method 0x750f0acc.
//
// Solidity: function withdraw(slot uint64) returns()
func (_RootChain *RootChainSession) Withdraw(slot uint64) (*types.Transaction, error) {
	return _RootChain.Contract.Withdraw(&_RootChain.TransactOpts, slot)
}

// Withdraw is a paid mutator transaction binding the contract method 0x750f0acc.
//
// Solidity: function withdraw(slot uint64) returns()
func (_RootChain *RootChainTransactorSession) Withdraw(slot uint64) (*types.Transaction, error) {
	return _RootChain.Contract.Withdraw(&_RootChain.TransactOpts, slot)
}

// RootChainChallengedExitIterator is returned from FilterChallengedExit and is used to iterate over the raw logs and unpacked data for ChallengedExit events raised by the RootChain contract.
type RootChainChallengedExitIterator struct {
	Event *RootChainChallengedExit // Event containing the contract specifics and raw log

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
func (it *RootChainChallengedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainChallengedExit)
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
		it.Event = new(RootChainChallengedExit)
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
func (it *RootChainChallengedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainChallengedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainChallengedExit represents a ChallengedExit event raised by the RootChain contract.
type RootChainChallengedExit struct {
	Slot uint64
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterChallengedExit is a free log retrieval operation binding the contract event 0x7866a4e258966d935dedb7b3a912110c549aba222af732ac3a2efb5609cb9792.
//
// Solidity: e ChallengedExit(slot indexed uint64)
func (_RootChain *RootChainFilterer) FilterChallengedExit(opts *bind.FilterOpts, slot []uint64) (*RootChainChallengedExitIterator, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "ChallengedExit", slotRule)
	if err != nil {
		return nil, err
	}
	return &RootChainChallengedExitIterator{contract: _RootChain.contract, event: "ChallengedExit", logs: logs, sub: sub}, nil
}

// WatchChallengedExit is a free log subscription operation binding the contract event 0x7866a4e258966d935dedb7b3a912110c549aba222af732ac3a2efb5609cb9792.
//
// Solidity: e ChallengedExit(slot indexed uint64)
func (_RootChain *RootChainFilterer) WatchChallengedExit(opts *bind.WatchOpts, sink chan<- *RootChainChallengedExit, slot []uint64) (event.Subscription, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "ChallengedExit", slotRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainChallengedExit)
				if err := _RootChain.contract.UnpackLog(event, "ChallengedExit", log); err != nil {
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

// RootChainDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the RootChain contract.
type RootChainDepositIterator struct {
	Event *RootChainDeposit // Event containing the contract specifics and raw log

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
func (it *RootChainDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainDeposit)
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
		it.Event = new(RootChainDeposit)
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
func (it *RootChainDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainDeposit represents a Deposit event raised by the RootChain contract.
type RootChainDeposit struct {
	Slot               uint64
	DepositBlockNumber *big.Int
	Denomination       uint64
	From               common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x51f61c9bfc0c0871dbb1aa0ecfc2166a0b5f6e158a489d17454f1ef618ba5eea.
//
// Solidity: e Deposit(slot indexed uint64, depositBlockNumber uint256, denomination uint64, from indexed address)
func (_RootChain *RootChainFilterer) FilterDeposit(opts *bind.FilterOpts, slot []uint64, from []common.Address) (*RootChainDepositIterator, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "Deposit", slotRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &RootChainDepositIterator{contract: _RootChain.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x51f61c9bfc0c0871dbb1aa0ecfc2166a0b5f6e158a489d17454f1ef618ba5eea.
//
// Solidity: e Deposit(slot indexed uint64, depositBlockNumber uint256, denomination uint64, from indexed address)
func (_RootChain *RootChainFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RootChainDeposit, slot []uint64, from []common.Address) (event.Subscription, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "Deposit", slotRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainDeposit)
				if err := _RootChain.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// RootChainFinalizedExitIterator is returned from FilterFinalizedExit and is used to iterate over the raw logs and unpacked data for FinalizedExit events raised by the RootChain contract.
type RootChainFinalizedExitIterator struct {
	Event *RootChainFinalizedExit // Event containing the contract specifics and raw log

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
func (it *RootChainFinalizedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainFinalizedExit)
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
		it.Event = new(RootChainFinalizedExit)
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
func (it *RootChainFinalizedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainFinalizedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainFinalizedExit represents a FinalizedExit event raised by the RootChain contract.
type RootChainFinalizedExit struct {
	Slot  uint64
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFinalizedExit is a free log retrieval operation binding the contract event 0x432647a5fb9bdea356d78f8e3d83b6ddc2e78b4e4a702ac7eb968f7fe03d9dda.
//
// Solidity: e FinalizedExit(slot indexed uint64, owner address)
func (_RootChain *RootChainFilterer) FilterFinalizedExit(opts *bind.FilterOpts, slot []uint64) (*RootChainFinalizedExitIterator, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "FinalizedExit", slotRule)
	if err != nil {
		return nil, err
	}
	return &RootChainFinalizedExitIterator{contract: _RootChain.contract, event: "FinalizedExit", logs: logs, sub: sub}, nil
}

// WatchFinalizedExit is a free log subscription operation binding the contract event 0x432647a5fb9bdea356d78f8e3d83b6ddc2e78b4e4a702ac7eb968f7fe03d9dda.
//
// Solidity: e FinalizedExit(slot indexed uint64, owner address)
func (_RootChain *RootChainFilterer) WatchFinalizedExit(opts *bind.WatchOpts, sink chan<- *RootChainFinalizedExit, slot []uint64) (event.Subscription, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "FinalizedExit", slotRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainFinalizedExit)
				if err := _RootChain.contract.UnpackLog(event, "FinalizedExit", log); err != nil {
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

// RootChainStartedExitIterator is returned from FilterStartedExit and is used to iterate over the raw logs and unpacked data for StartedExit events raised by the RootChain contract.
type RootChainStartedExitIterator struct {
	Event *RootChainStartedExit // Event containing the contract specifics and raw log

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
func (it *RootChainStartedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainStartedExit)
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
		it.Event = new(RootChainStartedExit)
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
func (it *RootChainStartedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainStartedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainStartedExit represents a StartedExit event raised by the RootChain contract.
type RootChainStartedExit struct {
	Slot      uint64
	Owner     common.Address
	CreatedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStartedExit is a free log retrieval operation binding the contract event 0xb2451f3e7d6dff5753a128f9a3982f46e0c287d3d3f71c87f13cc63a6081283d.
//
// Solidity: e StartedExit(slot indexed uint64, owner indexed address, created_at uint256)
func (_RootChain *RootChainFilterer) FilterStartedExit(opts *bind.FilterOpts, slot []uint64, owner []common.Address) (*RootChainStartedExitIterator, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "StartedExit", slotRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &RootChainStartedExitIterator{contract: _RootChain.contract, event: "StartedExit", logs: logs, sub: sub}, nil
}

// WatchStartedExit is a free log subscription operation binding the contract event 0xb2451f3e7d6dff5753a128f9a3982f46e0c287d3d3f71c87f13cc63a6081283d.
//
// Solidity: e StartedExit(slot indexed uint64, owner indexed address, created_at uint256)
func (_RootChain *RootChainFilterer) WatchStartedExit(opts *bind.WatchOpts, sink chan<- *RootChainStartedExit, slot []uint64, owner []common.Address) (event.Subscription, error) {

	var slotRule []interface{}
	for _, slotItem := range slot {
		slotRule = append(slotRule, slotItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "StartedExit", slotRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainStartedExit)
				if err := _RootChain.contract.UnpackLog(event, "StartedExit", log); err != nil {
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
