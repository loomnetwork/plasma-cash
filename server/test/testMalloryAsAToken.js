const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");
const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
import {increaseTimeTo, duration} from './helpers/increaseTime'
import assertRevert from './helpers/assertRevert.js';

const txlib = require('./UTXO.js')

contract("Plasma ERC721 - Multiple Validators and ERC721 tokens", async function(accounts) {

    const t1 = 3600 * 24 * 3; // 3 days later
    const t2 = 3600 * 24 * 5; // 5 days later

    // Alice registers and has 5 coins, and she deposits 3 of them.
    const INITIAL_COINS = 5;
    const DEPOSITED_COINS = 3;
    const COINS = [1, 2, 3];

    let cards, cards2;
    let plasma;
    let vmc;
    let events;
    let t0;

    let [owner, validator1, validator2, alice, bob, charlie, random_guy, mallory] = accounts;

    beforeEach(async function() {
        vmc = await ValidatorManagerContract.new({from: owner});
        plasma = await RootChain.new(vmc.address, {from: owner});
        cards = await CryptoCards.new(plasma.address);

        await vmc.toggleToken(cards.address);

        await cards.register({from: alice});
        assert.equal(await cards.balanceOf.call(alice), 5);

        let ret;
        // Deposit all of Alice's cards
        for (let i = 0; i < DEPOSITED_COINS; i ++) {
            await cards.depositToPlasma(COINS[i], {from: alice});
        }

        assert.equal((await cards.balanceOf.call(alice)).toNumber(), INITIAL_COINS - DEPOSITED_COINS);
        assert.equal((await cards.balanceOf.call(plasma.address)).toNumber(), DEPOSITED_COINS);

        const depositEvent = plasma.Deposit({}, {fromBlock: 0, toBlock: 'latest'});
        events = await txlib.Promisify(cb => depositEvent.get(cb));

        // Check that events were emitted properly
        let coin;
        for (let i = 0; i < DEPOSITED_COINS ; i++ ) {
            coin = events[i].args;
            assert.equal(coin.blockNumber.toNumber(), i+1);
            assert.equal(coin.denomination.toNumber(), 1);
            assert.equal(coin.from, alice);
        }
    });

    describe('Mallory tries to steal coins from Alice', function() {
        // This set of test(s) proves that even if the operator 
        // decides to enable arbitrary accounts as token accounts, 
        // they are not able to abuse the system and steal other user's coins
        // due to the `contractAddress` parameter in a coin's state.

        it("Mallory manually calls the onERC721Received and is able to generate a fake coin. Exits it, nobody can challenge. Even so, cannot steal the coin", async function() {
            // Validator triggers mallory as a token contract, 
            // allowing mallory to call the receiver function
            await vmc.toggleToken(mallory);

            let blk = web3.eth.blocknumber;
            await plasma.onERC721Received(mallory, COINS[1], '0x0', {from: mallory});

            let malloryEvent = plasma.Deposit({from: mallory}, {fromBlock: 0, toBlock: 'latest'});
            events = await txlib.Promisify(cb => malloryEvent.get(cb));

            let UTXO = {'slot': events[0]['args'].slot, 'block': events[0]['args'].blockNumber.toNumber()};

            // Mallory now has a UTXO of `COINS[1]` which is owned by Alice.
            // Tries to exit it.
            // Note however, that the `contractAddress` in the coin's state is not 
            // the ERC721, but Mallory's address!

            let ret = txlib.createUTXO(UTXO.slot, 0, mallory, mallory);
            let utxo = ret.tx;
            let sig = ret.sig;

            // The exit works, since it is a valid attestation
            await plasma.startExit(
                     UTXO.slot,
                    '0x', utxo,
                    '0x0', '0x0',
                     sig,
                     [0, UTXO.block],
                     {'from': mallory, 'value': web3.toWei(0.1, 'ether')}
            );
            t0 = await web3.eth.getBlock('latest').timestamp;

            await increaseTimeTo(t0 + t1 + t2);
            await plasma.finalizeExit(UTXO.slot, {from: random_guy});

            // However mallory cannot withdraw the coins since `contractAddress` 
            // of the coin is actually mallory's address and NOT the actual token address
            assertRevert(plasma.withdraw(UTXO.slot, {from: mallory}));

            // Nonetheless, Mallory's exit was actually valid and can withdraw the bond
            await txlib.withdrawBonds(plasma, mallory, 0.1)
        });
    });

});
