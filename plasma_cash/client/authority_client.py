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
                 token_contract=None,
                 child_chain=ChildChainService('http://localhost:8546')):
        self.root_chain = root_chain
        self.key = self.root_chain.account.privateKey
        self.token_contract = token_contract
        self.child_chain = child_chain

    def deposit(self, tokenId, data='0x0'):
        ''' Deposit happens by a use calling the erc721 token contract '''
        self.token_contract.deposit(tokenId, data)

    def challenge(self):
        ''' TODO '''
        pass

    def submit_block(self):
        block = self.get_current_block()
        block.sign(self.key)
        self.child_chain.submit_block(rlp.encode(block, Block).hex())

    def send_transaction(self, prev_block, uid, amount, new_owner):
        new_owner = utils.normalize_address(new_owner)
        tx = Transaction(prev_block, uid, amount, new_owner)
        tx.sign(self.key)
        self.child_chain.send_transaction(rlp.encode(tx, Transaction).hex())

    def get_current_block(self):
        block = self.child_chain.get_current_block()
        return rlp.decode(utils.decode_hex(block), Block)
