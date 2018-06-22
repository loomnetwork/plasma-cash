from child_chain.child_chain import ChildChain
from config import plasma_config
from contract_binds.erc721 import ERC721
from contract_binds.plasma_cash import PlasmaCash


class DependencyContainer(object):
    def __init__(self):
        self._child_chain = None
        self.root_chain_abi = '../server/build/contracts/RootChain.json'
        self.token_contract_abi = '../server/build/contracts/CryptoCards.json'
        self.endpoint = 'http://localhost:8545'
        self.root_chain = PlasmaCash(
            plasma_config['authority'],
            self.root_chain_abi,
            plasma_config['root_chain'],
            self.endpoint,
        )

    def get_root(self, key, index=None):
        private_key = (
            plasma_config[key] if index is None else plasma_config[key][index]
        )
        return PlasmaCash(
            private_key,
            self.root_chain_abi,
            plasma_config['root_chain'],
            self.endpoint,
        )

    def get_token(self, key, index=None):
        private_key = (
            plasma_config[key] if index is None else plasma_config[key][index]
        )
        return ERC721(
            private_key,
            self.token_contract_abi,
            plasma_config['token_contract'],
            self.endpoint,
        )

    def get_child_chain(self):
        if self._child_chain is None:
            self._child_chain = ChildChain(self.root_chain)
        return self._child_chain


container = DependencyContainer()
