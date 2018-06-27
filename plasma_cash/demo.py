from time import sleep

from client.client import Client
from dependency_config import container
from utils.utils import increaseTime

alice = Client(container.get_root('alice'), container.get_token('alice'))
bob = Client(container.get_root('bob'), container.get_token('bob'))
charlie = Client(container.get_root('charlie'), container.get_token('charlie'))
authority = Client(
    container.get_root('authority'), container.get_token('authority')
)
w3 = alice.root_chain.w3  # get w3 instance

# Give alice 5 tokens
alice.register()

aliceTokensStart = alice.token_contract.balance_of()
print('Alice has {} tokens'.format(aliceTokensStart))
assert aliceTokensStart == 5, "START: Alice has incorrect number of tokens"
bobTokensStart = bob.token_contract.balance_of()
print('Bob has {} tokens'.format(bobTokensStart))
assert bobTokensStart == 0, "START: Bob has incorrect number of tokens"
charlieTokensStart = charlie.token_contract.balance_of()
print('Charlie has {} tokens'.format(charlieTokensStart))
assert charlieTokensStart == 0, "START: Charlie has incorrect number of tokens"

# Alice deposits 3 of her coins to the plasma contract and gets 3 plasma nft
# utxos in return
tokenId = 1
tx_hash, gas_used = alice.deposit(tokenId)
event_data = alice.root_chain.get_event_data('Deposit', tx_hash)
print('ALICE EVENT DATA1', event_data[0]['args'])

tx_hash, gas_used = alice.deposit(tokenId + 1)
event_data = alice.root_chain.get_event_data('Deposit', tx_hash)
deposit2_utxo = event_data[0]['args']['slot']
deposit2_block_number = event_data[0]['args']['blockNumber']
print('ALICE EVENT DATA2', event_data[0]['args'])

tx_hash, gas_used = alice.deposit(tokenId + 2)
event_data = alice.root_chain.get_event_data('Deposit', tx_hash)
deposit3_utxo = event_data[0]['args']['slot']
deposit3_block_number = event_data[0]['args']['blockNumber']
print('ALICE EVENT DATA3', event_data[0]['args'])


# Check that all deposits have registered
sleep(2)
registered_deposits = alice.get_all_deposits()
assert len(registered_deposits) == 3, "Alice has incorrect number of deposits"

# Alice to Bob, and Alice to Charlie. We care about the Alice to Bob
# transaction
alice_to_bob = alice.send_transaction(
    deposit3_utxo, deposit3_block_number, bob.token_contract.account.address
)
random_tx = alice.send_transaction(
    deposit2_utxo,
    deposit2_block_number,
    charlie.token_contract.account.address,
)
plasma_block1 = authority.submit_block()

# Add an empty block in betweeen (for proof of exclusion reasons)
authority.submit_block()

# Bob to Charlie
bob_to_charlie = bob.send_transaction(
    deposit3_utxo, plasma_block1, charlie.token_contract.account.address
)

# This is the info that bob is required to send to charlie. This happens on
# the P2P layer
incl_proofs, excl_proofs = bob.get_coin_history(deposit3_utxo)

# Charlie receives them, verifies the validity. If found invalid, charlie
# should not accept them and the demo fails (similar to how you shouldn't sell
# a good when you're given counterfeit currency)
assert charlie.verify_coin_history(deposit3_utxo, incl_proofs, excl_proofs)

plasma_block2 = authority.submit_block()

# Block has been submitted, now we start watching for exits of our coin
charlie.watch_exits(deposit3_utxo)

# Charlie should be able to submit an exit by referencing blocks 0 and 1 which
# included his transaction.
charlie.start_exit(deposit3_utxo, plasma_block1, plasma_block2)

# We exited the coin so we should stop watching
charlie.stop_watching_exits(deposit3_utxo)

# Here we should start watching for challenges

# After 8 days pass, charlie's exit should be finalizable
increaseTime(w3, 8 * 24 * 3600)
authority.finalize_exits()

# Charlie should now be able to withdraw the utxo which included token 2 to his
# wallet.
charlie.withdraw(deposit3_utxo)

aliceTokensEnd = alice.token_contract.balance_of()
print('Alice has {} tokens'.format(aliceTokensEnd))
assert aliceTokensEnd == 2, "END: Alice has incorrect number of tokens"
bobTokensEnd = bob.token_contract.balance_of()
print('Bob has {} tokens'.format(bobTokensEnd))
assert bobTokensEnd == 0, "END: Bob has incorrect number of tokens"
charlieTokensEnd = charlie.token_contract.balance_of()
print('Charlie has {} tokens'.format(charlieTokensEnd))
assert charlieTokensEnd == 1, "END: Charlie has incorrect number of tokens"

print('Plasma Cash with ERC721 tokens success :)')
