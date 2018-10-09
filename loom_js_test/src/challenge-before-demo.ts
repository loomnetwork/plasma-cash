import test from 'tape'
import BN from 'bn.js'
import Web3 from 'web3'
import { PlasmaUser } from 'loom-js'

import { increaseTime, getEthBalanceAtAddress } from './ganache-helpers'
import { sleep, ADDRESSES, ACCOUNTS, setupContracts } from './config'

export async function runChallengeBeforeDemo(t: test.Test) {
  const web3Endpoint = 'ws://127.0.0.1:8545'
  const dappchainEndpoint = 'http://localhost:46658'
  const web3 = new Web3(new Web3.providers.WebsocketProvider(web3Endpoint))
  const { cards } = setupContracts(web3)

  const authority = PlasmaUser.createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.authority)
  const dan  = PlasmaUser.createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.dan)
  const trudy = PlasmaUser.createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.trudy)
  const mallory = PlasmaUser.createUser(web3Endpoint, ADDRESSES.root_chain, dappchainEndpoint, ACCOUNTS.mallory)

  // Give Dan 5 tokens
  await cards.registerAsync(dan.ethAddress)
  let balance = await cards.balanceOfAsync(dan.ethAddress)
  t.equal(balance.toNumber(), 6)

  const startBlockNum = await web3.eth.getBlockNumber()

  // Dan deposits a coin
  await cards.depositToPlasmaAsync({ tokenId: 16, from: dan.ethAddress })
  const deposits = await dan.deposits()
  t.equal(deposits.length, 1, 'All deposit events accounted for')

  await sleep(8000)

  const plasmaBlock1 = await authority.submitPlasmaBlockAsync()
  const plasmaBlock2 = await authority.submitPlasmaBlockAsync()
  const deposit1Slot = deposits[0].slot

  // Dan starts watching
  const coin = await dan.getPlasmaCoinAsync(deposit1Slot)
  const danCoin = dan.watchExit(deposit1Slot, coin.depositBlockNum)

  // Trudy creates an invalid spend of the coin to Mallory
  // Low level call since trudy doesn't actually have the data for this transfer in her state
  await trudy.transferTokenAsync({
    slot: deposit1Slot,
    prevBlockNum: coin.depositBlockNum,
    denomination: 1,
    newOwner: mallory.ethAddress
  })

  // Operator includes it
  const trudyToMalloryBlock = await authority.submitPlasmaBlockAsync()

  
  // Low level call for the malicious transfers
  await mallory.transferTokenAsync({
    slot: deposit1Slot,
    prevBlockNum: trudyToMalloryBlock,
    denomination: 1,
    newOwner: trudy.ethAddress
  })

  // Operator includes it
  const malloryToTrudyBlock = await authority.submitPlasmaBlockAsync()

  // Low level call for the malicious exit
  await trudy.startExitAsync({
    slot: deposit1Slot,
    prevBlockNum: trudyToMalloryBlock,
    exitBlockNum: malloryToTrudyBlock
  })
  // Dan challenges with his coin that hasn't moved

  await sleep(2000)

  // 8 days pass without any response to the challenge
  await increaseTime(web3, 8 * 24 * 3600)
  await authority.finalizeExitsAsync()

  // Having successufly challenged Mallory's exit Dan should be able to exit the coin
  await dan.startExitAsync({
    slot: deposit1Slot,
    prevBlockNum: new BN(0),
    exitBlockNum: coin.depositBlockNum
  })
  dan.stopWatching(danCoin)

  // Jump forward in time by 8 days
  await increaseTime(web3, 8 * 24 * 3600)

  await authority.finalizeExitsAsync()

  await dan.withdrawAsync(deposit1Slot)

  const danBalanceBefore = await getEthBalanceAtAddress(web3, dan.ethAddress)
  await dan.withdrawBondsAsync()
  const danBalanceAfter = await getEthBalanceAtAddress(web3, dan.ethAddress)
  t.ok(danBalanceBefore.cmp(danBalanceAfter) < 0, 'END: Dan withdrew his bonds')

  const danTokensEnd = await cards.balanceOfAsync(dan.ethAddress)
  t.equal(danTokensEnd.toNumber(), 6, 'END: Dan has correct number of tokens')

  // Close the websocket, hacky :/
  // @ts-ignore
  web3.currentProvider.connection.close()
  authority.disconnect()
  dan.disconnect()
  trudy.disconnect()
  mallory.disconnect()
  t.end()
}
