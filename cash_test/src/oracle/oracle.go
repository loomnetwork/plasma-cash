package oracle

import (
	"context"
	"crypto/ecdsa"
	"ethcontract"
	"log"
	"math/big"
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
	// Each Plasma block number must be a multiple of this value
	PlasmaBlockInterval uint32
	// URI of an Ethereum node
	EthereumURI string
	// Plasma contract address on Ethereum
	PlasmaHexAddress string
	ChainID          string
	WriteURI         string
	ReadURI          string
	// Used to sign txs sent to Loom DAppChain
	Signer auth.Signer
	// Private key that should be used to sign txs sent to Ethereum
	EthPrivateKey *ecdsa.PrivateKey
	// Override default gas computation when sending txs to Ethereum
	OverrideGas bool
}

type Oracle struct {
	cfg           OracleConfig
	solPlasma     *ethcontract.RootChain
	goPlasma      *client.Contract
	startEthBlock uint64 // Eth block from which the oracle should start looking for deposits
	ethClient     *ethclient.Client
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

		if orc.cfg.EthPrivateKey != nil {
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

	breq := &pctypes.SubmitBlockToMainnetRequest{}
	if _, err := orc.goPlasma.Call("SubmitBlockToMainnet", breq, orc.cfg.Signer, nil); err != nil {
		return errors.Wrap(err, "failed to commit SubmitBlockToMainnet tx")
	}

	return nil
}

// Ethereum -> Plasma Deposits -> DAppChain
func (orc *Oracle) sendPlasmaDepositsToDAppChain() error {
	// TODO: get start block from Plasma Go contract, like the Transfer Gateway Oracle
	startEthBlock := orc.startEthBlock

	// TODO: limit max block range per batch
	latestEthBlock, err := orc.getLatestEthBlockNumber()
	if err != nil {
		return errors.Wrap(err, "failed to obtain latest Ethereum block number")
	}

	if latestEthBlock < startEthBlock {
		// Wait for Ethereum to produce a new block...
		return nil
	}

	deposits, err := orc.fetchDeposits(startEthBlock, latestEthBlock)
	if err != nil {
		return errors.Wrap(err, "failed to fetch Plasma deposits from Ethereum")
	}

	for _, deposit := range deposits {
		if _, err := orc.goPlasma.Call("DepositRequest", deposit, orc.cfg.Signer, nil); err != nil {
			return errors.Wrap(err, "failed to commit DepositRequest tx")
		}
	}

	orc.startEthBlock = latestEthBlock + 1
	return nil
}

// Send any finalized but unsubmitted plasma blocks from the DAppChain to Ethereum.
func (orc *Oracle) syncPlasmaBlocksWithEthereum() error {
	curEthPlasmaBlockNum, err := orc.solPlasma.CurrentBlock(nil)
	if err != nil {
		return errors.Wrap(err, "failed to obtain current plasma block from Ethereum")
	}
	log.Printf("solPlasma.CurrentBlock: %s", curEthPlasmaBlockNum.String())

	req := &pctypes.GetCurrentBlockRequest{}
	resp := &pctypes.GetCurrentBlockResponse{}
	caller := loom.Address{
		ChainID: orc.cfg.ChainID,
		Local:   loom.LocalAddressFromPublicKey(orc.cfg.Signer.PublicKey()),
	}
	if _, err := orc.goPlasma.StaticCall("GetCurrentBlockRequest", req, caller, resp); err != nil {
		return errors.Wrap(err, "failed to call GetCurrentBlockRequest")
	}
	curLoomPlasmaBlockNum := resp.BlockHeight.Value.Int

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

		req := &pctypes.GetBlockRequest{
			BlockHeight: &ltypes.BigUInt{Value: *loom.NewBigUInt(unsubmittedPlasmaBlockNum)},
		}
		resp := &pctypes.GetBlockResponse{}
		if _, err := orc.goPlasma.StaticCall("GetBlockRequest", req, caller, resp); err != nil {
			return errors.Wrap(err, "failed to obtain plasma block from DAppChain")
		}
		if resp.Block == nil {
			return errors.New("DAppChain returned empty plasma block")
		}

		// TODO: Will this block until the tx is confirmed? If not then should wait until it is
		//       confirmed before submitting another block.
		if err := orc.submitPlasmaBlockToEthereum(unsubmittedPlasmaBlockNum, resp.Block.MerkleHash); err != nil {
			return errors.Wrap(err, "failed to submit plasma block to Ethereum")
		}

		unsubmittedPlasmaBlockNum = nextPlasmaBlockNum(unsubmittedPlasmaBlockNum, plasmaBlockInterval)
	}

	return nil
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
				// TODO: store ev.Hash... it's always a hash of ev.Slot, so a bit redundant
			})
		} else {
			err := it.Error()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get event data for Plasma deposit")
			}
			it.Close()
			break
		}
	}

	return deposits, nil
}

// Submits a Plasma block (or rather its merkle root) to the Plasma Solidity contract on Ethereum
func (orc *Oracle) submitPlasmaBlockToEthereum(plasmaBlockNum *big.Int, merkleRoot []byte) error {
	curEthPlasmaBlockNum, err := orc.solPlasma.CurrentBlock(nil)
	if err != nil {
		return errors.Wrap(err, "failed to obtain current plasma block from Ethereum")
	}

	// Try to avoid submitting the same plasma blocks multiple times
	if plasmaBlockNum.Cmp(curEthPlasmaBlockNum) <= 0 {
		return nil
	}

	if len(merkleRoot) != 32 {
		return errors.New("invalid merkle root size")
	}
	auth := bind.NewKeyedTransactor(orc.cfg.EthPrivateKey)
	if orc.cfg.OverrideGas {
		auth.GasPrice = big.NewInt(20000)
		auth.GasLimit = uint64(3141592)
	}

	var root [32]byte
	copy(root[:], merkleRoot)
	_, err = orc.solPlasma.SubmitBlock(auth, root)
	return err
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
