package client

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type EncodingTestSuite struct{}

var _ = Suite(&EncodingTestSuite{})

func (s *EncodingTestSuite) TestUnsignedTxRlpEncode(c *C) {
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(GetTestAccountHexKey("alice"), "0x"))
	if err != nil {
		c.Fatal(err)
	}
	ownerAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	tx := &LoomTx{
		Slot:         5,
		PrevBlock:    big.NewInt(0),
		Denomination: 1,
		Owner:        ownerAddr,
	}
	txBytes, err := tx.RlpEncode()
	if err != nil {
		c.Fatal(err)
	}
	hexStr := common.Bytes2Hex(txBytes)
	c.Assert(hexStr, Equals, "d8058001945194b63f10691e46635b27925100cfc0a5ceca62")
}
