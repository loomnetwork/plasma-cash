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

// PlasmaBlockWorker sends non-deposit Plasma block from the DAppChain to Ethereum.
type PlasmaBlockWorker struct {
	ethPlasmaClient     EthPlasmaClient
	dappPlasmaClient    DAppChainPlasmaClient
	plasmaBlockInterval uint32
}

func NewPlasmaBlockWorker(cfg *OracleConfig) *PlasmaBlockWorker {
	return &PlasmaBlockWorker{
		ethPlasmaClient:     &EthPlasmaClientImpl{EthPlasmaClientConfig: cfg.EthClientCfg},
		dappPlasmaClient:    &DAppChainPlasmaClientImpl{DAppChainPlasmaClientConfig: cfg.DAppChainClientCfg},
		plasmaBlockInterval: cfg.PlasmaBlockInterval,
	}
}

func (w *PlasmaBlockWorker) Init() error {
	if err := w.ethPlasmaClient.Init(); err != nil {
		return err
	}
	return w.dappPlasmaClient.Init()
}

func (w *PlasmaBlockWorker) Run() {
	go runWithRecovery(func() {
		loopWithInterval(w.sendPlasmaBlocksToEthereum, 5*time.Second)
	})
}

// DAppChain -> Plasma Blocks -> Ethereum
func (w *PlasmaBlockWorker) sendPlasmaBlocksToEthereum() error {
	if err := w.syncPlasmaBlocksWithEthereum(); err != nil {
		return errors.Wrap(err, "failed to sync plasma blocks with mainnet")
	}

	return w.dappPlasmaClient.FinalizeCurrentPlasmaBlock()
}

// Send any finalized but unsubmitted plasma blocks from the DAppChain to Ethereum.
func (w *PlasmaBlockWorker) syncPlasmaBlocksWithEthereum() error {
	curEthPlasmaBlockNum, err := w.ethPlasmaClient.CurrentPlasmaBlockNum()
	if err != nil {
		return err
	}
	log.Printf("solPlasma.CurrentBlock: %s", curEthPlasmaBlockNum.String())

	curLoomPlasmaBlockNum, err := w.dappPlasmaClient.CurrentPlasmaBlockNum()
	if err != nil {
		return err
	}

	if curLoomPlasmaBlockNum.Cmp(curEthPlasmaBlockNum) == 0 {
		// DAppChain and Ethereum both have all the finalized Plasma blocks
		return nil
	}

	plasmaBlockInterval := big.NewInt(int64(w.plasmaBlockInterval))
	unsubmittedPlasmaBlockNum := nextPlasmaBlockNum(curEthPlasmaBlockNum, plasmaBlockInterval)

	for {
		log.Printf("unsubmittedPlasmaBlockNum: %s", unsubmittedPlasmaBlockNum.String())
		if unsubmittedPlasmaBlockNum.Cmp(curLoomPlasmaBlockNum) > 0 {
			// All the finalized plasma blocks in the DAppChain have been submitted to Ethereum
			break
		}

		block, err := w.dappPlasmaClient.PlasmaBlockAt(unsubmittedPlasmaBlockNum)
		if err != nil {
			return err
		}

		if err := w.submitPlasmaBlockToEthereum(unsubmittedPlasmaBlockNum, block.MerkleHash); err != nil {
			return err
		}

		unsubmittedPlasmaBlockNum = nextPlasmaBlockNum(unsubmittedPlasmaBlockNum, plasmaBlockInterval)
	}

	return nil
}

// Submits a Plasma block (or rather its merkle root) to the Plasma Solidity contract on Ethereum.
// This function will block until the tx is confirmed, or times out.
func (w *PlasmaBlockWorker) submitPlasmaBlockToEthereum(plasmaBlockNum *big.Int, merkleRoot []byte) error {
	curEthPlasmaBlockNum, err := w.ethPlasmaClient.CurrentPlasmaBlockNum()
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
	return w.ethPlasmaClient.SubmitPlasmaBlock(plasmaBlockNum, root)
}

