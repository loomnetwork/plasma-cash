package client

type Block struct {
}

type Proof struct {
}

type ChainServiceClient interface {
	CurrentBlock() (error, *Block)
	BlockNumber() int

	Block(blknum int) (error, *Block)
	Proof(blknum int, uid int) (error, *Proof) //TODO what is the uid?

	SubmitBlock(*Block) error

	SendTransaction() error
}
