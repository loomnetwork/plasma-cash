package client

import (
	"crypto/ecdsa"

	loom "github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"

	ethcommon "github.com/ethereum/go-ethereum/common"
	amtypes "github.com/loomnetwork/go-loom/builtin/types/address_mapper"
	ssha "github.com/miguelmota/go-solidity-sha3"

	"github.com/loomnetwork/go-loom/common/evmcompat"
)

const (
	AddressMapperContractName = "addressmapper"
)

type AddressMapperClient struct {
	contract *client.Contract
}

func (a *AddressMapperClient) HasMapping(from, caller loom.Address) (*amtypes.AddressMapperHasMappingResponse, error) {
	addressMapperResponse := amtypes.AddressMapperHasMappingResponse{}

	_, err := a.contract.StaticCall("HasMapping", &amtypes.AddressMapperHasMappingRequest{
		From: from.MarshalPB(),
	}, caller, &addressMapperResponse)

	if err != nil {
		return nil, err
	}

	return &addressMapperResponse, nil
}

func (a *AddressMapperClient) AddIdentityMapping(from, to loom.Address, dappchainTxSigner auth.Signer, ethKey *ecdsa.PrivateKey) error {
	addressMappingSig, err := a.generateAddressMappingSignature(from, to, ethKey)
	if err != nil {
		return err
	}

	_, err = a.contract.Call("AddIdentityMapping", &amtypes.AddressMapperAddIdentityMappingRequest{
		From:      from.MarshalPB(),
		To:        to.MarshalPB(),
		Signature: addressMappingSig,
	}, dappchainTxSigner, nil)

	return err
}

func (a *AddressMapperClient) generateAddressMappingSignature(from, to loom.Address, key *ecdsa.PrivateKey) ([]byte, error) {
	hash := ssha.SoliditySHA3(
		ssha.Address(ethcommon.BytesToAddress(from.Local)),
		ssha.Address(ethcommon.BytesToAddress(to.Local)),
	)
	sig, err := evmcompat.SoliditySign(hash, key)
	if err != nil {
		return nil, err
	}
	// Prefix the sig with a single byte indicating the sig type, in this case EIP712
	return append(make([]byte, 1, 66), sig...), nil
}

func NewAddressMapperClient(chainID, writeUri, readUri string) (*AddressMapperClient, error) {
	rpcClient := client.NewDAppChainRPCClient(chainID, writeUri, readUri)

	contractAddr, err := rpcClient.Resolve(AddressMapperContractName)
	if err != nil {
		return nil, err
	}

	return &AddressMapperClient{contract: client.NewContract(rpcClient, contractAddr.Local)}, nil
}
