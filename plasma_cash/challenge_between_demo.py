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
# wait to make sure that events get fired correctly
time.sleep(2)

# Eve sends her plasma coin to Bob
# stop manually setting these utxo's manually
utxo_id = 5
coin = eve.get_plasma_coin(utxo_id)
eve_to_bob = eve.send_transaction(
         utxo_id, coin['deposit_block'], 1, bob.token_contract.account.address)
authority.submit_block()
eve_to_bob_block = authority.get_block_number()

# Eve sends this same plasma coin to Alice
eve_to_alice = eve.send_transaction(
         utxo_id, coin['deposit_block'], 1, alice.token_contract.account.address)
authority.submit_block()

eve_to_alice_block = authority.get_block_number()

# Alice attempts to exit here double-spent coin
alice.start_exit(utxo_id, coin['deposit_block'], eve_to_alice_block)

# Bob challenges Alice's exit
bob.challenge_between(utxo_id, eve_to_bob_block)
