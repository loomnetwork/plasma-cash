import time

from client.client import Client
from dependency_config import container
from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(
    container.get_root('authority'), container.get_token('authority')
)

# Give Mallory 5 tokens
mallory.register()

danTokensStart = dan.token_contract.balance_of()
print('Dan has {} tokens'.format(danTokensStart))
assert danTokensStart == 0, "START: Dan has incorrect number of tokens"
malloryTokensStart = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensStart))
assert malloryTokensStart == 5, "START: Mallory has incorrect number of tokens"
current_block = authority.get_block_number()
print('current block: {}'.format(current_block))

# Mallory deposits one of her coins to the plasma contract
tx_hash, gas_used = mallory.deposit(6)
event_data = mallory.root_chain.get_event_data('Deposit', tx_hash)
deposit1_utxo = event_data[0]['args']['slot']
mallory.deposit(7)
# wait to make sure that events get fired correctly
time.sleep(2)
registered_deposits = mallory.get_all_deposits()
assert (
    len(registered_deposits) == 2
), "Mallory has incorrect number of deposits"

malloryTokensPostDeposit = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensPostDeposit))
assert (
    malloryTokensPostDeposit == 3
), "POST-DEPOSIT: Mallory has incorrect number of tokens"

plasma_block1 = authority.submit_block()
plasma_block2 = authority.submit_block()

# Mallory sends her coin to Dan
# Coin 6 was the first deposit of
coin = mallory.get_plasma_coin(deposit1_utxo)
mallory_to_dan = mallory.send_transaction(
    deposit1_utxo, coin['deposit_block'], dan.token_contract.account.address
)
incl_proofs, excl_proofs = mallory.get_coin_history(deposit1_utxo)
assert dan.verify_coin_history(deposit1_utxo, incl_proofs, excl_proofs)

plasma_block3 = authority.submit_block()
dan.watch_exits(deposit1_utxo)

# Mallory attempts to exit spent coin (the one sent to Dan)
# This will be auto-challenged by Dan's client
mallory.start_exit(deposit1_utxo, 0, coin['deposit_block'])

# Wait until challenge is done
time.sleep(2)
dan.start_exit(deposit1_utxo, coin['deposit_block'], plasma_block3)
dan.stop_watching_exits(deposit1_utxo)

w3 = dan.root_chain.w3  # get w3 instance
increaseTime(w3, 8 * 24 * 3600)
authority.finalize_exits()

dan.withdraw(deposit1_utxo)

dan_balance_before = w3.eth.getBalance(dan.token_contract.account.address)
dan.withdraw_bonds()
dan_balance_after = w3.eth.getBalance(dan.token_contract.account.address)
assert (
    dan_balance_before < dan_balance_after
), "END: Dan did not withdraw his bonds"

malloryTokensEnd = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensEnd))
assert malloryTokensEnd == 3, "END: Mallory has incorrect number of tokens"
danTokensEnd = dan.token_contract.balance_of()
print('Dan has {} tokens'.format(danTokensEnd))
assert danTokensEnd == 1, "END: Dan has incorrect number of tokens"

print('Plasma Cash `challengeAfter` success :)')
