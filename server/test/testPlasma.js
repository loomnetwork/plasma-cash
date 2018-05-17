const _ = require('lodash')

const RLP = require('rlp')
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");

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
    let utxo_slot; // the utxo that alice will attempt to withdraw

    let [authority, alice, bob, random_guy, random_guy2] = accounts;


    before("Deploys the contracts", async function() {
        plasma = await RootChain.new({from: authority});
        cards = await CryptoCards.new(plasma.address);
        plasma.setCryptoCards(cards.address);
    });

    describe("Tests that alice can deposit and withdraw a coin", function() {

        it("Registers alice charlie bob", async function() {
            cards.register({from: alice});
            assert.equal(await cards.balanceOf.call(alice), 5);
        });

        it("Transfers NFT 1,2,3 from Alice to Plasma contract", async function() {
            await cards.depositToPlasma(1, {from: alice});
            await cards.depositToPlasma(2, {from: alice});
            await cards.depositToPlasma(3, {from: alice});
            assert.equal(await cards.balanceOf.call(alice), 2);
            assert.equal(await cards.balanceOf.call(plasma.address), 3);

        });

        it("Submits an exit for the UTXO of Coin 2", async function() {
            const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
            const events = await Promisify(cb => depositEvent.get(cb));
            let coin2 = events[1].args;

            // data for hash
            utxo_slot = coin2.slot.toNumber();
            let txHash = utils.soliditySha3(utxo_slot);

            let prevBlock = coin2.depositBlockNumber.toNumber();
            let denom = coin2.denomination.toNumber();
            let from  = coin2.from;
            let data = [
                utxo_slot, // exit coin id
                prevBlock, // this is a deposit block
                denom,
                from
            ];
            let utxo = '0x' + RLP.encode(data).toString('hex');

            // sign the exit with Alice's key since she is the owner of the private key that corresponds to `from`. If it's not signed by alice this should fail. 
            let sig = web3.eth.sign(alice, txHash); 
           //  console.log('RLP Encoded:', utxo);
           //  console.log('TXHash:', txHash);
           //  console.log('Signature', sig);
            let rootBlkHash = await plasma.childChain.call(prevBlock);
            // console.log('Block Hash:', rootBlkHash[0]);

            // Since we are exiting a deposit UTXO, there is no need to provide a previous tx or proof.
            await plasma.startExit(
                    '0x', utxo,  // prevTx, exitingTx
                    '0x0', '0x0', // inclusion proofs
                    sig // msg signed by alice
                    );
            start = (await web3.eth.getBlock('latest')).timestamp;
        });

        it("A random guy finalizes the exits 1 week later. No challlenge has happened this time period", async function() {
            let expire = start + 3600 * 24 * 8; // 8 days pass, can finalize exit
            await increaseTimeTo(expire);
            await plasma.finalizeExits({from: random_guy2 });
        });

        it("Exit is finalized, Alice can withdraw her tokens", async function() {
            plasma.withdraw(utxo_slot, {from : alice });
            assert.equal((await cards.balanceOf.call(alice)).toNumber(), 3);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 2);
        })
    });

});
