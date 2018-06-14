package client

type Block struct {
}

type Proof struct {
}

type TokenContract interface {
	Register()
	Deposit(int)
	BalanceOf() int
}

type RootChainClient interface {
	FinalizeExits()
	Withdraw(int)
	WithdrawBonds()
	PlasmaCoin(int)
}

type ChainServiceClient interface {
	CurrentBlock() (error, *Block)
	BlockNumber() int

	Block(blknum int) (error, *Block)
	Proof(blknum int, uid int) (error, *Proof) //TODO what is the uid?

	SubmitBlock(*Block) error

	SendTransaction() error
}
