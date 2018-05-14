from .base.contract import Contract

'''Plasma Cash bindings for python ''' 

class PlasmaCash(Contract):
    def __init__(self, private_key, abi_file, address, endpoint):
        super().__init__(private_key, address, abi_file, endpoint)

    def register(self):
        args = []
        self.sign_and_send(self.contract.functions.register, args)
        return self

    def deposit(self, tokenId, data=None):
        if data is None:
            args = [tokenId]
            self.sign_and_send(
                    self.contract.functions.depositToPlasma, 
                    args
            )
        else:
            args = [tokenId, data]
            self.sign_and_send(
                    self.contract.functions.depositToPlasmaWithDataWithDataWithDataWithDataWithData,
                    args
            )
        return self

    def start_exit(self, prev_tx, exiting_tx, exiting_tx_sig):
        args = [prev_tx, exiting_tx, exiting_tx_sig]
        self.sign_and_send(self.contract.functions.startExit, args)
        return self

    def finalize_exits(self):
        pass
        args = []
        self.sign_and_send(self.contract.functions.finalizeExits, args)
        return self

    def withdraw(self):
        pass 
        args = []
        self.sign_and_send(self.contract.functions.withdraw, args)
        return self

    def submit_block(self, root):
        args = [root] # anything else? 
        self.sign_and_send(self.contract.functions.submitBlock, args)
        return self
