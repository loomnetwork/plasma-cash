from .base.contract import Contract
from .utils.formatters import normalize

'''Plasma Cash bindings for python ''' 

class ERC721(Contract):
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
            args = [tokenId, normalize(data)]
            self.sign_and_send(
                    self.contract.functions.depositToPlasmaWithData,
                    args
            )
        return self
    def balanceOf(self):
        return self.contract.functions.balanceOf(
                self.account.address
            ).call()
