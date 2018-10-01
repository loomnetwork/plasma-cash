import test from 'tape'
import BN from 'bn.js'
import Web3 from 'web3'
import { PlasmaDB, SignedContract, IPlasmaDeposit, marshalDepositEvent } from 'loom-js'

import { increaseTime } from './ganache-helpers'
import { sleep, createTestEntity, ADDRESSES, ACCOUNTS } from './config'
import { EthCardsContract } from './cards-contract'

// Alice registers and has 5 coins, and she deposits 3 of them.
const ALICE_INITIAL_COINS = 5
const ALICE_DEPOSITED_COINS = 3
const COINS = [1, 2, 3]

// All the contracts are expected to have been deployed to Ganache when this function is called.
function setupContracts(web3: Web3): { cards: EthCardsContract } {
  const abi = require('./contracts/cards-abi.json')
  const cards = new EthCardsContract(new web3.eth.Contract(abi, ADDRESSES.token_contract))
  return { cards }
}

export async function runDemo(t: test.Test) {
  const endpoint = 'ws://127.0.0.1:8545'
  const web3 = new Web3(new Web3.providers.WebsocketProvider(endpoint))
  const { cards } = setupContracts(web3)
  const database = new PlasmaDB(endpoint, 'localhost:45578', '0x', ACCOUNTS.charlie) // Demo values to store in the db

  const authority = createTestEntity(web3, ACCOUNTS.authority)
  const alice = createTestEntity(web3, ACCOUNTS.alice)
  const bob = createTestEntity(web3, ACCOUNTS.bob)
  const charlie = createTestEntity(web3, ACCOUNTS.charlie, database)

  await cards.registerAsync(alice.ethAddress)
  let balance = await cards.balanceOfAsync(alice.ethAddress)
  t.equal(balance.toNumber(), 5)

  const startBlockNum = await web3.eth.getBlockNumber()

  for (let i = 0; i < ALICE_DEPOSITED_COINS; i++) {
    await cards.depositToPlasmaAsync({ tokenId: COINS[i], from: alice.ethAddress })
  }

  // Get deposit events for all
  const deposits: IPlasmaDeposit[] = await authority.getDepositEvents(new BN(0), true)
  t.equal(deposits.length, ALICE_DEPOSITED_COINS, 'All deposit events accounted for')

  for (let i = 0; i < deposits.length; i++) {
    const deposit = deposits[i]
    t.equal(deposit.blockNumber.toNumber(), i + 1, `Deposit ${i + 1} block number is correct`)
    t.equal(deposit.denomination.toNumber(), 1, `Deposit ${i + 1} denomination is correct`)
    t.equal(deposit.from, alice.ethAddress, `Deposit ${i + 1} sender is correct`)
  }

  balance = await cards.balanceOfAsync(alice.ethAddress)
  t.equal(
    balance.toNumber(),
    ALICE_INITIAL_COINS - ALICE_DEPOSITED_COINS,
    'alice should have 2 tokens in cards contract'
  )
  balance = await cards.balanceOfAsync(ADDRESSES.root_chain)
  t.equal(
    balance.toNumber(),
    ALICE_DEPOSITED_COINS,
    'plasma contract should have 3 tokens in cards contract'
  )

  // NOTE: In practice the Plasma Cash Oracle will submit the deposits to the DAppChain,
  // we're doing it here manually to simplify the test setup.
  for (let i = 0; i < deposits.length; i++) {
    await authority.submitPlasmaDepositAsync(deposits[i])
  }
  await sleep(2000)


  const coins = await alice.getUserCoinsAsync()
  t.ok(coins[0].slot.eq(deposits[0].slot), 'got correct deposit coins 1')
  t.ok(coins[1].slot.eq(deposits[1].slot), 'got correct deposit coins 2')
  t.ok(coins[2].slot.eq(deposits[2].slot), 'got correct deposit coins 3')

  // Alice to Bob, and Alice to Charlie. We care about the Alice to Bob
  // transaction
  const deposit3 = deposits[2]
  const deposit2 = deposits[1]
  // Alice -> Bob
  await alice.transferTokenAsync({
    slot: deposit3.slot,
    prevBlockNum: deposit3.blockNumber,
    denomination: 1,
    newOwner: bob
  })
  // Alice -> Charlie
  await alice.transferTokenAsync({
    slot: deposit2.slot,
    prevBlockNum: deposit2.blockNumber,
    denomination: 1,
    newOwner: charlie
  })
  const plasmaBlockNum1 = await authority.submitPlasmaBlockAsync()

  // Add an empty block in between (for proof of exclusion)
  await authority.submitPlasmaBlockAsync()

  // Bob -> Charlie
  await bob.transferTokenAsync({
    slot: deposit3.slot,
    prevBlockNum: new BN(1000),
    denomination: 1,
    newOwner: charlie
  })
  const plasmaBlockNum2 = await authority.submitPlasmaBlockAsync()

  const coin = await charlie.getPlasmaCoinAsync(deposit3.slot)
  const blocks = await bob.getBlockNumbersAsync(coin.depositBlockNum)

  const proofs = await bob.getCoinHistoryAsync(deposit3.slot, blocks)
  t.equal(await charlie.verifyCoinHistoryAsync(deposit3.slot, proofs), true)
  let charlieCoin = charlie.watchExit(deposit3.slot, coin.depositBlockNum)


  await charlie.startExitAsync({
    slot: deposit3.slot,
    prevBlockNum: plasmaBlockNum1,
    exitBlockNum: plasmaBlockNum2
  })
  charlie.stopWatching(charlieCoin)

  // Jump forward in time by 8 days
  await increaseTime(web3, 8 * 24 * 3600)
  // Charlie's exit should be finalizable...
  await authority.finalizeExitsAsync()
  // Charlie should now be able to withdraw the UTXO (plasma token) which contains ERC721 token #2
  // into his wallet.
  await charlie.withdrawAsync(deposit3.slot)

  balance = await cards.balanceOfAsync(alice.ethAddress)
  t.equal(balance.toNumber(), 2, 'alice should have 2 tokens in cards contract')
  balance = await cards.balanceOfAsync(bob.ethAddress)
  t.equal(balance.toNumber(), 0, 'bob should have no tokens in cards contract')
  balance = await cards.balanceOfAsync(charlie.ethAddress)
  t.equal(balance.toNumber(), 1, 'charlie should have 1 token in cards contract')

  // Close the websocket, hacky :/
  // @ts-ignore
  web3.currentProvider.connection.close()

  t.end()
}
