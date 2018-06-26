package client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Client struct {
	childChain         ChainServiceClient
	RootChain          RootChainClient
	TokenContract      TokenContract
	childBlockInterval int64
}

// Token Functions

// Register a new player and grant 5 cards, for demo purposes
func (c *Client) Register() {
	c.TokenContract.Register()
}

// Deposit happens by a use calling the erc721 token contract
func (c *Client) Deposit(tokenId int64) {
	c.TokenContract.Deposit(tokenId)
}

// Plasma Functions

//Placeholder
func Transaction(slot uint64, prevTxBlkNum int64, domination uint32, address string) Tx {
	panic(address)
	return &LoomTx{Slot: slot,
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
		panic("TODO")

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
	fmt.Printf("NOT exiting deposit transaction\n")

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
	fmt.Printf("exitingTx-%x exitingTxProof-%x\n", exitingTx, exitingTxProof)
	fmt.Printf("prevTx-%v prevTxProof-%v\n", prevTx, prevTxProof)
	fmt.Printf("prevTxBlkNum-%d-txBlkNum-%d\n", prevTxBlkNum, txBlkNum)
	fmt.Printf("exitingTxIncBlock MOD childBlockInterval %d\n", txBlkNum%1000)
	sig := exitingTx.Sig()
	fmt.Printf("exitingTx.Sig() %x\n", sig)
	fmt.Printf("len(sig)- %d\n", len(sig))
	fmt.Printf("byte 0(sig)- %d\n", sig[0])
	fmt.Printf("prevTx-Owner -%v\n", prevTx.NewOwner().String())
	fmt.Printf("len(exitingTxProof)-%d\n", len(exitingTxProof))

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
		proof)
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
		proof)
	return txHash, err
}

// ChallengeAfter - `Exit Spent Coin Challenge`: Challenge an exit with a spend
// after the exit's blocks
func (c *Client) ChallengeAfter(slot uint64, challengingBlockNumber int64) ([]byte, error) { //
	challengingTx, proof, err := c.getTxAndProof(challengingBlockNumber,
		slot)
	if err != nil {
		return nil, err
	}

	txHash, err := c.RootChain.ChallengeAfter(
		slot, challengingBlockNumber,
		challengingTx,
		proof)
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

func (c *Client) PlasmaCoin(slot uint64) (*PlasmaCoin, error) {
	return c.RootChain.PlasmaCoin(slot)
}

func (c *Client) DebugCoinMetaData() {
	c.RootChain.DebugCoinMetaData()
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
	fmt.Printf("newowner -%s\n", ethAddress)

	tx := &LoomTx{
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

func (c *Client) getTxAndProof(blkHeight int64, slot uint64) (Tx, []byte, error) {
	block, err := c.childChain.Block(blkHeight)
	if err != nil {
		return nil, nil, err
	}
	tx, err := block.TxFromSlot(slot)
	if err != nil {
		return nil, nil, err
	}

	return tx, tx.Proof(), nil
}

func (c *Client) GetBlockNumber() (int64, error) {
	return c.childChain.BlockNumber()
}

func (c *Client) GetBlock(blkHeight int64) (Block, error) {
	return c.childChain.Block(blkHeight)
	//return rlp.decode(utils.decode_hex(block), Block)
}

/*
//These methods exist in python but are unused so we dont need them
func (c *Client) CurrentBlock() (Block, error) {
	return c.childChain.CurrentBlock()
	//	return rlp.decode(utils.decode_hex(block), Block)
}

func (c *Client) Proof(blkHeight int64, slot uint64) (Proof, error) {
	return c.childChain.Proof(blkHeight, slot)
	//	return base64.b64decode(c.childChain.get_proof(blknum, slot))
}
*/

func NewClient(childChainServer ChainServiceClient, rootChain RootChainClient, tokenContract TokenContract) *Client {
	return &Client{childChain: childChainServer, childBlockInterval: 1000, RootChain: rootChain, TokenContract: tokenContract}
}
