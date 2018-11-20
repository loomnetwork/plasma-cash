package client

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"

	"github.com/loomnetwork/go-loom/client/plasma_cash"
	"github.com/loomnetwork/go-loom/client/plasma_cash/eth"
)

func PollForBlockChange(c *Client, currentBlockNumber *big.Int, maxIteration int, sleepPerIteration time.Duration) (*big.Int, error) {
	currentIteration := 0
	var err error
	var updatedBlockNumber = big.NewInt(0)
	for {
		time.Sleep(sleepPerIteration)

		fmt.Printf("Polling, CurrentBlockNumber: %s\n", currentBlockNumber.String())

		updatedBlockNumber, err = c.GetBlockNumber()
		if err != nil {
			err = fmt.Errorf("Error while polling for current block number")
			break
		}

		if updatedBlockNumber.Cmp(currentBlockNumber) != 0 {
			fmt.Printf("Updated BlockNumber to: %s\n", updatedBlockNumber.String())
			break
		}

		currentIteration += 1
		if currentIteration >= maxIteration {
			err = fmt.Errorf("Maximum iteration exceeded but, block didnt change")
			break
		}
	}

	return updatedBlockNumber, err
}

type Client struct {
	childChain         plasma_cash.ChainServiceClient
	RootChain          plasma_cash.RootChainClient
	TokenContract      plasma_cash.TokenContract
	childBlockInterval int64
	blocks             map[string]plasma_cash.Block
	plasmaEthClient    eth.EthPlasmaClient
}

const ChildBlockInterval = 1000

// Token Functions

// Register a new player and grant 5 cards, for demo purposes
func (c *Client) Register() {
	c.TokenContract.Register()
}

// Deposit happens by a use calling the erc721 token contract
func (c *Client) Deposit(tokenID *big.Int) common.Hash {
	txHash, err := c.TokenContract.Deposit(tokenID)
	if err != nil {
		panic(err)
	}
	return txHash
}

// Plasma Functions

func Transaction(slot uint64, prevTxBlkNum *big.Int, denomination *big.Int, address string) plasma_cash.Tx {
	return &plasma_cash.LoomTx{
		Slot:         slot,
		PrevBlock:    prevTxBlkNum,
		Denomination: denomination,
		Owner:        common.HexToAddress(address), //TODO: 0x?
	}
}

func (c *Client) StartExit(slot uint64, prevTxBlkNum *big.Int, txBlkNum *big.Int) ([]byte, error) {
	// As a user, you declare that you want to exit a coin at slot `slot`
	//at the state which happened at block `txBlkNum` and you also need to
	// reference a previous block

	// TODO The actual proof information should be passed to a user from its
	// previous owners, this is a hacky way of getting the info from the
	// operator which sould be changed in the future after the exiting
	// process is more standardized
	account, err := c.TokenContract.Account()
	if err != nil {
		return nil, err
	}

	blkModInterval := new(big.Int)
	blkModInterval = blkModInterval.Mod(txBlkNum, big.NewInt(c.childBlockInterval))
	if blkModInterval.Cmp(big.NewInt(0)) != 0 {
		// In case the sender is exiting a Deposit transaction, they should
		// just create a signed transaction to themselves. There is no need
		// for a merkle proof.
		fmt.Printf("exiting deposit transaction\n")

		// prev_block = 0 , denomination = 1
		exitingTx := Transaction(slot, big.NewInt(0), big.NewInt(1), account.Address)
		exitingTxSig, err := exitingTx.Sign(account.PrivateKey)
		if err != nil {
			return nil, err
		}

		txHash, err := c.RootChain.StartExit(
			slot,
			nil, exitingTx,
			nil, nil, //proofs?
			exitingTxSig,
			big.NewInt(0), txBlkNum)
		if err != nil {
			return nil, err
		}
		return txHash, nil
	}

	// Otherwise, they should get the raw tx info from the block
	// And the merkle proof and submit these
	exitingTx, exitingTxProof, err := c.getTxAndProof(txBlkNum, slot)
	if err != nil {
		return nil, err
	}

	prevTx, prevTxProof, err := c.getTxAndProof(prevTxBlkNum, slot)
	if err != nil {
		return nil, err
	}
	sig := exitingTx.Sig()

	return c.RootChain.StartExit(
		slot,
		prevTx, exitingTx,
		prevTxProof, exitingTxProof,
		sig,
		prevTxBlkNum, txBlkNum)
}

func (c *Client) ChallengeBefore(slot uint64, txBlkNum *big.Int) ([]byte, error) {
	account, err := c.TokenContract.Account()
	if err != nil {
		return nil, err
	}

	blkModInterval := new(big.Int)
	blkModInterval = blkModInterval.Mod(txBlkNum, big.NewInt(c.childBlockInterval))
	if blkModInterval.Cmp(big.NewInt(0)) != 0 {
		// If the client is challenging an exit with a deposit they can create a signed transaction themselves.
		// There is no need for a merkle proof.
		exitingTx := Transaction(slot, big.NewInt(0), big.NewInt(1), account.Address)
		exitingTxSig, err := exitingTx.Sign(account.PrivateKey)
		if err != nil {
			return nil, err
		}

		txHash, err := c.RootChain.ChallengeBefore(
			slot,
			exitingTx,
			nil,
			exitingTxSig,
			txBlkNum)
		if err != nil {
			return nil, err
		}
		return txHash, nil
	}

	// Otherwise, they should get the raw tx info from the block
	// And the merkle proof and submit these
	exitingTx, exitingTxProof, err := c.getTxAndProof(txBlkNum, slot)
	if err != nil {
		return nil, err
	}

	exitingTxSig, err := exitingTx.Sign(account.PrivateKey)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.ChallengeBefore(
		slot,
		exitingTx,
		exitingTxProof,
		exitingTxSig,
		txBlkNum)
	return txHash, err

}

