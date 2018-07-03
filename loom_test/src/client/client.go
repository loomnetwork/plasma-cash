package client

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/loomnetwork/go-loom/client/plasma_cash"
	"github.com/loomnetwork/loomchain/builtin/plugins/plasma_cash/oracle"
)

type Client struct {
	childChain         plasma_cash.ChainServiceClient
	RootChain          plasma_cash.RootChainClient
	TokenContract      plasma_cash.TokenContract
	childBlockInterval int64
	blocks             map[string]plasma_cash.Block
	plasmaEthClient    oracle.EthPlasmaClient
}

const ChildBlockInterval = 1000

// Token Functions

// Register a new player and grant 5 cards, for demo purposes
func (c *Client) Register() {
	c.TokenContract.Register()
}

// Deposit happens by a use calling the erc721 token contract
func (c *Client) Deposit(tokenID int64) common.Hash {
	txHash, err := c.TokenContract.Deposit(tokenID)
	if err != nil {
		panic(err)
	}

	//To prevent us to having to run the oracle, we are going to run the oracle manually here
	//Normally this would run as a seperate process, in future tests we can spin it up independantly
	deposits, err := c.plasmaEthClient.FetchDeposits(0, 1000)
	if err != nil {
		panic(errors.Wrap(err, "failed to fetch Plasma deposits from Ethereum"))
	}

	for _, deposit := range deposits {
		if err := c.childChain.Deposit(deposit); err != nil {
			panic(err)
		}
	}

	return txHash
}

// Plasma Functions

//Placeholder
func Transaction(slot uint64, prevTxBlkNum int64, domination uint32, address string) plasma_cash.Tx {
	return &plasma_cash.LoomTx{Slot: slot,
		PrevBlock:    big.NewInt(prevTxBlkNum),
		Denomination: domination,
		Owner:        common.HexToAddress(address), //TODO: 0x?
	}
}

func (c *Client) StartExit(slot uint64, prevTxBlkNum int64, txBlkNum int64) ([]byte, error) {
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

	if txBlkNum%c.childBlockInterval != 0 {
		// In case the sender is exiting a Deposit transaction, they should
		// just create a signed transaction to themselves. There is no need
		// for a merkle proof.
		fmt.Printf("exiting deposit transaction\n")

		// prev_block = 0 , denomination = 1
		exitingTx := Transaction(slot, 0, 1, account.Address)
		exitingTxSig, err := exitingTx.Sign(account.PrivateKey)
		if err != nil {
			return nil, err
		}

		txHash, err := c.RootChain.StartExit(
			slot,
			nil, exitingTx,
			nil, nil, //proofs?
			exitingTxSig,
			0, txBlkNum)
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

func (c *Client) ChallengeBefore(slot uint64, prevTxBlkNum int64, txBlkNum int64) ([]byte, error) {
	account, err := c.TokenContract.Account()
	if err != nil {
		return nil, err
	}

	if txBlkNum%c.childBlockInterval != 0 {
		// In case the sender is exiting a Deposit transaction, they should
		// just create a signed transaction to themselves. There is no need
		// for a merkle proof.
		panic("TODO")

		//  prev_block = 0 , denomination = 1
		exitingTx := Transaction(slot, 0, 1, account.Address)
		exitingTxSig, err := exitingTx.Sign(account.PrivateKey)
		if err != nil {
			return nil, err
		}
		txHash, err := c.RootChain.ChallengeBefore(
			slot,
			nil, exitingTx,
			nil, nil,
			exitingTxSig,
			0, txBlkNum)

		return txHash, err
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

	exitingTxSig, err := exitingTx.Sign(account.PrivateKey)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.ChallengeBefore(
		slot,
		prevTx, exitingTx,
		prevTxProof, exitingTxProof,
		exitingTxSig,
		prevTxBlkNum, txBlkNum)
	return txHash, err

}

// RespondChallengeBefore - Respond to an exit with invalid history challenge by proving that
// you were given the coin under question
func (c *Client) RespondChallengeBefore(slot uint64, challengingBlockNumber int64) ([]byte, error) {
	challengingTx, proof, err := c.getTxAndProof(challengingBlockNumber,
		slot)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.RespondChallengeBefore(slot,
		challengingBlockNumber,
		challengingTx,
		proof,
		challengingTx.Sig())
	return txHash, err
}

// ChallengeBetween - `Double Spend Challenge`: Challenge a double spend of a coin
// with a spend between the exit's blocks
func (c *Client) ChallengeBetween(slot uint64, challengingBlockNumber int64) ([]byte, error) {
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
func (c *Client) ChallengeAfter(slot uint64, challengingBlockNumber int64) ([]byte, error) { //
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

func (c *Client) FinalizeExits() error {
	return c.RootChain.FinalizeExits()
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

func (c *Client) SubmitBlock() error {
	if err := c.childChain.SubmitBlock(); err != nil {
		return err
	}

	blockNum, err := c.childChain.BlockNumber()
	if err != nil {
		return err
	}

	block, err := c.childChain.Block(blockNum)
	if err != nil {
		return err
	}

	var root [32]byte
	copy(root[:], block.MerkleHash())
	return c.RootChain.SubmitBlock(big.NewInt(blockNum), root)
}

func (c *Client) SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) error {
	ethAddress := common.HexToAddress(newOwner)

	tx := &plasma_cash.LoomTx{
		Slot:         slot,
		Denomination: uint32(denomination),
		Owner:        ethAddress,
		PrevBlock:    big.NewInt(prevBlock),
	}

	account, err := c.TokenContract.Account()
	if err != nil {
		return err
	}

	sig, err := tx.Sign(account.PrivateKey)
	if err != nil {
		return err
	}

	return c.childChain.SendTransaction(slot, prevBlock, denomination, newOwner, sig)
}

func (c *Client) getTxAndProof(blkHeight int64, slot uint64) (plasma_cash.Tx, []byte, error) {
	block, err := c.childChain.Block(blkHeight)
	if err != nil {
		return nil, nil, err
	}

	tx, err := block.TxFromSlot(slot)
	if err != nil {
		return nil, nil, err
	}

	// server should handle this
	/*
		if blkHeight%ChildBlockInterval != 0 {
			proof := []byte{00000000}
		} else {

		}
	*/

	return tx, tx.Proof(), nil
}

func (c *Client) GetBlockNumber() (int64, error) {
	return c.childChain.BlockNumber()
}

func (c *Client) GetBlock(blkHeight int64) (plasma_cash.Block, error) {
	return c.childChain.Block(blkHeight)
}

func NewClient(childChainServer plasma_cash.ChainServiceClient, rootChain plasma_cash.RootChainClient, tokenContract plasma_cash.TokenContract) *Client {
	ethPrivKeyHexStr := GetTestAccountHexKey("authority")
	ethPrivKey, err := crypto.HexToECDSA(strings.TrimPrefix(ethPrivKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	ethCfg := oracle.EthPlasmaClientConfig{
		EthereumURI:      "http://localhost:8545",
		PlasmaHexAddress: GetContractHexAddress("root_chain"),
		PrivateKey:       ethPrivKey,
		OverrideGas:      true,
	}

	pbc := &oracle.EthPlasmaClientImpl{EthPlasmaClientConfig: ethCfg}
	err = pbc.Init()
	if err != nil {
		panic(err) //todo return
	}

	return &Client{childChain: childChainServer, childBlockInterval: 1000, RootChain: rootChain, TokenContract: tokenContract, plasmaEthClient: pbc}
}
