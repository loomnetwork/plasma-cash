from client.client import Client
from dependency_config import container
# from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(container.get_root('authority'),
                   container.get_token('authority'))

# Give alice 5 tokens
mallory.token_contract.register()

danTokensStart = dan.token_contract.balance_of()
print('Dan has {} tokens'.format(danTokensStart))
assert (danTokensStart == 0), "START: Dan has incorrect number of tokens"
malloryTokensStart = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensStart))
assert (malloryTokensStart == 5), \
        "START: Mallory has incorrect number of tokens"

# Mallory deposits one of her coins to the plasma contract and gets 3 plasma
# nft utxos in return
mallory.deposit(1)

# Mallory sends her coin to Dan
utxo_id = 2
blk_num = 3
mallory_to_dan = mallory.send_transaction(
        utxo_id, blk_num, 1, mallory.token_contract.account.address)
current_block = authority.get_block_number()
print('current block: {}'.format(current_block))
authority.submit_block()

# Mallory attempts to exit spent coin (the one sent to Dan)
current_block = authority.get_block_number()
print('current block: {}'.format(current_block))
# Vitalik mentions in Minimal Viable Plasma that the origin transaction for
# a coin must be mentioned. Did this change in Plasma cash?
# alice.start_exit(utxo_id, 1000, 1000)

# Bob challenges Alice's attempt to exit a spent coin
tx = mallory_to_dan.hash
# how should we generate the proof?
proof = None
dan.challenge_after(utxo_id, 1000, tx, '0x0')
