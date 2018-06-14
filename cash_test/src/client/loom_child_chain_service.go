package client

import "fmt"

// LoomChildChainService child client to reference server
type LoomChildChainService struct {
	url string
}

func (c *LoomChildChainService) CurrentBlock() (error, *Block) {
	return nil, nil
}

func (c *LoomChildChainService) BlockNumber() int {
	return 0
}

func (c *LoomChildChainService) Block(blknum int) (error, *Block) {
	return nil, nil
}

func (c *LoomChildChainService) Proof(blknum int, uid int) (error, *Proof) {
	return nil, nil

}

func (c *LoomChildChainService) SubmitBlock(*Block) error {
	return nil
}

func (c *LoomChildChainService) SendTransaction() error {
	return nil
}

func NewLoomChildChainService(url string) ChainServiceClient {
	fmt.Printf("Using Loom Service as Plasma Chain\n")
	return &LoomChildChainService{url}
}
