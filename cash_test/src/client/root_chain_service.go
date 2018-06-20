package client

import (
	"crypto/ecdsa"
	"ethcontract"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RootChainService struct {
	Name           string
	plasmaContract *ethcontract.RootChain
	callerKey      *ecdsa.PrivateKey
	callerAddr     common.Address
}

func (d *RootChainService) PlasmaCoin(slot uint64) (*PlasmaCoin, error) {
	uid, depositBlockNum, denom, ownerAddr, state, err := d.plasmaContract.GetPlasmaCoin(&bind.CallOpts{
		From: d.callerAddr,
	}, slot)
	if err != nil {
		return nil, err
	}
	return &PlasmaCoin{
		UID:             uid,
		DepositBlockNum: depositBlockNum.Int64(),
		Denomination:    denom,
		Owner:           ownerAddr.Hex(),
		State:           PlasmaCoinState(state),
	}, nil
}

func (d *RootChainService) Withdraw(slot uint64) error {
	auth := bind.NewKeyedTransactor(d.callerKey)
	_, err := d.plasmaContract.Withdraw(auth, slot)
	return err
}

func (d *RootChainService) ChallengeBefore(slot uint64, prevTxBytes []byte, exitingTxBytes []byte,
	prevTxInclusionProof Proof, exitingTxInclusionProof Proof,
	sig []byte, prevTxBlockNum int64, exitingTxBlockNum int64) ([]byte, error) {

	//return self.sign_and_send(self.contract.functions.challengeBefore, args,	  value=self.BOND)
	return []byte{}, nil
}

func (d *RootChainService) RespondChallengeBefore(slot uint64, challengingBlockNumber int64,
	challenging_transaction Tx, proof Proof) ([]byte, error) {

	//return self.sign_and_send(self.contract.functions.respond_challenge_before,					args)
	return []byte{}, nil
}

func (d *RootChainService) ChallengeBetween(slot uint64, challengingBlockNumber int64,
	challengingTransaction Tx, proof Proof) ([]byte, error) {
	//return self.sign_and_send(self.contract.functions.challengeBetween, args)
	return []byte{}, nil
}

func (d *RootChainService) ChallengeAfter(slot uint64, challengingBlockNumber int64,
	challengingTransaction Tx, proof Proof) ([]byte, error) {
	//return self.sign_and_send(self.contract.functions.challengeAfter, args)
	return []byte{}, nil
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
		callerAddr:     crypto.PubkeyToAddress(callerKey.PublicKey),
		plasmaContract: boundContract,
	}
}
