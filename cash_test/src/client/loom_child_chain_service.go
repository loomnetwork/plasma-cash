package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// LoomChildChainService child client to reference server
type LoomChildChainService struct {
	url string
}

func (c *LoomChildChainService) CurrentBlock() (error, *Block) {
	return nil, nil
}

func (c *LoomChildChainService) BlockNumber() int {

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
