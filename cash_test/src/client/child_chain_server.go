package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

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

	req, err := http.NewRequest("GET", "http://localhost:46657/abci_info", nil)
	if err != nil {
		fmt.Print(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
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
