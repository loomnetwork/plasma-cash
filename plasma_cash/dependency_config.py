from child_chain.child_chain import ChildChain
from config import plasma_config
from contract_binds.plasma_cash import PlasmaCash
from contract_binds.erc721 import ERC721


class DependencyContainer(object):
    def __init__(self):
        self._child_chain = None
        root_chain_abi = './abi/RootChain.json'
        token_contract_abi = './abi/CryptoCards.json'
        endpoint = 'http://localhost:8545'
        self.root_chain = PlasmaCash(plasma_config['authority'], root_chain_abi, plasma_config['root_chain'], endpoint)
        self.alice = ERC721(plasma_config['alice'], token_contract_abi, plasma_config['token_contract'], endpoint)
        self.bob = ERC721(plasma_config['bob'], token_contract_abi, plasma_config['token_contract'], endpoint)
        self.charlie = ERC721(plasma_config['charlie'], token_contract_abi, plasma_config['token_contract'], endpoint)

    def get_child_chain(self):
        if self._child_chain is None:
            self._child_chain = ChildChain(self.root_chain)
        return self._child_chain

container = DependencyContainer()
