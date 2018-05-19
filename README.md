# Plasma Cash - ERC721 Version

## Architecture

Client connects to the Child Chain Service, the Root Chain Plasma Contract and the Root Chain Token Contract

Send tx:
Client -> Child Chain Service sendtx -> Child Chain server listens to that request and calls child chain sendtransaction, tx gets added to the block 

Child Chain ALWAYS listens for events on the RootChain contract and acts on them accordingly


## Demo

1. `npm run migrate:dev` on the server directory, contracts are now deployed
3. `python demo.py` on plasma-cash directory

1. Alice registers and is given coins 1-5 
2. Alice deposits `Coin 1`, `Coin 2`, `Coin 3` in the plasma chain
3. 3 Deposit Blocks have been generated in the child chain, each one having 1 UTXO at slots 0,1,2 repsectively
4. Alice sends a Coin 1 to Bob, adding a transaction to the current block. 
5. Operator calls submitBlock, checkpointing the block merkle root which includes the transaction that gives ownership to charlie. At this point, both the child chain and the root chain, have checkpointed alice's transaction at block number 1000.
6. Bob transfers the previous UTXO to Charlie and the operator submits that block as well
5. Charlie tries to exit coin 1, Alice & Bob should be unable to challenge
6. After challenge period passes, charlie should be able to withdraw his coin

## TODO

- Challenges in contract and client
- Verify exiting / optimize for size



## Dev notes
### Using monkeypatched web3.py 4.2.1 version for ganache issues
https://github.com/jordanjambazov/web3.py/commit/a61b382dd5c11de1102779868d377408c7590839
Also https://github.com/ethereum/web3.py/pull/827

###Signing locally
http://web3py.readthedocs.io/en/latest/web3.eth.account.html#prepare-message-for-ecrecover-in-solidity

