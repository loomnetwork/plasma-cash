package main

import (
	"client"
	"log"
	"oracle"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/go-loom/auth"
	"golang.org/x/crypto/ed25519"
)

func main() {
	_, loomPrivKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("failed to gnerate private key for authority: %v", err)
	}
	ethPrivKeyHexStr := client.GetTestAccountHexKey("authority")
	ethPrivKey, err := crypto.HexToECDSA(strings.TrimPrefix(ethPrivKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	plasmaOrc := oracle.NewOracle(oracle.OracleConfig{
		PlasmaBlockInterval: 1000,
		EthereumURI:         "http://localhost:8545",
		PlasmaHexAddress:    client.GetContractHexAddress("root_chain"),
		ChainID:             "default",
		WriteURI:            "http://localhost:46658/rpc",
		ReadURI:             "http://localhost:46658/query",
		Signer:              auth.NewEd25519Signer(loomPrivKey),
		EthPrivateKey:       ethPrivKey,
		OverrideGas:         true,
	})
	if err := plasmaOrc.Init(); err != nil {
		log.Fatal(err)
	}
	plasmaOrc.Run()
}
