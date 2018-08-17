const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const LoomToken = artifacts.require("LoomToken");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
const SparseMerkleTree = require('./SparseMerkleTree.js');
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract("Plasma Debit - All In One", async function(accounts) {

    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    // Alice registers and has 5 coins, and she deposits 3 of them.
    const ALICE_INITIAL_COINS = 5;
    const ALICE_DEPOSITED_COINS = 3;
    const coins = [1, 2, 3];
    const ETHER = 10 ** 18;

    let erc20;
    let erc721;
    let plasma;
    let vmc;
    let events;
    let t0;

    let [authority, alice, bob, charlie, dylan, elliot, random_guy, random_guy2, challenger] = accounts;

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

    beforeEach(async function() {
        vmc = await ValidatorManagerContract.new({from: authority});
        plasma = await RootChain.new(vmc.address, {from: authority});
        erc20 = await LoomToken.new(plasma.address, {from: authority});
        erc721 = await CryptoCards.new(plasma.address, {from: authority});

        await vmc.toggleToken(erc20.address, {from: authority});
        await vmc.toggleToken(erc721.address, {from: authority});

        await erc20.transfer(alice, 10000 * DECIMALS, {from: authority});
        await erc721.register({from: alice});

        for (let i = 0; i < denominations.length; i ++) {
            await web3.eth.sendTransaction({from: alice, to: plasma.address, value: ethers[i], gas: 250000 });
            await erc20.depositToPlasma(denominations[i], {from: alice});
            await erc721.depositToPlasma(coins[i], {from: alice});
        }
        assert.equal(await erc20.balanceOf.call(alice), 1000 * DECIMALS);
        assert.equal(await erc20.balanceOf.call(plasma.address), 9000 * DECIMALS);

        assert.equal(await erc721.balanceOf.call(plasma.address), 3);
        assert.equal(await erc721.balanceOf.call(alice), 2);

        assert.equal(await web3.eth.getBalance(plasma.address), web3.toWei(10, 'ether'));

        const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
        events = await txlib.Promisify(cb => depositEvent.get(cb));
    });

    describe('Plasma Debit', function() {

		it('Operator provides liquidity!', async function() {
			let UTXO = [
                {'slot': events[0]['args'].slot, 'block': events[0]['args'].blockNumber.toNumber()},
                {'slot': events[1]['args'].slot, 'block': events[1]['args'].blockNumber.toNumber()},
            ]
            // Fill up the ETH token, had 1 ether
            await plasma.provideLiquidity(UTXO[0].slot, 0, {'value':  web3.toWei(14, 'ether') });
            // Fill up the ERC20 token, had 3000 erc20 coins
            await erc20.approve(plasma.address, 4000 * DECIMALS, {from: authority});
            await plasma.provideLiquidity(UTXO[1].slot, 4000 * DECIMALS, {from: authority});

            // TODO Improve ux, if user does not provide a value the contract
            // should be checking and automatically giving the user the default
            // balance value
            let values = [ 1 * ETHER, 3000 * DECIMALS ];
            let prevBlock = 0;
            for (let i in UTXO) {
                let aUTXO = UTXO[i];
                let ret = txlib.createDebitUTXO(aUTXO.slot, prevBlock, alice, alice, values[i], 0);
                let utxo = ret.tx;
                let sig = ret.sig;

                await plasma.startExit(
                    aUTXO.slot,
                    '0x', utxo,
                    '0x0', '0x0',
                    sig,
                    [prevBlock, aUTXO.block],
                    {'from': alice, 'value': web3.toWei(0.1, 'ether')}
                );
            }
            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo(t0 + t1);
            await plasma.finalizeExits({from: random_guy2 });

            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });
            for (let i in UTXO) {
                let aUTXO = UTXO[i];
                await plasma.withdraw(aUTXO.slot, {from : alice });
            }

            const withdrewEvent = plasma.Withdrew({}, {fromBlock: 0, toBlock: 'latest'});
            const withdrew = await txlib.Promisify(cb => withdrewEvent.get(cb));

            // The authority should have got the liquidity provided back.
            assert.equal(withdrew[0]['args'].denomination, web3.toWei(1, 'ether'));
            assert.equal(withdrew[0]['args'].toOperator, web3.toWei(14, 'ether'));
            assert.equal(withdrew[1]['args'].denomination, 3000 * DECIMALS);
            assert.equal(withdrew[1]['args'].toOperator, 4000 * DECIMALS);
            // Alice has her coins back.
            await txlib.withdrawBonds(plasma, alice, 0.2);
        });

		it('User exits a partial coin, cooperative case', async function() {
			let UTXO =
                {'slot': events[0]['args'].slot, 'block': events[0]['args'].blockNumber.toNumber()};
            let prevBlock = 0;

            // Auth and Alice sign nonce 0
            let changeBalance = txlib.createDebitUTXO(
                UTXO.slot,
                UTXO.block,
                alice,
                alice,
                0.75 * ETHER,
                0
            ); // Alice has now signed the TXO.

            let auth_sig = txlib.signHash(authority, changeBalance.hash);

            // Auth and alice now sign nonce 1
            changeBalance = txlib.createDebitUTXO(
                UTXO.slot,
                UTXO.block,
                alice,
                alice,
                0.6 * ETHER,
                1
            );

            auth_sig = txlib.signHash(authority, changeBalance.hash);

            // Alice should be able to exit this TXO and get 0.5 out of the 1
            // ether.

            let utxo = changeBalance.tx;
            let sig = changeBalance.sig;
            let txs = [changeBalance.leaf]

            // Authority submits a block to plasma with that transaction included
            let tree = await txlib.submitTransactions(authority, plasma, txs);
            let submittedBlock = 1000;

            let exiting_tx_proof = tree.createMerkleProof(UTXO.slot)

            // Auth and alice now sign nonce 1
            let prev_tx = txlib.createDebitUTXO(
                UTXO.slot,
                0,
                alice,
                alice,
                1 * ETHER,
                0
            ).tx;

            await plasma.startExit(
                UTXO.slot,
                prev_tx, utxo,
                '0x0', exiting_tx_proof,
                sig,
                [UTXO.block, submittedBlock],
                {'from': alice, 'value': web3.toWei(0.1, 'ether')}
            );

            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExits({from: random_guy2 });

            await plasma.withdraw(UTXO.slot, {from : alice });

            const withdrewEvent = plasma.Withdrew({}, {fromBlock: 0, toBlock: 'latest'});
            const withdrew = await txlib.Promisify(cb => withdrewEvent.get(cb));

            // The authority should have got the liquidity provided back.
            assert.equal(withdrew[0]['args'].denomination, web3.toWei(0.6, 'ether'));
            assert.equal(withdrew[0]['args'].toOperator, web3.toWei(0.4, 'ether'));
            await txlib.withdrawBonds(plasma, alice, 0.1);
        });

		it('No withholding - Operator tries to exit with an earlier nonce', async function() {
			let UTXO =
                {'slot': events[0]['args'].slot, 'block': events[0]['args'].blockNumber.toNumber()};
            let prevBlock = 0;

            // Auth and Alice sign nonce 0
            let changeBalance_0 = txlib.createDebitUTXO(
                UTXO.slot,
                UTXO.block,
                alice,
                alice,
                0.5 * ETHER,
                0
            );

            // Operator signs nonce 0
            let auth_sig_0 = txlib.signHash(authority, changeBalance_0.hash);

            // Auth and alice now sign nonce 1
            let changeBalance_1 = txlib.createDebitUTXO(
                UTXO.slot,
                UTXO.block,
                alice,
                alice,
                0.75 * ETHER,
                1
            );

            // Operator has now signed a later state, contract will accept an
            // exit involving nonce 0, but a rational user will not submit that
            // exit. They will always submit 0.
            let auth_sig_1 = txlib.signHash(authority, changeBalance_1.hash);

            let utxo_0 = changeBalance_0.tx;
            let sig_0 = changeBalance_0.sig;
            let txs_0 = [changeBalance_0.leaf]

            // The operator maliciously submits the transaction at an earlier
            // nonce, the user should be able to prove that the operator messed
            // up.
            let tree_0 = await txlib.submitTransactions(authority, plasma, txs_0);
            let submittedBlock = 1000;

            // Previous tx that is used for the exit
            let prev_tx = txlib.createDebitUTXO(
                UTXO.slot,
                0,
                alice,
                alice,
                1 * ETHER,
                0
            ).tx;

            let txs_1 = [changeBalance_1.leaf];
            let leaves = {};
            for (let l in txs_0) {
                leaves[txs_0[l].slot] = txs_0[l].hash;
            }
            let tree_1 = new SparseMerkleTree(64, leaves);
            let honest_proof = tree_1.createMerkleProof(UTXO.slot);
            let utxo_1 = changeBalance_1.tx;
            let sig_1 = changeBalance_1.sig;

            // The exit by the user should fail, because the transaction that
            // they are expecting to see included was not actually included by
            // the operator
            assertRevert(plasma.startExit(
                UTXO.slot,
                prev_tx, utxo_1,
                '0x0', honest_proof,
                sig_1,
                [UTXO.block, submittedBlock],
                {'from': alice, 'value': web3.toWei(0.1, 'ether')}
            ));

            // Since user is unable to exit, they will try to challenge and
            // claim that the operator has included a transaction which is not
            // valid, according to the user.
            await plasma.challengeChannel(
                UTXO.slot,
                submittedBlock, // the block used
                utxo_1, // the later transaction
                sig_1, // the user's sig
                auth_sig_1, // the operator's sig
                {'from': alice, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = (await web3.eth.getBlock('latest')).timestamp;

            // The operator must respond by revealign that the transaction they
            // included had a bigger OR equal nonce. In this case they cannot
            // reveal within the set challenge period (1 week again?!) and so
            // they lose the whole coin and the user gets it.
            // In this case, they cannot
            // let included_proof = tree_0.createMerkleProof(UTXO.slot);
            assertRevert(plasma.respondChallengeChannelWithSignature(
                UTXO.slot,
            //     submittedBlock, // the block used
                utxo_0, // the responding transaction

                sig_1, // the user's sig
            ));

            // Challenge period passes, all outstanding challenges on channels
            // will be finalized now, and an exit will be initiated (?) for
            // each of these coins?
            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeChannel(UTXO.slot, {from: random_guy2 });

            // Let the exit mature
            t0 = (await web3.eth.getBlock('latest')).timestamp;
            await increaseTimeTo(t0 + t1 + t2);

            await plasma.finalizeExits({from: random_guy2});
            await plasma.withdraw(UTXO.slot, {from : alice });

            const withdrewEvent = plasma.Withdrew({}, {fromBlock: 0, toBlock: 'latest'});
            const withdrew = await txlib.Promisify(cb => withdrewEvent.get(cb));

            // // The authority should have got nothing since they cheated
            assert.equal(withdrew[0]['args'].denomination, web3.toWei(1, 'ether'));
            assert.equal(withdrew[0]['args'].toOperator, web3.toWei(0, 'ether'));
            await txlib.withdrawBonds(plasma, alice, 0.1);
        });

    });

});
