package client

type Tx struct {
}

type Block struct {
	blockId string
}

type Proof struct {
}

type Account struct {
	Address string
}

type TokenContract interface {
	Register()
	Deposit(int)
	BalanceOf() int

	Account() *Account
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
