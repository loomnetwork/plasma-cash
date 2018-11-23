const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract("Plasma ERC721 - Invalid History Challenge / `challengeBefore`", async function(accounts) {

    const RESPONSE_PERIOD = 3600 * 24 * 3.5; // 3.5 days later
    const MATURITY_PERIOD = 3600 * 24 * 7; // 3.5 days later
    const e = 3600;

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
    const blk_5 = 5000


    let [authority, alice, bob, charlie, dylan, elliot, fred, random_guy, random_guy2, challenger] = accounts;


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
            assert.equal(coin.blockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }
        UTXO = {'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()};

    });

    describe('Without Responses', function() {
        it("Challenge with C = Deposit (no response)", async function() {

            // The authority submits a block, but there is no transaction from Alice to Bob
            let tree_bob = await txlib.submitTransactions(authority, plasma, 1000);

            // Nevertheless, Bob pretends he received the coin, and by
            // colluding with the chain operator he is able to include his
            // invalid transaction in a block.
            let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
            let txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000, txs);

            // Charlie having received the coin, gives it to Dylan.
            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

            // Dylan normally should be always checking the coin's history and
            // not accepting the payment if it's invalid like in this case, but
            // it is considered that they are all colluding together to steal
            // Alice's coin.  Dylan has all the info required to submit
            // an exit, even if one of the transactions in the coin's history
            // were invalid.
            t0 = await txlib.exit(plasma, dylan,
               UTXO.slot,
                
                { 'block': blk_3, 'tx': charlie_to_dylan },
                tree_dylan.createMerkleProof(UTXO.slot),

                { 'block': blk_2, 'tx': bob_to_charlie },
                tree_charlie.createMerkleProof(UTXO.slot),
            )

            await txlib.challengeBefore(plasma, challenger,
                UTXO.slot,
                { 'block': UTXO.block, 'tx': txlib.createUTXO(UTXO.slot, 0, alice, alice) },
                '0x0'
            )

            // Go a litle after the maturity period
            await increaseTimeTo(t0 + MATURITY_PERIOD + e);
            await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0, "State should be reset")

        });

        it("Challenge with C = Non-Deposit (no response)", async function() {
            // Alice gives her coin legitimately to Bob
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, 1000, txs);

            // The authority submits a block, but there is no transaction from Bob to Charlie
            let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000);

            // Nevertheless, Charlie pretends he received the coin, and by
            // colluding with the chain operator he is able to include his
            // invalid transaction in a block.
            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

            // Dylan having received the coin, gives it to Elliot.
            let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
            txs = [dylan_to_elliot.leaf]
            let tree_elliot = await txlib.submitTransactions(authority, plasma, 4000, txs);

            // Elliot normally should be always checking the coin's history and
            // not accepting the payment if it's invalid like in this case, but
            // it is considered that they are all colluding together to steal
            // Bob's coin.  Elliot has all the info required to submit
            // an exit, even if one of the transactions in the coin's history
            // were invalid.
            t0 = await txlib.exit(plasma, elliot,
               UTXO.slot,
                
                { 'block': blk_4, 'tx': dylan_to_elliot },
                tree_elliot.createMerkleProof(UTXO.slot),

                { 'block': blk_3, 'tx': charlie_to_dylan },
                tree_dylan.createMerkleProof(UTXO.slot),
            )

            await txlib.challengeBefore(plasma, challenger,
                UTXO.slot,
                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot) 
            )

            await increaseTimeTo(t0 + MATURITY_PERIOD + e);
            await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0, "State should be reset")
        });
    });

    describe('With Responses', function() {
        it("Cannot respond with an earlier spend", async function() {
            // alice gives coin to bob
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, 1000, txs);

            // bob gives it back to alice
            let bob_to_alice = txlib.createUTXO(UTXO.slot, 1000, bob, alice);
            txs = [bob_to_alice.leaf]
            let tree_alice = await txlib.submitTransactions(authority, plasma, 2000, txs);

            // The authority submits a block, but there is no transaction from Bob to Charlie
            let tree_charlie = await txlib.submitTransactions(authority, plasma, 3000);

            // Nevertheless, Charlie pretends he received the coin, and by
            // colluding with the chain operator he is able to include his
            // invalid transaction in a block.
            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, 4000, txs);

            // Dylan having received the coin, gives it to Elliot.
            let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
            txs = [dylan_to_elliot.leaf]
            let tree_elliot = await txlib.submitTransactions(authority, plasma, 5000, txs);

            t0 = await txlib.exit(plasma, elliot,
               UTXO.slot,
                
                { 'block': blk_5, 'tx': dylan_to_elliot },
                tree_elliot.createMerkleProof(UTXO.slot),

                { 'block': blk_4, 'tx': charlie_to_dylan },
                tree_dylan.createMerkleProof(UTXO.slot),
            )

            await txlib.challengeBefore(plasma, challenger,
                UTXO.slot,
                { 'block': blk_2, 'tx': bob_to_alice },
                tree_alice.createMerkleProof(UTXO.slot) 
            )

            // Elliot tries to make an invalid response by providing
            // the initial transfer from Alice to Bob. A valid response
            // would involve a later spend
            try { 
                await txlib.respondChallengeBefore(plasma, challenger,
                    UTXO.slot,
                    { 'block': blk_1, 'tx': alice_to_bob },
                    tree_bob.createMerkleProof(UTXO.slot) 
                )
            } catch (e) {
                assert.ok(e !== undefined)
            }

            await increaseTimeTo(t0 + MATURITY_PERIOD + e);
            await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
            assert.equal(await txlib.getState(plasma, UTXO.slot), 0)
        });

        it("Can respond with a later spend", async function() {
            let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
            let txs = [alice_to_bob.leaf]
            let tree_bob = await txlib.submitTransactions(authority, plasma, 1000, txs);

            let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
            txs = [bob_to_charlie.leaf]
            let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000, txs);

            let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
            txs = [charlie_to_dylan.leaf]
            let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

            let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
            txs = [dylan_to_elliot.leaf]
            let tree_elliot = await txlib.submitTransactions(authority, plasma, 4000, txs);

            t0 = await txlib.exit(plasma, elliot,
               UTXO.slot,
                
                { 'block': blk_4, 'tx': dylan_to_elliot },
                tree_elliot.createMerkleProof(UTXO.slot),

                { 'block': blk_3, 'tx': charlie_to_dylan },
                tree_dylan.createMerkleProof(UTXO.slot),
            )

            await txlib.challengeBefore(plasma, challenger,
                UTXO.slot,
                { 'block': blk_1, 'tx': alice_to_bob },
                tree_bob.createMerkleProof(UTXO.slot) 
            )

            await txlib.respondChallengeBefore(plasma, elliot,
                UTXO.slot,
                { 'block': blk_2, 'tx': bob_to_charlie },
                tree_charlie.createMerkleProof(UTXO.slot),
                alice_to_bob.leaf.hash
            )
            await increaseTimeTo(t0 + MATURITY_PERIOD + e);
            await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
            assert.equal(await txlib.getState(plasma, UTXO.slot), 2)
        });

        describe("Multiple Challenges", async function() {
            // This scenario can occur as follows when the operator is byzaantine
            // 1) Operator makes an Invalid exit
            // 2) Operator challenges the invalid exit
            // 3) Operator responds with a response that looks valid (but in fact was forged)
            // If multiple challenges were not allowed, the operator would
            // challenge their exit right before the response period starts, which would
            // disallow any valid responses, thus allowing them to finalize an invalid exit
            // By allowing multiple challenges, we require that during finalizations there are
            // no pending challenges which mitigates the above vector.
            // Splitting the maturity period in 2 was introduced to avoid griefing attacks as
            // discussed in issue #103
            it("Invalid exit - Multichallenged in time", async function() {
                // Bob/Operator forges a transaction and gives the coin all the way to Elliot
                // Elliot exits, Bob challenges,
                // The authority submits a block, but there is no transaction from Alice to Bob
                let tree_bob = await txlib.submitTransactions(authority, plasma, 1000);

                // Nevertheless, Bob pretends he received the coin, and by
                // colluding with the chain operator he is able to include his
                // invalid transaction in a block.
                let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
                let txs = [bob_to_charlie.leaf]
                let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000, txs);

                // Charlie having received the coin, gives it to Dylan.
                let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
                txs = [charlie_to_dylan.leaf]
                let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

                // Dylan having received the coin, gives it to Elliot.
                let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
                txs = [dylan_to_elliot.leaf]
                let tree_elliot = await txlib.submitTransactions(authority, plasma, 4000,txs);
                //
                // Dylan having received the coin, gives it to Elliot.
                let elliot_to_fred = txlib.createUTXO(UTXO.slot, 4000, elliot, fred);
                txs = [elliot_to_fred.leaf]
                let tree_fred = await txlib.submitTransactions(authority, plasma, 5000, txs);

                // Elliot normally should be always checking the coin's history and
                // not accepting the payment if it's invalid like in this case, but
                // it is considered that they are all colluding together to steal
                // Bob's coin.  Elliot has all the info required to submit
                // an exit, even if one of the transactions in the coin's history
                // were invalid.
                let sig = elliot_to_fred.sig;
                let prev_tx_proof = tree_elliot.createMerkleProof(UTXO.slot)
                let exiting_tx_proof = tree_fred.createMerkleProof(UTXO.slot)
                let prev_tx = dylan_to_elliot.tx;
                let exiting_tx = elliot_to_fred.tx;

                // Fred submits the exit which is valid as far as the chain is concerned
                await plasma.startExit(
                    UTXO.slot,
                    prev_tx, exiting_tx,
                    prev_tx_proof, exiting_tx_proof,
                    sig,
                    [4000, 5000],
                    {'from': fred, 'value': web3.toWei(0.1, 'ether')}
                );
                t0 = (await web3.eth.getBlock('latest')).timestamp;

                // Dylan challenges the exit
                sig = charlie_to_dylan.sig;
                let tx_proof = tree_dylan.createMerkleProof(UTXO.slot)
                prev_tx_proof = tree_charlie.createMerkleProof(UTXO.slot)
                prev_tx = bob_to_charlie.tx;
                let tx = charlie_to_dylan.tx;
                await plasma.challengeBefore(
                    UTXO.slot,
                    tx,
                    tx_proof,
                    sig,
                    3000,
                    {'from': dylan, 'value': web3.toWei(0.1, 'ether')}
                );

                // Alice ALSO challenges the exit! Note that Dylan's challenge is an artificial one
                // that can be responded to by a response which holds on to invalid transactions
                let alice_to_alice = txlib.createUTXO(UTXO.slot, 0, alice, alice);
                await plasma.challengeBefore(
                    UTXO.slot,
                    alice_to_alice.tx,
                    '0x0',
                    alice_to_alice.sig,
                    UTXO.block,
                    {'from': alice, 'value': web3.toWei(0.1, 'ether')}
                );

                // No more challenges can be issued
                await increaseTimeTo(t0 + RESPONSE_PERIOD);

                // Elliot responds to the challenge (what are they even trying to do?!)
                let challengingTxHash = charlie_to_dylan.leaf.hash;
                let responseTx = dylan_to_elliot.tx;
                sig = dylan_to_elliot.sig;
                let responseProof = tree_elliot.createMerkleProof(UTXO.slot);
                await plasma.respondChallengeBefore(
                    UTXO.slot, challengingTxHash, 4000, responseTx, responseProof, sig,
                    {'from': elliot}
                );

                await increaseTimeTo(t0 + MATURITY_PERIOD + e);
                await plasma.finalizeExit(UTXO.slot, {from: random_guy2});

                // Fred shouldn't be able to withdraw the coin.
                assertRevert(plasma.withdraw(UTXO.slot, {from : fred}));

                assert.equal(await cards.balanceOf.call(alice), 2);
                assert.equal(await cards.balanceOf.call(bob), 0);
                assert.equal(await cards.balanceOf.call(charlie), 0);
                assert.equal(await cards.balanceOf.call(dylan), 0);
                assert.equal(await cards.balanceOf.call(elliot), 0);
                assert.equal(await cards.balanceOf.call(plasma.address), 3);


                // Bond withdrawals
                await plasma.withdrawBonds({from: elliot});
                await plasma.withdrawBonds({from: alice});
                let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
                let eventLog = await txlib.Promisify(cb => withdrewBonds.get(cb));

                // Elliot is able to get the challenger's bond
                // (nothing gained since they are all colluding)
                let elliot_bond = eventLog[0].args;
                assert.equal(elliot_bond.from, elliot);
                assert.equal(elliot_bond.amount, web3.toWei(0.1, 'ether'));

                // Alice gets her own bond back and the exitor's bond
                let alice_bond = eventLog[1].args;
                assert.equal(alice_bond.from, alice);
                assert.equal(alice_bond.amount, web3.toWei(0.2, 'ether'));
            });

            it("Cannot challenge after the first half of the MATURITY_PERIOD", async function() {
                // Bob/Operator forges a transaction and gives the coin all the way to Elliot
                // Elliot exits, Bob challenges,
                // The authority submits a block, but there is no transaction from Alice to Bob
                let tree_bob = await txlib.submitTransactions(authority, plasma, 1000);

                // Nevertheless, Bob pretends he received the coin, and by
                // colluding with the chain operator he is able to include his
                // invalid transaction in a block.
                let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
                let txs = [bob_to_charlie.leaf]
                let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000, txs);

                // Charlie having received the coin, gives it to Dylan.
                let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
                txs = [charlie_to_dylan.leaf]
                let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

                // Dylan having received the coin, gives it to Elliot.
                let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
                txs = [dylan_to_elliot.leaf]
                let tree_elliot = await txlib.submitTransactions(authority, plasma, 4000, txs);
                //
                // Dylan having received the coin, gives it to Elliot.
                let elliot_to_fred = txlib.createUTXO(UTXO.slot, 4000, elliot, fred);
                txs = [elliot_to_fred.leaf]
                let tree_fred = await txlib.submitTransactions(authority, plasma, 5000, txs);

                // Elliot normally should be always checking the coin's history and
                // not accepting the payment if it's invalid like in this case, but
                // it is considered that they are all colluding together to steal
                // Bob's coin.  Elliot has all the info required to submit
                // an exit, even if one of the transactions in the coin's history
                // were invalid.
                t0 = await txlib.exit(plasma, fred,
                    UTXO.slot,

                    { 'block': blk_5, 'tx': elliot_to_fred },
                    tree_fred.createMerkleProof(UTXO.slot),

                    { 'block': blk_4, 'tx': dylan_to_elliot },
                    tree_elliot.createMerkleProof(UTXO.slot),
                )

                // No more challenges can be issued
                await increaseTimeTo(t0 + RESPONSE_PERIOD + e);

                // Alice can no longer challenge the exit! :(
                try { 
                    await txlib.challengeBefore(plasma, challenger,
                        UTXO.slot,
                        { 'block': UTXO.block, 'tx': txlib.createUTXO(UTXO.slot, 0, alice, alice) },
                        '0x0'
                    )
                } catch (e) { 
                    assert.ok(e !== undefined)
                }

                await increaseTimeTo(t0 + MATURITY_PERIOD + e);
                await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
                assert.equal(await txlib.getState(plasma, UTXO.slot), 2, "State should be finalized")
            });

            it("Can have multiple challenges per exit", async function() {
                // Bob/Operator forges a transaction and gives the coin all the way to Elliot
                // Elliot exits, Bob challenges,
                // The authority submits a block, but there is no transaction from Alice to Bob
                let tree_bob = await txlib.submitTransactions(authority, plasma, 1000);

                // Nevertheless, Bob pretends he received the coin, and by
                // colluding with the chain operator he is able to include his
                // invalid transaction in a block.
                let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
                let txs = [bob_to_charlie.leaf]
                let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000, txs);

                // Charlie having received the coin, gives it to Dylan.
                let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
                txs = [charlie_to_dylan.leaf]
                let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

                // Dylan having received the coin, gives it to Elliot.
                let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
                txs = [dylan_to_elliot.leaf]
                let tree_elliot = await txlib.submitTransactions(authority, plasma, 4000, txs);
                //
                // Dylan having received the coin, gives it to Elliot.
                let elliot_to_fred = txlib.createUTXO(UTXO.slot, 4000, elliot, fred);
                txs = [elliot_to_fred.leaf]
                let tree_fred = await txlib.submitTransactions(authority, plasma, 5000, txs);

                t0 = await txlib.exit(plasma, fred,
                    UTXO.slot,

                    { 'block': blk_5, 'tx': elliot_to_fred },
                    tree_fred.createMerkleProof(UTXO.slot),

                    { 'block': blk_4, 'tx': dylan_to_elliot },
                    tree_elliot.createMerkleProof(UTXO.slot),
                )

                // 2 challenges
                await txlib.challengeBefore(plasma, dylan,
                    UTXO.slot,
                    { 'block': blk_3, 'tx': charlie_to_dylan },
                    tree_dylan.createMerkleProof(UTXO.slot)
                )

                await txlib.challengeBefore(plasma, challenger,
                    UTXO.slot,
                    { 'block': UTXO.block, 'tx': txlib.createUTXO(UTXO.slot, 0, alice, alice) },
                    '0x0'
                )

                await increaseTimeTo(t0 + MATURITY_PERIOD + e);
                await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
                assert.equal(await txlib.getState(plasma, UTXO.slot), 0, "State should be reset")
                
                // Only the first challenger gets the bond from the exit!

                // Bond withdrawals
                await plasma.withdrawBonds({from: dylan});
                await plasma.withdrawBonds({from: challenger});

                let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
                let eventLog = await txlib.Promisify(cb => withdrewBonds.get(cb));

                let dylan_bond = eventLog[0].args;
                assert.equal(dylan_bond.from, dylan);
                assert.equal(dylan_bond.amount, web3.toWei(0.2, 'ether'));

                let alice_bond = eventLog[1].args;
                assert.equal(alice_bond.from, challenger);
                assert.equal(alice_bond.amount, web3.toWei(0.1, 'ether'));
            });

            it("All challenges must be responded to", async function() {
                let alice_to_bob = txlib.createUTXO(UTXO.slot, UTXO.block, alice, bob);
                let txs = [alice_to_bob.leaf]
                let tree_bob = await txlib.submitTransactions(authority, plasma, 1000, txs);

                let bob_to_charlie = txlib.createUTXO(UTXO.slot, 1000, bob, charlie);
                txs = [bob_to_charlie.leaf]
                let tree_charlie = await txlib.submitTransactions(authority, plasma, 2000, txs);

                let charlie_to_dylan = txlib.createUTXO(UTXO.slot, 2000, charlie, dylan);
                txs = [charlie_to_dylan.leaf]
                let tree_dylan = await txlib.submitTransactions(authority, plasma, 3000, txs);

                let dylan_to_elliot = txlib.createUTXO(UTXO.slot, 3000, dylan, elliot);
                txs = [dylan_to_elliot.leaf]
                let tree_elliot = await txlib.submitTransactions(authority, plasma, 4000, txs);

                let elliot_to_fred = txlib.createUTXO(UTXO.slot, 4000, elliot, fred);
                txs = [elliot_to_fred.leaf]
                let tree_fred = await txlib.submitTransactions(authority, plasma, 5000, txs);

                t0 = await txlib.exit(plasma, fred,
                    UTXO.slot,

                    { 'block': blk_5, 'tx': elliot_to_fred },
                    tree_fred.createMerkleProof(UTXO.slot),

                    { 'block': blk_4, 'tx': dylan_to_elliot },
                    tree_elliot.createMerkleProof(UTXO.slot),
                )

                await txlib.challengeBefore(plasma, dylan,
                    UTXO.slot,
                    { 'block': blk_3, 'tx': charlie_to_dylan },
                    tree_dylan.createMerkleProof(UTXO.slot)
                )

                let tx = txlib.createUTXO(UTXO.slot, 0, alice, alice)
                await txlib.challengeBefore(plasma, challenger,
                    UTXO.slot,
                    { 'block': UTXO.block, 'tx': tx },
                    '0x0'
                )

                // 2 responses
                await txlib.respondChallengeBefore(plasma, fred,
                    UTXO.slot,
                    { 'block': blk_4, 'tx': dylan_to_elliot },
                    tree_dylan.createMerkleProof(UTXO.slot),
                    charlie_to_dylan.leaf.hash
                )

                await txlib.respondChallengeBefore(plasma, fred,
                    UTXO.slot,
                    { 'block': blk_1, 'tx': alice_to_bob },
                    tree_bob.createMerkleProof(UTXO.slot),
                    tx.leaf.hash
                )


                await increaseTimeTo(t0 + MATURITY_PERIOD + e);
                await plasma.finalizeExit(UTXO.slot, {from: random_guy2});
                assert.equal(await txlib.getState(plasma, UTXO.slot), 2)
                await txlib.withdrawBonds(plasma, fred, 0.3)
            });

        });

    });
});
