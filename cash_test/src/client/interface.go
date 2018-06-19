package client

type Tx struct {
}

type Block struct {
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
	CurrentBlock() (error, *Block)
	BlockNumber() int64

	Block(blknum int64) (error, *Block)
	Proof(blknum int64, uid int64) (error, *Proof) //TODO what is the uid?

	SubmitBlock(*Block) error

	SendTransaction() error
}
