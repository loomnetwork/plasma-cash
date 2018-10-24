import test from 'tape'

import { runDemo } from './demo'
import { complexDemo } from './complex-demo'
import { runChallengeAfterDemo } from './challenge-after-demo'
import { PlasmaUser } from 'loom-js'

// TODO: Redeploy the Solidity contracts before each demo so the demos don't share any state.

PlasmaUser.contractName = 'plasmacash'
test('Plasma Cash ETH - Complex Demo', complexDemo)
test('Plasma Cash with ERC721 Demo', runDemo)
test('Plasma Cash Challenge After Demo', runChallengeAfterDemo)
