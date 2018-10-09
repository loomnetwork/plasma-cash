import test from 'tape'
import Web3 from 'web3'
import { createUser } from 'loom-js'

import { increaseTime, getEthBalanceAtAddress } from './ganache-helpers'
import { sleep, ADDRESSES, ACCOUNTS, setupContracts } from './config'

export async function runChallengeBetweenDemo(t: test.Test) {
  const web3Endpoint = 'ws://127.0.0.1:8545'
  const dappchainEndpoint = 'http://localhost:46658'
  const web3 = new Web3(new Web3.providers.WebsocketProvider(web3Endpoint))
  const { cards } = setupContracts(web3)

  const authority = createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.authority)
  const alice = createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.alice)
  const bob = createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.bob)
  const eve = createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.eve)

  const bobTokensStart = await cards.balanceOfAsync(bob.ethAddress)

  // Give Eve 5 tokens
  await cards.registerAsync(eve.ethAddress)

  // Eve deposits a coin
  await cards.depositToPlasmaAsync({ tokenId: 11, from: eve.ethAddress })
  const deposits = await eve.deposits()
  t.equal(deposits.length, 1, 'Eve has correct number of deposits')

  await sleep(8000)

  const deposit1Slot = deposits[0].slot

  // Eve sends her plasma coin to Bob
  const coin = await eve.getPlasmaCoinAsync(deposit1Slot)
  await eve.transfer(deposit1Slot, bob.ethAddress)

  const eveToBobBlockNum = await authority.submitPlasmaBlockAsync()

  const blocks = await eve.getBlockNumbersAsync(coin.depositBlockNum)
  const proofs = await eve.getCoinHistoryAsync(deposit1Slot, blocks)
  t.equal(await bob.verifyCoinHistoryAsync(deposit1Slot, proofs), true)
  const bobCoin = bob.watchExit(deposit1Slot, coin.depositBlockNum)

  // Eve sends this same plasma coin to Alice
  await eve.transfer(deposit1Slot, alice.ethAddress)

  const eveToAliceBlockNum = await authority.submitPlasmaBlockAsync()

  // Alice attempts to exit her double-spent coin
  // Low level call to exit the double spend
  await alice.startExitAsync({
    slot: deposit1Slot,
    prevBlockNum: coin.depositBlockNum,
    exitBlockNum: eveToAliceBlockNum
  })
  // Bob challenges here

  await sleep(2000)

  await bob.exit(deposit1Slot)
  bob.stopWatching(bobCoin)

  // Jump forward in time by 8 days
  await increaseTime(web3, 8 * 24 * 3600)

  await authority.finalizeExitsAsync()

  await bob.withdrawAsync(deposit1Slot)

  const bobBalanceBefore = await getEthBalanceAtAddress(web3, bob.ethAddress)

  await bob.withdrawBondsAsync()

  const bobBalanceAfter = await getEthBalanceAtAddress(web3, bob.ethAddress)

  t.ok(bobBalanceBefore.cmp(bobBalanceAfter) < 0, 'END: Bob withdrew his bonds')

  const bobTokensEnd = await cards.balanceOfAsync(bob.ethAddress)

  t.equal(
    bobTokensEnd.toNumber(),
    bobTokensStart.toNumber() + 1,
    'END: Bob has correct number of tokens'
  )

  // Close the websocket, hacky :/
  // @ts-ignore
  authority.web3.currentProvider.connection.close()
  // @ts-ignore
  alice.web3.currentProvider.connection.close()
  // @ts-ignore
  bob.web3.currentProvider.connection.close()
  // @ts-ignore
  eve.web3.currentProvider.connection.close()
  // @ts-ignore
  web3.currentProvider.connection.close()

  t.end()
}
