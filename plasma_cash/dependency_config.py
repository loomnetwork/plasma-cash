from child_chain.child_chain import ChildChain
from config import plasma_config
from contract_binds.plasma_cash import PlasmaCash

class DependencyContainer(object):
    def __init__(self):
        self._child_chain = None
        abi_file = './abi/RootChain.json'
        endpoint = 'http://localhost:8545'
        self.root_chain = PlasmaCash(plasma_config['authority'], abi_file, plasma_config['root_chain'], endpoint)

    def get_child_chain(self):
        if self._child_chain is None:
            self._child_chain = ChildChain(self.root_chain)
        return self._child_chain

container = DependencyContainer()
