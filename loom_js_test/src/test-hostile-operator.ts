import test from 'tape'

import { enableHostilePlasmaCashOperator } from './config'
import { runDemo } from './demo'
import { runChallengeAfterDemo } from './challenge-after-demo'
import { runChallengeBetweenDemo } from './challenge-between-demo'
import { runChallengeBeforeDemo } from './challenge-before-demo'
import { runRespondChallengeBeforeDemo } from './respond-challenge-before-demo'

// TODO: Redeploy the Solidity contracts before each demo so the demos don't share any state.

enableHostilePlasmaCashOperator(true)
test('Plasma Cash with ERC721 Demo', runDemo)
test('Plasma Cash Challenge After Demo', runChallengeAfterDemo)
test('Plasma Cash Challenge Between Demo', runChallengeBetweenDemo)
test('Plasma Cash Challenge Before Demo', runChallengeBeforeDemo)
test('Plasma Cash Respond Challenge Before Demo', runRespondChallengeBeforeDemo)
