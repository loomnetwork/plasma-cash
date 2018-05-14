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

tokenId = 2
alice.token_contract.deposit(tokenId, data = '12345789')
tx1 = alice.send_transaction(0, tokenId, bob.token_contract.account.address)
tx2 = bob.send_transaction(0, tokenId, charlie.token_contract.account.address)

print ('Current block has {} transactions'.format(alice.get_current_block().transaction_set))

authority.submit_block()

# Block has been submitted, transactions have been checkpointed. Charlie needs to send the last UTXO which will be sufficient proof that he can claim the token 

# If alice wants to exit her deposit, still WIP
# deposit_tx = Transaction(1000, 2, alice.token_contract.account.address)
# deposit_bytes = rlp.encode(deposit_tx, UnsignedTransaction)
# alice.exit(deposit_bytes, deposit_bytes, deposit_bytes)

# If charlie wants to exit the tx he received, UNIMPLEMENTED
# tx1_bytes = rlp.encode(tx1, UnsignedTransaction)
# tx2_bytes = rlp.encode(tx2, UnsignedTransaction)
# charlie.exit(tx1_bytes, tx2_bytes, tx2.sig)


