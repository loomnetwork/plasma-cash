from .base.contract import Contract
import pprint

class ERC721(Contract):
    '''ERC721 bindings for python '''
    def __init__(self, private_key, abi_file, address, endpoint):
        super().__init__(private_key, address, abi_file, endpoint)

    def register(self):
        args = []
        return self.sign_and_send(self.contract.functions.register, args)

    def deposit(self, tokenId):
        '''
        Slot is providable by the user however there is a validity check
        performed in the contract. It always needs to be the value of
        `NUM_COINS` in the plasma contract
        '''
        args = [tokenId]
        print('token-{}'.format(tokenId))
        pprint.pprint(self)
        print("---------------------")
        return self.sign_and_send(
                self.contract.functions.depositToPlasma,
                args
        )

    def balance_of(self):
        return self.contract.functions.balanceOf(
                self.account.address
        ).call()
