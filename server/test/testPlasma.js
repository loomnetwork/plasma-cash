const _ = require('lodash')

const RLP = require('rlp')
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");

const SparseMerkleTree = require('./SparseMerkleTree.js');

import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';
const utils = require('web3-utils');

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

contract("Plasma ERC721", async function(accounts) {

    let cards;
    let plasma;
    let start;
    let utxo_slot; // the utxo that alice will attempt to withdrawo

    const UTXO_ID = 2

    let [authority, alice, bob, charlie, random_guy, random_guy2] = accounts;

    let exit_coin;
    let data;
    let rawdata;

    function signPrefixed(from, hash) {
        let prefix = "\u0019Ethereum Signed Message:\n32";
        let prefixedHash = utils.soliditySha3(prefix, hash);
        // BUG! Recovering the sender from the signing does not work for some reason. WIP
        let sig = web3.eth.sign(from, prefixedHash);
        return sig
    }

    function createUTXO(slot, prevBlock, from, to) {
        let data = [ slot, prevBlock, 1, to ];
        data = '0x' + RLP.encode(data).toString('hex');
        let txHash = utils.soliditySha3(data);
        let sig = signPrefixed(from, txHash); // prefixed signature on the hash
        return [data, sig];
    }
    //
    async function submitUTXO(slot, tx) {
        // Create merkle Tree from A SINGLE UTXO and submit it.
        // Returns merkle tree that was created
        let leaves = {}
        leaves[slot] = utils.soliditySha3(tx);

        let tree = new SparseMerkleTree(64, leaves);
        await plasma.submitBlock(tree.root, {'from': authority});

        return tree;
    }


    beforeEach("Deploys the contracts, Registers Alice and deposits her coins", async function() {
        plasma = await RootChain.new({from: authority});
        cards = await CryptoCards.new(plasma.address);
        plasma.setCryptoCards(cards.address);
        cards.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);


        let ret = createUTXO(0, 0, alice, alice); data = ret[0];
        await cards.depositToPlasmaWithData(1, data, {from: alice});

        ret = createUTXO(1, 0, alice, alice); data = ret[0];
        await cards.depositToPlasmaWithData(2, data, {from: alice});

        ret = createUTXO(2, 0, alice, alice); data = ret[0];
        await cards.depositToPlasmaWithData(3, data, {from: alice});

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

    it('Tests that Merkle Proofs work', async function() {
        let slot = 1500;
        let tx = createUTXO(slot, 1000, alice, bob);
        let txHash = utils.soliditySha3(tx[0]);

        let leaves = {};
        leaves[slot] = txHash;

        // BUG! proofbits are not properly returned. WIP
        // slot = 63;
        // data = [slot, 2000, bob, charlie];
        // tx = '0x' + RLP.encode(data).toString('hex');
        // txHash = utils.soliditySha3(tx);
        // leaves[slot] = txHash;

        // This will be happening on the Plasma Cash client
        // The root will be submitted by the authority
        let tree = new SparseMerkleTree(64, leaves);
        let proof = tree.createMerkleProof(slot);

        let ret = await plasma.checkMembership(txHash, tree.root, slot, proof);
        assert.equal(ret, true);

    
    });

    it("Submits an exit for the UTXO of Coin 3 (utxo id 2)  directly after depositing it", async function() {
        utxo_slot = exit_coin.slot.toNumber();
        let includedBlock = exit_coin.depositBlockNumber.toNumber();
        let denom = exit_coin.denomination.toNumber();
        let from  = exit_coin.from;

        let ret = createUTXO(2, 0, alice, alice); 
        let utxo = ret[0];
        let sig = ret[1];

        await plasma.startExit(
                 utxo_slot,
                '0x', utxo,  // prevTx, exitingTx
                '0x0', '0x0', // inclusion proofs
                 sig,
                 0, includedBlock, 
                 {'from': alice}
        );

        start = (await web3.eth.getBlock('latest')).timestamp;
        let expire = start + 3600 * 24 * 8; // 8 days pass, can finalize exit
        await increaseTimeTo(expire);

        await plasma.finalizeExits({from: random_guy2 });

        plasma.withdraw(utxo_slot, {from : alice });
        assert.equal((await cards.balanceOf.call(alice)).toNumber(), 3);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);
    });

    it("Transfers Coin 2 from Alice to Bob and then to Charlie who tries to exit it", async function() {
        let leaves = {};
        let utxo_slot = 2;

        // Block 1000: Transaction from Alice to Bob
        // Block 2000: Transaction from Bob to Charlie
        
        let to_bob = createUTXO(utxo_slot, 3, alice, bob);
        // submits tree root frmo authority
        let tree_bob = await submitUTXO(utxo_slot, to_bob[0]);

        // Tx to Charlie from Bob referencing Bob's UTXO at block 1000
        let to_charlie = createUTXO(utxo_slot, 1000, bob, charlie);
        let tree_charlie = await submitUTXO(utxo_slot, to_charlie[0]);
       
        // Concatenate the 2 signatures
        let sigs = to_bob[1] + to_charlie[1].substr(2, 132);

        let prev_tx_proof = tree_bob.createMerkleProof(utxo_slot)
        let exiting_tx_proof = tree_charlie.createMerkleProof(utxo_slot)

        let prev_tx = to_bob[0];
        let exiting_tx = to_charlie[0];

        plasma.startExit(
                utxo_slot,
                prev_tx, exiting_tx, // rlp encoded
                prev_tx_proof, exiting_tx_proof, // proofs from the tree
                sigs, // concatenated signatures
                1000, 2000, // 1000 is when alice->bob got included, 2000 for bob->charlie
                {'from': charlie }
        );

        start = (await web3.eth.getBlock('latest')).timestamp;
        let expire = start + 3600 * 24 * 8; // 8 days pass, can finalize exit
        await increaseTimeTo(expire);

        await plasma.finalizeExits({from: random_guy2 });

        plasma.withdraw(utxo_slot, {from : charlie });
        assert.equal((await cards.balanceOf.call(alice)).toNumber(), 2);
        assert.equal((await cards.balanceOf.call(bob)).toNumber(), 0);
        assert.equal((await cards.balanceOf.call(charlie)).toNumber(), 1);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);
    });


});
