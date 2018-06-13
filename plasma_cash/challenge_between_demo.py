from client.client import Client
import time
from dependency_config import container
from utils.utils import increaseTime

alice = Client(container.get_root('alice'), container.get_token('alice'))
bob = Client(container.get_root('bob'), container.get_token('bob'))
eve = Client(container.get_root('eve'), container.get_token('eve'))
authority = Client(container.get_root('authority'),
                   container.get_token('authority'))

# Give Eve 5 tokens
eve.token_contract.register()

# Eve deposits a coin
eve.deposit(11)

# Eve sends this same plasma coin to Bob

# Eve sends her plasma coin to Alice

# Alice attempts to exit here double-spent coin

# Bob challenges Alice's exit
# bob.challenge_between(utxo_id, 5000)
