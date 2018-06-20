package client

import (
	"encoding/base64"
	"errors"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/cli"
	"github.com/loomnetwork/go-loom/client"
)

type LoomContract struct {
	WriteURI     string
	ReadURI      string
	ContractAddr string
	ChainID      string
	PrivFile     string
}

func NewLoomContract(readAddr string, writeAddr string, contractAddr string) *LoomContract {
	return &LoomContract{
		WriteURI:     readAddr,
		ReadURI:      writeAddr,
		ContractAddr: contractAddr,
		ChainID:      "",
		PrivFile:     "test.key",
	}
}

func (l LoomContract) ResolveAddress(s string) (loom.Address, error) {
	rpcClient := client.NewDAppChainRPCClient(l.ChainID, l.WriteURI, l.ReadURI)
	contractAddr, err := cli.ParseAddress(s)
	if err != nil {
		// if address invalid, try to resolve it using registry
		contractAddr, err = rpcClient.Resolve(s)
		if err != nil {
			return loom.Address{}, err
		}
	}

	return contractAddr, nil
}

func (l LoomContract) contract() (*client.Contract, error) {
	if l.ContractAddr == "" {
		return nil, errors.New("contract address or name required")
	}

	contractAddr, err := l.ResolveAddress(l.ContractAddr)
	if err != nil {
		return nil, err
	}

	// create rpc client
	rpcClient := client.NewDAppChainRPCClient(l.ChainID, l.WriteURI, l.ReadURI)
	// create contract
	contract := client.NewContract(rpcClient, contractAddr.Local)
	return contract, nil
}

func (l LoomContract) CallContract(method string, params proto.Message, result interface{}) error {
	if l.PrivFile == "" {
		return errors.New("private key required to call contract")
	}

	privKeyB64, err := ioutil.ReadFile(l.PrivFile)
	if err != nil {
		return err
	}

	privKey, err := base64.StdEncoding.DecodeString(string(privKeyB64))
	if err != nil {
		return err
	}

	signer := auth.NewEd25519Signer(privKey)

	contract, err := l.contract()
	if err != nil {
		return err
	}
	_, err = contract.Call(method, params, signer, result)
	return err
}

func (l LoomContract) StaticCallContract(method string, params proto.Message, result interface{}) error {
	contract, err := l.contract()
	if err != nil {
		return err
	}

	_, err = contract.StaticCall(method, params, loom.RootAddress(l.ChainID), result)
	return err
}
