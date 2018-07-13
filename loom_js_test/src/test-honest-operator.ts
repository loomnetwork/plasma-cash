import test from 'tape'

import { enableHostilePlasmaCashOperator } from './config'
import { runDemo } from './demo'
import { runChallengeAfterDemo } from './challenge-after-demo'

// TODO: Redeploy the Solidity contracts before each demo so the demos don't share any state.

enableHostilePlasmaCashOperator(false)
test('Plasma Cash with ERC721 Demo', runDemo)
test('Plasma Cash Challenge After Demo', runChallengeAfterDemo)
