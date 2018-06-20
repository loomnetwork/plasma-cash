import rlp
from client.client import Client
from dependency_config import container
from child_chain.transaction import Transaction, UnsignedTransaction

token_contract = container.get_token('alice')
tx = Transaction(5, 0, 1, token_contract.account.address)
tx_hex = rlp.encode(tx, UnsignedTransaction).hex()
print('Transaction(5, 0, 1, token_contract.account.address): {}'.format(tx_hex))
