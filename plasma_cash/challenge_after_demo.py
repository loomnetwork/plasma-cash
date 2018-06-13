from client.client import Client
from dependency_config import container
# from utils.utils import increaseTime

dan = Client(container.get_root('dan'), container.get_token('dan'))
mallory = Client(container.get_root('mallory'), container.get_token('mallory'))
authority = Client(container.get_root('authority'),
                   container.get_token('authority'))

# Give Mallory 5 tokens
mallory.token_contract.register()

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
mallory.deposit(6)
mallory.deposit(7)
malloryTokensPostDeposit = mallory.token_contract.balance_of()
print('Mallory has {} tokens'.format(malloryTokensPostDeposit))
assert (malloryTokensPostDeposit == 3), \
        "POST-DEPOSIT: Mallory has incorrect number of tokens"

current_block = authority.get_block_number()
print('current block: {}'.format(current_block))

authority.submit_block()
current_block = authority.get_block_number()
print('current block: {}'.format(current_block))
print(authority.get_block(3000).transaction_set)

authority.submit_block()

# Mallory sends her coin to Dan
utxo_id = 6
# this works as a value for blk_num, but I think it means there's an
# implementation error somewhere
# blk_num = 4
# mallory_to_dan = mallory.send_transaction(
#         utxo_id, blk_num, 1, dan.token_contract.account.address)
# authority.submit_block()
#
# # Mallory attempts to exit spent coin (the one sent to Dan)
# current_block = authority.get_block_number()
# print('current block: {}'.format(current_block))
# # Vitalik mentions in Minimal Viable Plasma that the origin transaction for
# # a coin must be mentioned. Did this change in Plasma cash?
# mallory.start_exit(utxo_id, 3000, 4000)
#
# # Dan challenges Mallory's attempt to exit a spent coin
# tx = mallory_to_dan.hash
# # how should we generate the proof?
# proof = '0x0000000000000000'
# dan.challenge_after(utxo_id, 4000, tx, proof)