// RespondChallengeBefore - Respond to an exit with invalid history challenge by proving that
// you were given the coin under question
func (c *Client) RespondChallengeBefore(slot uint64, respondingBlockNumber *big.Int, challengingTxHash [32]byte) ([]byte, error) {
	respondingTx, proof, err := c.getTxAndProof(respondingBlockNumber, slot)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.RespondChallengeBefore(slot,
		challengingTxHash,
		respondingBlockNumber,
		respondingTx,
		proof,
		respondingTx.Sig())
	return txHash, err
}

// ChallengeBetween - `Double Spend Challenge`: Challenge a double spend of a coin
// with a spend between the exit's blocks
func (c *Client) ChallengeBetween(slot uint64, challengingBlockNumber *big.Int) ([]byte, error) {
	challengingTx, proof, err := c.getTxAndProof(challengingBlockNumber, slot)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.ChallengeBetween(
		slot,
		challengingBlockNumber,
		challengingTx,
		proof,
		challengingTx.Sig(),
	)
	return txHash, err
}

// ChallengeAfter - `Exit Spent Coin Challenge`: Challenge an exit with a spend
// after the exit's blocks
func (c *Client) ChallengeAfter(slot uint64, challengingBlockNumber *big.Int) ([]byte, error) { //
	fmt.Printf("Challenege after getting block-%d - slot %d\n", challengingBlockNumber, slot)
	challengingTx, proof, err := c.getTxAndProof(challengingBlockNumber,
		slot)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.ChallengeAfter(
		slot, challengingBlockNumber,
		challengingTx,
		proof,
		challengingTx.Sig())
	return txHash, err
}

func (c *Client) CancelExit(slot uint64) error {
	return c.RootChain.CancelExit(slot)
}

func (c *Client) CancelExits(slots []uint64) error {
	return c.RootChain.CancelExits(slots)
}

func (c *Client) FinalizeExit(slot uint64) error {
	return c.RootChain.FinalizeExit(slot)
}

func (c *Client) FinalizeExits(slots []uint64) error {
	return c.RootChain.FinalizeExits(slots)
}

func (c *Client) Withdraw(slot uint64) error {
	return c.RootChain.Withdraw(slot)
}

func (c *Client) WithdrawBonds() error {
	return c.RootChain.WithdrawBonds()
}

func (c *Client) PlasmaCoin(slot uint64) (*plasma_cash.PlasmaCoin, error) {
	return c.RootChain.PlasmaCoin(slot)
}

func (c *Client) DebugCoinMetaData(slots []uint64) {
	c.RootChain.DebugCoinMetaData(slots)
}

// Child Chain Functions

func (c *Client) SendTransaction(slot uint64, prevBlock *big.Int, denomination *big.Int, newOwner string) error {
	ethAddress := common.HexToAddress(newOwner)

	tx := &plasma_cash.LoomTx{
		Slot:         slot,
		Denomination: denomination,
		Owner:        ethAddress,
		PrevBlock:    prevBlock,
	}

	account, err := c.TokenContract.Account()
	if err != nil {
		return err
	}

	sig, err := tx.Sign(account.PrivateKey)
	if err != nil {
		return err
	}

	return c.childChain.SendTransaction(slot, prevBlock, denomination, newOwner, account.Address, sig)
}

func (c *Client) getTxAndProof(blkHeight *big.Int, slot uint64) (plasma_cash.Tx, []byte, error) {
	tx, err := c.childChain.PlasmaTx(blkHeight, slot)
	if err != nil {
		return nil, nil, err
	}
	return tx, tx.Proof(), nil
}

func (c *Client) WatchExits(slot uint64) error {
	panic("TODO")
}

func (c *Client) StopWatchingExits(slot uint64) error {
	panic("TODO")
}

func (c *Client) GetBlockNumber() (*big.Int, error) {
	return c.childChain.BlockNumber()
}

func (c *Client) GetBlock(blkHeight *big.Int) (plasma_cash.Block, error) {
	return c.childChain.Block(blkHeight)
}

func NewClient(cfg *viper.Viper, childChainServer plasma_cash.ChainServiceClient, rootChain plasma_cash.RootChainClient, tokenContract plasma_cash.TokenContract) *Client {
	ethPrivKeyHexStr := cfg.GetString("authority")
	ethPrivKey, err := crypto.HexToECDSA(strings.TrimPrefix(ethPrivKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	ethCfg := eth.EthPlasmaClientConfig{
		EthereumURI:      "http://localhost:8545",
		PlasmaHexAddress: cfg.GetString("root_chain"),
		PrivateKey:       ethPrivKey,
		OverrideGas:      true,
	}

	pbc := eth.NewEthPlasmaClient(ethCfg)
	err = pbc.Init()
	if err != nil {
		panic(err) //todo return
	}

	return &Client{childChain: childChainServer, childBlockInterval: 1000, RootChain: rootChain, TokenContract: tokenContract, plasmaEthClient: pbc}
}
