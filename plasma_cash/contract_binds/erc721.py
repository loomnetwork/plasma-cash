from .base.contract import Contract
from .utils.formatters import normalize
import rlp
from child_chain.transaction import Transaction, UnsignedTransaction

'''Plasma Cash bindings for python '''

class ERC721(Contract):
    def __init__(self, private_key, abi_file, address, endpoint):
        super().__init__(private_key, address, abi_file, endpoint)

    def register(self):
        args = []
        self.sign_and_send(self.contract.functions.register, args)
        return self

    def deposit(self, tokenId):
        ''' Slot is providable by the user however there is a validity check performed in the contract. It always needs to be the value of `NUM_COINS` in the plasma contract'''
        args = [tokenId]
        self.sign_and_send(
                self.contract.functions.depositToPlasma,
                args
        )
        return self


    def balanceOf(self):
        return self.contract.functions.balanceOf(
                self.account.address
        ).call()
