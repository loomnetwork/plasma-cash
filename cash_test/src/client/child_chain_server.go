package client

// ChildChainService child client to reference server
type ChildChainService struct {
	url string
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

func NewChildChainService(url string) ChainServiceClient {
	return &ChildChainService{url}
}
