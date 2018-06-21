package oracle

import (
	"context"
	"ethcontract"
	"log"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/loomnetwork/go-loom"
	"github.com/loomnetwork/go-loom/auth"
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
	"github.com/loomnetwork/go-loom/client"
	ltypes "github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
)

type OracleConfig struct {
	// URI of an Ethereum node
	EthereumURI string
	// Plasma contract address on Ethereum
	PlasmaHexAddress string
	ChainID          string
	WriteURI         string
	ReadURI          string
	Signer           auth.Signer
}

type Oracle struct {
	cfg        OracleConfig
	solPlasma  *ethcontract.RootChain
	goPlasma   *client.Contract
	startBlock uint64
	ethClient  *ethclient.Client
}

func NewOracle(cfg OracleConfig) *Oracle {
	return &Oracle{cfg: cfg}
}

func (orc *Oracle) Init() error {
	cfg := &orc.cfg
	var err error
	orc.ethClient, err = ethclient.Dial(cfg.EthereumURI)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Ethereum")
	}

	orc.solPlasma, err = ethcontract.NewRootChain(common.HexToAddress(cfg.PlasmaHexAddress), orc.ethClient)
	if err != nil {
		return errors.Wrap(err, "failed to bind Plasma Solidity contract")
	}

	dappClient := client.NewDAppChainRPCClient(cfg.ChainID, cfg.WriteURI, cfg.ReadURI)
	contractAddr, err := dappClient.Resolve("plasmacash")
	if err != nil {
		return errors.Wrap(err, "failed to resolve Plasma Go contract address")
	}
	orc.goPlasma = client.NewContract(dappClient, contractAddr.Local)
	return nil
}

// RunWithRecovery should run in a goroutine, it will ensure the oracle keeps on running as long
// as it doesn't panic due to a runtime error.
func (orc *Oracle) RunWithRecovery() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic in Gateway Oracle: %v\n", r)
			// Unless it's a runtime error restart the goroutine
			if _, ok := r.(runtime.Error); !ok {
				time.Sleep(30 * time.Second)
				log.Printf("Restarting Gateway Oracle...")
				go orc.RunWithRecovery()
			}
		}
	}()
	orc.Run()
}

// TODO: Graceful shutdown
func (orc *Oracle) Run() {
	skipSleep := true
	for {
		if !skipSleep {
			// TODO: should be configurable
			time.Sleep(5 * time.Second)
		} else {
			skipSleep = false
		}

		// TODO: get start block from Plasma Go contract, like the Transfer Gateway Oracle
		startBlock := orc.startBlock

		// TODO: limit max block range per batch
		latestBlock, err := orc.getLatestEthBlockNumber()
		if err != nil {
			log.Printf("failed to obtain latest Ethereum block number: %v", err)
			continue
		}

		if latestBlock < startBlock {
			// Wait for Ethereum to produce a new block...
			continue
		}

		deposits, err := orc.fetchDeposits(startBlock, latestBlock)
		if err != nil {
			log.Printf("failed to fetch events from Ethereum: %v", err)
			continue
		}

		for _, deposit := range deposits {
			if _, err := orc.goPlasma.Call("DepositRequest", deposit, orc.cfg.Signer, nil); err != nil {
				log.Printf("failed to commit DepositRequest tx: %v", err)
				continue
			}
		}

		orc.startBlock = latestBlock + 1
	}
}

func (orc *Oracle) getLatestEthBlockNumber() (uint64, error) {
	blockHeader, err := orc.ethClient.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		return 0, err
	}
	return blockHeader.Number.Uint64(), nil
}

// Fetches all deposit events from an Ethereum node from startBlock to endBlock (inclusive)
func (orc *Oracle) fetchDeposits(startBlock, endBlock uint64) ([]*pctypes.DepositRequest, error) {
	// NOTE: Currently either all blocks from w.StartBlock are processed successfully or none are.
	filterOpts := &bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}
	deposits := []*pctypes.DepositRequest{}

	it, err := orc.solPlasma.FilterDeposit(filterOpts, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Plasma deposit logs")
	}
	for {
		ok := it.Next()
		if ok {
			ev := it.Event
			fromAddr, err := loom.LocalAddressFromHexString(ev.From.Hex())
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse Plasma deposit 'from' address")
			}
			deposits = append(deposits, &pctypes.DepositRequest{
				Slot:         ev.Slot,
				DepositBlock: &ltypes.BigUInt{Value: *loom.NewBigUInt(ev.BlockNumber)},
				Denomination: &ltypes.BigUInt{Value: *loom.NewBigUIntFromInt(1)}, // TODO: ev.Denomination
				From:         loom.Address{ChainID: "eth", Local: fromAddr}.MarshalPB(),
				// TODO: store ev.Hash ... should it's always a hash of ev.Slot, so a bit redundant
			})
		} else {
			err := it.Error()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get event data for ETHReceived")
			}
			it.Close()
			break
		}
	}

	return deposits, nil
}
