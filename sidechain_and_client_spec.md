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

For each block storing (i) the Merkle root, (ii) the time the Merkle root was
submitted.

    fields = [
        ('transaction_set', CountableList(Transaction)),
        ('sig', binary)
    ]

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

### Transactions

Plasma transactions are initiated by ... the sidechain client. The Plasma
chain operator accumulates transactions submitted by clients and can decide to
submit a block containing those transactions whenever he sees fit. Adding
a block to the plasma chain requires registering that block's merkle root with
the `RootChain.sol` contract.

All blocks filled with transactions are assigned a number divisible by
`childBlockInterval` (currently 1000) in the `RootChain.sol` contract.

## Sidechain Client

The sidechain client interacts with both the Plasma sidechain and the Ethereum
rootchain.

### Sidechain interactions

#### 1. submit_transaction

#### 2. submit_block (authority only)

#### 3. get_current_block

Get the latest block in the Plasma chain.

#### 4. get_block

Get a block with a user-specified number.

#### 5. get_proof

Get a merkle proof of inclusion of a particular uid in a particular block.

### Rootchain interactions

#### 1. challengeBefore

#### 2. challengeBetween

#### 3. challengeAfter

#### 4. startExit

#### 5. finalizeExits

#### 6. withdraw

Withdraws a particular UTXO from `RootChain.sol` contract that has been exited
from the Plasma chain.

#### 7. withdrawBonds (to be implemented)

Allows a user to withdraw any bonds he was forced to submit along when
initiating an exit.

## Sidechain & Sidechain Client Integration plans

Re-write all truffle tests as Python tests using 1) an Ethereum blockchain 2) the Plasma sidechain and 3) the sidechain client.
