package main

import (
	"client"
	"log"
	"oracle"

	"github.com/loomnetwork/go-loom/auth"
	"golang.org/x/crypto/ed25519"
)

func main() {
	_, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("failed to gnerate private key for authority: %v", err)
	}
	plasmaOrc := oracle.NewOracle(oracle.OracleConfig{
		EthereumURI:      "http://localhost:8545",
		PlasmaHexAddress: client.GetContractHexAddress("root_chain"),
		ChainID:          "default",
		WriteURI:         "http://localhost:46658/rpc",
		ReadURI:          "http://localhost:46658/query",
		Signer:           auth.NewEd25519Signer(privKey),
	})
	if err := plasmaOrc.Init(); err != nil {
		log.Fatal(err)
	}
	plasmaOrc.Run()
}
