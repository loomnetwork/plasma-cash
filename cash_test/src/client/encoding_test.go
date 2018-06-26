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

	tx.PrevBlock = big.NewInt(85478557858583)
	txBytes, err = tx.RlpEncode()
	if err != nil {
		c.Fatal(err)
	}
	hexStr = common.Bytes2Hex(txBytes)
	c.Assert(hexStr, Equals, "de05864dbe0713bb1701945194b63f10691e46635b27925100cfc0a5ceca62")
}

func (s *EncodingTestSuite) TestTxSignature(c *C) {
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(GetTestAccountHexKey("alice"), "0x"))
	if err != nil {
		c.Fatal(err)
	}
	ownerAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	tx := &LoomTx{
		Slot:         5,
		PrevBlock:    big.NewInt(85478557858583),
		Denomination: 1,
		Owner:        ownerAddr,
	}
	txBytes, err := tx.Sign(privKey)
	if err != nil {
		c.Fatal(err)
	}
	hexStr := common.Bytes2Hex(txBytes)
	c.Assert(hexStr, Equals, "00b0e4901dc74b9851dba3c52406e1325c2ac9c4fe9f4d0379099a3357b763c96c104d3fffb78e99515db2e583568588d740b743ad3105d63fb252014f806fd06b1b")
}
