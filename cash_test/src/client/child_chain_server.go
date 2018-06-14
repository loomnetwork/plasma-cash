package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ChildChainService child client to reference server
type ChildChainService struct {
	url string
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

	type jResponse struct {
		Last_block_height   int
		Last_block_app_hash string
	}

	type jResult struct {
		Response jResponse
	}

	type jBlock struct {
		Jsonrpc string
		Id      string
		Result  jResult
	}

	var jblock jBlock

	err = json.Unmarshal([]byte(body), &jblock)

	if err != nil {
		fmt.Print(err)
	}
	return jblock.Result.Response.Last_block_height
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
