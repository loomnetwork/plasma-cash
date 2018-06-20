package client

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
func rlpEncode(i interface{}, t string) Tx {
	return &LoomTx{}
}

//Placeholder
func Transaction(slot uint64, prevTxBlkNum int64, domination int64, address string) Tx {
	return &LoomTx{}
}

func (c *Client) StartExit(slot uint64, prevTxBlkNum int64, txBlkNum int64) ([]byte, error) {
	// As a user, you declare that you want to exit a coin at slot `slot`
	//at the state which happened at block `txBlkNum` and you also need to
	// reference a previous block

	// TODO The actual proof information should be passed to a user from its
	// previous owners, this is a hacky way of getting the info from the
	// operator which sould be changed in the future after the exiting
	// process is more standardized
	var txHash []byte

	if txBlkNum%c.childBlockInterval != 0 {
		// In case the sender is exiting a Deposit transaction, they should
		// just create a signed transaction to themselves. There is no need
		// for a merkle proof.

		account, err := c.TokenContract.Account()
		if err != nil {
			return nil, err
		}
		// prev_block = 0 , denomination = 1
		exitingTx := Transaction(slot, 0, 1, account.Address)
		//		exiting_tx.sign(c.key)  //????
		txHash, err = c.RootChain.StartExit(
			slot,
			nil, rlpEncode(exitingTx, "UnsignedTransaction"),
			nil, nil, //proofs?
			exitingTx.Sig(),
			0, txBlkNum)
		if err != nil {
			return nil, err
		}
		return txHash, nil

	}

	// Otherwise, they should get the raw tx info from the block
	// And the merkle proof and submit these
	exitingTx, exitingTxProof, err := c.getTxAndProof(txBlkNum,
		slot)
	if err != nil {
		return nil, err
	}
	prevTx, prevTxProof, err := c.getTxAndProof(prevTxBlkNum,
		slot)
	if err != nil {
		return nil, err
	}

	txHash, err = c.RootChain.StartExit(
		slot,
		rlpEncode(prevTx, "UnsignedTransaction"),
		rlpEncode(exitingTx, "UnsignedTransaction"),
		prevTxProof, exitingTxProof,
		exitingTx.Sig(),
		prevTxBlkNum, txBlkNum)
	return txHash, nil

}

func (c *Client) ChallengeBefore() { //slot, prev_txBlkNum, txBlkNum
	/*
	   		block = c.get_block(txBlkNum)
	           tx = block.get_tx_by_uid(slot)
	           tx_proof = c.get_proof(txBlkNum, slot)

	           # If the referenced transaction is a deposit transaction then no need
	           prev_tx = Transaction(0, 0, 0,
	                                 0x0000000000000000000000000000000000000000)
	           prev_tx_proof = '0x0000000000000000'
	           if prev_txBlkNum % c.child_block_interval == 0:
	               prev_block = c.get_block(prev_txBlkNum)
	               prev_tx = prev_block.get_tx_by_uid(slot)
	               prev_tx_proof = c.get_proof(prev_txBlkNum, slot)

	           c.RootChain.challenge_before(
	               slot, rlp.encode(prev_tx, UnsignedTransaction),
	               rlp.encode(tx, UnsignedTransaction), prev_tx_proof,
	               tx_proof, tx.sig, prev_txBlkNum, txBlkNum
	           )
	   		return
	*/
}

// RespondChallengeBefore - Respond to an exit with invalid history challenge by proving that
// you were given the coin under question
func (c *Client) RespondChallengeBefore() { //slot, challenging_block_number

	/*
	   block = c.get_block(challenging_block_number)
	   challenging_tx = block.get_tx_by_uid(slot)
	   proof = c.get_proof(challenging_block_number, slot)

	   c.RootChain.respond_challenge_before(
	       slot, challenging_block_number,
	       rlp.encode(challenging_tx, UnsignedTransaction), proof
	   )
	   return
	*/
}

// ChallengeBetween - `Double Spend Challenge`: Challenge a double spend of a coin
// with a spend between the exit's blocks
func (c *Client) ChallengeBetween() { //slot, challenging_block_number
	/*
		        block = c.get_block(challenging_block_number)
		        challenging_tx = block.get_tx_by_uid(slot)
		        proof = c.get_proof(challenging_block_number, slot)

		        c.RootChain.challenge_between(
		            slot, challenging_block_number,
		            rlp.encode(challenging_tx, UnsignedTransaction), proof
		        )
		        return self
			}

		        // `Exit Spent Coin Challenge`: Challenge an exit with a spend
		        // after the exit's blocks
				func (c *Client) challenge_after(self, slot, challenging_block_number) {
		        block = c.get_block(challenging_block_number)
		        challenging_tx = block.get_tx_by_uid(slot)
		        proof = c.get_proof(challenging_block_number, slot)

		        c.RootChain.challenge_after(
		            slot, challenging_block_number,
		            rlp.encode(challenging_tx, UnsignedTransaction), proof
		        )
				return self
	*/
}

func (c *Client) FinalizeExits() {
	c.RootChain.FinalizeExits()
}

func (c *Client) Withdraw(slot uint64) {
	c.RootChain.Withdraw(slot)
}

func (c *Client) WithdrawBonds() {
	c.RootChain.WithdrawBonds()
}

func (c *Client) PlasmaCoin(slot uint64) {
	c.RootChain.PlasmaCoin(slot)
}

// Child Chain Functions

func (c *Client) SubmitBlock() error {
	return c.childChain.SubmitBlock()
}

func (c *Client) SendTransaction(slot uint64, prevBlock int64, denomination int64, newOwner string) (Tx, error) {
	return c.childChain.SendTransaction(slot, prevBlock, denomination, newOwner)
}

func (c *Client) BlockNumber() int64 {
	return c.childChain.BlockNumber()
}

func (c *Client) getTxAndProof(blknum int64, slot uint64) (Tx, Proof, error) {
	//	data = json.loads(self.child_chain.get_tx_and_proof(blknum, slot))
	//	tx = rlp.decode(utils.decode_hex(data['tx']), Transaction)
	//	proof = utils.decode_hex(data['proof'])
	return nil, &SimpleProof{}, nil
}

func (c *Client) CurrentBlock() {
	//	block = c.childChain.CurrentBlock()
	//	return rlp.decode(utils.decode_hex(block), Block)
}

func (c *Client) Block(blkHeight int) {
	//	block = c.childChain.get_block(blknum)
	//	return rlp.decode(utils.decode_hex(block), Block)
}

func (c *Client) Proof(blkHeight int, slot int) {
	//	return base64.b64decode(c.childChain.get_proof(blknum, slot))
}

func NewClient(childChainServer ChainServiceClient, rootChain RootChainClient, tokenContract TokenContract) *Client {
	return &Client{childChain: childChainServer, childBlockInterval: 1000, RootChain: rootChain, TokenContract: tokenContract}
}
