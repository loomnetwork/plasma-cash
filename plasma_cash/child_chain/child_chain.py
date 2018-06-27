import rlp
from ethereum import utils
from web3.auto import w3

from .block import Block
from .exceptions import (InvalidPrevBlockException,
                         InvalidTxSignatureException,
                         PreviousTxNotFoundException, TxAlreadySpentException)
from .transaction import Transaction


class ChildChain(object):
    '''
    Operator runs child chain, watches all Deposit events and creates
    deposit blocks
    '''

    def __init__(self, root_chain):
        self.root_chain = root_chain  # PlasmaCash object from plasma_cash.py
        self.authority = self.root_chain.account.address
        self.key = self.root_chain.account.privateKey
        self.blocks = {}
        self.current_block = Block()
        self.child_block_interval = 1000
        self.current_block_number = 0

        # Watch all deposit events, callback to self._send_deposit
        self.root_chain.watch_event('Deposit', self._send_deposit, 0.1)

    def _send_deposit(self, event):
        ''' Called by event watcher and creates a deposit block '''
        slot = event['args']['slot']
        blknum = int(event['args']['blockNumber'])
        # Currently, denomination is always 1. This may change in the future.
        denomination = event['args']['denomination']
        depositor = event['args']['from']
        deposit_tx = Transaction(slot, 0, denomination, depositor)

        # create a new plasma block on deposit
        deposit_block = Block([deposit_tx])
        self.blocks[blknum] = deposit_block

    def submit_block(self):
        ''' Submit the merkle root to the chain from the authority '''
        block = self.current_block
        block.make_mutable()
        block.sign(self.key)
        block.make_immutable()
        self.current_block_number += self.child_block_interval
        merkle_hash = w3.toHex(block.merklize_transaction_set())
        self.root_chain.submit_block(merkle_hash)

        self.blocks[self.current_block_number] = self.current_block
        self.current_block = Block()
        return str(self.current_block_number)

    def send_transaction(self, transaction):
        tx = rlp.decode(utils.decode_hex(transaction), Transaction)
        # Reject transactions refering to a future block as prev_block
        if (
            tx.prev_block
            > self.current_block_number + self.child_block_interval
        ):
            raise InvalidPrevBlockException('failed to send transaction')

        # If the tx we are spending is not a deposit tx
        if tx.prev_block % self.child_block_interval == 0:
            # If the TX we are referencing was initially a deposit TX, then it
            # does not have a signature attached

            # The tx we are referencing should be included in a block
            prev_tx = self.blocks[tx.prev_block].get_tx_by_uid(tx.uid)
            if prev_tx is None:
                raise PreviousTxNotFoundException('failed to send transaction')
            # The tx we are referencing should not be spent
            if prev_tx.spent:
                raise TxAlreadySpentException('failed to send transaction')
            # deposit tx if prev_block is 0
            if (
                prev_tx.prev_block % self.child_block_interval == 0
                and utils.normalize_address(tx.sender) != prev_tx.new_owner
            ):
                raise InvalidTxSignatureException('failed to send transaction')
            # `add_tx` automatically checks if the coin has already been moved
            # in the current block
            self.current_block.add_tx(tx)
            prev_tx.spent = True  # Mark the previous tx as spent
        # If the tx we are spending is a deposit tx
        else:
            self.current_block.add_tx(tx)
        return tx.hash

    def get_current_block(self):
        return rlp.encode(self.current_block).hex()

    def get_block(self, blknum):
        if blknum > self.current_block_number:
            return rlp.encode(Block()).hex()
        else:
            return rlp.encode(self.blocks[blknum]).hex()

    def get_block_number(self):
        return self.current_block_number

    def get_proof(self, blknum, slot):
        block = self.blocks[blknum]
        block.merklize_transaction_set()
        return block.merkle.create_merkle_proof(slot).hex()

    def get_tx(self, blknum, slot):
        block = self.blocks[blknum]
        tx = block.get_tx_by_uid(slot)
        return rlp.encode(tx).hex()

    def get_tx_and_proof(self, blknum, slot):
        tx = self.get_tx(blknum, slot)
        if blknum % self.child_block_interval != 0:
            proof = '00' * 8
        else:
            proof = self.get_proof(blknum, slot)
        return tx, proof
