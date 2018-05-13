from client.authority_client import Client
from dependency_config import container

from child_chain.transaction import Transaction
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

tokenId = 2
alice.token_contract.deposit(tokenId, data = '12345789')
# tx1 = alice.send_transaction(0, tokenId, bob.token_contract.account.address)
# alice.submit_block() # signed by authority account
tx2 = bob.send_transaction(0, tokenId, charlie.token_contract.account.address)

print ('Current block has {} transactions'.format(alice.get_current_block().transaction_set))

blk = alice.get_current_block()
authority.submit_block()

# Block has been submitted, transactions have been checkpointed. Charlie needs to send the last UTXO which will be sufficient proof that he can claim the token 

