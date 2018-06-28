package client

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/client/plasma_cash"
)

func NewLoomChildChainService(writeuri, readuri string) (plasma_cash.ChainServiceClient, error) {
	privFile := "test.key" // hard coded for integration tests

	privKeyB64, err := ioutil.ReadFile(privFile)
	if err != nil {
		return nil, err
	}

	privKey, err := base64.StdEncoding.DecodeString(string(privKeyB64))
	if err != nil {
		return nil, err
	}

	signer := auth.NewEd25519Signer(privKey)

	return client.NewPlasmaCashClient(signer, "", writeuri, readuri)
}
