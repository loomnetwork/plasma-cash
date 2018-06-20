package client

type Tx interface {
}

type Block interface {
}

type DummyBlock struct {
	blockId string
}

type Proof struct {
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
	FinalizeExits()
	Withdraw(uint64)
	WithdrawBonds()
	PlasmaCoin(uint64)
}

type ChainServiceClient interface {
	CurrentBlock() (Block, error)
	BlockNumber() int64

	Block(blknum int64) (Block, error)
	Proof(blknum int64, uid int64) (*Proof, error) //TODO what is the uid?

	SubmitBlock() error

	SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) (Tx, error)
}
