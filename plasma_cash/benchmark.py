from client.client import Client
from dependency_config import container
# from time import sleep

authority = Client(container.get_root('authority'),
                   container.get_token('authority'))
w3 = authority.root_chain.w3  # get w3 instance

coins_per_register = 5
number_of_extras = 8
extras_indices = range(0, number_of_extras)

extras = list(map(lambda index: Client(container.get_root('extras', index),
                                       container.get_token('extras', index)),
                  extras_indices))

deposits = {}
for index in extras_indices:
    extras[index].register()
    for coin_index in range(0, coins_per_register):
        extras[index].deposit(index * coins_per_register + coin_index + 1)
    deposits[index] = list(map(lambda event: event['args'],
                               extras[index].get_all_deposits()))


for index in extras_indices:
    print(deposits[index])
    extras[index].send_transaction(
        deposits[index][0]['slot'], deposits[index][0]['blockNumber'], 1,
        extras[(index + 1) % number_of_extras].token_contract.account.address)

authority.submit_block()
