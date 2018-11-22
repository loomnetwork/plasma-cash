const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract.only("Plasma ERC721 - Exit Spent Coin Challenge / `challengeAfter`", async function(accounts) {

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
    let UTXO;

    const blk_1 = 1000
    const blk_2 = 2000
    const blk_3 = 3000
    const blk_4 = 4000

    let [authority, alice, bob, charlie, dylan, elliot, random_guy, random_guy2, challenger] = accounts;


    beforeEach(async function() {
        vmc = await ValidatorManagerContract.new({from: authority});
        plasma = await RootChain.new(vmc.address, {from: authority});
        cards = await CryptoCards.new(plasma.address);
        await vmc.toggleToken(cards.address);
        cards.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);

        let ret;
        for (let i = 0; i < ALICE_DEPOSITED_COINS; i ++) {
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
            // assert.equal(coin.slot.toNumber(), i);
            assert.equal(coin.blockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }

        UTXO = {'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()};

    });

    describe('Unit Tests', function() {

        it("Cannot challenge a coin not being exited", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let challengeTx = alice_to_bob.tx
            let sig = alice_to_bob.sig
            let proof = tree_bob.createMerkleProof(UTXO.slot)

            // State before must be 0
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
            try { 
                await txlib.challengeAfter(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_1, 'tx': alice_to_bob },
                    tree_bob.createMerkleProof(UTXO.slot)
                )
            } catch (e) {
            }

            // State after must be 0
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
        })

        it("Bonds get slashed correctly", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            await txlib.exitDeposit(plasma, alice, UTXO)

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)

            await txlib.challengeAfter(plasma, challenger,
                UTXO.slot,
                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot)
            )
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
            await txlib.withdrawBonds(plasma, challenger, 0.1)
        })

        it("Cannot provide an invalid signature", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            await txlib.exitDeposit(plasma, alice, UTXO)

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)

            try { 
                alice_to_bob.sig = "0x1234" // make the sig invalid
                await txlib.challengeAfter(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_1, 'tx': alice_to_bob },
                    tree_bob.createMerkleProof(UTXO.slot)
                )
            } catch (e) {
                // console.log(e.reason)
            }

            // State after must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
        })

        it("Cannot provide an invalid merkle proof", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            await txlib.exitDeposit(plasma, alice, UTXO)

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
            try {
                await txlib.challengeAfter(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_1, 'tx': alice_to_bob },
                    "0x12345678"
                )
            } catch (e) {
                // console.log(e)
            }
            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
        })

        it("Cannot provide an invalid transaction", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            // Alice exits at [0,UTXO.block]
            await txlib.exitDeposit(plasma, alice, UTXO)

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)

            try { 
            await txlib.challengeAfter(plasma, challenger,
                UTXO.slot,
                { 'block': blk_1, 'tx': bob_to_charlie },
                tree_bob.createMerkleProof(UTXO.slot)
            )
            } catch (e) {
                // console.log(e)
            }
            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
        })

        it("Exit of a spent coin is successful if left unchallenged", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, blk_2, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, blk_3, txs);

            let dylan_to_elliot = txlib.createUTXO(UTXO.slot, blk_3, dylan, elliot);
            txs = [dylan_to_elliot.leaf]
            let tree_elliot = await txlib.submitTransactions(authority, plasma, blk_4, txs);

            // Bob exits at [UTXO.block,blk_1]
            await txlib.exit(plasma, charlie,
                UTXO.slot,

                { 'block': blk_2, 'tx': bob_to_charlie },
                tree_charlie.createMerkleProof(UTXO.slot),

                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot),
            )
            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo( t0 + t1 + t2);
            await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
            assert.equal(await txlib.getState(plasma, UTXO.slot), 2)
        });
    })

    describe('C = Deposit, PC = Deposit', function() {

        it("Can challenge with a direct spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            // Alice exits at [0,UTXO.block]
            await txlib.exitDeposit(plasma, alice, UTXO)

            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
            await txlib.challengeAfter(plasma, challenger,
                UTXO.slot,
                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot)
            )

            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
        })


        it("Cannot challenge with a non-direct spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);
            
            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);
            // Alice exits at [0,UTXO.block]
            await txlib.exitDeposit(plasma, alice, UTXO)

            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)

            try { 
                await txlib.challengeAfter(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_2, 'tx': bob_to_charlie },
                    tree_charlie.createMerkleProof(UTXO.slot)
                )
            } catch (e) { }

            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
        })

    })

    describe('C = Non-Deposit, PC = Deposit', function() {
        it("Can challenge with a direct spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            // Bob exits at [UTXO.block,blk_1]
            await txlib.exit(plasma, bob,
                UTXO.slot,
                
                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot),

                { 'block': UTXO.block, 'tx': txlib.createUTXO(UTXO.slot, 0, alice, alice) },
                '0x',
            )

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
            await txlib.challengeAfter(plasma, challenger,
                UTXO.slot,
                { 'block': blk_2, 'tx': bob_to_charlie },
                tree_charlie.createMerkleProof(UTXO.slot)
            )
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
        })

        it("Cannot challenge with a non-direct spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, blk_2, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, blk_3, txs);

            // Bob exits at [UTXO.block,blk_1]
            await txlib.exit(plasma, bob,
                UTXO.slot,
                
                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot),

                { 'block': UTXO.block, 'tx': txlib.createUTXO(UTXO.slot, 0, alice, alice) },
                '0x',
            )

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
            try {
                await txlib.challengeAfter(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_3, 'tx': charlie_to_dylan },
                    tree_dylan.createMerkleProof(UTXO.slot)
                )
            } catch (e) { 

            }
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
        })

    })

    describe('C = Non-Deposit, PC = Non-Deposit', function() {
        it("Can challenge with a direct spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, blk_2, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, blk_3, txs);

            // Bob exits at [UTXO.block,blk_1]
            await txlib.exit(plasma, charlie,
                UTXO.slot,
                
                { 'block': blk_2, 'tx': bob_to_charlie },
                tree_charlie.createMerkleProof(UTXO.slot),

                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot),
            )

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
            await txlib.challengeAfter(plasma, challenger,
                UTXO.slot,
                { 'block': blk_3, 'tx': charlie_to_dylan },
                tree_dylan.createMerkleProof(UTXO.slot)
            )
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
        })

        it("Cannot challenge with a non-direct spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, blk_1, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, blk_1, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, blk_2, txs);

            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, blk_2, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, blk_3, txs);

            let dylan_to_elliot = txlib.createUTXO(UTXO.slot, blk_3, dylan, elliot);
            txs = [dylan_to_elliot.leaf]
            let tree_elliot = await txlib.submitTransactions(authority, plasma, blk_4, txs);

            // Bob exits at [UTXO.block,blk_1]
            await txlib.exit(plasma, charlie,
                UTXO.slot,
                
                { 'block': blk_2, 'tx': bob_to_charlie },
                tree_charlie.createMerkleProof(UTXO.slot),

                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot),
            )

            // State before must be 1
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
            try {
                await txlib.challengeAfter(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_4, 'tx': dylan_to_elliot },
                    tree_elliot.createMerkleProof(UTXO.slot)
                )
            } catch (e) { 

            }
            assert.equal(await txlib.getState(plasma, UTXO.slot), 1)
        })

    })

});
