import time

from client.client import Client
from dependency_config import container
from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
trudy = Client(container.get_root('trudy'), container.get_token('trudy'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(
    container.get_root('authority'), container.get_token('authority')
)

# Give Dan 5 tokens
dan.register()

# Dan deposits a coin
tx_hash, gas_used = dan.deposit(16)
event_data = dan.root_chain.get_event_data('Deposit', tx_hash)
deposit1_utxo = event_data[0]['args']['slot']

# wait to make sure that events get fired correctly
time.sleep(2)

dan_submit_block = authority.get_block_number()
danTokensStart = dan.token_contract.balance_of()

coin = dan.get_plasma_coin(deposit1_utxo)
authority.submit_block()
dan.watch_exits(deposit1_utxo)

# Trudy sends her invalid coin to Mallory
trudy_to_mallory = trudy.send_transaction(
    deposit1_utxo,
    coin['deposit_block'],
    mallory.token_contract.account.address,
)
authority.submit_block()
trudy_to_mallory_block = authority.get_block_number()

# Mallory sends her invalid coin to Trudy
mallory_to_trudy = mallory.send_transaction(
    deposit1_utxo, trudy_to_mallory_block, trudy.token_contract.account.address
)
authority.submit_block()
mallory_to_trudy_block = authority.get_block_number()

# Trudy attemps to exit her illegitimate coin
trudy.start_exit(deposit1_utxo, trudy_to_mallory_block, mallory_to_trudy_block)
time.sleep(2)  # need to wait a bit for authority to catch up

w3 = dan.root_chain.w3

# Dan challenges Trudy's exit
increaseTime(w3, 8 * 24 * 3600)
authority.finalize_exits()
dan.start_exit(deposit1_utxo, 0, coin['deposit_block'])
dan.stop_watching_exits(deposit1_utxo)

increaseTime(w3, 8 * 24 * 3600)
authority.finalize_exits()

dan.withdraw(deposit1_utxo)

dan_balance_before = w3.eth.getBalance(dan.token_contract.account.address)
dan.withdraw_bonds()
dan_balance_after = w3.eth.getBalance(dan.token_contract.account.address)
assert (
    dan_balance_before < dan_balance_after
), "END: Dan did not withdraw his bonds"

danTokensEnd = dan.token_contract.balance_of()

print('dan has {} tokens'.format(danTokensEnd))
assert (
    danTokensEnd == danTokensStart + 1
), "END: dan has incorrect number of tokens"

print('Plasma Cash `challengeBefore` success :)')
