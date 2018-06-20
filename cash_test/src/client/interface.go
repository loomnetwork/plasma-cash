package client

type Tx interface {
	Sig() []byte
}

type Block interface {
}

type Proof interface {
}

//TODO not sure what this struct looks like
type SimpleProof struct {
	proofstring string
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

type RootChainClient interface {
	FinalizeExits() error
	Withdraw(uint64)
	WithdrawBonds()
	PlasmaCoin(uint64)
	StartExit(uid uint64, prevTx Tx, exiting_tx Tx, prevTxProof Proof,
		exitingTxProof Proof, sigs []byte, prevTxBlkNum int64, txBlkNum int64) ([]byte, error)
}

type ChainServiceClient interface {
	CurrentBlock() (Block, error)
	BlockNumber() int64

	Block(blknum int64) (Block, error)
	Proof(blknum int64, uid uint64) (Proof, error) //TODO what is the uid?

	SubmitBlock() error

	SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) (Tx, error)
}
