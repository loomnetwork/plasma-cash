package client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
)

// Plasma Block
type PbBlock struct {
	block *pctypes.PlasmaBlock
}

func NewClientBlock(pb *pctypes.PlasmaBlock) Block {
	return &PbBlock{pb}
}

func (p *PbBlock) MerkleHash() []byte {
	return p.block.MerkleHash
}

func (p *PbBlock) TxFromSlot(slot uint64) (Tx, error) {
	var tx *pctypes.PlasmaTx

	for _, v := range p.block.Transactions {
		fmt.Printf("testing transaction at slot-%d\n", v.Slot)
		if v.Slot == slot {
			tx = v
		}
	}
	if tx == nil {
		return nil, fmt.Errorf("can't find transaction at slot %d. We had %d Transactions\n", slot, len(p.block.Transactions))
	}

	address := tx.NewOwner.Local.String()
	fmt.Printf("ethaddress-b-%s\n", address)
	ethAddress := common.HexToAddress(address)
	fmt.Printf("ethaddress-%s previous block -%d   Denomination-%d\n", ethAddress, tx.GetPreviousBlock().Value.Int64(), tx.Denomination.Value.Uint64())

	return &LoomTx{Slot: slot,
		PrevBlock:    big.NewInt(tx.GetPreviousBlock().Value.Int64()), //TODO ugh bad casting
		Denomination: uint32(tx.Denomination.Value.Uint64()),          //TODO get this from somewhere
		Owner:        ethAddress,
		Signature:    tx.Signature}, nil
}
