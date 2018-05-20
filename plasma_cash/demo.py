from client.authority_client import Client
from dependency_config import container

from child_chain.transaction import UnsignedTransaction, Transaction
import rlp


alice = Client(container.root_chain)
bob = Client(container.root_chain, container.bob)
charlie = Client(container.root_chain, container.charlie)
authority = alice
authority.key = authority.root_chain.account.privateKey # hack to give proper private key to authority

# Give alice 5 tokens
print ('Current block has {} transactions'.format(alice.get_current_block().transaction_set))

alice.token_contract.register()

print('Alice has {} tokens'.format(alice.token_contract.balanceOf()))
print('Bob has {} tokens'.format(bob.token_contract.balanceOf()))
print('Charlie has {} tokens'.format(bob.token_contract.balanceOf()))

# Alice deposits 3 of her coins to the plasma contract and gets 3 plasma nft utxos in return 
tokenId = 1
alice.token_contract.deposit(tokenId)
alice.token_contract.deposit(tokenId+1)
alice.token_contract.deposit(tokenId+2)

# Alice's UTXOs are with id 0, 1 and 2.
utxo_id = 2
blk_num = 3

tx1 = alice.send_transaction(utxo_id, blk_num, 1, bob.token_contract.account.address)

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

# Charlie should be able to submit an exit by referencing blocks 0 and 1 which included his transaction. 
# charlie.start_exit(utxo_id, 1000, 2000)


# # After 7 days pass, charlie's exit should be finalizable
# 
# authority.finalize_exits()
# 
# # Charlie should now be able to withdraw the utxo which included token 2 to his wallet.
# 
# charlie.withdraw(tokenId)

# Plasma Cash with ERC721 tokens success :) 
