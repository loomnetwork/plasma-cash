package gocontracts

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/loomnetwork/transfer-gateway/gocontracts/sha3"
)

//TODO in future all interfaces and not do conversions from strings
type Pair struct {
	Type  string
	Value string
}

func SoliditySHA3(pairs []*Pair) (error, []byte) {
	//convert to packed bytes like solidity
	err, data := SolidityPackedBytes(pairs)
	if err != nil {
		return err, nil
	}

	d := sha3.NewKeccak256()
	d.Write(data)
	return nil, d.Sum(nil)
}

func SolidityPackedBytes(pairs []*Pair) (error, []byte) {
	var b bytes.Buffer

	for _, pair := range pairs {
		fmt.Printf("%v\n", pair)
		switch strings.ToLower(pair.Type) {
		case "address":
			decoded, err := hex.DecodeString(pair.Value)
			if err != nil {
				return err, nil
			}
			if len(decoded) != 20 {
				return fmt.Errorf("we don't support partial addresses, the len was %d we wanted 20", len(decoded)), nil
			}
			b.Write(decoded)
		case "uint16": //"uint", "uint16", "uint64":
			//pack integers
			u, err := strconv.ParseUint(pair.Value, 10, 32)
			if err != nil {
				return err, nil
			}
			var bTest []byte = make([]byte, 2)
			//			binary.LittleEndian.PutUint32(bTest, uint32(u))
			//			fmt.Printf("little-%v\n", bTest)
			binary.BigEndian.PutUint16(bTest, uint16(u))
			fmt.Printf("big-%v\n", bTest)
			b.Write(bTest)
		case "uint32": //"uint", "uint16", "uint64":
			//pack integers
			u, err := strconv.ParseUint(pair.Value, 10, 32)
			if err != nil {
				return err, nil
			}
			var bTest []byte = make([]byte, 4)
			//			binary.LittleEndian.PutUint32(bTest, uint32(u))
			//			fmt.Printf("little-%v\n", bTest)
			binary.BigEndian.PutUint32(bTest, uint32(u))
			fmt.Printf("big-%v\n", bTest)
			b.Write(bTest)
		case "uint64": //"uint", "uint16", "uint64":
			//pack integers
			u, err := strconv.ParseUint(pair.Value, 10, 32)
			if err != nil {
				return err, nil
			}
			var bTest []byte = make([]byte, 4)
			//			binary.LittleEndian.PutUint32(bTest, uint32(u))
			//			fmt.Printf("little-%v\n", bTest)
			binary.BigEndian.PutUint32(bTest, uint32(u))
			fmt.Printf("big-%v\n", bTest)
			b.Write(bTest)
		case "uint256":
			n := new(big.Int)
			_, valid := n.SetString(pair.Value, 10)
			if !valid {
				return errors.New("invalid big int"), nil
			}

			bytes := n.Bytes()
			padlen := 32 - len(bytes)
			if padlen < 0 {
				return errors.New("big int byte length too large"), nil
			}
			pad := make([]byte, padlen, padlen)
			b.Write(pad)
			b.Write(bytes)
		}
	}

	return nil, b.Bytes()
}
