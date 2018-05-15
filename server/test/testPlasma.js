const _ = require('lodash')

const RLP = require('rlp')
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");

import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';
const utils = require('web3-utils');

contract("Plasma ERC721 WIP", async function(accounts) {

    let cards;
    let plasma;
    let start;

    let [authority, alice, bob, random_guy, random_guy2] = accounts;

    // UTXO = (prev_block, uid, new_owner)
    let aliceUTXO = '0x' + RLP.encode([0, 1, alice]).toString('hex');
    // Alice's UTXO Is included in deposit block 1.
    let bobUTXO = '0x' + RLP.encode([1, 1, bob]).toString('hex')

    before("Deploys the contracts", async function() {
        plasma = await RootChain.new({from: authority});
        cards = await CryptoCards.new(plasma.address);
        plasma.setCryptoCards(cards.address);
    });

    describe("Tests Plasma Deposits through the linked ERC721 contract", function() {

        it("Registers alice charlie bob", async function() {
            cards.register({from: alice});
            assert.equal(await cards.balanceOf.call(alice), 5);
        });

        it("Transfers NFT 1 from Alice to Plasma contract", async function() {

            // Call without extra data
            await cards.depositToPlasmaWithData(1, aliceUTXO, {from: alice});
            assert.equal(await cards.balanceOf.call(alice), 4);
            assert.equal(await cards.balanceOf.call(plasma.address), 1);

        });

        it("Checks that events were emitted correctly", async function() {
            plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'}).get((error, res) => {
               let aliceLogs = res[0].args; // Why only 1 event gets emitted?
                assert.equal(aliceLogs.depositor, alice);
                assert.equal(aliceLogs.tokenId.toNumber(), 1);
                assert.equal(aliceLogs.data, aliceUTXO);
            });
        });
    });

    describe("Exit of Alice's UTXO1 by Bob", function() {

        it("Alice signs a utxo for coin1, operator submits root proof", async function() {
            let txHash = utils.soliditySha3(bobUTXO);
            let sig = web3.eth.sign(alice, txHash);
            await plasma.startExit(aliceUTXO, bobUTXO, sig);
            start = (await web3.eth.getBlock('latest')).timestamp;
        });

        it("Another random guy finalizes the exits 1 week later. No challlenge has happened this time period", async function() {
            let expire = start + 3600 * 24 * 8; // 8 days pass, can finalize exit
            await increaseTimeTo(expire);
            await plasma.finalizeExits({from: random_guy2 });
        });

        it("Exit is finalized, Bob can withdraw his tokens", async function() {
            plasma.withdraw({from : bob});
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 1);
            assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), 0);
        })
    });
});
