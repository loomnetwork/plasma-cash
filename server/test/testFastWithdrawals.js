const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
const FastWithdrawal = artifacts.require("FastWithdrawal");
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
    let fast;
    let events;
    let t0;

    let [owner, validator1, validator2, alice, bob, charlie, random_guy] = accounts;

    beforeEach(async function() {
        vmc = await ValidatorManagerContract.new({from: owner});
        plasma = await RootChain.new(vmc.address, {from: owner});
        cards = await CryptoCards.new(plasma.address);
        cards2 = await CryptoCards.new(plasma.address);
        fast  = await FastWithdrawal.new(plasma.address);

        await vmc.toggleToken(cards.address);
        await vmc.toggleToken(cards2.address);

        await vmc.toggleValidator(validator1);
        await vmc.toggleValidator(validator2);

        await cards.register({from: alice});
        await cards2.register({from: bob});
        assert.equal(await cards.balanceOf.call(alice), 5);
        assert.equal(await cards2.balanceOf.call(bob), 5);

        let ret;
        // Deposit all of Alice's cards
        for (let i = 0; i < DEPOSITED_COINS; i ++) {
            await cards.depositToPlasma(COINS[i], {from: alice});
        }

        // Deposit all of Alice's cards from the other erc721
        for (let i = 0; i < DEPOSITED_COINS; i ++) {
            await cards2.depositToPlasma(COINS[i], {from: bob});
        }

        assert.equal((await cards.balanceOf.call(alice)).toNumber(), INITIAL_COINS - DEPOSITED_COINS);
        assert.equal((await cards2.balanceOf.call(bob)).toNumber(), INITIAL_COINS - DEPOSITED_COINS);
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
            assert.equal(coin.from, bob);
        }
    });

    describe('Fast exit of ERC721 where buyer waits 7 days and gives the receiver one of the coins they wanted', function() {
        it("Alice gives coin to Bob who gives it to Charlie. Charlie fast exits and demands a (address, uid) pair", async function() {
            let UTXO = {'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()}

            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf];
            let tree_1000 = await txlib.submitTransactions(validator1, plasma, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
            txs = [bob_to_charlie.leaf];
            let tree_2000 = await txlib.submitTransactions(validator2, plasma, txs);

            // Charlie now exits both coins.

            let prev_tx = alice_to_bob.tx;
            let exiting_tx = bob_to_charlie.tx;
            let prev_tx_proof = tree_1000.createMerkleProof(UTXO.slot)
            let exiting_tx_proof = tree_2000.createMerkleProof(UTXO.slot)
            let sig = bob_to_charlie.sig;
            // for coins 1,2,3 from the other cards contract. 
            // Atomic swap during exit:)
            // let whitelist_coins = [ 4, 5];
            let whitelist_coins = 4;
            let contractAddress = cards2.address

            let blocks = new Array(1000, 2000);

            // Essentially the owner of the exit is the contract
            // And charlie is the owner of the exit in the contract
            await fast.startExit(
                contractAddress,
                whitelist_coins,
                UTXO.slot,
                prev_tx, exiting_tx,
                prev_tx_proof, exiting_tx_proof,
                sig,
                blocks,
                {'from': charlie, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            let buyout = whitelist_coins;
            // // Approve and buy - this can be done in 1 transaction 
            // // by making FastWithdrawal an erc721 receiver
            await cards2.approve(fast.address, buyout, {from: bob});
            await fast.buyExit(
                contractAddress,
                UTXO.slot, 
                buyout,
                {'from': bob}
            );

            // // Bob bought charlie's exit with `buyout`
            // // Charlie now owns Bob's coin in the cards2 contract
            assert.equal(await cards2.ownerOf(buyout), charlie);
            assert.equal(await cards2.balanceOf(charlie), 1);

            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy});

            await fast.withdraw(cards.address, UTXO.slot, bob, {from: bob});
            assert.equal(await cards.balanceOf(bob), 1);
            // // Bob should now own a coin in the cards contract

            // TODO
            // Charlie is also able to withdraw his bond from the fast exit
            // Small optimization issue that 2 withdrawals need to be made,
            // Plasma -> Fast -> Charlie
            // await fast.withdrawBonds({from: charlie})
        });
    });

});
