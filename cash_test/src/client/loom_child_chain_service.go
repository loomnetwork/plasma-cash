package client

import (
	"fmt"
	"log"

	loom "github.com/loomnetwork/go-loom"
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
	"github.com/loomnetwork/go-loom/types"
)

// LoomChildChainService child client to reference server
type LoomChildChainService struct {
	url          string
	ChainID      string
	WriteURI     string
	ReadURI      string
	contractAddr string
	loomcontract *LoomContract
}

func (c *LoomChildChainService) CurrentBlock() (error, *Block) {
	return nil, nil
}

func (c *LoomChildChainService) BlockNumber() int64 {
	return int64(0)
}

func (c *LoomChildChainService) Block(blknum int64) (error, *Block) {
	fmt.Printf("trying to get Block data\n")
	blk := loom.NewBigUIntFromInt(blknum)

	var result pctypes.GetBlockResponse
	params := &pctypes.GetBlockRequest{
		BlockHeight: &types.BigUInt{*blk},
	}

	if err := c.loomcontract.StaticCallContract("GetBlockRequest", params, &result); err != nil {
		log.Fatalf("failed getting Block data - %v\n", err)

		return err, nil
	}

	log.Printf("get block value %v '\n", result)

	return nil, nil
}

func (c *LoomChildChainService) Proof(blknum int64, uid int64) (error, *Proof) {

	return nil, nil
}

func (c *LoomChildChainService) SubmitBlock(*Block) error {
	return nil
}

func (c *LoomChildChainService) SendTransaction() error {
	return nil
}

func NewLoomChildChainService(readuri, writeuri string) ChainServiceClient {
	fmt.Printf("Using Loom Service as Plasma Chain\n")
	return &LoomChildChainService{loomcontract: NewLoomContract(readuri, writeuri, "plasmacash")}
}
