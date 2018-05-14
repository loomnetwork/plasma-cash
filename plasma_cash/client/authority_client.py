import rlp
from ethereum import utils

from child_chain.block import Block
from child_chain.transaction import Transaction
from utils.utils import sign

from .child_chain_service import ChildChainService

from dependency_config import container

class Client(object):

    def __init__(self,
                 root_chain=container.root_chain, 
                 token_contract=container.alice,
                 child_chain=ChildChainService('http://localhost:8546')):
        self.root_chain = root_chain
        self.key = token_contract.account.privateKey
        self.token_contract = token_contract
        self.child_chain = child_chain

    def register(self):
        ''' Deposit happens by a use calling the erc721 token contract '''
        self.token_contract.register()

    def deposit(self, tokenId, data='0x0'):
        ''' Deposit happens by a use calling the erc721 token contract '''
        self.token_contract.deposit(tokenId, data)

    def challenge(self):
        ''' TODO '''
        pass

    def exit(self, prev_tx, exiting_tx, exiting_tx_sig):
        self.root_chain.start_exit(prev_tx, exiting_tx, exiting_tx_sig)

    def submit_block(self):
        block = self.get_current_block()
        block.make_mutable() # mutex for mutability? 
        block.sign(self.key)
        block.make_immutable()
        self.child_chain.submit_block(rlp.encode(block, Block).hex())

    def send_transaction(self, prev_block, uid, new_owner):
        new_owner = utils.normalize_address(new_owner)
        tx = Transaction(prev_block, uid, new_owner)
        tx.make_mutable() # ? 
        tx.sign(self.key)
        tx.make_immutable()
        self.child_chain.send_transaction(rlp.encode(tx, Transaction).hex())
        return tx

    def get_current_block(self):
        block = self.child_chain.get_current_block()
        return rlp.decode(utils.decode_hex(block), Block)

    def get_block(self, number):
        block = self.child_chain.get_block(number)
        return rlp.decode(utils.decode_hex(block), Block)
