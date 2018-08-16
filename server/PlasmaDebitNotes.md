# Plasma Debit

Previously the denomination field in a transaction was always set to 1, because
we were always transferring the whole coin. Now a UTXO off-chain will also have the denomination field set, to some number that must be less than the coin's . 

Plasma Debit's assumption is that a coin can have any value between 0 and its
initial deposit value. A 5/5 coin is owned by the user, a 3/5 coin allows the
user to withdraw 3/5 and the 2/5 go to the operator, while a 0/5 coin is woned
100% by the operator. 0/5 coins can be generated whne the operator makes a
deposit and can be used to bootstrap channels between the operator and users,
WITHOUT HAVING THE USERES TO MAKE DEPOSITS!

There are 2 kinds of transactions: 
- Ownership transfer transactions, which transfer the whole coin
- Denomination change transactions, which just alter the balance of the coin
  (which can be thought of paying fees to an operator)

A transaction that changes the denominations of two coins must meet the
following:
- Coin slots must be different
- Flow is always from address 1 to address 2 
- Amount >0 indicates that address1's coin gets reduced by amount, and address
  2's coin gets increased by the same amount, opposite applies for amount <0.

- Whenever the `authority` address makes a deposit, a new coin is minted to
  their address, as usual, however the coin's denomination is set to 0, to
indicate that they own the whole coin. 

When a coin is being exited, 

# Providing liquidity to a coin!

The operator can deposit into any coin from the main chain, which increments v
but not a. So if Bob has a 5 ETH coin and a current balance of 5 ETH (i.e. v is
5 and a is 5), then the operator can deposit 1 ETH into the coin, making it a 6
ETH coin where Bob has a balance of 5 ETH (i.e. v is 6 and a is 5), which gives
the operator an implied balance of 1 ETH.


## Sending atomic transactions

Ok, with proofs extended with these joint signatures, how can some other
entity that is receiving a partially spent Plasma Debit coin know that it is
receiving the highest nonce joint signature between the coinholder and the
Plasma operator?

Nope, because notarized transactions (those included in the Plasma Cash chain)
take strict precedence over the intermediate payment-channel-like transactions. <-- OK?

The transaction that transfers to Bob doesn’t have to mention or acknowledge
any of the intermediate states between Alice and the operator, because those
are all obsoleted as soon as the transfer to Bob is included on the Plasma Cash
chain. 

The transfer to Bob does have to state what the balance of the coin is
that Bob receives, though. <-- BALANCE 

Updates to the balance would not actually need to be included in the Plasma
blocks, since they require only the mutual consent of the coinholder and
operator. The exit game could be altered relatively easily to allow the
operator and current coinholder to instantly update their balance by exchanging <-- NEED TO ADD EXTRA STUFF TO THE EXIT GAME + NONCES
signatures on a state update (which would then be almost exactly equivalent to
a payment channel). 

> The only transactions that need to be notarized (i.e.
included in Plasma blocks) are those that either change the owner of the
channel, or involve multiple coins.

No, because you could instead just have the operator and user both sign the new
balance, along with an incrementing nonce representing which intermediate
balance it was. 

1) When withdrawing, the most recent committed Plasma Debit coin
would take priority, 
2) after that, the most recent (i.e. highest-nonce) set
of balances that was signed by both the operator and user would take priority.
This is the same basic mechanism as state channels.


^ NEED A CHALLENGE FOR AN EXIT WHERE WE REVEAL A SPEND WITH A HIGHER NONCE



## Liquidity Injections

Have a staking contract that handles all liquuidity. users deposit funds to the
contract in order to have a big reserve. We consider the liquidity reserve to
be governed by the validators. If we have 10 eth inside it, any liquidity
injection should be voted on. Plasma Cash/Deit withdrawals incur a X% fee. That
fee is paid out to the staking contract after each withdrawal (maybe make it
just an sstore and then make it a full withdrawal in contract terms - sounds
better!). After fees have benen withdrawn from Plasma to the dao, the dao will
give out the dividends. Users will get more fees based on the % of the coin
that the liquidity was provided.

T0 -> 30% 50% 20%, 10 coins
T1000000 -> 100 coins 
30%Plasma Cash/Deit withdrawals incur a X% fee. That fee is paid out to the
staking contract after each withdrawal (maybe make it just an sstore and then
make it a full withdrawal in contract terms - sounds better!). After fees have
benen withdrawn from Plasma to the dao, the dao will give out the dividends.
Users will get more fees based on the % of the coin that the liquidity was
provided.

T0 -> 30% 50% 20%, 10 coins, 3, 5, 2
T100 -> 100 coins, 30, 50, 20

1) U2 can withdraw up to STAKED/TOTAL = 50%
2) U2 withdraws 20 coins, new balance is 30
3) Reduce total stake by 20, new stake = 80
4) Recalculate percentages: 30/80 30/80 20/80

--

U1 can now withdrawl up to 3/8 of new stake: 30 coins

