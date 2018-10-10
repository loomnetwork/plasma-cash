import test from 'tape'
import Web3 from 'web3'
import BN from 'bn.js'
import { PlasmaUser } from 'loom-js'

import { increaseTime, getEthBalanceAtAddress } from './ganache-helpers'
import { sleep, ADDRESSES, ACCOUNTS, setupContracts } from './config'

export async function runChallengeAfterDemo(t: test.Test) {
  const web3Endpoint = 'ws://127.0.0.1:8545'
  const dappchainEndpoint = 'http://localhost:46658'
  const web3 = new Web3(new Web3.providers.WebsocketProvider(web3Endpoint))
  const { cards } = setupContracts(web3)

  const authority = PlasmaUser.createUser(
    web3Endpoint,
    ADDRESSES.root_chain,
    dappchainEndpoint,
    ACCOUNTS.authority
  )
  const mallory = PlasmaUser.createUser(
    web3Endpoint,
    ADDRESSES.root_chain,
    dappchainEndpoint,
    ACCOUNTS.mallory
  )
  const dan = PlasmaUser.createUser(
    web3Endpoint,
    ADDRESSES.root_chain,
    dappchainEndpoint,
    ACCOUNTS.dan
  )

  // Give Mallory 5 tokens
  await cards.registerAsync(mallory.ethAddress)

  const danTokensStart = await cards.balanceOfAsync(dan.ethAddress)
  t.equal(danTokensStart.toNumber(), 0, 'START: Dan has correct number of tokens')
  const malloryTokensStart = await cards.balanceOfAsync(mallory.ethAddress)
  t.equal(malloryTokensStart.toNumber(), 5, 'START: Mallory has correct number of tokens')

  // Mallory deposits one of her coins to the plasma contract
  await cards.depositToPlasmaAsync({ tokenId: 6, from: mallory.ethAddress })
  await cards.depositToPlasmaAsync({ tokenId: 7, from: mallory.ethAddress })

  const deposits = await mallory.deposits()
  t.equal(deposits.length, 2, 'Mallory has correct number of deposits')

  const malloryTokensPostDeposit = await cards.balanceOfAsync(mallory.ethAddress)
  t.equal(
    malloryTokensPostDeposit.toNumber(),
    3,
    'POST-DEPOSIT: Mallory has correct number of tokens'
  )

  await sleep(8000)

  await authority.submitPlasmaBlockAsync()
  await authority.submitPlasmaBlockAsync()

  const deposit1Slot = deposits[0].slot

  // Mallory -> Dan
  const coin = await mallory.getPlasmaCoinAsync(deposit1Slot)
  await mallory.transferAsync(deposit1Slot, dan.ethAddress)
  await authority.submitPlasmaBlockAsync()
  t.equal(await dan.checkHistoryAsync(coin), true, 'Coin history verified')
  const danCoin = dan.watchExit(deposit1Slot, coin.depositBlockNum)

  // Mallory attempts to exit spent coin (the one sent to Dan)
  // Needs to use the low level API to make an invalid tx
  await mallory.startExitAsync({
    slot: deposit1Slot,
    prevBlockNum: new BN(0),
    exitBlockNum: coin.depositBlockNum
  })

  // Having successufly challenged Mallory's exit Dan should be able to exit the coin
  await sleep(2000)
  await dan.exitAsync(deposit1Slot)

  dan.stopWatching(danCoin)

  // Jump forward in time by 8 days
  await increaseTime(web3, 8 * 24 * 3600)

  await authority.finalizeExitsAsync()

  await dan.withdrawAsync(deposit1Slot)

  const danBalanceBefore = await getEthBalanceAtAddress(web3, dan.ethAddress)
  await dan.withdrawBondsAsync()
  const danBalanceAfter = await getEthBalanceAtAddress(web3, dan.ethAddress)
  t.ok(danBalanceBefore.cmp(danBalanceAfter) < 0, 'END: Dan withdrew his bonds')

  const malloryTokensEnd = await cards.balanceOfAsync(mallory.ethAddress)
  t.equal(malloryTokensEnd.toNumber(), 3, 'END: Mallory has correct number of tokens')
  const danTokensEnd = await cards.balanceOfAsync(dan.ethAddress)
  t.equal(danTokensEnd.toNumber(), 1, 'END: Dan has correct number of tokens')

  // Close the websocket, hacky :/
  // @ts-ignore
  web3.currentProvider.connection.close()
  authority.disconnect()
  dan.disconnect()
  mallory.disconnect()
  t.end()
}
