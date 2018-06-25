package client

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/loomnetwork/go-loom/common/evmcompat"
)

type LoomTx struct {
	Slot         uint64
	Denomination uint32
	Owner        common.Address
	PrevBlock    *big.Int
	//IncludeBlock  *big.Int // TODO
}

//Python version signs here
func (l *LoomTx) Sig() []byte {
	return l.Hash()
}

func (l *LoomTx) RlpEncode() ([]byte, error) {
	logdebug("RlpEncode()")

	return rlp.EncodeToBytes([]interface{}{
		l.Slot,
		l.PrevBlock,
		l.Denomination,
		l.Owner,
	})
}

func (l *LoomTx) Hash() []byte {
	//    if l.IncludeBlock.Mod(1000) == 0 {
	//            ret = w3.sha3(rlp.encode(self, UnsignedTransaction))
	//   }

	//      else
	data, err := soliditySha3(l.Slot)
	if err != nil {
		panic(err)
	}
	return data
}

func (l *LoomTx) MerkleHash() []byte {
	//        return w3.sha3(rlp.encode(self))
	panic("TODO")
}

func soliditySha3(data uint64) ([]byte, error) {
	pairs := []*evmcompat.Pair{&evmcompat.Pair{"uint64", strconv.FormatUint(data, 10)}}
	hash, err := evmcompat.SoliditySHA3(pairs)
	if err != nil {
		return []byte{}, err
	}
	return hash, err
}
