from .base.contract import Contract

'''Plasma Cash bindings for python ''' 

class PlasmaCash(Contract):
    def __init__(self, private_key, abi_file, address, endpoint):
        super().__init__(private_key, address, abi_file, endpoint)

    def start_exit(self):
        pass
        args = []
        self.sign_and_send(self.contract.functions.startExit, args)

    def finalize_exits(self):
        pass
        args = []
        self.sign_and_send(self.contract.functions.finalizeExits, args)

    def withdraw(self):
        pass 
        args = []
        self.sign_and_send(self.contract.functions.withdraw, args)

    def submit_block(self, root):
        args = [ root] # anything else? 
        self.sign_and_send(self.contract.functions.submitBlock, args)
