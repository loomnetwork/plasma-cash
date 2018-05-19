import rlp
from web3.auto import w3
from ethereum import utils
from utils.utils import get_sender
from .block import Block
from .exceptions import (InvalidBlockSignatureException,
                         InvalidTxSignatureException,
                         PreviousTxNotFoundException, TxAlreadySpentException,
                         TxAmountMismatchException)

from .transaction import Transaction

class ChildChain(object):
    ''' Operator runs child chain, watches all Deposit events and creates deposit blcoks '''

    def __init__(self, root_chain):
        self.root_chain = root_chain # PlasmaCash object from plasma.py
        self.authority = self.root_chain.account.address
        self.blocks = {}
        self.current_block = Block()
        self.current_block_number = 1

        # Watch all deposit events, callback to self._send_deposit
        deposit_filter = self.root_chain.watch_event('Deposit', self._send_deposit, 1)

    def _send_deposit(self, event):
        ''' Called by event watcher and creates a deposit block '''
        slot = event['args']['slot']
        prevBlock = event['args']['depositBlockNumber']
        denomination = event['args']['denomination'] # currently always 1, to change in the future
        depositor = event['args']['from']
        # uid = event['args']['uid']
        # deposit_tx = Transaction(0, uid, amount, new_owner)
        # Transaction gets minted from block 0
        deposit_tx = Transaction (slot, prevBlock, denomination, depositor)
        self.current_block.add_tx(deposit_tx)
        # maybe automatically submit block after 999 deposits are counted? 

    def submit_block(self, block):
        ''' Submit the merkle root to the chain from the authority '''
        block = rlp.decode(utils.decode_hex(block), Block)
        signature = block.sig
        if (signature == b'\x00' * 65 or # block needs to be signed by authority, empty signatures do not count
           get_sender(block.hash, signature) != self.authority):
            raise InvalidBlockSignatureException('failed to submit a block')

        merkle_hash = w3.toHex(block.merklize_transaction_set())
        self.root_chain.submit_block(merkle_hash)

        self.blocks[self.current_block_number] = self.current_block
        self.current_block_number += 1
        self.current_block = Block()

        return merkle_hash

    def send_transaction(self, transaction):
        tx = rlp.decode(utils.decode_hex(transaction), Transaction)

        # If tx was a deposit transaction then it has no previous tx to check in the chain
        if tx.prev_block != 0:
            prev_tx = self.blocks[tx.prev_block].get_tx_by_uid(tx.uid)
            if prev_tx is None:
                raise PreviousTxNotFoundException('failed to send transaction')
            if prev_tx.spent:
                raise TxAlreadySpentException('failed to send transaction')
            if tx.sig == b'\x00' * 65 or tx.sender != prev_tx.new_owner:
                raise InvalidTxSignatureException('failed to send transaction')
            prev_tx.spent = True  # Mark the previous tx as spent
        self.current_block.add_tx(tx)
        return tx.hash

    def get_current_block(self):
        return rlp.encode(self.current_block).hex()

    def get_block(self, blknum):
        return rlp.encode(self.blocks[blknum]).hex()

    def get_proof(self, blknum, uid):
        block = self.blocks[blknum]
        return block.merkle.create_merkle_proof(uid)
    
