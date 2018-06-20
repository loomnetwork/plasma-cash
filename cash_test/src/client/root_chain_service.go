package client

import (
	"crypto/ecdsa"
	"ethcontract"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RootChainService struct {
	Name           string
	plasmaContract *ethcontract.RootChain
	callerKey      *ecdsa.PrivateKey
}

// TODO: implement for challenge_after_demo
func (d *RootChainService) PlasmaCoin(uint64) {
}

func (d *RootChainService) Withdraw(slot uint64) error {
	auth := bind.NewKeyedTransactor(d.callerKey)
	_, err := d.plasmaContract.Withdraw(auth, slot)
	return err
}

func (d *RootChainService) StartExit(
	slot uint64, prevTx Tx, exitingTx Tx, prevTxInclusion Proof, exitingTxInclusion Proof,
	sigs []byte, prevTxIncBlock int64, exitingTxIncBlock int64) ([]byte, error) {
	auth := bind.NewKeyedTransactor(d.callerKey)
	// TODO: encode params into bytes...
	var prevTxBytes, exitingTxBytes, prevTxInclusionProof, exitingTxInclusionProof []byte
	_, err := d.plasmaContract.StartExit(
		auth, slot,
		prevTxBytes, exitingTxBytes, prevTxInclusionProof, exitingTxInclusionProof,
		sigs, big.NewInt(prevTxIncBlock), big.NewInt(exitingTxIncBlock))
	return []byte{}, err
}

func (d *RootChainService) FinalizeExits() error {
	_, err := d.plasmaContract.FinalizeExits(nil)
	return err
}

// TODO: implement for challenge_after_demo
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

func NewRootChainService(callerName string, callerKey *ecdsa.PrivateKey, boundContract *ethcontract.RootChain) *RootChainService {
	return &RootChainService{
		Name:           callerName,
		callerKey:      callerKey,
		plasmaContract: boundContract,
	}
}
