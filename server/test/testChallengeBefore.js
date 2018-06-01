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

contract("Plasma ERC721 - Invalid History Challenge / `challengeBefore`", async function(accounts) {

    const UTXO_SLOT = 2;
    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    let cards;
    let plasma;
    let t0;

    let [authority, alice, bob, charlie, dylan, elliot, random_guy, random_guy2, challenger] = accounts;

    let exit_coin;
    let data;
    let txs;
    let to_alice;

    beforeEach(async function() {
        plasma = await RootChain.new({from: authority});
        cards = await CryptoCards.new(plasma.address);
        plasma.setCryptoCards(cards.address);
        cards.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);


        let ret = UTXO.createUTXO(0, 0, alice, alice); data = ret[0];
        await cards.depositToPlasmaWithData(1, data, {from: alice});

        ret = UTXO.createUTXO(1, 0, alice, alice); data = ret[0];
        await cards.depositToPlasmaWithData(2, data, {from: alice});

        ret = UTXO.createUTXO(2, 0, alice, alice); data = ret[0];
        await cards.depositToPlasmaWithData(3, data, {from: alice});

        to_alice = ret;

        assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 3);

        const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
        const events = await Promisify(cb => depositEvent.get(cb));
        exit_coin = events[2].args;
        assert.equal(exit_coin.slot.toNumber(), 2);
        assert.equal(exit_coin.depositBlockNumber.toNumber(), 3);
        assert.equal(exit_coin.denomination.toNumber(), 1);
        assert.equal(exit_coin.from, alice);
    });

    describe('Invalid Exit of UTXO 2', function() {

        it("Elliot tries to exit a coin that has invalid history. Elliot's exit gets challenged with challengeBefore w/o response as there is no valid transaction to respond with", async function() {
            let UTXO_SLOT = 2;
            let bobTx = await elliotInvalidHistoryExit();
            let to_bob = bobTx[0];
            let tree_bob = bobTx[1];


            // Concatenate the 2 signatures
            let sigs = to_alice[1] + to_bob[1].replace('0x', '')
            let tx_proof = tree_bob.createMerkleProof(UTXO_SLOT)

            let prev_tx = to_alice[0];
            let tx = to_bob[0];

            // Challenge before is essentially a challenge where the challenger submits the proof required to exit a coin, claiming that this is the last valid state of a coin. Due to bonds the challenger will only do this when he actually knows that there was an invalid spend. If the challenger is a rational player, there should be no case where respondChallengeBefore should succeed.
            await plasma.challengeBefore(
                UTXO_SLOT,
                prev_tx , tx, // rlp encoded
                '0x0', tx_proof, // proofs from the tree
                sigs, // concatenated signatures
                3, 1000,
                {'from': challenger, 'value': web3.toWei(0.1, 'ether')}
            );

            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo( t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });

            // Charlie shouldn't be able to withdraw the coin.
            assertRevert( plasma.withdraw(UTXO_SLOT, {from : elliot }));

            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(dylan)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(elliot)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 3);

            // On the contrary, his bond must be slashed, and `challenger` must be able to claim it
            await plasma.withdrawBonds({from: challenger});
            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, challenger);
            assert.equal(withdraw.amount, web3.toWei(0.1, 'ether'));
        });

        it("Elliot makes a valid exit which gets challenged, however he responds with `respondChallengeBefore`", async function() {
            let bobTx = await elliotValidHistoryExit();
            let to_bob = bobTx[0];
            let tree_bob = bobTx[1];

            let to_charlie = bobTx[2];
            let tree_charlie = bobTx[3];

            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // Concatenate the 2 signatures
            let sigs = to_alice[1] + to_bob[1].substr(2,132);
            let proof = tree_bob.createMerkleProof(UTXO_SLOT)

            let prev_tx = to_alice[0];
            let tx = to_bob[0];

            // Challenge before is essentially a challenge where the challenger submits the proof required to exit a coin, claiming that this is the last valid state of a coin. Due to bonds the challenger will only do this when he actually knows that there was an invalid spend. If the challenger is a rational player, there should be no case where respondChallengeBefore should succeed.
            await plasma.challengeBefore(
                UTXO_SLOT,
                prev_tx , tx, // rlp encoded
                '0x0', proof, // proofs from the tree
                sigs, // concatenated signatures
                3, 1000,
                {'from': challenger, 'value': web3.toWei(0.1, 'ether')}
            );


            let responseTx = to_charlie[0];
            let responseProof = tree_charlie.createMerkleProof(UTXO_SLOT);

            await plasma.respondChallengeBefore(
                UTXO_SLOT, 2000, responseTx, responseProof,
                {'from': elliot }
            );


            await increaseTimeTo( t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2});
            await plasma.withdraw(UTXO_SLOT, {from : elliot});

            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(dylan)).toNumber(), 0);
            assert.equal((await cards.balanceOf.call(elliot)).toNumber(), 1);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);

            // On the contrary, his bond must be slashed, and `challenger` must be able to claim it
            await plasma.withdrawBonds({from: elliot});

            let withdrewBonds = plasma.WithdrewBonds({}, {fromBlock: 0, toBlock: 'latest'});
            let e = await Promisify(cb => withdrewBonds.get(cb));
            let withdraw = e[0].args;
            assert.equal(withdraw.from, elliot);
            assert.equal(withdraw.amount, web3.toWei(0.2, 'ether'));
        });

        async function elliotInvalidHistoryExit() {
            let to_bob = UTXO.createUTXO(UTXO_SLOT, 3, alice, bob);
            txs = [ to_bob[2] ]
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            // The authority submits a block, but there is no transaction from Bob to Charlie
            let tree_charlie = await UTXO.submitTransactions(authority, plasma);

            // Nevertheless, Charlie pretends he received the coin, and by colluding with the chain operator he is able to include his invalid transaction in a block.
            let to_dylan = UTXO.createUTXO(UTXO_SLOT, 2000, charlie, dylan);
            txs = [ to_dylan[2] ]
            let tree_dylan = await UTXO.submitTransactions(authority, plasma, txs);

            // Dylan having received the coin, gives it to Elliot. 
            let to_elliot = UTXO.createUTXO(UTXO_SLOT, 3000, dylan, elliot);
            txs = [ to_elliot[2] ]
            let tree_elliot = await UTXO.submitTransactions(authority, plasma, txs);

            // Elliot normally should be always checking the coin's history and not accepting the payment if it's invalid like in this case, but it is considered that they are all colluding together to steal Bob's coin.;
            // Elliot actually has all the info required to submit an exit, even if one of the transactions in the coin's history were invalid. 
            let sigs = to_dylan[1] + to_elliot[1].substr(2, 132);

            let prev_tx_proof = tree_dylan.createMerkleProof(UTXO_SLOT)
            let exiting_tx_proof = tree_elliot.createMerkleProof(UTXO_SLOT)

            let prev_tx = to_dylan[0];
            let exiting_tx = to_elliot[0]; 

            plasma.startExit(
                UTXO_SLOT,
                prev_tx, exiting_tx, 
                prev_tx_proof, exiting_tx_proof, 
                sigs, 
                3000, 4000, 
                {'from': elliot, 'value': web3.toWei(0.1, 'ether')}
            );

            return [to_bob, tree_bob];

        }

        async function elliotValidHistoryExit() {
            let to_bob = UTXO.createUTXO(UTXO_SLOT, 3, alice, bob);
            txs = [ to_bob[2] ]
            let tree_bob = await UTXO.submitTransactions(authority, plasma, txs);

            // The authority submits a block, but there is no transaction from Bob to Charlie
            let to_charlie = UTXO.createUTXO(UTXO_SLOT, 1000, bob, charlie);
            txs = [ to_charlie[2] ]
            let tree_charlie = await UTXO.submitTransactions(authority, plasma, txs);

            // Nevertheless, Charlie pretends he received the coin, and by colluding with the chain operator he is able to include his invalid transaction in a block.
            let to_dylan = UTXO.createUTXO(UTXO_SLOT, 2000, charlie, dylan);
            txs = [ to_dylan[2] ]
            let tree_dylan = await UTXO.submitTransactions(authority, plasma, txs);

            // Dylan having received the coin, gives it to Elliot. 
            let to_elliot = UTXO.createUTXO(UTXO_SLOT, 3000, dylan, elliot);
            txs = [ to_elliot[2] ]
            let tree_elliot = await UTXO.submitTransactions(authority, plasma, txs);

            // Elliot normally should be always checking the coin's history and not accepting the payment if it's invalid like in this case, but it is considered that they are all colluding together to steal Bob's coin.;
            // Elliot actually has all the info required to submit an exit, even if one of the transactions in the coin's history were invalid. 
            let sigs = to_dylan[1] + to_elliot[1].replace('0x', '')

            let prev_tx_proof = tree_dylan.createMerkleProof(UTXO_SLOT)
            let exiting_tx_proof = tree_elliot.createMerkleProof(UTXO_SLOT)

            let prev_tx = to_dylan[0];
            let exiting_tx = to_elliot[0]; 

            plasma.startExit(
                UTXO_SLOT,
                prev_tx, exiting_tx, 
                prev_tx_proof, exiting_tx_proof, 
                sigs, 
                3000, 4000, 
                {'from': elliot, 'value': web3.toWei(0.1, 'ether')}
            );

            return [to_bob, tree_bob, to_charlie, tree_charlie];

        }

    })
});
