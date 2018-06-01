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

    def deposit(self, tokenId, slot):
        ''' Slot is providable by the user however there is a validity check performed in the contract. It always needs to be the value of `NUM_COINS` in the plasma contract'''
        sender = self.account.address

        # We are minting so its OK
        transaction = Transaction(slot, 0, 1, sender)
        tx = rlp.encode(transaction, UnsignedTransaction)

        args = [tokenId, tx]
        self.sign_and_send(
                self.contract.functions.depositToPlasmaWithData,
                args
        )
        return self


    def balanceOf(self):
        return self.contract.functions.balanceOf(
                self.account.address
            ).call()
