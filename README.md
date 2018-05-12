# Plasma Cash - ERC721 Version

## Architecture

Client connects to the Child Chain Service, the Root Chain Plasma Contract and the Root Chain Token Contract

Send tx:
Client -> Child Chain Service sendtx -> Child Chain server listens to that request and calls child chain sendtransaction, tx gets added to the block 

Child Chain ALWAYS listens for events on the RootChain contract and acts on them accordingly


## Demo

1. `npm run migrate:dev` on the server directory, contracts are now deployed
2. `./init` on server directory 
3. `python demo.py`

1. Alice registers and is given coins 1-5 
2. Alice deposits `Coin 1`, `Coin 2`, `Coin 3` in the plasma chain
3. Current block now has 3 transactions, all generated from thin air
3. Alice sends coin 1 to Bob, current block has 4 transactions
4. Bob sends coin 1 to Charlie, current block has 5 transactions
5. Operator calls submitBlock, checkpointing the block merkle root which includes the transaction that gives ownership to charlie
5. Charlie tries to exit coin 1, Alice & Bob should be unable to challenge
6. Challenge period passes (simulate with evmAdvancetime), exits get finalized
7. Charlie is able to withdraw coin 1.

## TODO

Challenges in contract and client, 

## Dev notes
### Using monkeypatched web3.py 4.2.1 version for ganache issues
https://github.com/jordanjambazov/web3.py/commit/a61b382dd5c11de1102779868d377408c7590839
Also https://github.com/ethereum/web3.py/pull/827

###Signing locally
http://web3py.readthedocs.io/en/latest/web3.eth.account.html#prepare-message-for-ecrecover-in-solidity

