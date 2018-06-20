package client

import (
	"ethcontract"

	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RootChainService struct {
	Name              string
	rootchainContract *ethcontract.RootChain
}

func (d *RootChainService) PlasmaCoin(uint64) {
}

func (d *RootChainService) Withdraw(uint64) {
}

func (d *RootChainService) FinalizeExits() {
}
func (d *RootChainService) WithdrawBonds() {
}

var conn *ethclient.Client

func InitClients(connStr string) {
	var err error
	conn, err = ethclient.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}

func GetRootChain(name string) RootChainClient {

	// Instantiate the contract and display its name
	rootchainContract, err := ethcontract.NewRootChain(common.HexToAddress("0x21e6fc92f93c8a1bb41e2be64b4e1f88a54d3576"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	/*
		name, err := token.Name(nil)
		if err != nil {
			log.Fatalf("Failed to retrieve token name: %v", err)
		}
		fmt.Println("Token name:", name)
	*/
	// (plasma_config[key], self.root_chain_abi,
	//	plasma_config['root_chain'], self.endpoint)

	return &RootChainService{name, rootchainContract}
}
