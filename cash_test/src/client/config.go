package client

import (
	"ethcontract"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

// Loads plasma-config.yml or equivalent from the cwd
func parseConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("plasma-config")
	v.AddConfigPath(".")
	return v, v.ReadInConfig()
}

func GetTokenContract(name string) TokenContract {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	tokenAddr := common.HexToAddress(cfg.GetString("token_contract"))
	privKeyHexStr := cfg.GetString(name)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key for %s: %v", name, err)
	}
	tokenContract, err := ethcontract.NewCards(tokenAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	return NewTokenContract(name, privKey, tokenContract)
}

func GetRootChain(name string) RootChainClient {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	contractAddr := common.HexToAddress(cfg.GetString("root_chain"))
	privKeyHexStr := cfg.GetString(name)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key for %s: %v", name, err)
	}
	plasmaContract, err := ethcontract.NewRootChain(contractAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	return NewRootChainService(name, privKey, plasmaContract)
}
