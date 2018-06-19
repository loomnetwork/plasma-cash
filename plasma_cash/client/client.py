import rlp
from ethereum import utils
from child_chain.block import Block
from child_chain.transaction import Transaction, UnsignedTransaction
from .child_chain_service import ChildChainService
import json


class Client(object):

    def __init__(self,
                 root_chain,
                 token_contract,
                 child_chain=ChildChainService('http://localhost:8546')):
        self.root_chain = root_chain
        self.key = token_contract.account.privateKey
        self.token_contract = token_contract
        self.child_chain = child_chain
        self.child_block_interval = 1000
        self.proofs = {}

    # Token Functions

    def register(self):
        ''' Register a new player and grant 5 cards, for demo purposes'''
        tx_hash = self.token_contract.register()
        return tx_hash

    def deposit(self, tokenId):
        ''' Deposit happens by a use calling the erc721 token contract '''
        tx_hash = self.token_contract.deposit(tokenId)
        return tx_hash

    # Plasma Functions

    def start_exit(self, slot, prev_tx_blk_num, tx_blk_num):
        '''
        As a user, you declare that you want to exit a coin at slot `slot`
        at the state which happened at block `tx_blk_num` and you also need to
        reference a previous block
        '''
        # TODO The actual proof information should be passed to a user from its
        # previous owners, this is a hacky way of getting the info from the
        # operator which sould be changed in the future after the exiting
        # process is more standardized

        if (tx_blk_num % self.child_block_interval != 0):
            # In case the sender is exiting a Deposit transaction, they should
            # just create a signed transaction to themselves. There is no need
            # for a merkle proof.

            # prev_block = 0 , denomination = 1
            exiting_tx = Transaction(slot, 0, 1,
                                     self.token_contract.account.address)
            exiting_tx.make_mutable()
            exiting_tx.sign(self.key)
            exiting_tx.make_immutable()
            tx_hash = self.root_chain.start_exit(
                    slot,
                    b'0x0', rlp.encode(exiting_tx, UnsignedTransaction),
                    b'0x0', b'0x0',
                    exiting_tx.sig,
                    0, tx_blk_num)
        else:
            # Otherwise, they should get the raw tx info from the block
            # And the merkle proof and submit these
            exiting_tx, exiting_tx_proof = self.get_tx_and_proof(tx_blk_num,
                                                                 slot)
            prev_tx, prev_tx_proof = self.get_tx_and_proof(prev_tx_blk_num,
                                                           slot)

            tx_hash = self.root_chain.start_exit(
                    slot,
                    rlp.encode(prev_tx, UnsignedTransaction),
                    rlp.encode(exiting_tx, UnsignedTransaction),
                    prev_tx_proof, exiting_tx_proof,
                    exiting_tx.sig,
                    prev_tx_blk_num, tx_blk_num)
        return tx_hash

    def challenge_before(self, slot, prev_tx_blk_num, tx_blk_num):
        if (tx_blk_num % self.child_block_interval != 0):
            # In case the sender is exiting a Deposit transaction, they should
            # just create a signed transaction to themselves. There is no need
            # for a merkle proof.

            # prev_block = 0 , denomination = 1
            exiting_tx = Transaction(slot, 0, 1,
                                     self.token_contract.account.address)
            exiting_tx.make_mutable()
            exiting_tx.sign(self.key)
            exiting_tx.make_immutable()
            tx_hash = self.root_chain.challenge_before(
                    slot,
                    b'0x0', rlp.encode(exiting_tx, UnsignedTransaction),
                    b'0x0', b'0x0',
                    exiting_tx.sig,
                    0, tx_blk_num)
        else:
            # Otherwise, they should get the raw tx info from the block
            # And the merkle proof and submit these
            exiting_tx, exiting_tx_proof = self.get_tx_and_proof(tx_blk_num,
                                                                 slot)
            prev_tx, prev_tx_proof = self.get_tx_and_proof(prev_tx_blk_num,
                                                           slot)

            tx_hash = self.root_chain.challenge_before(
                    slot,
                    rlp.encode(prev_tx, UnsignedTransaction),
                    rlp.encode(exiting_tx, UnsignedTransaction),
                    prev_tx_proof, exiting_tx_proof,
                    exiting_tx.sig,
                    prev_tx_blk_num, tx_blk_num)
        return tx_hash

    def respond_challenge_before(self, slot, challenging_block_number):
        '''
        Respond to an exit with invalid history challenge by proving that
        you were given the coin under question
        '''
        challenging_tx, proof = self.get_tx_and_proof(challenging_block_number,
                                                      slot)

        tx_hash = self.root_chain.respond_challenge_before(
            slot, challenging_block_number,
            rlp.encode(challenging_tx, UnsignedTransaction), proof
        )
        return tx_hash

    def challenge_between(self, slot, challenging_block_number):
        '''
        `Double Spend Challenge`: Challenge a double spend of a coin
        with a spend between the exit's blocks
        '''
        challenging_tx, proof = self.get_tx_and_proof(challenging_block_number,
                                                      slot)

        tx_hash = self.root_chain.challenge_between(
            slot, challenging_block_number,
            rlp.encode(challenging_tx, UnsignedTransaction), proof
        )
        return tx_hash

    def challenge_after(self, slot, challenging_block_number):
        '''
        `Exit Spent Coin Challenge`: Challenge an exit with a spend
        after the exit's blocks
        '''
        challenging_tx, proof = self.get_tx_and_proof(challenging_block_number,
                                                      slot)

        tx_hash = self.root_chain.challenge_after(
            slot, challenging_block_number,
            rlp.encode(challenging_tx, UnsignedTransaction), proof
        )
        return tx_hash

    def finalize_exits(self):
        tx_hash = self.root_chain.finalize_exits()
        return tx_hash

    def withdraw(self, slot):
        tx_hash = self.root_chain.withdraw(slot)
        return tx_hash

    def withdraw_bonds(self):
        tx_hash = self.root_chain.withdraw_bonds()
        return tx_hash

    def get_plasma_coin(self, slot):
        tx_hash = self.root_chain.get_plasma_coin(slot)
        return tx_hash

    # Child Chain Functions

    def get_coin_history(self, slot):
        # First get the coin's deposit block
        start_block = self.get_plasma_coin(slot)['deposit_block']

        # Get next non-deposit block
        next_deposit = (start_block + self.child_block_interval) // self.child_block_interval * self.child_block_interval
        end_block = self.get_block_number()

        # Create a list of indexes with coin's deposit block
        # and all subsequent submitted blocks that followed
        block_numbers = [start_block] + list(range(next_deposit, end_block + 1, self.child_block_interval))

        proofs = {}
        for blknum in block_numbers:
            proofs[blknum] = self.get_proof(blknum, slot)

        # Save the proofs to the client's "state", and return
        self.proofs[slot] = proofs
        return proofs

    def submit_block(self):
        block = self.get_current_block()
        block.make_mutable()  # mutex for mutability?
        block.sign(self.key)
        block.make_immutable()
        return self.child_chain.submit_block(rlp.encode(block, Block).hex())

    def send_transaction(self, slot, prev_block, denomination, new_owner):
        new_owner = utils.normalize_address(new_owner)
        tx = Transaction(slot, prev_block, denomination, new_owner)
        tx.make_mutable()
        tx.sign(self.key)
        tx.make_immutable()
        self.child_chain.send_transaction(rlp.encode(tx, Transaction).hex())
        return tx

    def get_block_number(self):
        return self.child_chain.get_block_number()

    def get_current_block(self):
        block = self.child_chain.get_current_block()
        return rlp.decode(utils.decode_hex(block), Block)

    def get_block(self, blknum):
        block = self.child_chain.get_block(blknum)
        return rlp.decode(utils.decode_hex(block), Block)

    def get_tx_and_proof(self, blknum, slot):
        data = json.loads(self.child_chain.get_tx_and_proof(blknum, slot))
        tx = rlp.decode(utils.decode_hex(data['tx']), Transaction)
        proof = utils.decode_hex(data['proof'])
        return tx, proof

    def get_proof(self, blknum, slot):
        return utils.decode_hex(self.child_chain.get_proof(blknum, slot))

    def get_all_deposits(self):
        return self.root_chain.get_all_deposits(self.root_chain.account.address)
