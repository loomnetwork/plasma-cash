package client

import (
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
)

// Plasma Block
type PbBlock struct {
	block *pctypes.PlasmaBlock
}

func NewClientBlock(pb *pctypes.PlasmaBlock) Block {
	return &PbBlock{pb}
}
