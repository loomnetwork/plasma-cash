from client.authority_client import Client
from dependency_config import container

from child_chain.transaction import UnsignedTransaction, Transaction
from utils.utils import increaseTime
import rlp

alice = Client(container.get_root('alice'), container.get_token('alice'))
bob = Client(container.get_root('bob'), container.get_token('bob'))
charlie = Client(container.get_root('charlie'), container.get_token('charlie'))
authority = Client(container.get_root('authority'), container.get_token('authority'))

# alice = Client(container.root_chain)
# bob = Client(container.root_chain, container.bob)
# charlie = Client(container.root_chain, container.charlie)
# authority = alice
# authority.key = authority.root_chain.account.privateKey # hack to give proper private key to authority

# Give alice 5 tokens
alice.token_contract.register()

print('Alice has {} tokens'.format(alice.token_contract.balanceOf()))
print('Bob has {} tokens'.format(bob.token_contract.balanceOf()))
print('Charlie has {} tokens'.format(bob.token_contract.balanceOf()))

# Alice deposits 3 of her coins to the plasma contract and gets 3 plasma nft utxos in return 
tokenId = 1
alice.deposit(tokenId)
alice.deposit(tokenId+1)
alice.deposit(tokenId+2)

# Alice's UTXOs are with id 0, 1 and 2.
utxo_id = 2
blk_num = 3

tx1 = alice.send_transaction(utxo_id, blk_num, 1, bob.token_contract.account.address)
tx4 = alice.send_transaction(utxo_id-1, blk_num-1, 1, charlie.token_contract.account.address)

# Alice's 3 coins are now checkpointed in block `1000` in the root chain
authority.submit_block() 

# Block has been submitted, transactions have been checkpointed. Charlie needs to send the last UTXO which will be sufficient proof that he can claim the token 

blk_num = 1000
submitted_block = alice.get_block(blk_num)
onchain_block = alice.root_chain.contract.functions.childChain(blk_num).call()[0].hex() # block 1000 is the block that the authority submitted
assert submitted_block.merklize_transaction_set().hex() == '0x' + onchain_block # check that child-chain and root-chain are in sync

# Bob is now checkpointed as the owner of `utxo_id`. He makes a spend to charlie by referencing block 1 which was in the child chain.
tx2 = bob.send_transaction(utxo_id, blk_num, 1, charlie.token_contract.account.address)

authority.submit_block()

prev_tx_blk_num = 1000
exiting_tx_blk_num = 2000
uid = utxo_id

# Charlie should be able to submit an exit by referencing blocks 0 and 1 which included his transaction. 
utxo_id = 2
prev_tx_blk_num = 1000
exiting_tx_blk_num = 2000


charlie.start_exit(utxo_id, prev_tx_blk_num, exiting_tx_blk_num)
# 
# # # After 8 days pass, charlie's exit should be finalizable

w3 = charlie.root_chain.w3 # get w3 instance
increaseTime(w3, 8 * 24 * 3600)
authority.finalize_exits()
# 
# # Charlie should now be able to withdraw the utxo which included token 2 to his wallet.
charlie.withdraw(2)
print('Alice has {} tokens'.format(alice.token_contract.balanceOf()))
print('Bob has {} tokens'.format(bob.token_contract.balanceOf()))
print('Charlie has {} tokens'.format(charlie.token_contract.balanceOf()))
# 
# # Plasma Cash with ERC721 tokens success :) 
