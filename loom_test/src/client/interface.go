package client

import (
	"crypto/ecdsa"
	"ethcontract"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/loomnetwork/go-loom/client/plasma_cash"
)

type Account struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

type TokenContract interface {
	Register() error
	Deposit(int64) (common.Hash, error)
	BalanceOf() (int64, error)

	Account() (*Account, error)
}

type PlasmaCoinState uint8

const (
	PlasmaCoinDeposited PlasmaCoinState = iota
	PlasmaCoinExiting
	PlasmaCoinChallenged
	PlasmaCoinResponded
	PlasmaCoinExited
)

type PlasmaCoin struct {
	UID             uint64
	DepositBlockNum int64
	Denomination    uint32
	Owner           string
	State           PlasmaCoinState
}

type RootChainClient interface {
	FinalizeExits() error
	Withdraw(slot uint64) error
	WithdrawBonds() error
	PlasmaCoin(slot uint64) (*PlasmaCoin, error)
	StartExit(slot uint64, prevTx Tx, exitingTx Tx, prevTxProof Proof,
		exitingTxProof Proof, sigs []byte, prevTxBlkNum int64, txBlkNum int64) ([]byte, error)

	ChallengeBefore(slot uint64, prevTx Tx, exitingTx Tx,
		prevTxInclusionProof Proof, exitingTxInclusionProof Proof,
		sig []byte, prevTxBlockNum int64, exitingTxBlockNum int64) ([]byte, error)

	RespondChallengeBefore(slot uint64, challengingBlockNumber int64,
		challengingTransaction Tx, proof Proof, sig []byte) ([]byte, error)

	ChallengeBetween(slot uint64, challengingBlockNumber int64,
		challengingTransaction Tx, proof Proof, sig []byte) ([]byte, error)

	ChallengeAfter(slot uint64, challengingBlockNumber int64,
		challengingTransaction Tx, proof Proof, sig []byte) ([]byte, error)

	SubmitBlock(blockNum *big.Int, merkleRoot [32]byte) error

	DebugCoinMetaData(slots []uint64)
	DepositEventData(txHash common.Hash) (*ethcontract.RootChainDeposit, error)
}
