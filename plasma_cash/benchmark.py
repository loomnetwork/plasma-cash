from client.client import Client
from dependency_config import container

authority = Client(container.get_root('authority'),
                   container.get_token('authority'))
w3 = authority.root_chain.w3  # get w3 instance

coins_per_register = 5
coin_indices = range(0, coins_per_register)
number_of_extras = 100
extras_indices = range(0, number_of_extras)

extras = list(map(lambda index: Client(container.get_root('extras', index),
                                       container.get_token('extras', index)),
                  extras_indices))

deposits = {}
for index in extras_indices:
    extras[index].register()
    for coin_index in coin_indices:
        extras[index].deposit(index * coins_per_register + coin_index + 1)
    deposits[index] = list(map(lambda event: event['args'],
                               extras[index].get_all_deposits()))


for index in extras_indices:
    print(deposits[index])
    for coin_index in coin_indices:
        neighbor_index = (index + 1) % number_of_extras
        extras[index].send_transaction(
            deposits[index][coin_index]['slot'],
            deposits[index][coin_index]['blockNumber'], 1,
            extras[neighbor_index].token_contract.account.address)

authority.submit_block()
