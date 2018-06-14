package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ChildChainService child client to reference server
type ChildChainService struct {
	url string
}

func (c *ChildChainService) CurrentBlock() (error, *Block) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/block", c.url), nil)
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
	block := Block{blockId: string(body)}
	return nil, &block
}

func (c *ChildChainService) BlockNumber() int {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/blocknumber", c.url), nil)
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

	i, _ := strconv.Atoi(string(body))
	return i
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
