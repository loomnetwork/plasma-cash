# Sidechain implementation

Currently, both the Plasma Cash sidechain and the sidechain client are
implemented in Python though the final implementation will be in Go. Below is
a description of the current sidechain functionality. Based on this
description, it will be decided what exactly should be ported to Go.

## Plasma Cash Sidechain

The state maintained by the sidechain is described followed by its
functionality.

### Plasma Cash state

The Plamsa chain maintains several pieces of state during operation:

- authority address
    + the only address allowed to submit blocks
- interface to rootchain
    + url and portnumber over which to communicate with Ethereum blockchain
    + abi of RootChain.sol
- address of the ERC721 token contract on the rootchain
- list of plasma blocks
    + a record of all transactions on the plasma chain, described in more
        detail below

#### List of Plasma Blocks

For each block of the Plasma chain, the `RootChain.sol` contract stores the
Merkle root of the block's transactions and the time the Merkle root was
submitted. The Plasma chain maintains a list of blocks each of which contains
a set of a transactions and a signature.

### Deposits

Deposits are initiated automatically on the Plasma chain when a `Deposit` event
is registered on the rootchain. All deposit blocks are assigned a number not
divisible by `childBlockInterval` (currently 1000) in the `RootChain.sol` contract.

This is the event as it appears in the `RootChain.sol` contract:

```javascript
    event Deposit(uint64 indexed slot, uint256 depositBlockNumber, uint64 denomination, address indexed from);
```

And this is the python source that causes `_send_deposit` to be invoked as
a callback on every `Deposit` event emitted by the `RootChain.sol` contract:

```python
    # Watch all deposit events, callback to self._send_deposit
    deposit_filter = self.root_chain.watch_event('Deposit', self._send_deposit, 1)

def _send_deposit(self, event):
    ''' Called by event watcher and creates a deposit block '''
    slot = event['args']['slot']
    blknum = event['args']['depositBlockNumber']
    denomination = event['args']['denomination'] # currently always 1, to change in the future
    depositor = event['args']['from']
    deposit_tx = Transaction (slot, blknum, denomination, depositor)
    deposit_block = Block( [ deposit_tx ] ) # create a new plasma block on deposit

    self.blocks[blknum] = deposit_block
```

### Transactions & Non-deposit Blocks

Plasma transactions are initiated by the sidechain client on behalf of any user
of the Plasma chain. The Plasma chain operator receives transactions submitted
by clients, submitting a block of the accumulated transactions at will. Adding
a block to the plasma chain requires registering that block's merkle root with
the `RootChain.sol` contract by calling the `submitBlock` function.

A Plasma block can contain only a single spend of a particular coin.

All blocks filled with transactions are assigned a number divisible by
`childBlockInterval` (currently 1000) in the `RootChain.sol` contract.

## Sidechain Client

The sidechain client interacts with both the Plasma sidechain and the Ethereum
rootchain.

### Sidechain interactions

In the current implementation of sidechain interactions, the client calls
functions of the Plasma sidechain via an http interface.

#### 1. send_transaction

A Plasma chain user submits a transaction for inclusion in a Plasma block.
Transactions are currently limited to spends of a particular coin.

#### 2. submit_block (authority only)

Initiates a block submission from the Plasma Chain onto the rootchain. A block
submission consists of only a Plasma block's transaction merkle root.

#### 3. get_current_block

Get the latest block in the Plasma chain.

#### 4. get_block

Get a block with a user-specified number.

#### 5. get_proof

Get a merkle proof of inclusion of a particular uid in a particular block.

### Rootchain interactions

For rootchain interactions, the client calls functions of the `RootChain.sol`
contract running on the Ethereum rootchain.

#### 1. challengeBefore

Submit a fraud proof showing that someone is attempting to exit an
coin with an invalid history.

#### 2. challengeBetween

Submit a fraud proof showing that an exited coin was double spent.

#### 3. challengeAfter

Submit a fraud proof showing that an exited coin was already spent.

#### 4. startExit

Start exiting a coin at a particular UTXO.

#### 5. finalizeExits

Exit all coins whose exists are older than 7 days and have not been
successfully challenged. Any successfully challenged coins should have their
states changed to DEPOSITED and have their owner's bonds slashed.

#### 6. withdraw

Withdraws a particular UTXO from `RootChain.sol` contract that has been exited
from the Plasma chain.

#### 7. withdrawBonds (to be implemented)

Allows a user to withdraw any bonds he was forced to submit along when
initiating an exit.

## Sidechain & Sidechain Client Integration plans

Re-write all truffle tests as Python tests using 1) an Ethereum blockchain 2) the Plasma sidechain and 3) the sidechain client.
