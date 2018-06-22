package client

type Tx interface {
	Sig() []byte
	RlpEncode() ([]byte, error)
}

type Block interface {
}

type Proof interface {
	Bytes() []byte
}

//TODO not sure what this struct looks like
type SimpleProof struct {
	proofdata []byte
}

func (s SimpleProof) Bytes() []byte {
	return s.proofdata
}

type Account struct {
	Address string
}

type TokenContract interface {
	Register() error
	Deposit(int64) error
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
	StartExit(uid uint64, prevTx Tx, exitingTx Tx, prevTxProof Proof,
		exitingTxProof Proof, sigs []byte, prevTxBlkNum int64, txBlkNum int64) ([]byte, error)

	ChallengeBefore(slot uint64, prevTx Tx, exitingTx Tx,
		prevTxInclusionProof Proof, exitingTxInclusionProof Proof,
		sig []byte, prevTxBlockNum int64, exitingTxBlockNum int64) ([]byte, error)

	RespondChallengeBefore(slot uint64, challengingBlockNumber int64,
		challengingTransaction Tx, proof Proof) ([]byte, error)

	ChallengeBetween(slot uint64, challengingBlockNumber int64,
		challengingTransaction Tx, proof Proof) ([]byte, error)

	ChallengeAfter(slot uint64, challengingBlockNumber int64,
		challengingTransaction Tx, proof Proof) ([]byte, error)
}

type ChainServiceClient interface {
	CurrentBlock() (Block, error)
	BlockNumber() (int64, error)

	Block(blknum int64) (Block, error)
	//Proof(blknum int64, uid uint64) (Proof, error) //TODO what is the uid?

	SubmitBlock() error

	SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) (Tx, error)
}
