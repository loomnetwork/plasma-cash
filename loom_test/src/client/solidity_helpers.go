package client

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SoliditySign signs the given data with the specified private key and returns the 65-byte signature.
// The signature is in a format that's compatible with the ecverify() Solidity function.
func SoliditySign(data []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(data, privKey)
	if err != nil {
		return nil, err
	}

	v := sig[len(sig)-1]
	sig[len(sig)-1] = v + 27
	return sig, nil
}

// SolidityRecover recovers the Ethereum address from the signed hash and the 65-byte signature.
func SolidityRecover(hash []byte, sig []byte) (common.Address, error) {
	stdSig := make([]byte, 65)
	copy(stdSig[:], sig[:])
	stdSig[len(sig)-1] -= 27

	var signer common.Address
	pubKey, err := crypto.Ecrecover(hash, stdSig)
	if err != nil {
		return signer, err
	}

	copy(signer[:], crypto.Keccak256(pubKey[1:])[12:])
	return signer, nil
}
