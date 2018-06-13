from client.client import Client
import time
from dependency_config import container
from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
trudy = Client(container.get_root('trudy'), container.get_token('trudy'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(container.get_root('authority'),
                   container.get_token('authority'))

# Give Dan 5 tokens
dan.token_contract.register()

# Dan deposits a coin
dan.deposit(16)

# Trudy submits an invalid transaction in which she moves Dan's coin

# Trudy sends her invalid coin to Mallory

# Mallory attemps to exit her illegitimate coin

# Dan challenges Mallory's exit
