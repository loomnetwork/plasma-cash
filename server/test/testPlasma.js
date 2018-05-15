const _ = require('lodash')

const RLP = require('rlp')
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");

import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

contract("Plasma ERC721 WIP", async function(accounts) {

    let cards;
    let plasma;

    let [authority, alice, bob, random_guy, random_guy2] = accounts;

    before("Deploys the contracts", async function() {
        plasma = await RootChain.new();
        cards = await CryptoCards.new(plasma.address);
        plasma.setCryptoCards(cards.address);
    });

    describe("Tests Plasma Deposits through the linked ERC721 contract", function() {

        it("Registers alice charlie bob", async function() {
            cards.register({from: alice});
            cards.register({from: bob});

            // Each should own 5 NFTs
            assert.equal(await cards.balanceOf.call(alice), 5);
            assert.equal(await cards.balanceOf.call(bob), 5);
        });

        it("Transfers NFT 1 from Alice to Plasma contract", async function() {

            // Call without extra data
            await cards.depositToPlasma(1, {from: alice});
            //await cards.safeTransferFrom(alice, plasma.address, 1, {from: alice});
            assert.equal(await cards.balanceOf.call(alice), 4);
            assert.equal(await cards.balanceOf.call(plasma.address), 1);

        });

        it("Transfers NFT 7 + data from Bob to Plasma contract", async function() {
           
            // Call with extra data, format: [chaindId][Address][Metadata]
            // Can use RLP format here for encoding UTXO Information
            await cards.depositToPlasmaWithData(7, '150x123WTF', {from: bob});
            //await cards.safeTransferFrom(bob, plasma.address, 7, '150x123WTF', {from: bob});
            assert.equal(await cards.balanceOf.call(bob), 4);
            assert.equal(await cards.balanceOf.call(plasma.address), 2);

        });

        it("Checks that events were emitted correctly", async function() {
            plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'}).get((error, res) => {
               let aliceLogs = res[0].args; // Why only 1 event gets emitted?
                assert.equal(aliceLogs.depositor, alice);
                assert.equal(aliceLogs.tokenId.toNumber(), 1);
                assert.equal(aliceLogs.data, '0x');

               let bobLogs = res[1].args; // Why only 1 event gets emitted?
                assert.equal(bobLogs.depositor, bob);
                assert.equal(bobLogs.tokenId.toNumber(), 7);
                assert.equal(bobLogs.data, '0x31353078313233575446');
            });
        });

    });

    describe("Exit Mechanism", function() {
        let start;
        let blk; 
        it("Random guy submits an exit for Bob's coins (exiter as a service)", async function() {
            // Bob wants to get his coins back after doing whatever on the plasma chain so exits
            let exitingTxBytes = '0x' + RLP.encode([1001, 7, bob]).toString('hex');
            let prevTxbytes = '0x' + RLP.encode([0, 0, 0]).toString('hex');  // doesn't matter
            let sig = exitingTxBytes;
            await plasma.startExit(prevTxbytes, exitingTxBytes, sig); // anyone can submit an exit for someone else
            start = (await web3.eth.getBlock('latest')).timestamp;
        });

        it("Another random guy finalizes the exits 1 week later. No challlenge has happened this time period", async function() {
            let expire = start + 3600 * 24 * 8; // 8 days pass, can finalize exit
            await increaseTimeTo(expire);
            await plasma.finalizeExits({from: random_guy2 });
        });

        it("Exit is finalized, Bob can withdraw his tokens", async function() {
            plasma.withdraw({from : bob});
            assert.equal((await cards.balanceOf.call(bob)).toNumber(), 5);
            assert.equal(await cards.balanceOf.call(plasma.address), 1);
        })
    });
});
