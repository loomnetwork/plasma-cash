package client

import (
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
)

// Plasma Block
type PbBlock struct {
	block *pctypes.PlasmaBlock
	proof []byte
}

func (pb *PbBlock) Proof() []byte {
	return pb.proof
}
func NewClientBlock(pb *pctypes.PlasmaBlock) Block {
	return &PbBlock{}
	//fmt.Printf("proof---%v\n", pb.Proof)
	//return &PbBlock{pb, pb.Proof}
}
