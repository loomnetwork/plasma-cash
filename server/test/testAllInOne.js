const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const LoomToken = artifacts.require("LoomToken");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract("Plasma Cash - All In One", async function(accounts) {

    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    // Alice registers and has 5 coins, and she deposits 3 of them.
    const ALICE_INITIAL_COINS = 5;
    const ALICE_DEPOSITED_COINS = 3;
    const coins = [1, 2, 3];

    let erc20;
    let erc721;
    let plasma;
    let vmc;
    let events;
    let t0;

    let [authority, alice, bob, charlie, dylan, elliot, random_guy, random_guy2, challenger] = accounts;

    const blk_1 = 1000
    const blk_2 = 2000
    const blk_3 = 3000
    const blk_4 = 4000

    const DECIMALS = 10 ** 18;
    const denominations = [
        3000 * DECIMALS, 
        2000 * DECIMALS, 
        4000 * DECIMALS
    ];

    const ethers = [
        web3.toWei(1, 'ether'),
        web3.toWei(4, 'ether'),
        web3.toWei(5, 'ether')
    ];
    let UTXO
    let slots

    it('Does not accept more than MAX_VALUE = 10 Ether', async function() {
        vmc = await ValidatorManagerContract.new({from: authority});
        plasma = await RootChain.new(vmc.address, {from: authority});
        await plasma.sendTransaction({from: alice, value: ethers[2], gas: 220000 });
        await plasma.sendTransaction({from: alice, value: ethers[2], gas: 220000 });
        await assertRevert(plasma.sendTransaction({from: alice, value: ethers[2], gas: 220000 }))
        assert.equal((await web3.eth.getBalance(plasma.address)).toNumber(), web3.toWei(10, 'ether'));
    })

    it('Does not accept ether when paused', async function() {
        vmc = await ValidatorManagerContract.new({from: authority});
        plasma = await RootChain.new(vmc.address, {from: authority});
        await plasma.sendTransaction({from: alice,  value: web3.toWei(1, 'ether'), gas: 220000 });
        await plasma.pause({from: authority})
        await assertRevert(plasma.sendTransaction({from: alice, value: web3.toWei(1, 'ether'), gas: 220000 }))
        await plasma.unpause({from: authority})
        await plasma.sendTransaction({from: alice, value: ethers[2], gas: 220000 })
    })

    describe('Multideposit tests', function() {
        beforeEach(async function() {
            vmc = await ValidatorManagerContract.new({from: authority});
            plasma = await RootChain.new(vmc.address, {from: authority});
            erc20 = await LoomToken.new(plasma.address, {from: authority});
            erc721 = await CryptoCards.new(plasma.address, {from: authority});

            await vmc.toggleToken(erc20.address, {from: authority});
            await vmc.toggleToken(erc721.address, {from: authority});

            await erc20.transfer(alice, 10000 * DECIMALS, {from: authority});
            await erc721.register({from: alice});

            for (let i = 0; i < denominations.length - 1; i ++) {
                await web3.eth.sendTransaction({from: alice, to: plasma.address, value: ethers[i], gas: 220000 });
                await erc20.depositToPlasma(denominations[i], {from: alice});
                await erc721.depositToPlasma(coins[i], {from: alice});
            }

            // Make the last transfer come from approve/deposit pattern
            let ind = denominations.length - 1;
            await web3.eth.sendTransaction({from: alice, to: plasma.address, value: ethers[ind], gas: 220000 });
            await erc20.approve(plasma.address, denominations[ind], { 'from': alice })
            await plasma.depositERC20(denominations[ind], erc20.address, { 'from': alice})

            await erc721.approve(plasma.address, coins[ind], { 'from': alice })
            await plasma.depositERC721(coins[ind], erc721.address, { 'from': alice})

            assert.equal((await erc20.balanceOf.call(alice)).toNumber(), 1000 * DECIMALS);
            assert.equal((await erc20.balanceOf.call(plasma.address)).toNumber(), 9000 * DECIMALS);
            assert.equal(await erc721.balanceOf.call(plasma.address), 3);
            assert.equal(await erc721.balanceOf.call(alice), 2);
            assert.equal((await web3.eth.getBalance(plasma.address)).toNumber(), web3.toWei(10, 'ether'));

            const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
            events = await txlib.Promisify(cb => depositEvent.get(cb));

            UTXO = [
                {'slot': events[0]['args'].slot, 'block': events[0]['args'].blockNumber.toNumber()},
                {'slot': events[1]['args'].slot, 'block': events[1]['args'].blockNumber.toNumber()},
                {'slot': events[2]['args'].slot, 'block': events[2]['args'].blockNumber.toNumber()},
            ]
            slots = UTXO.map(u => u.slot)
        });

        it('Cancel exits / getExit', async function() {
            for (let i in UTXO) {
                t0 = await txlib.exitDeposit(plasma, alice, UTXO[i])
            }
            await increaseTimeTo(t0 + t1 + t2);

            // Fails to cancel exit from anyone other than the exitor
            assertRevert(plasma.cancelExits(slots, {from: random_guy2}))
            await plasma.cancelExits(slots, {from: alice})

            // Nothing should happen
            await plasma.finalizeExits(slots, {from: random_guy2 });

            // State of both coins should be 0 since their exits got cancelled
            for (let i in slots) {
                assert.equal(await txlib.getState(plasma, slots[i]), 0)
            }
        });

            it('C = Deposit, PC = Null', async function() {
                for (let i in UTXO) {
                    t0 = await txlib.exitDeposit(plasma, alice, UTXO[i])
                }
                await increaseTimeTo(t0 + t1 + t2);
                await plasma.finalizeExits(slots, {from: random_guy2 });
                for (let i in slots) {
                    assert.equal(await txlib.getState(plasma, slots[i]), 2)
                }
            });

            it('C = Non-Deposit, PC = Deposit', async function() {
                let leaves = [];
                let txs = [];
                for (let i in UTXO) {
                    let aUTXO = UTXO[i];
                    let tx = txlib.createUTXO(aUTXO.slot, 0, alice, bob);
                    leaves.push(tx.leaf)
                    txs.push(tx)
                }
                let tree_1000 = await txlib.submitTransactions(authority, plasma, 1000, leaves);

                for (let i in UTXO) {
                    t0 = await txlib.exit(plasma, bob,
                        UTXO[i].slot,

                        { 'block': blk_1, 'tx': txs[i] },
                        tree_1000.createMerkleProof(UTXO[i].slot),

                        { 'block': UTXO[i].block, 'tx': txlib.createUTXO(UTXO[i].slot, 0, alice, alice) },
                        '0x',
                    )
                }
                await increaseTimeTo(t0 + t1 + t2);
                await plasma.finalizeExits(slots, {from: random_guy2 });
                for (let i in slots) {
                    assert.equal(await txlib.getState(plasma, slots[i]), 2)
                }
            });

            it('C = Non-Deposit, PC = Non-Deposit', async function() {
                let prevBlock = 0;
                let leaves = [];
                let prev_txs = [];
                for (let i in UTXO) {
                    let tx = txlib.createUTXO(UTXO[i].slot, prevBlock, alice, bob);
                    leaves.push(tx.leaf)
                    prev_txs.push(tx)
                }
                let tree_1000 = await txlib.submitTransactions(authority, plasma, 1000, leaves);

                prevBlock = 1000;
                leaves = [];
                let txs = [];
                for (let i in UTXO) {
                    let aUTXO = UTXO[i];
                    let tx = txlib.createUTXO(aUTXO.slot, prevBlock, bob, charlie);
                    leaves.push(tx.leaf)
                    txs.push(tx)
                }
                let tree_2000 = await txlib.submitTransactions(authority, plasma, 2000, leaves);


                for (let i in UTXO) {
                    t0 = await txlib.exit(plasma, charlie,
                        UTXO[i].slot,

                        { 'block': blk_2, 'tx': txs[i] },
                        tree_2000.createMerkleProof(UTXO[i].slot),

                        { 'block': blk_1, 'tx': prev_txs[i] },
                        tree_1000.createMerkleProof(UTXO[i].slot),
                    )
                }
                await increaseTimeTo(t0 + t1 + t2);
                await plasma.finalizeExits(slots, {from: random_guy2 });
                for (let i in slots) {
                    assert.equal(await txlib.getState(plasma, slots[i]), 2)
                }
            });

        });

});
