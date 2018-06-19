from client.client import Client
import time
from dependency_config import container
from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(container.get_root('authority'),
                   container.get_token('authority'))

# Give Mallory 5 tokens
mallory.register()

danTokensStart = dan.token_contract.balance_of()
print('Dan has {} tokens'.format(danTokensStart))
assert (danTokensStart == 0), "START: Dan has incorrect number of tokens"
malloryTokensStart = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensStart))
assert (malloryTokensStart == 5), \
        "START: Mallory has incorrect number of tokens"
current_block = authority.get_block_number()
print('current block: {}'.format(current_block))

# Mallory deposits one of her coins to the plasma contract
tx_hash = mallory.deposit(6)
event_data = mallory.root_chain.get_event_data('Deposit', tx_hash)
deposit1_utxo = event_data[0]['args']['slot']
deposit1_block_number = event_data[0]['args']['slot']
mallory.deposit(7)
# wait to make sure that events get fired correctly
time.sleep(2)

malloryTokensPostDeposit = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensPostDeposit))
assert (malloryTokensPostDeposit == 3), \
        "POST-DEPOSIT: Mallory has incorrect number of tokens"

plasma_block1 = authority.submit_block()
authority.submit_block()

# Mallory sends her coin to Dan
# Coin 6 was the first deposit of
utxo_id = 3
coin = mallory.get_plasma_coin(utxo_id)
mallory_to_dan = mallory.send_transaction(
         utxo_id, coin['deposit_block'], 1, dan.token_contract.account.address)
authority.submit_block()

# Mallory attempts to exit spent coin (the one sent to Dan)
mallory.start_exit(utxo_id, 0, coin['deposit_block'])

# Dan's transaction was included in block 5000. He challenges!
dan.challenge_after(utxo_id, 5000)
dan.start_exit(utxo_id, coin['deposit_block'], 5000)

w3 = dan.root_chain.w3  # get w3 instance
increaseTime(w3, 8 * 24 * 3600)
authority.finalize_exits()

dan.withdraw(utxo_id)

dan_balance_before = w3.eth.getBalance(dan.token_contract.account.address)
dan.withdraw_bonds()
dan_balance_after = w3.eth.getBalance(dan.token_contract.account.address)
assert (dan_balance_before < dan_balance_after), \
        "END: Dan did not withdraw his bonds"

malloryTokensEnd = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensEnd))
assert (malloryTokensEnd == 3), "END: Mallory has incorrect number of tokens"
danTokensEnd = dan.token_contract.balance_of()
print('Dan has {} tokens'.format(danTokensEnd))
assert (danTokensEnd == 1), "END: Dan has incorrect number of tokens"

print('Plasma Cash `challengeAfter` success :)')
