package client

type Tx struct {
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
	Deposit(int) error
	BalanceOf() (int, error)

	Account() (*Account, error)
}

type RootChainClient interface {
	FinalizeExits()
	Withdraw(int)
	WithdrawBonds()
	PlasmaCoin(int)
}

type ChainServiceClient interface {
	CurrentBlock() (Block, error)
	BlockNumber() int64

	Block(blknum int64) (Block, error)
	Proof(blknum int64, uid int64) (*Proof, error) //TODO what is the uid?

	SubmitBlock() error

	SendTransaction() error
}