// PlasmaDepositWorker sends Plasma deposits from Ethereum to the DAppChain.
type PlasmaDepositWorker struct {
	ethPlasmaClient  EthPlasmaClient
	dappPlasmaClient DAppChainPlasmaClient
	startEthBlock    uint64 // Eth block from which the oracle should start looking for deposits
}

func NewPlasmaDepositWorker(cfg *OracleConfig) *PlasmaDepositWorker {
	return &PlasmaDepositWorker{
		ethPlasmaClient:  &EthPlasmaClientImpl{EthPlasmaClientConfig: cfg.EthClientCfg},
		dappPlasmaClient: &DAppChainPlasmaClientImpl{DAppChainPlasmaClientConfig: cfg.DAppChainClientCfg},
	}
}

func (w *PlasmaDepositWorker) Init() error {
	if err := w.ethPlasmaClient.Init(); err != nil {
		return err
	}
	return w.dappPlasmaClient.Init()
}

func (w *PlasmaDepositWorker) Run() {
	go runWithRecovery(func() {
		loopWithInterval(w.sendPlasmaDepositsToDAppChain, 5*time.Second)
	})
}

// Ethereum -> Plasma Deposits -> DAppChain
func (w *PlasmaDepositWorker) sendPlasmaDepositsToDAppChain() error {
	// TODO: get start block from Plasma Go contract, like the Transfer Gateway Oracle
	startEthBlock := w.startEthBlock

	// TODO: limit max block range per batch
	latestEthBlock, err := w.ethPlasmaClient.LatestEthBlockNum()
	if err != nil {
		return errors.Wrap(err, "failed to obtain latest Ethereum block number")
	}

	if latestEthBlock < startEthBlock {
		// Wait for Ethereum to produce a new block...
		return nil
	}

	deposits, err := w.ethPlasmaClient.FetchDeposits(startEthBlock, latestEthBlock)
	if err != nil {
		return errors.Wrap(err, "failed to fetch Plasma deposits from Ethereum")
	}

	for _, deposit := range deposits {
		if err := w.dappPlasmaClient.Deposit(deposit); err != nil {
			return err
		}
	}

	w.startEthBlock = latestEthBlock + 1
	return nil
}

type Oracle struct {
	cfg           *OracleConfig
	depositWorker *PlasmaDepositWorker
	blockWorker   *PlasmaBlockWorker
}

func NewOracle(cfg *OracleConfig) *Oracle {
	return &Oracle{
		cfg:           cfg,
		depositWorker: NewPlasmaDepositWorker(cfg),
		blockWorker:   NewPlasmaBlockWorker(cfg),
	}
}

func (orc *Oracle) Init() error {
	if err := orc.depositWorker.Init(); err != nil {
		return err
	}
	return orc.blockWorker.Init()
}

// TODO: Graceful shutdown
func (orc *Oracle) Run() {
	if orc.cfg.EthClientCfg.PrivateKey != nil {
		orc.blockWorker.Run()
	}
	orc.depositWorker.Run()
}

// runWithRecovery should run in a goroutine, it will ensure the given function keeps on running in
// a goroutine as long as it doesn't panic due to a runtime error.
func runWithRecovery(run func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic in a Plasma Oracle worker: %v\n", r)
			// Unless it's a runtime error restart the goroutine
			if _, ok := r.(runtime.Error); !ok {
				time.Sleep(30 * time.Second)
				log.Printf("Restarting Plasma Oracle worker...")
				go runWithRecovery(run)
			}
		}
	}()
	run()
}

// loopWithInterval will execute the step function in an endless loop while ensuring that each
// loop iteration takes up the minimum specified duration.
func loopWithInterval(step func() error, minStepDuration time.Duration) {
	for {
		start := time.Now()
		if err := step(); err != nil {
			log.Println(err)
		}
		diff := time.Now().Sub(start)
		if diff < minStepDuration {
			time.Sleep(minStepDuration - diff)
		}
	}
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
