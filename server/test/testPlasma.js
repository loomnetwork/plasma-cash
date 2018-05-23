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

contract("Plasma ERC721 WIP", async function(accounts) {

    let cards;
    let plasma;
    let start;
    let utxo_slot; // the utxo that alice will attempt to withdrawo

    const UTXO_ID = 2

    let [authority, alice, bob, charlie, random_guy, random_guy2] = accounts;

    let exit_coin;
    let data;
    let rawdata;

    function createUTXO(slot, prevBlock, from, to) {
        let data = [ slot, prevBlock, 1, to ]
        data = '0x' + RLP.encode(data).toString('hex')
        // console.log('Created UTXO for slot', slot);
        // console.log(data)
        let txHash = utils.soliditySha3(data);
        let sig = web3.eth.sign(from, txHash); 
        return [data, sig];
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
        let slot = 60
        let prevblock = 1000;
        let denom = 1;
        let newowner = bob;
        let data = [slot, prevblock, denom, newowner];
        let tx = '0x' + RLP.encode(data).toString('hex');
        let txHash = utils.soliditySha3(tx);

        let leaves = {};
        leaves[slot] = txHash;

        // slot = 63;
        // data = [slot, prevblock, denom, newowner];
        // tx = '0x' + RLP.encode(data).toString('hex');
        // txHash = utils.soliditySha3(tx);
        // leaves[slot] = txHash;

        let tree = new SparseMerkleTree(64, leaves);
        // tree.root will be submited to `submitBlock`
        let proof = tree.createMerkleProof(slot);

        let ret = await plasma.checkMembership(txHash, tree.root, slot, proof);
        console.log('Sent:', txHash, tree.root, slot, proof);
        console.log(ret);

    
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
        // We submit 2 precomputed block roots which represent:
        // Block 1000: Transaction from Alice to Bob
        // Block 2000: Transaction from Bob to Charlie
        // These were precomputed from the Python client for a Sparse Merkle Tree of Depth 32.
        
        let utxo_slot = 2;
        // Tx to Bob from Alice referencing Alice's UTXO at block 3
        let to_bob = createUTXO(utxo_slot, 3, alice, bob);
        let block_1000_root = '0x20d4251cbfd45b0f2c708e63f9c96ddd8cd832d087bacfa241716e41ddbe136b';
        // plasma.submitBlock(block_1000_root, {'from': authority});

        // Tx to Charlie from Bob referencing Bob's UTXO at block 1000
        let to_charlie = createUTXO(utxo_slot, 1000, bob, charlie);
        let block_2000_root = '0xe9ac5d9bc7cb7cc7c00d06b80a4e2f0a40a3803d3a84968f403aa312474e1ca6';

        plasma.submitBlock(block_2000_root, {'from': authority});
       
        // Concatenate the 2 signatures
        let sigs = to_bob[1] + to_charlie[1].substr(2, 132);

        // Merkle branches -> TODO
        let prev_tx_proof = '0x';
        let exiting_tx_proof = '0x'; // To add the valid proofs of inclusion

        let prev_tx = to_bob[0];
        let exiting_tx = to_charlie[0];

        plasma.startExit(
                utxo_slot,
                prev_tx, exiting_tx, 
                prev_tx_proof, exiting_tx_proof, 
                sigs,
                1000, 2000,
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
