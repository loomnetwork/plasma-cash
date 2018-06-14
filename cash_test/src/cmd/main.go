package main

import (
	"client"
)

func main() {
	client := client.NewClient(client.NewChildChainService("http://localhost:8546"))
	/*
		alice = Client(container.get_root('alice'), container.get_token('alice'))
		bob = Client(container.get_root('bob'), container.get_token('bob'))
		charlie = Client(container.get_root('charlie'), container.get_token('charlie'))
		authority = Client(container.get_root('authority'),
						   container.get_token('authority'))

		// Give alice 5 tokens
		alice.token_contract.register()

		aliceTokensStart = alice.token_contract.balance_of()
		log.Printf("Alice has %d tokens", aliceTokensStart)
		assert (aliceTokensStart == 5), "START: Alice has incorrect number of tokens"
		bobTokensStart = bob.token_contract.balance_of()
		log.Printf('Bob has {} tokens'.format(bobTokensStart))
		assert (bobTokensStart == 0), "START: Bob has incorrect number of tokens"
		charlieTokensStart = charlie.token_contract.balance_of()
		log.Printf('Charlie has {} tokens'.format(charlieTokensStart))
		assert (charlieTokensStart == 0), \
				"START: Charlie has incorrect number of tokens"

	*/
	/*

		// Alice deposits 3 of her coins to the plasma contract and gets 3 plasma nft
		// utxos in return
		tokenId = 1
		alice.deposit(tokenId)
		alice.deposit(tokenId+1)
		alice.deposit(tokenId+2)
		``

		 //Alice to Bob, and Alice to Charlie. We care about the Alice to Bob
		// transaction
		utxo_id := 2
		blk_num := 3
		alice_to_bob = alice.send_transaction(utxo_id, blk_num, 1,
											  bob.token_contract.account.address)
		random_tx = alice.send_transaction(utxo_id-1, blk_num-1, 1,
										   charlie.token_contract.account.address)
		authority.submit_block()

		// Bob to Charlie
		blk_num = 1000  // the prev transaction was included in block 1000
		bob_to_charlie = bob.send_transaction(utxo_id, blk_num, 1,
											  charlie.token_contract.account.address)
		authority.submit_block()

		// Charlie should be able to submit an exit by referencing blocks 0 and 1 which
		// included his transaction.
		utxo_id = 2
		prev_tx_blk_num = 1000
		exiting_tx_blk_num = 2000
		charlie.start_exit(utxo_id, prev_tx_blk_num, exiting_tx_blk_num)

		// After 8 days pass, charlie's exit should be finalizable
		w3 = charlie.root_chain.w3  // get w3 instance
		increaseTime(w3, 8 * 24 * 3600)
		authority.finalize_exits()
		// Charlie should now be able to withdraw the utxo which included token 2 to his
		// wallet.

		charlie.withdraw(utxo_id)

		aliceTokensEnd = alice.token_contract.balance_of()
		log.Printf('Alice has {} tokens'.format(aliceTokensEnd))
		assert (aliceTokensEnd == 2), "END: Alice has incorrect number of tokens"
		bobTokensEnd = bob.token_contract.balance_of()
		log.Printf('Bob has {} tokens'.format(bobTokensEnd))
		assert (bobTokensEnd == 0), "END: Bob has incorrect number of tokens"
		charlieTokensEnd = charlie.token_contract.balance_of()
		log.Printf('Charlie has {} tokens'.format(charlieTokensEnd))
		assert (charlieTokensEnd == 1), "END: Charlie has incorrect number of tokens"

		log.Printf('Plasma Cash with ERC721 tokens success :)')
	*/
}
