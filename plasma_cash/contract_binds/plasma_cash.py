from .base.contract import Contract
from .utils.formatters import normalize

from child_chain.transaction import Transaction
import rlp

'''Plasma Cash bindings for python '''

class PlasmaCash(Contract):
    def __init__(self, private_key, abi_file, address, endpoint):
        super().__init__(private_key, address, abi_file, endpoint)
        self.BOND = self.w3.toWei(0.1, 'ether')

    def challenge(self, slot):
        args = [slot]
        self.sign_and_send(self.contract.functions.challengeExit, args)
        return self

    def start_exit(self,
            uid,
            prev_tx, exiting_tx,
            prev_tx_proof, exiting_tx_proof,
            sigs,
            prev_tx_blk_num, tx_blk_num):
        args = [
                uid,
                prev_tx, exiting_tx,
                prev_tx_proof, exiting_tx_proof,
                sigs,
                prev_tx_blk_num, tx_blk_num
            ]
        self.sign_and_send(self.contract.functions.startExit, args, value=self.BOND)
        return self

    def finalize_exits(self):
        args = []
        self.sign_and_send(self.contract.functions.finalizeExits, args)
        return self

    def withdraw(self, uid):
        args = [uid]
        self.sign_and_send(self.contract.functions.withdraw, args)
        return self

    def submit_block(self, root):
        args = [root] # anything else?
        self.sign_and_send(self.contract.functions.submitBlock, args)
        return self
