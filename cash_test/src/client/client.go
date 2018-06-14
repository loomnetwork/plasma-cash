package client

type Client struct {
	/*
			        c.rootChain = rootChain
		        c.key = token_contract.account.privateKey
		        c.token_contract = token_contract
	*/
	childChain         ChainServiceClient
	rootChain          RootChainClient
	TokenContract      TokenContract
	childBlockInterval int
}

// Token Functions

// Register a new player and grant 5 cards, for demo purposes
func (c *Client) Register() {

	c.TokenContract.Register()
}

// Deposit happens by a use calling the erc721 token contract
func (c *Client) Deposit(tokenId int) {
	c.TokenContract.Deposit(tokenId)
}

// Plasma Functions

func (c *Client) StartExit() { //slot, prev_tx_blk_num, tx_blk_num
	// As a user, you declare that you want to exit a coin at slot `slot`
	//at the state which happened at block `tx_blk_num` and you also need to
	// reference a previous block

	// TODO The actual proof information should be passed to a user from its
	// previous owners, this is a hacky way of getting the info from the
	// operator which sould be changed in the future after the exiting
	// process is more standardized
	/*
	   if (tx_blk_num % c.child_block_interval != 0 and
	           prev_tx_blk_num == 0):
	       # In case the sender is exiting a Deposit transaction, they should
	       # just create a signed transaction to themselves. There is no need
	       # for a merkle proof.

	       # prevBlockehw = 0 , denomination = 1
	       exiting_tx = Transaction(slot, 0, 1,
	                                c.token_contract.account.address,
	                                incl_block=tx_blk_num)
	       exiting_tx.make_mutable()
	       exiting_tx.sign(c.key)
	       exiting_tx.make_immutable()
	       c.rootChain.start_exit(
	               slot,
	               b'0x0', rlp.encode(exiting_tx, UnsignedTransaction),
	               b'0x0', b'0x0',
	               exiting_tx.sig,
	               0, tx_blk_num
	       )
	   else:
	       # Otherwise, they should get the raw tx info from the block
	       # And the merkle proof and submit these
	       block = c.get_block(tx_blk_num)
	       exiting_tx = block.get_tx_by_uid(slot)
	       exiting_tx_proof = c.get_proof(tx_blk_num, slot)

	       prev_block = c.get_block(prev_tx_blk_num)
	       prev_tx = prev_block.get_tx_by_uid(slot)
	       if (prev_tx_blk_num % c.child_block_interval != 0):
	           # After 1 off-chain transfer, referencing a deposit
	           # transaction, no need for proof
	           prev_tx_proof = b'0x0000000000000000'
	       else:
	           prev_tx_proof = c.get_proof(prev_tx_blk_num, slot)
	       c.rootChain.start_exit(
	               slot,
	               rlp.encode(prev_tx, UnsignedTransaction),
	               rlp.encode(exiting_tx, UnsignedTransaction),
	               prev_tx_proof, exiting_tx_proof,
	               exiting_tx.sig,
	               prev_tx_blk_num, tx_blk_num
	       )
	*/
	return
}

func (c *Client) ChallengeBefore() { //slot, prev_tx_blk_num, tx_blk_num
	/*
	   		block = c.get_block(tx_blk_num)
	           tx = block.get_tx_by_uid(slot)
	           tx_proof = c.get_proof(tx_blk_num, slot)

	           # If the referenced transaction is a deposit transaction then no need
	           prev_tx = Transaction(0, 0, 0,
	                                 0x0000000000000000000000000000000000000000)
	           prev_tx_proof = '0x0000000000000000'
	           if prev_tx_blk_num % c.child_block_interval == 0:
	               prev_block = c.get_block(prev_tx_blk_num)
	               prev_tx = prev_block.get_tx_by_uid(slot)
	               prev_tx_proof = c.get_proof(prev_tx_blk_num, slot)

	           c.rootChain.challenge_before(
	               slot, rlp.encode(prev_tx, UnsignedTransaction),
	               rlp.encode(tx, UnsignedTransaction), prev_tx_proof,
	               tx_proof, tx.sig, prev_tx_blk_num, tx_blk_num
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

	   c.rootChain.respond_challenge_before(
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

		        c.rootChain.challenge_between(
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

		        c.rootChain.challenge_after(
		            slot, challenging_block_number,
		            rlp.encode(challenging_tx, UnsignedTransaction), proof
		        )
				return self
	*/
}

func (c *Client) FinalizeExits() {
	c.rootChain.FinalizeExits()
}

func (c *Client) Withdraw(slot int) {
	c.rootChain.Withdraw(slot)
}

func (c *Client) WithdrawBonds() {
	c.rootChain.WithdrawBonds()
}

func (c *Client) PlasmaCoin(slot int) {
	c.rootChain.PlasmaCoin(slot)
}

// Child Chain Functions

func (c *Client) SubmitBlock() {
	/*
		block = c.GetCurrentBlock()
		block.make_mutable() // mutex for mutability?
		block.sign(c.key)
		block.make_immutable()
		return c.childChain.submit_block(rlp.encode(block, Block).hex())
	*/
}

func (c *Client) SendTransaction(slot int, prevBlock int, denomination int, newOwner string) *Tx {
	/*
		        new_owner = utils.normalize_address(new_owner)
		        incl_block = c.BlockNumber()
		        tx = Transaction(slot, prev_block, denomination, new_owner,
		                         incl_block=incl_block)
		        tx.make_mutable()
		        tx.sign(c.key)
		        tx.make_immutable()
		        c.childChain.SendTransaction(rlp.encode(tx, Transaction).hex())
				return tx
	*/
	return &Tx{}
}

func (c *Client) BlockNumber() int {
	return c.childChain.BlockNumber()
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
	return &Client{childChain: childChainServer, childBlockInterval: 1000, rootChain: rootChain, TokenContract: tokenContract}
}
