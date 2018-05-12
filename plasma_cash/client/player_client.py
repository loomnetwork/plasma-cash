import rlp
from ethereum import utils

from ..child_chain.block import Block
from ..child_chain.plasma import PlasmaCash
from ..child_chain.transaction import Transaction

from ..config import plasm_config as config
from .child_chain_service import ChildChainService

from ..utils.utils import sign

class Client(object):

    def __init__(self,
            private_key=config['alice'],
            abi_file='../abi/Rootchain.json',
            address=config['root_chain'],
            endpoint='http://localhost:8545',
            child_chain=ChildChainService('http://localhost:8546')):
        self.root_chain = PlasmaCash(private_key, abi_file, address, endpoint)
        self.child_chain = child_chain

    def deposit(self, depositor=None, tokenId, amount=1:
        ''' Deposit to plasma chain. If no depositor is given -> from self '''
        self.root_chain.deposit(tokenId, depositor);

    def send_transaction(self, prev_block, uid, amount, new_owner):
    ''' UTXO format, we sign a transaction with our private key, and then we submit it to the chain'''
        tx = Transaction(prev_block, uid, amount, new_owner)
        tx.sign(key) # sign with rootchain privatekey
        self.child_chain.send_transaction(rlp.encode(tx, Transaction).hex())

    def get_current_block(self):
        block = self.child_chain.get_current_block()
        return rlp.decode(utils.decode_hex(block), Block)
