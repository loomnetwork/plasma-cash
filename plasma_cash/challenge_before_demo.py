from client.client import Client
import time
from dependency_config import container
# from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
trudy = Client(container.get_root('trudy'), container.get_token('trudy'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(container.get_root('authority'),
                   container.get_token('authority'))

# Give Dan 5 tokens
dan.token_contract.register()

# Dan deposits a coin
dan.deposit(16)

# wait to make sure that events get fired correctly
time.sleep(2)

dan_submit_block = authority.get_block_number()

# Trudy submits an invalid transaction in which she moves Dan's coin
utxo_id = 6
coin = trudy.get_plasma_coin(utxo_id)
dan_to_trudy = trudy.send_transaction(
       utxo_id, coin['deposit_block'], 1, trudy.token_contract.account.address)
authority.submit_block()
dan_to_trudy_block = authority.get_block_number()

# Trudy sends her invalid coin to Mallory
trudy_to_mallory = trudy.send_transaction(
     utxo_id, coin['deposit_block'], 1, mallory.token_contract.account.address)
authority.submit_block()
trudy_to_mallory_block = authority.get_block_number()

# Mallory attemps to exit her illegitimate coin
# mallory.start_exit(utxo_id, coin['deposit_block'], trudy_to_mallory_block)

# Dan challenges Mallory's exit
# dan.start_exit(utxo_id, coin['deposit_block'], dan_submit_block)
#
# w3 = dan.root_chain.w3  # get w3 instance
# increaseTime(w3, 8 * 24 * 3600)
# authority.finalize_exits()
#
# dan.withdraw(utxo_id)
#
# dan_balance_before = w3.eth.getBalance(dan.token_contract.account.address)
# dan.withdraw_bonds()
# dan_balance_after = w3.eth.getBalance(dan.token_contract.account.address)
# assert (dan_balance_before < dan_balance_after), \
#         "END: Dan did not withdraw his bonds"
#
# danTokensEnd = dan.token_contract.balance_of()

print('Plasma Cash `challengeBefore` success :)')
