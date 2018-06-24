package oracle

import (
	"log"
	"math/big"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

type OracleConfig struct {
	// Each Plasma block number must be a multiple of this value
	PlasmaBlockInterval uint32
	DAppChainClientCfg  DAppChainPlasmaClientConfig
	EthClientCfg        EthPlasmaClientConfig
}

type Oracle struct {
	cfg              OracleConfig
	ethPlasmaClient  EthPlasmaClient
	dappPlasmaClient DAppChainPlasmaClient
	startEthBlock    uint64 // Eth block from which the oracle should start looking for deposits
}

func NewOracle(cfg OracleConfig) *Oracle {
	return &Oracle{
		cfg:              cfg,
		ethPlasmaClient:  &EthPlasmaClientImpl{EthPlasmaClientConfig: cfg.EthClientCfg},
		dappPlasmaClient: &DAppChainPlasmaClientImpl{DAppChainPlasmaClientConfig: cfg.DAppChainClientCfg},
	}
}

func (orc *Oracle) Init() error {
	if err := orc.ethPlasmaClient.Init(); err != nil {
		return err
	}
	return orc.dappPlasmaClient.Init()
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

		if orc.cfg.EthClientCfg.PrivateKey != nil {
			if err := orc.sendPlasmaBlocksToEthereum(); err != nil {
				log.Println(err.Error())
			}
		}

		if err := orc.sendPlasmaDepositsToDAppChain(); err != nil {
			log.Println(err.Error())
		}
	}
}

// DAppChain -> Plasma Blocks -> Ethereum
func (orc *Oracle) sendPlasmaBlocksToEthereum() error {
	if err := orc.syncPlasmaBlocksWithEthereum(); err != nil {
		return errors.Wrap(err, "failed to sync plasma blocks with mainnet")
	}

	return orc.dappPlasmaClient.FinalizeCurrentPlasmaBlock()
}

// Ethereum -> Plasma Deposits -> DAppChain
func (orc *Oracle) sendPlasmaDepositsToDAppChain() error {
	// TODO: get start block from Plasma Go contract, like the Transfer Gateway Oracle
	startEthBlock := orc.startEthBlock

	// TODO: limit max block range per batch
	latestEthBlock, err := orc.ethPlasmaClient.LatestEthBlockNum()
	if err != nil {
		return errors.Wrap(err, "failed to obtain latest Ethereum block number")
	}

	if latestEthBlock < startEthBlock {
		// Wait for Ethereum to produce a new block...
		return nil
	}

	deposits, err := orc.ethPlasmaClient.FetchDeposits(startEthBlock, latestEthBlock)
	if err != nil {
		return errors.Wrap(err, "failed to fetch Plasma deposits from Ethereum")
	}

	for _, deposit := range deposits {
		if err := orc.dappPlasmaClient.Deposit(deposit); err != nil {
			return err
		}
	}

	orc.startEthBlock = latestEthBlock + 1
	return nil
}

// Send any finalized but unsubmitted plasma blocks from the DAppChain to Ethereum.
func (orc *Oracle) syncPlasmaBlocksWithEthereum() error {
	curEthPlasmaBlockNum, err := orc.ethPlasmaClient.CurrentPlasmaBlockNum()
	if err != nil {
		return err
	}
	log.Printf("solPlasma.CurrentBlock: %s", curEthPlasmaBlockNum.String())

	curLoomPlasmaBlockNum, err := orc.dappPlasmaClient.CurrentPlasmaBlockNum()
	if err != nil {
		return err
	}

	if curLoomPlasmaBlockNum.Cmp(curEthPlasmaBlockNum) == 0 {
		// DAppChain and Ethereum both have all the finalized Plasma blocks
		return nil
	}

	plasmaBlockInterval := big.NewInt(int64(orc.cfg.PlasmaBlockInterval))
	unsubmittedPlasmaBlockNum := nextPlasmaBlockNum(curEthPlasmaBlockNum, plasmaBlockInterval)

	for {
		log.Printf("unsubmittedPlasmaBlockNum: %s", unsubmittedPlasmaBlockNum.String())
		if unsubmittedPlasmaBlockNum.Cmp(curLoomPlasmaBlockNum) > 0 {
			// All the finalized plasma blocks in the DAppChain have been submitted to Ethereum
			break
		}

		block, err := orc.dappPlasmaClient.PlasmaBlockAt(unsubmittedPlasmaBlockNum)
		if err != nil {
			return err
		}

		// TODO: Will this block until the tx is confirmed? If not then should wait until it is
		//       confirmed before submitting another block.
		if err := orc.submitPlasmaBlockToEthereum(unsubmittedPlasmaBlockNum, block.MerkleHash); err != nil {
			return errors.Wrap(err, "failed to submit plasma block to Ethereum")
		}

		unsubmittedPlasmaBlockNum = nextPlasmaBlockNum(unsubmittedPlasmaBlockNum, plasmaBlockInterval)
	}

	return nil
}

// Submits a Plasma block (or rather its merkle root) to the Plasma Solidity contract on Ethereum
func (orc *Oracle) submitPlasmaBlockToEthereum(plasmaBlockNum *big.Int, merkleRoot []byte) error {
	curEthPlasmaBlockNum, err := orc.ethPlasmaClient.CurrentPlasmaBlockNum()
	if err != nil {
		return err
	}

	// Try to avoid submitting the same plasma blocks multiple times
	if plasmaBlockNum.Cmp(curEthPlasmaBlockNum) <= 0 {
		return nil
	}

	if len(merkleRoot) != 32 {
		return errors.New("invalid merkle root size")
	}

	var root [32]byte
	copy(root[:], merkleRoot)
	return orc.ethPlasmaClient.SubmitPlasmaBlock(root)
}

// TODO: This function should be moved to loomchain/builtin/plasma_cash when the Oracle is
//       integrated into loomchain.
// Computes the block number of the next non-deposit Plasma block.
// The current Plasma block number can be for a deposit or non-deposit Plasma block.
// Plasma block numbers of non-deposit blocks are expected to be multiples of the specified interval.
func nextPlasmaBlockNum(current *big.Int, interval *big.Int) *big.Int {
	if current.Cmp(new(big.Int)) == 0 {
		return new(big.Int).Set(interval)
	}
	if current.Cmp(interval) == 0 {
		return new(big.Int).Add(current, interval)
	}
	r := new(big.Int).Add(current, new(big.Int).Sub(interval, big.NewInt(1)))
	r.Div(r, interval)        // (current + (interval - 1)) / interval
	return r.Mul(r, interval) // ((current + (interval - 1)) / interval) * interval
}
