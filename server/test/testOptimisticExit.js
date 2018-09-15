const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract.only("Plasma Cash - Optimistic Exits", async function(accounts) {

    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    // Alice registers and has 5 coins, and she deposits 3 of them.
    const ALICE_INITIAL_COINS = 5;
    const ALICE_DEPOSITED_COINS = 3;
    const COINS = [1, 2, 3];

    let cards;
    let plasma;
    let vmc;
    let events;
    let t0;

    let [authority, alice, bob, charlie, dylan, elliot, random_guy, random_guy2, challenger] = accounts;

    beforeEach(async function() {
        vmc = await ValidatorManagerContract.new({from: authority});
        plasma = await RootChain.new(vmc.address, {from: authority});
        cards = await CryptoCards.new(plasma.address);
        await vmc.toggleToken(cards.address);
        cards.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);

        let ret;
        for (let i = 0; i < ALICE_DEPOSITED_COINS; i++) {
            await cards.depositToPlasma(COINS[i], {from: alice});
        }


        assert.equal((await cards.balanceOf.call(alice)).toNumber(), ALICE_INITIAL_COINS - ALICE_DEPOSITED_COINS);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), ALICE_DEPOSITED_COINS);

        const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
        events = await txlib.Promisify(cb => depositEvent.get(cb));

        // Check that events were emitted properly
        let coin;
        for (let i = 0; i < events.length; i++) {
            coin = events[i].args;
            assert.equal(coin.blockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }

    });


        it("Challenge wrong parent optimistic exit", async function() {
            let UTXO = {'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()};


            // Alice gives the coin to dylan
            let alice_to_dylan = txlib.createUTXO(UTXO.slot, UTXO.block, alice, dylan);
            let txs = [ alice_to_dylan.leaf ];
            let tree_dylan = await txlib.submitTransactions(authority, plasma, txs);

            // Tx to Charlie from Bob referencing Bob's UTXO at block 1000
            let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
            txs = [ bob_to_charlie.leaf ];
            let tree_charlie = await txlib.submitTransactions(authority, plasma, txs);
            let exiting_tx_proof = tree_charlie.createMerkleProof(UTXO.slot)

            await plasma.startExit(
                UTXO.slot,
                bob_to_charlie.tx,
                exiting_tx_proof,
                bob_to_charlie.sig,
                bob,
                [1000, 2000],
                {'from': charlie, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo(t0 + t1 + t2);

            // Challenge with the invalid parent! There was no tx to Bob at the block!!
            let proof = tree_dylan.createMerkleProof(UTXO.slot);
            await plasma.challengeWrongParent(UTXO.slot, 1000, alice_to_dylan.tx, proof, {from: challenger});

            await plasma.finalizeExits({from: random_guy2});

            assertRevert(plasma.withdraw(UTXO.slot, {from : charlie}));

            assert.equal(await cards.balanceOf.call(alice), 2);
            assert.equal(await cards.balanceOf.call(bob), 0);
            assert.equal(await cards.balanceOf.call(charlie), 0);
            assert.equal(await cards.balanceOf.call(dylan), 0);
            assert.equal(await cards.balanceOf.call(plasma.address), 3);

            await txlib.withdrawBonds(plasma, challenger, 0.1);
        });

        it("Optimistic exit where the only data provided is block height and address", async function() {
            let UTXO = {'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()};
            // instead of providing this transaction in the exit, we just give its block number. anybody is welcome to challenge:) 
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [ alice_to_bob.leaf ];
            let tree_bob = await txlib.submitTransactions(authority, plasma, txs);

            // Tx to Charlie from Bob referencing Bob's UTXO at block 1000
            let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
            txs = [ bob_to_charlie.leaf ];
            let tree_charlie = await txlib.submitTransactions(authority, plasma, txs);

            let exiting_tx_proof = tree_charlie.createMerkleProof(UTXO.slot)

            await plasma.startExit(
                UTXO.slot,
                bob_to_charlie.tx,
                exiting_tx_proof,
                bob_to_charlie.sig,
                bob,
                [1000, 2000],
                {'from': charlie, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2});

            await plasma.withdraw(UTXO.slot, {from : charlie});

            assert.equal(await cards.balanceOf.call(alice), 2);
            assert.equal(await cards.balanceOf.call(bob), 0);
            assert.equal(await cards.balanceOf.call(charlie), 1);
            assert.equal(await cards.balanceOf.call(dylan), 0);
            assert.equal(await cards.balanceOf.call(plasma.address), 2);

            await txlib.withdrawBonds(plasma, charlie, 0.1);
        });

});
