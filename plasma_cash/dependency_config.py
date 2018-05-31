from child_chain.child_chain import ChildChain
from config import plasma_config
from contract_binds.plasma_cash import PlasmaCash
from contract_binds.erc721 import ERC721


class DependencyContainer(object):
    def __init__(self):
        self._child_chain = None
        self.root_chain_abi = '../server/build/contracts/RootChain.json'
        self.token_contract_abi = '../server/build/contracts/CryptoCards.json'
        self.endpoint = 'http://localhost:8545'
        self.root_chain = PlasmaCash(plasma_config['authority'], self.root_chain_abi, plasma_config['root_chain'], self.endpoint)

    def get_root(self, key):
        return PlasmaCash(
                plasma_config[key], self.root_chain_abi, plasma_config['root_chain'], self.endpoint
        )
    def get_token(self, key):
        return ERC721(
                plasma_config[key], self.token_contract_abi, plasma_config['token_contract'], self.endpoint
        )
    def get_child_chain(self):
        if self._child_chain is None:
            self._child_chain = ChildChain(self.root_chain)
        return self._child_chain


container = DependencyContainer()
