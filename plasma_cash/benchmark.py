from client.client import Client
from dependency_config import container
from utils.utils import increaseTime

authority = Client(
    container.get_root("authority"), container.get_token("authority")
)
w3 = authority.root_chain.w3  # get w3 instance

child_block_interval = 1000
block_iterations = 4
coins_per_register = 5
num_deposited_coins = 2
coin_indices = range(0, num_deposited_coins)
number_of_players = 5
players_indices = range(0, number_of_players)

players = list(
    map(
        lambda index: Client(
            container.get_root("players", index),
            container.get_token("players", index),
        ),
        players_indices,
    )
)

# Step 1: All players register to the game and are given 5 coins each.
# Step 2: They then proceed to deposit 2 of their coins to Plasma Cash.
deposits = {}
for index in players_indices:
    players[index].register()
    for coin_index in coin_indices:
        players[index].deposit(index * coins_per_register + coin_index + 1)
    deposits[index] = list(
        map(lambda event: event["args"], players[index].get_all_deposits())
    )
    print("Player {} deposited coins: {}".format(index, deposits[index]))


print("STEP 3")
# Step 3: Each player gives their deposited coins to the next player
# 1000 players * 2 coins = 2k Plasma transactions per round
# This loops `block_iterations` times.
for iteration in range(block_iterations):
    for index in players_indices:
        neighbor_index = (index + 1) % number_of_players
        deposit_index = (index - iteration) % number_of_players
        for coin_index in coin_indices:
            prev_block = (
                deposits[index][coin_index]["blockNumber"]
                if iteration == 0
                else iteration * child_block_interval
            )
            print(
                '{}: PLAYER {} to {} : Coin {} from block {}'.format(
                    iteration, index, neighbor_index,
                    deposits[deposit_index][coin_index]['slot'],
                    prev_block
                )
            )
            players[index].send_transaction(
                deposits[deposit_index][coin_index]["slot"],
                # prev_block
                deposits[deposit_index][coin_index]["blockNumber"],
                players[neighbor_index].token_contract.account.address
            )
    authority.submit_block()

print("STEP 4")
# Step 4: All players initiate an exit for the coins they own.
# Since each player gave their coin to their neighbour, player `i`
# now owns the coins that player `(i-block_iterations) % num_players`
# initially had. Everyone initializes their exits
# by referencing the last 2 blocks
for index in players_indices:
    received = (index - block_iterations) % number_of_players
    for coin_index in coin_indices:
        print(
            'PLAYER {} exiting {} from {}'.format(index, coin_index, received)
        )
        slot = deposits[received][coin_index]["slot"]
        prev_block = (
            deposits[received][coin_index]["blockNumber"]
            if block_iterations == 1
            else (block_iterations - 1) * child_block_interval
        )
        players[index].start_exit(
            slot, prev_block, block_iterations * child_block_interval
        )


increaseTime(w3, 8 * 24 * 3600)
# Somebody can finalize all exits, or each user can finalize their own
# authority.finalize_exits()

# Final step: Each user finalizes their exit after challenge period
# has passed and then withdraws their coins.
for index in players_indices:
    received = (index - block_iterations) % number_of_players
    for coin_index in coin_indices:
        slot = deposits[received][coin_index]["slot"]
        players[index].finalize_exit(slot)
        players[index].withdraw(slot)
        print("Player {} withdrew coin: {}".format(index, slot))

print('Benchmarking done :)')
