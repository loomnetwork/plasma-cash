# Loom SDK - Plasma Cash Integration Spec

## Glossary:

1. Root Chain: Ethereum Mainnet
2. RootChain.sol: Plasma Cash smart contract, deployed on the Root Chain
3. Plasma Chain: Ethereum sidechain, in our case Loom DAppChain
4. Authority: Owner of the keypair which allows calling privileged functions such as `submitBlock` in RootChain.sol, as well as the operator of the Plasma Chain
5. UTXO: Unspent Transaction Output. When referring to the slot of a UTXO, we refer to the position of the UTXO in the Merkle Tree that includes it.

## Root Chain Interaction

Rootchain interactions come in two varieties: (1) the client calls functions of the `RootChain.sol` contract running on the Ethereum rootchain, or (2) events on the rootchain trigger actions on the Plasma chain.

1. Bindings to all contract functions of [RootChain.sol](https://github.com/loomnetwork/plasma-erc721/blob/master/server/contracts/Core/RootChain.sol). Given the ABI, Address and mainnet node endpoint (infura, localhost:8545 etc.) it should be able to call all the functions of the RootChain.sol contract successfully.

    - `challengeBefore(slot, prev_tx_bytes, exiting_tx_bytes, prev_tx_inclusion_proof, exiting_tx_inclusion_proof, sig, prev_tx_block_num, exiting_tx_block_num)`: Submit a fraud proof showing that someone is attempting to exit an coin with an invalid history.

    - `respondChallengeBefore(uint64 slot, uint challengingBlockNumber, bytes challengingTransaction, bytes proof)`: Allows a user to respond to a `challengeBefore` by providing information about the child of the transaction which the challenger claims is the last valid transaction involving a particular coin.

    - `challengeBetween(uint64 slot, uint challengingBlockNumber, bytes challengingTransaction, bytes proof)`: This function is used to challenge exists of double-spent coins. A cheat spends a coin deposited in block `n` in both blocks `n + 1` and `n + 2`. If the cheat attempts to exit the coin by referencing the transaction in block `n + 2`, a challenger must provide proof that the coin was also spent in block `n + 1` which is positioned between the deposit block `n` and the block with a second spend `n + 2`. Hence, the function is called `challengeBetween`.

    - `challengeAfter(uint64 slot, uint challengingBlockNumber, bytes challengingTransaction, bytes proof)`: This function is used to challenge exits of already-spent coins. A cheat will try to exit a coin included in block `n` and a challenger must produce a proof that the coin was spent in a later block, e.g. `n + 1` or `n + 2`. Therefore, the function is called `challengeAfter` because it requires a challenger to find a later block containing a spend of particular coin than the one indicated by the cheat.

    - `startExit(slot, previous_tx_blk_num, latest_tx_blk_num)`: Start exiting a coin at a particular UTXO.

    - `finalizeExits()`: Exit all coins whose exists are older than 7 days and have not been successfully challenged. Any successfully challenged coins should have their states changed to DEPOSITED and have their owner's bonds slashed.

    - `withdraw(uint64 slot)`: Withdraws a particular UTXO from `RootChain.sol` contract that has been exited from the Plasma chain.

    - `withdrawBonds()`: Allows a user to withdraw any bonds he was forced to submit along when initiating an exit.

2. Event handlers: Functionality to make a callback to a function whenever a certain type of event gets emitted from RootChain.sol.

This is discussed in more detail in the [Block Structure](#block-structure) section below.

## Data Types

_All data type sizes (i.e uint64, uint8 etc) are preliminary and may change_.

### `Transaction`:

Whenever a transaction is created, it gets signed by its owner and gets passed on to a new owner. Each Transaction/UTXO contains the following (ref impl in Python):
1. `uid(uint64)`:The slot of the UTXO - Currently uint64, subject to change.
2. `previousBlock(uint256)`: Each time a transaction is created, it MUST refer to a previous block which also included that transaction. A transaction is considered a “deposit transaction”, if it’s the first UTXO after a user deposits their coin in the Plasma Chain. This transaction mints coins from nowhere in the Plasma Chain and as a result its previous block is 0.
3. `denomination(uint256)`: How many coins are included in that UTXO. Currently this is always 1 since we’re using ERC721 tokens which are unique, however in future iterations this can be any number.
4. `new_owner(address)`: The new owner of the transaction.
5. `signature(bytes)`: Signature on the transaction’s hash
6. `hash(byte32)`: The hash of the RLP encoded unsigned transaction’s bytes. If the transaction is a deposit transaction (its prevblock is 0), its hash is the hash of its uid
7. `merkle_hash(byte32)`: The hash of the RLP encoded signed transaction’s bytes
8. `sender(address)`: The transaction’s sender, derived from the hash and the signature

### `Block`:

1. Transaction Set(`Transaction[]`): List of transactions included in Block
2. Merklized Transaction Set: All transactions in the block get sorted in a Sparse Merkle Tree of depth 64. The block’s merkle root gets generated, and is to be submitted at a future point to the `RootChain.sol` contract
3. `hash(byte32)`: The hash of the RLP encoded unsigned block’s bytes.
3. `merkle_hash(byte32)`: The block's merkle root from its included transactions
3. `signature(bytes)`: Signature on the block’s hash

Multiple spends of the same coin cannot be included in a single block since the Sparse Merkle Tree structure allocates a unique slot for each coin on a plasma chain.

### Sparse Merkle Tree (SMT)

A N-depth merkle tree where all of its leaves are initialized to H(0), where H is Ethereum’s keccak256. Full description in [ethresearch](https://ethresear.ch/t/plasma-cash-with-sparse-merkle-trees-bloom-filters-and-probabilistic-transfers/2006) and in the [relevant paper](https://eprint.iacr.org/2016/683.pdf). We have 3 implementations, in [Python](https://github.com/loomnetwork/plasma-erc721/blob/master/plasma_cash/utils/merkle/sparse_merkle_tree.py) and [Javascript](https://github.com/loomnetwork/plasma-erc721/blob/master/server/test/SparseMerkleTree.js) for the creation/validation, and additionally a [Solidity Smart Contract](https://github.com/loomnetwork/plasma-erc721/blob/master/server/contracts/Core/SparseMerkleTree.sol) for the on-chain validation.

## Plasma Chain

### Block Structure

The Plasma Chain supports multiple validators. It keeps the full state of the chain and is responsible for providing clients with any data regarding the state of the chain they ask for (data availability).

There are 2 kind of blocks:

1. **Deposit Blocks:** They get created in the `RootChain.sol` contract after a deposit, and include **only 1** transaction, namely the deposit transaction. The `RootChain.sol` contract stores the Block's hash only, which is the hash of the included transaction. The Plasma Chain, stores the whole Block, including the transaction's data.
2. **Plasma Blocks:** These blocks contain multiple transactions. Each time a transaction is created in the Plasma Chain, it gets added to the Plasma Block. The Plasma Block's root is the root of the Sparse Merkle Tree of N-depth, which gets generated from all its included transactions. Periodically, the Plasma Chain operator submits the current Plasma Block's root to the blockchain. The frequency of block submissions can be either some time constant, or can be depending on a metric such as Merkle Tree sparseness.

In other words, the Plasma Chain generates Deposit Blocks, in response to deposits, while `RootChain.sol` creates blocks in response to block submissions. The Plasma Chain holds the full information for each block, while `RootChain.sol` just stores a hash, which in the case of a Deposit Block is the hash of the deposit transaction, while in the case of a Plasma Block is the merkle root of the Sparse Merkle Tree which gets created from the block's transactions in the Plasma Chain.

The Plasma Chain generates Deposit Blocks by listening to `Deposit` type events from `RootChain.sol`. This is the event as it appears in the `RootChain.sol` contract:

```solidity
event Deposit(uint64 indexed slot, uint256 depositBlockNumber, uint64 denomination, address indexed from);
```

As an example, this is the Python source code for the callback which will create a Deposit Block in the Plasma Chain:
```python
# Watch all deposit events, callback to self._send_deposit
    deposit_filter = self.root_chain.watch_event('Deposit', self._send_deposit, 1)

def _send_deposit(self, event):
    ''' Called by event watcher and creates a deposit block '''
    slot = event['args']['slot']
    blknum = event['args']['depositBlockNumber']
    denomination = event['args']['denomination'] # currently always 1, to change in the future
    depositor = event['args']['from']
    deposit_tx = Transaction (slot, 0, denomination, depositor, incl_block=blknum)
    deposit_block = Block( [ deposit_tx ] ) # create a new plasma block on deposit
    self.blocks[blknum] = deposit_block
```

After every deposit transaction (previous block = 0), a block is created on the Plasma Chain, based on the parameters of the event emitted from RootChain.sol. That way, both the Plasma Chain state and the `RootChain.sol` contract have the same state for deposit blocks. In order to have separate block numbering between the two, we consider that Plasma Blocks are separated by N blocks, where every block between 2 Plasma Block can be a Deposit Block. As an example for N=1000:

```
0 -> 1 -> 2 -> 3 -> 1000 -> 1001 -> 1002 -> ... -> 2000 -> 2001 ...
```
Deposit Blocks are the blocks which are non-multiples of the interval N.

### Communication Protocol

As with all blockchain nodes, the Plasma Chain should also have an API to accept communication. We define the following:

1. `send_transaction`: A Plasma chain user submits a transaction for inclusion in a Plasma block. Transactions are currently limited to spends of a particular coin which means every transaction involves a user transfering a single coin to another user.
2. `submit_block` (authority only): Submits the current block root to `RootChain.sol`. Can only be called from the authority.
3. `get_current_block`: Get the latest block in the Plasma chain.
4. `get_block`: Get a block with a user-specified number.
5. `get_proof`: Get a merkle proof of inclusion of a particular uid in a particular block.

Reference implementation in Python [[1](https://github.com/loomnetwork/plasma-erc721/blob/master/plasma_cash/child_chain/server.py)][[2](https://github.com/loomnetwork/plasma-erc721/blob/master/plasma_cash/child_chain/child_chain.py)]

In order for a client to communicate with the Plasma chain, the client should send requests that the Plasma Chain understands. We can think of the Plasma Chain an Ethereum Node, and the client as web3.js. Therefore, the above should be implemented for a client that wants to talk to the Plasma Chain.

Reference implementation in Python [[1](https://github.com/loomnetwork/plasma-erc721/blob/master/plasma_cash/client/child_chain_service.py)[[2](https://github.com/loomnetwork/plasma-erc721/blob/master/plasma_cash/client/client.py)]


> Note: Any kind of data that  gets passed around between clients is encoded in RLP (import "github.com/ethereum/go-ethereum/rlp”, wiki)

TODO: Further explain how exits should be implemented in clients, poll the authority for data, local db schema for storing data etc.
