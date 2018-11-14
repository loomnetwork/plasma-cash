package client

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"ethcontract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/go-loom/auth"
	"github.com/loomnetwork/go-loom/client"
	"github.com/loomnetwork/go-loom/client/plasma_cash"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	loom "github.com/loomnetwork/go-loom"
	loom_ethcontract "github.com/loomnetwork/go-loom/client/plasma_cash/eth/ethcontract"

	"crypto/ecdsa"
	"crypto/elliptic"
)

type TestContext struct {
	Alice   *Client
	Bob     *Client
	Eve     *Client
	Mallory *Client
	Charlie *Client
	Dan     *Client
	Trudy   *Client

	Authority *Client
}

func getDAppchainTxSigner(name string) (auth.Signer, error) {
	privFile := name + ".key"

	privKeyB64, err := ioutil.ReadFile(privFile)
	if err != nil {
		return nil, err
	}

	privKey, err := base64.StdEncoding.DecodeString(string(privKeyB64))
	if err != nil {
		return nil, err
	}

	signer := auth.NewEd25519Signer(privKey)
	return signer, nil
}

func getTokenContract(cfg *viper.Viper, name string, privKey *ecdsa.PrivateKey) (plasma_cash.TokenContract, error) {
	tokenAddr := common.HexToAddress(cfg.GetString("token_contract"))
	tokenContract, err := ethcontract.NewCards(tokenAddr, conn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to instantiate a Token contract")
	}
	return NewTokenContract(name, privKey, tokenContract), nil
}

func getRootChain(cfg *viper.Viper, name string) (plasma_cash.RootChainClient, error) {
	contractAddr := common.HexToAddress(cfg.GetString("root_chain"))
	privKeyHexStr := cfg.GetString(name)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyHexStr, "0x"))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load private key for %s", name)
	}
	plasmaContract, err := loom_ethcontract.NewRootChain(contractAddr, conn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to instantiate a Token contract")
	}
	return NewRootChainService(name, privKey, plasmaContract), nil
}

// Loads plasma-config.yml or equivalent from the cwd
func parseConfig() (*viper.Viper, error) {
	// When running "go test" the cwd is set to the package dir, not the root dir
	// where the config is, so gotta do a bit more work to figure out the config dir...
	_, filename, _, _ := runtime.Caller(0)
	cfgDir := filepath.Join(filepath.Dir(filename), "../..")

	v := viper.New()

	v.SetConfigName("plasma-config")
	v.AddConfigPath(cfgDir)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func deriveAddressFromECPubKey(pubKey *ecdsa.PublicKey) loom.LocalAddress {
	pubKeyInBinary := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	return crypto.Keccak256(pubKeyInBinary[1:])[12:]
}

func setupClient(cfg *viper.Viper, addressMapper *AddressMapperClient, hostile bool, entityName, readUri, writeUri string) (*Client, error) {
	signer, err := getDAppchainTxSigner(entityName)
	if err != nil {
		return nil, err
	}

	contractName := "plasmacash"
	if hostile {
		contractName = "hostileoperator"
	}

	chainServiceClient, err := client.NewPlasmaCashClient(contractName, signer, "default", writeUri, readUri)
	if err != nil {
		return nil, err
	}

	privKeyHexStr := cfg.GetString(entityName)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyHexStr, "0x"))
	if err != nil {
		return nil, err
	}

	from := loom.Address{
		ChainID: "default",
		Local:   loom.LocalAddressFromPublicKey(signer.PublicKey()),
	}

	to := loom.Address{
		ChainID: "eth",
		Local:   deriveAddressFromECPubKey(&privKey.PublicKey),
	}

	hasMappingResponse, err := addressMapper.HasMapping(from, from)
	if err != nil {
		return nil, err
	}

	if !hasMappingResponse.HasMapping {
		err = addressMapper.AddIdentityMapping(from, to, signer, privKey)
		if err != nil {
			return nil, err
		}
	}

	rootChainClient, err := getRootChain(cfg, entityName)
	if err != nil {
		return nil, err
	}

	tokenContract, err := getTokenContract(cfg, entityName, privKey)
	if err != nil {
		return nil, err
	}

	return NewClient(cfg, chainServiceClient, rootChainClient, tokenContract), nil

}

func SetupTest(hostile bool, readUri, writeUri string) (*TestContext, error) {
	var err error
	testCtx := TestContext{}

	addressMapper, err := NewAddressMapperClient("default", writeUri, readUri)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig()
	if err != nil {
		return nil, err
	}

	testCtx.Alice, err = setupClient(cfg, addressMapper, hostile, "alice", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Bob, err = setupClient(cfg, addressMapper, hostile, "bob", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Charlie, err = setupClient(cfg, addressMapper, hostile, "charlie", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Dan, err = setupClient(cfg, addressMapper, hostile, "dan", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Eve, err = setupClient(cfg, addressMapper, hostile, "eve", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Mallory, err = setupClient(cfg, addressMapper, hostile, "mallory", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Trudy, err = setupClient(cfg, addressMapper, hostile, "trudy", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	testCtx.Authority, err = setupClient(cfg, addressMapper, hostile, "authority", readUri, writeUri)
	if err != nil {
		return nil, err
	}

	return &testCtx, nil
}
