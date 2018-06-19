package client

import (
	"ethcontract"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TContract struct {
	Name          string
	tokenContract *ethcontract.Cards
}

func (d *TContract) Deposit(int) error {
	return nil
}

func (d *TContract) Register() error {
	_, err := d.tokenContract.Register(nil)
	return err
}

func (d *TContract) BalanceOf() (int, error) {
	return 0, nil
}

func (d *TContract) Account() (*Account, error) {
	return &Account{}, nil
}

var connToken *ethclient.Client

func InitTokenClient(connStr string) {
	var err error
	connToken, err = ethclient.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}

func GetTokenContract(name string) TokenContract {

	tokenContract, err := ethcontract.NewCards(common.HexToAddress("0x21e6fc92f93c8a1bb41e2be64b4e1f88a54d3576"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	return &TContract{name, tokenContract}
}
