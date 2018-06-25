package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ChildChainService child client to reference server
type ChildChainService struct {
	url string
}

type DummyBlock struct {
	blockId string
}

func (d *DummyBlock) TxFromSlot(slot uint64) (Tx, error) {
	return &LoomTx{}, nil
}
func (d *DummyBlock) MerkleHash() []byte {
	return []byte{}
}

func (c *ChildChainService) CurrentBlock() (Block, error) {
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
	// fmt.Println("CurrentBlock Body:", string(body))
	return &DummyBlock{blockId: string(body)}, nil
}

func (c *ChildChainService) BlockNumber() (int64, error) {
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
	// fmt.Println("BlockNumber Body:", string(body))

	i, _ := strconv.Atoi(string(body))
	return int64(i), nil
}

func (c *ChildChainService) Block(blknum int64) (Block, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/block/%d", c.url, blknum), nil)
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
	// fmt.Println("Block Body:", string(body))
	return &DummyBlock{blockId: string(body)}, nil
}

func (c *ChildChainService) Proof(blknum int64, uid uint64) (Proof, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/proof?blknum=%d&uid=%d", c.url, blknum, uid), nil)
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
	proof := SimpleProof{proofdata: body}
	return proof, nil
}

func (c *ChildChainService) SubmitBlock() error {

	data := url.Values{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/submit_block", c.url), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	return err
}

type ChildChainTx struct {
}

func (c *ChildChainTx) Sig() []byte {
	return []byte{}
}

func (c *ChildChainTx) RlpEncode() ([]byte, error) {
	return []byte{}, nil
}

func (c *ChildChainService) SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) (Tx, error) {
	/*
	   new_owner = utils.normalize_address(new_owner)
	   incl_block = c.BlockNumber()
	   tx = Transaction(slot, prev_block, denomination, new_owner,
	                    incl_block=incl_block)
	   tx.sign(c.key)
	*/

	return &ChildChainTx{}, nil
}

func NewChildChainService(url string) ChainServiceClient {
	return &ChildChainService{url}
}
