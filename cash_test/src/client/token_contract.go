package client

import (
	"crypto/ecdsa"
	"ethcontract"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/ethclient"
)

type TContract struct {
	Name          string
	tokenContract *ethcontract.Cards
	callerKey     *ecdsa.PrivateKey
	callerAddr    common.Address
}

func (d *TContract) Deposit(int) error {
	return nil
}

func (d *TContract) Register() error {
	auth := bind.NewKeyedTransactor(d.callerKey)
	_, err := d.tokenContract.Register(auth)
	return err
}

func (d *TContract) BalanceOf() (int64, error) {
	bal, err := d.tokenContract.BalanceOf(nil, d.callerAddr)
	if err != nil {
		return 0, err
	}
	return bal.Int64(), nil
}

func (d *TContract) Account() (*Account, error) {
	return &Account{
		Address: d.callerAddr.String(),
	}, nil
}

var connToken *ethclient.Client

func InitTokenClient(connStr string) {
	var err error
	connToken, err = ethclient.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}

func NewTokenContract(callerName string, callerKey *ecdsa.PrivateKey, boundContract *ethcontract.Cards) TokenContract {
	return &TContract{
		Name:          callerName,
		tokenContract: boundContract,
		callerKey:     callerKey,
		callerAddr:    crypto.PubkeyToAddress(callerKey.PublicKey),
	}
}
