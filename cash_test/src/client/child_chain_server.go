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

// ChildChainService child client to reference server
type ChildChainService struct {
}

func (c *ChildChainService) CurrentBlock() (error, *Block) {
	return nil, nil
}

func (c *ChildChainService) BlockNumber() int {
	return 0
}

func (c *ChildChainService) Block(blknum int) (error, *Block) {
	return nil, nil
}

func (c *ChildChainService) Proof(blknum int, uid int) (error, *Proof) {
	return nil, nil

}

func (c *ChildChainService) SubmitBlock(*Block) error {
	return nil
}

func (c *ChildChainService) SendTransaction() error {
	return nil
}

func NewChildChainServer() ChainServiceClient {
	return &ChildChainService{}
}
