package oracle

import (
	"context"
	"crypto/ecdsa"
	"ethcontract"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	loom "github.com/loomnetwork/go-loom"
	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
	ltypes "github.com/loomnetwork/go-loom/types"
	"github.com/pkg/errors"
)

type EthPlasmaClientConfig struct {
	// URI of an Ethereum node
	EthereumURI string
	// Plasma contract address on Ethereum
	PlasmaHexAddress string
	// Private key that should be used to sign txs sent to Ethereum
	PrivateKey *ecdsa.PrivateKey
	// Override default gas computation when sending txs to Ethereum
	OverrideGas bool
	// How often Ethereum should be polled for mined txs (defaults to 10 secs).
	TxPollInterval time.Duration
	// Maximum amount of time to way for a tx to be mined by Ethereum (defaults to 5 mins).
	TxTimeout time.Duration
}

type EthPlasmaClient interface {
	Init() error
	CurrentPlasmaBlockNum() (*big.Int, error)
	LatestEthBlockNum() (uint64, error)
	// SubmitPlasmaBlock will submit a Plasma block to Ethereum and wait until the tx is confirmed.
	// The maximum wait time can be specified via the TxTimeout option in the client config.
	SubmitPlasmaBlock(plasmaBlockNum *big.Int, merkleRoot [32]byte) error
	FetchDeposits(startBlock, endBlock uint64) ([]*pctypes.DepositRequest, error)
}

type EthPlasmaClientImpl struct {
	EthPlasmaClientConfig
	ethClient      *ethclient.Client
	plasmaContract *ethcontract.RootChain
}

func (c *EthPlasmaClientImpl) Init() error {
	var err error
	c.ethClient, err = ethclient.Dial(c.EthereumURI)
	if err != nil {
		return errors.Wrap(err, "failed to connect to Ethereum")
	}

	c.plasmaContract, err = ethcontract.NewRootChain(common.HexToAddress(c.PlasmaHexAddress), c.ethClient)
	if err != nil {
		return errors.Wrap(err, "failed to bind Plasma Solidity contract")
	}
	return nil
}

func (c *EthPlasmaClientImpl) CurrentPlasmaBlockNum() (*big.Int, error) {
	curEthPlasmaBlockNum, err := c.plasmaContract.CurrentBlock(nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain current plasma block from Ethereum")
	}
	return curEthPlasmaBlockNum, nil
}

func (c *EthPlasmaClientImpl) LatestEthBlockNum() (uint64, error) {
	blockHeader, err := c.ethClient.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		return 0, err
	}
	return blockHeader.Number.Uint64(), nil
}

// SubmitPlasmaBlock will submit a Plasma block to Ethereum and wait until the tx is confirmed.
func (c *EthPlasmaClientImpl) SubmitPlasmaBlock(blockNum *big.Int, merkleRoot [32]byte) error {
	failMsg := "failed to submit plasma block to Ethereum"
	auth := bind.NewKeyedTransactor(c.PrivateKey)
	if c.OverrideGas {
		auth.GasPrice = big.NewInt(20000)
		auth.GasLimit = uint64(3141592)
	}
	tx, err := c.plasmaContract.SubmitBlock(auth, merkleRoot)
	if err != nil {
		return errors.Wrap(err, failMsg)
	}
	receipt, err := c.waitForTxReceipt(context.TODO(), tx)
	if err != nil {
		return err
	}
	if receipt.Status == 0 {
		return errors.New(failMsg)
	}
	return nil
}

// FetchDeposits fetches all deposit events from an Ethereum node from startBlock to endBlock (inclusive).
func (c *EthPlasmaClientImpl) FetchDeposits(startBlock, endBlock uint64) ([]*pctypes.DepositRequest, error) {
	// NOTE: Currently either all blocks from w.StartBlock are processed successfully or none are.
	filterOpts := &bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	}
	deposits := []*pctypes.DepositRequest{}

	it, err := c.plasmaContract.FilterDeposit(filterOpts, nil, nil)
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

// waitForTxReceipt waits for a tx to be confirmed.
// It stops waiting if the context is canceled, or the tx hasn't been confirmed after TxTimeout.
func (c *EthPlasmaClientImpl) waitForTxReceipt(ctx context.Context, tx *etypes.Transaction) (*etypes.Receipt, error) {
	timeout := c.TxTimeout
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	interval := c.TxPollInterval
	if interval == 0 {
		interval = 10 * time.Second
	}

	timer := time.NewTimer(timeout)
	ticker := time.NewTicker(interval)

	defer timer.Stop()
	defer ticker.Stop()

	txHash := tx.Hash()
	for {
		receipt, err := c.ethClient.TransactionReceipt(ctx, txHash)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve tx receipt")
		}
		if receipt != nil {
			return receipt, nil
		}
		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			return nil, errors.New("timed out waiting for tx receipt")
		case <-ticker.C:
		}
	}
}
