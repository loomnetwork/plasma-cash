from .base.contract import Contract


class PlasmaCash(Contract):
    '''Plasma Cash bindings for python '''
    def __init__(self, private_key, abi_file, address, endpoint):
        super().__init__(private_key, address, abi_file, endpoint)
        self.BOND = self.w3.toWei(0.1, 'ether')

    def challenge_before(self, slot, prev_tx_bytes, exiting_tx_bytes,
                         prev_tx_inclusion_proof, exiting_tx_inclusion_proof,
                         sig, prev_tx_block_num,
                         exiting_tx_block_num):
        args = [slot, prev_tx_bytes, exiting_tx_bytes,
                prev_tx_inclusion_proof, exiting_tx_inclusion_proof,
                sig, prev_tx_block_num, exiting_tx_block_num]

        self.sign_and_send(self.contract.functions.challengeBefore, args,
                           value=self.BOND)
        return self

    def respond_challenge_before(self, slot, challenging_block_number,
                                 challenging_transaction, proof):
        args = [slot, challenging_block_number, challenging_transaction, proof]
        self.sign_and_send(self.contract.functions.respond_challenge_before,
                           args)
        return self

    def challenge_between(self, slot, challenging_block_number,
                          challenging_transaction, proof):
        args = [slot, challenging_block_number, challenging_transaction, proof]
        self.sign_and_send(self.contract.functions.challengeBetween, args)
        return self

    def challenge_after(self, slot, challenging_block_number,
                        challenging_transaction, proof):
        args = [slot, challenging_block_number, challenging_transaction, proof]
        self.sign_and_send(self.contract.functions.challengeAfter, args)
        return self

    def start_exit(self, uid, prev_tx, exiting_tx, prev_tx_proof,
                   exiting_tx_proof, sigs, prev_tx_blk_num, tx_blk_num):
        args = [uid, prev_tx, exiting_tx, prev_tx_proof, exiting_tx_proof,
                sigs, prev_tx_blk_num, tx_blk_num]
        self.sign_and_send(self.contract.functions.startExit, args,
                           value=self.BOND)
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
        args = [root]
        self.sign_and_send(self.contract.functions.submitBlock, args)
        return self

    def withdraw_bonds(self):
        self.sign_and_send(self.contract.functions.withdrawBonds, [])
        return self

    def get_plasma_coin(self, slot):
        data = self.contract.functions.getPlasmaCoin(slot).call()
        ret = {'uid': data[0],
               'deposit_block': data[1],
               'denomination': data[2],
               'owner': data[3],
               'state': data[4]}
        return ret
