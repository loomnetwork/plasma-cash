const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract("Plasma ERC721 - Multiple Validators and ERC721 tokens", async function(accounts) {

    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    // Alice registers and has 5 coins, and she deposits 3 of them.
    const INITIAL_COINS = 5;
    const DEPOSITED_COINS = 3;
    const COINS = [1, 2, 3];

    let cards, cards2;
    let plasma;
    let vmc;
    let events;
    let t0;
    let UTXO;
    let slots;

    const blk_1 = 1000
    const blk_2 = 2000

    let [owner, validator1, validator2, alice, bob, charlie, random_guy] = accounts;

    beforeEach(async function() {
        vmc = await ValidatorManagerContract.new({from: owner});
        plasma = await RootChain.new(vmc.address, {from: owner});
        cards = await CryptoCards.new(plasma.address);
        cards2 = await CryptoCards.new(plasma.address);

        await vmc.toggleToken(cards.address);
        await vmc.toggleToken(cards2.address);

        await vmc.toggleValidator(validator1);
        await vmc.toggleValidator(validator2);

        await cards.register({from: alice});
        await cards2.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);
        assert.equal(await cards2.balanceOf.call(alice), 5);

        // Deposit all of Alice's cards
        for (let i = 0; i < DEPOSITED_COINS; i ++) {
            await cards.depositToPlasma(COINS[i], {from: alice});
        }

        // Deposit all of Alice's cards from the other erc721
        for (let i = 0; i < DEPOSITED_COINS; i ++) {
            await cards2.depositToPlasma(COINS[i], {from: alice});
        }

        assert.equal((await cards.balanceOf.call(alice)).toNumber(), INITIAL_COINS - DEPOSITED_COINS);
        assert.equal((await cards2.balanceOf.call(alice)).toNumber(), INITIAL_COINS - DEPOSITED_COINS);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), DEPOSITED_COINS);
        assert.equal((await cards2.balanceOf.call(plasma.address)).toNumber(), DEPOSITED_COINS);

        const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
        events = await txlib.Promisify(cb => depositEvent.get(cb));

        // Check that events were emitted properly
        let coin;
        for (let i = 0; i < DEPOSITED_COINS ; i++ ) {
            coin = events[i].args;
            assert.equal(coin.blockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }

        for (let i = DEPOSITED_COINS; i < 2 * DEPOSITED_COINS ; i++ ) {
            coin = events[i].args;
            assert.equal(coin.blockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }
        UTXO = [{'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()}, 
            {'slot': events[5]['args'].slot, 'block': events[5]['args'].blockNumber.toNumber()}];
        slots = UTXO.map(u => u.slot)
    });

    it("Exit two coins from different contracts (submit block by 2 validators)", async function() {

        let alice_to_bob = txlib.createUTXO(UTXO[0].slot, UTXO[0].block, alice, bob);
        let alice_to_bob2 = txlib.createUTXO(UTXO[1].slot, UTXO[1].block, alice, bob);
        let txs = [alice_to_bob.leaf, alice_to_bob2.leaf];
        let tree_1000 = await txlib.submitTransactions(validator1, plasma, 1000,txs);

        let bob_to_charlie = txlib.createUTXO(UTXO[0].slot, 1000, bob, charlie);
        let bob_to_charlie2 = txlib.createUTXO(UTXO[1].slot, 1000, bob, charlie);
        txs = [bob_to_charlie.leaf, bob_to_charlie2.leaf];
        let tree_2000 = await txlib.submitTransactions(validator2, plasma, 2000, txs);


        t0 = await txlib.exit(plasma, charlie,
            UTXO[0].slot,

            { 'block': blk_2, 'tx': bob_to_charlie },
            tree_2000.createMerkleProof(UTXO[0].slot),

            { 'block': blk_1, 'tx': alice_to_bob },
            tree_1000.createMerkleProof(UTXO[0].slot),
        )

        t0 = await txlib.exit(plasma, charlie,
            UTXO[1].slot,

            { 'block': blk_2, 'tx': bob_to_charlie2 },
            tree_2000.createMerkleProof(UTXO[1].slot),

            { 'block': blk_1, 'tx': alice_to_bob2 },
            tree_1000.createMerkleProof(UTXO[1].slot),
        )
        await increaseTimeTo(t0 + t1 + t2);
        await plasma.finalizeExits(slots, {from: random_guy});
        await plasma.withdraw(UTXO[0].slot, {from: charlie});
        await plasma.withdraw(UTXO[1].slot, {from: charlie});

        assert.equal(await cards.balanceOf(charlie), 1);
        assert.equal(await cards2.balanceOf(charlie), 1);

        await txlib.withdrawBonds(plasma, charlie, 0.2)
    });

});
