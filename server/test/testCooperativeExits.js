const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");

const SparseMerkleTree = require('./SparseMerkleTree.js');

import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const UTXO = require('./UTXO.js')

const Promisify = (inner) =>
new Promise((resolve, reject) =>
        inner((err, res) => {
            if (err) {
                reject(err);
            } else {
                resolve(res);
            }
        })
);

contract("Plasma ERC721 - Cooperative Exits, no challenges", async function(accounts) {

    const UTXO_SLOT = 2;
    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    // Alice registers and has 5 coins, and she deposits 3 of them.
    const ALICE_INITIAL_COINS = 5;
    const ALICE_DEPOSITED_COINS = 3;

    let cards;
    let plasma;
    let t0;

    let [authority, alice, bob, charlie, dylan, elliot, random_guy, random_guy2, challenger] = accounts;

    let data;
    let to_alice = [];

    beforeEach(async function() {
        plasma = await RootChain.new({from: authority});
        cards = await CryptoCards.new(plasma.address);
        plasma.setCryptoCards(cards.address);
        cards.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);


        let ret;
        for (let i = 0; i < ALICE_DEPOSITED_COINS; i ++) {
            ret = UTXO.createUTXO(i, 0, alice, alice); data = ret[0];
            await cards.depositToPlasmaWithData(i+1, data, {from: alice});
            to_alice.push(ret);
        }


        assert.equal((await cards.balanceOf.call(alice)).toNumber(), ALICE_INITIAL_COINS - ALICE_DEPOSITED_COINS);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), ALICE_DEPOSITED_COINS);

        const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
        const events = await Promisify(cb => depositEvent.get(cb));

        // Check that events were emitted properly
        let coin;
        for (let i = 0; i < events.length; i++) {
            coin = events[i].args;
            assert.equal(coin.slot.toNumber(), i);
            assert.equal(coin.depositBlockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }

    });

    describe('Exit of UTXO 2 (Coin 3)', function() {
        it('Directly after its deposit', async function() {
            // Prevblock = 0, Exit block = deposit block
            let ret = UTXO.createUTXO(2, 0, alice, alice);
            let utxo = ret[0];
            let sig = ret[1];

            await plasma.startExit(
                     UTXO_SLOT,
                    '0x', utxo,
                    '0x0', '0x0', 
                     sig,
                     0, 3,
                     {'from': alice, 'value': web3.toWei(0.1, 'ether')}
            );

            t0 = (await web3.eth.getBlock('latest')).timestamp;

            await increaseTimeTo(t0 + t1);
            await plasma.finalizeExits({from: random_guy2 });
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 3);

            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });
            await plasma.withdraw(UTXO_SLOT, {from : alice });
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 3);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);
        });

        it('After 1 Plasma-Chain transfer', async function() {
            // Prevblock = deposit block, Exit block = normal block

            // Create a UTXO to Bob from Alice and sign it. Refer to Alice's deposit transaction at block 3
            let to_bob = UTXO.createUTXO(UTXO_SLOT, 3, alice, bob); 
            let txs = [ to_bob[2] ]
            // Authority submits a block to plasma with that transaction included
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            // Concatenate the 2 signatures. `to_alice` is the deposit transaction from `beforeEach` block.
            let sigs = to_alice[2][1] + to_bob[1].replace('0x','');
            let exiting_tx_proof = tree_bob.createMerkleProof(UTXO_SLOT)

            let prev_tx = to_alice[2][0];
            let exiting_tx = to_bob[0];

            plasma.startExit(
                    UTXO_SLOT,
                    prev_tx , exiting_tx,
                    '0x0', exiting_tx_proof,
                    sigs,
                    3, 1000, // Prev tx was included in block 3, exiting tx was included in block 1000
                     {'from': bob, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // Even though the coin still belongs to Alice, it is in the `EXITING` state so it shouldn't be possible for her to exit it
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));

            // If enough time hasn't passed, neither bob nor alice should be able to withdraw the coin
            await increaseTimeTo(t0 + t1);
            await plasma.finalizeExits({from: random_guy2 });
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : bob }));
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 3);

            // After the exit is matured and finalized, bob can withdraw the coin.
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));
            await plasma.withdraw(UTXO_SLOT, {from : bob });
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 1);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);

            // Bob is also able to withdraw his deposit bond of 0.1 ether
            await plasma.withdrawBonds({from: bob });
            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, bob);
            assert.equal(withdraw.amount, web3.toWei(0.1, 'ether'));
        });
        
        it("After 2 Plasma-Chain transfers", async function() {
            // Prevblock = normal block, Exit block = normal block
            let to_bob = UTXO.createUTXO(UTXO_SLOT, 3, alice, bob);
            let txs = [ to_bob[2] ];
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            // Tx to Charlie from Bob referencing Bob's UTXO at block 1000
            let to_charlie = UTXO.createUTXO(UTXO_SLOT, 1000, bob, charlie);
            txs = [ to_charlie[2] ];
            let tree_charlie = await UTXO.submitTransactions(authority, plasma, txs);

            // Concatenate the 2 signatures
            let sigs = to_bob[1] + to_charlie[1].replace('0x', '');

            let prev_tx_proof = tree_bob.createMerkleProof(UTXO_SLOT)
            let exiting_tx_proof = tree_charlie.createMerkleProof(UTXO_SLOT)

            let prev_tx = to_bob[0];
            let exiting_tx = to_charlie[0];

            plasma.startExit(
                    UTXO_SLOT,
                    prev_tx, exiting_tx,
                    prev_tx_proof, exiting_tx_proof, 
                    sigs,
                    1000, 2000,
                     {'from': charlie, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // Even though the coin still belongs to Alice, it is in the `EXITING` state so it shouldn't be possible for her to exit it
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));

            // If enough time hasn't passed, none of bob, alice or charlie should be able to withdraw the coin
            await increaseTimeTo(t0 + t1);
            await plasma.finalizeExits({from: random_guy2 });
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : bob }));
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : charlie }));
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 3);

            // After the exit is matured and finalized, Charlie can withdraw the coin.
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : alice }));
            assertRevert(plasma.withdraw(UTXO_SLOT, {from : bob }));
            await plasma.withdraw(UTXO_SLOT, {from : charlie });
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 1);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);

            // Charlie is also able to withdraw his deposit bond of 0.1 ether
            await plasma.withdrawBonds({from: charlie });
            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, charlie);
            assert.equal(withdraw.amount, web3.toWei(0.1, 'ether'));
        });

    });

    // If it works for 2 coins, proof by induction it will work for N coins >2
    describe('Exit of UTXO 1 and 2', function() {
        it('Alice gives Bob 2 coins who exits both', async function() {
            let to_bob = [];
            let txs = [];
            let tx;
            for (let i = 1 ; i < 3; i++) {
                tx = UTXO.createUTXO(i, i+1, alice, bob);
                to_bob.push(tx);
                txs.push(tx[2]); // push leaf
            }

            // Tree contains both transactions
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            // Block 1000 has now been checkpointed with both transactions that give ownership of the coins to Bob
            // UTXO 1 was deposited at Block 2, UTXO 2 was created at block 3, hence the i+1 in the for loop

            let prev_tx, exiting_tx, exiting_tx_proof, sigs;
            for (let i = 1 ; i < 3; i++) {
                sigs = to_alice[i][1] + to_bob[i-1][1].replace('0x','');
                exiting_tx_proof = tree_bob.createMerkleProof(i)
                prev_tx = to_alice[i][0];
                exiting_tx = to_bob[i-1][0];

                plasma.startExit(
                        i,
                        prev_tx , exiting_tx,
                        '0x0', exiting_tx_proof,
                        sigs,
                        i+1, 1000,
                         {'from': bob, 'value': web3.toWei(0.1, 'ether')}
                );
            }

            // Bob exit of UTXO 2
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // Even though the coin still belongs to Alice, it is in the `EXITING` state so it shouldn't be possible for her to exit it
            for (let i = 1 ; i < 3; i++) {
                assertRevert(plasma.withdraw(i, {from : alice }));
            }

            // If enough time hasn't passed, neither bob nor alice should be able to withdraw the coin
            await increaseTimeTo(t0 + t1);
            await plasma.finalizeExits({from: random_guy2 });
            for (let i = 1 ; i < 3; i++) {
                assertRevert(plasma.withdraw(i, {from : alice }));
                assertRevert(plasma.withdraw(i, {from : bob }));
            }
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 3);

            // After the exit is matured and finalized, bob can withdraw the coin.
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });
            for (let i = 1 ; i < 3; i++) {
                assertRevert(plasma.withdraw(i, {from : alice }));
                await plasma.withdraw(i, {from : bob });
            }
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 1);

            // Bob is also able to withdraw his deposit bonds of 0.2 ether for 2 exits
            await plasma.withdrawBonds({from: bob });
            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, bob);
            assert.equal(withdraw.amount, web3.toWei(0.1 * 2, 'ether'));
        });

        it('Alice gives Bob 2 coins, he exits 1 and gives another to Charlie who also exits it', async function() {
            let to_bob = [];
            let txs = [];
            let tx;
            for (let i = 1 ; i < 3; i++) {
                tx = UTXO.createUTXO(i, i+1, alice, bob);
                to_bob.push(tx);
                txs.push(tx[2]); // push leaf
            }

            // Tree contains both transactions
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            // Bob has ownership of the 2 coin's and now he gives 1 to Charlie
            let to_charlie = UTXO.createUTXO(2, 3, bob, charlie);
            txs = [ to_charlie[2] ];
            let tree_charlie = await UTXO.submitTransactions(authority, plasma, txs);

            // Bob exits UTXO 1
            let sigs = to_alice[1][1] + to_bob[0][1].replace('0x','');
            let exiting_tx_proof = tree_bob.createMerkleProof(1)
            let prev_tx = to_alice[1][0];
            let exiting_tx = to_bob[0][0];

            plasma.startExit(
                    1,
                    prev_tx , exiting_tx,
                    '0x0', exiting_tx_proof,
                    sigs,
                    2, 1000,
                     {'from': bob, 'value': web3.toWei(0.1, 'ether')}
            );

            // Charlie will now exit UTXO 2
            sigs = to_bob[1][1] + to_charlie[1].replace('0x', '');
            let prev_tx_proof = tree_bob.createMerkleProof(2)
            exiting_tx_proof = tree_charlie.createMerkleProof(2)
            prev_tx = to_bob[1][0];
            exiting_tx = to_charlie[0];

            plasma.startExit(
                    2,
                    prev_tx , exiting_tx,
                    prev_tx_proof, exiting_tx_proof,
                    sigs,
                    1000, 2000,
                     {'from': charlie, 'value': web3.toWei(0.1, 'ether')}
            );

            // Bob exit of UTXO 2
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // After the exit is matured and finalized, bob and charlie can withdraw the coin.
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });

            await plasma.withdraw(1, {from : bob });
            await plasma.withdraw(2, {from : charlie });
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 1);
            assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 1);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 1);

            // Bob is also able to withdraw his deposit bonds of 0.2 ether for 2 exits
            await plasma.withdrawBonds({from: bob });
            await plasma.withdrawBonds({from: charlie });
            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, bob);
            assert.equal(withdraw.amount, web3.toWei(0.1, 'ether'));
            withdraw = e[1].args;
            assert.equal(withdraw.from, charlie);
            assert.equal(withdraw.amount, web3.toWei(0.1, 'ether'));

        });

        it('Alice gives Bob 2 coins, who gives both to Charlie who exits both', async function() {
            let to_bob = [];
            let txs = [];
            let tx;
            for (let i = 1 ; i < 3; i++) {
                tx = UTXO.createUTXO(i, i+1, alice, bob);
                to_bob.push(tx);
                txs.push(tx[2]); // push leaf
            }

            // Tree contains both transactions
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            let to_charlie = [];
            txs = [];
            for (let i = 1 ; i < 3; i++) {
                tx = UTXO.createUTXO(i, 1000, bob, charlie);
                to_charlie.push(tx);
                txs.push(tx[2]); // push leaf
            }

            let tree_charlie = await UTXO.submitTransactions(authority, plasma, txs);

            let prev_tx, exiting_tx, prev_tx_proof, exiting_tx_proof, sigs;
            for (let i = 1 ; i < 3; i++) {
                sigs = to_bob[i-1][1] + to_charlie[i-1][1].replace('0x','');
                prev_tx_proof = tree_bob.createMerkleProof(i)
                exiting_tx_proof = tree_charlie.createMerkleProof(i)
                prev_tx = to_bob[i-1][0];
                exiting_tx = to_charlie[i-1][0];

                plasma.startExit(
                        i,
                        prev_tx , exiting_tx,
                        prev_tx_proof, exiting_tx_proof,
                        sigs,
                        1000, 2000,
                         {'from': charlie, 'value': web3.toWei(0.1, 'ether')}
                );
            }
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // After the exit is matured and finalized, Charlie can withdraw the coin.
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });
            for (let i = 1 ; i < 3; i++) {
                await plasma.withdraw(i, {from : charlie });
            }
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 1);

            // Charlie is also able to withdraw his deposit bonds of 0.2 ether for 2 exits
            await plasma.withdrawBonds({from: charlie });
            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, charlie);
            assert.equal(withdraw.amount, web3.toWei(0.1 * 2, 'ether'));

        });

        it('Alice gives Bob and Charlie 1 coin, they exit each one', async function() {
        });

        it('Alice gives Bob and Charlie 1 coin, Bob gives Charlie his coin, charlie exits it', async function() {
        });
    });

});


