package main

import (
	"client"
	"fmt"
	"log"
)

func main() {

	svc := client.NewChildChainService("http://localhost:8546")
	alice := client.NewClient(svc, client.GetRootChain("alice"), client.GetTokenContract("alice"))

	bob := client.NewClient(svc, client.GetRootChain("bob"), client.GetTokenContract("bob"))
	charlie := client.NewClient(svc, client.GetRootChain("charlie"), client.GetTokenContract("charlie"))
	authority := client.NewClient(svc, client.GetRootChain("authority"),
		client.GetTokenContract("authority"))

	// Give alice 5 tokens
	alice.TokenContract.Register()

	aliceTokensStart := alice.TokenContract.BalanceOf()
	log.Printf("Alice has %d tokens\n", aliceTokensStart)

	if aliceTokensStart != 5 {
		log.Fatalf("START: Alice has incorrect number of tokens")
	}
	bobTokensStart := bob.TokenContract.BalanceOf()
	log.Printf("Bob has %d tokens\n", bobTokensStart)
	if bobTokensStart != 0 {
		log.Fatalf("START: Bob has incorrect number of tokens")
	}
	charlieTokensStart := charlie.TokenContract.BalanceOf()
	log.Printf("Charlie has %d tokens\n", charlieTokensStart)
	if charlieTokensStart != 0 {
		log.Fatalf("START: Charlie has incorrect number of tokens")
	}

	fmt.Printf("initialized %v\n", authority)

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
											  bob.TokenContract.account.address)
		random_tx = alice.send_transaction(utxo_id-1, blk_num-1, 1,
										   charlie.TokenContract.account.address)
		authority.submit_block()

		// Bob to Charlie
		blk_num = 1000  // the prev transaction was included in block 1000
		bob_to_charlie = bob.send_transaction(utxo_id, blk_num, 1,
											  charlie.TokenContract.account.address)
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

		aliceTokensEnd = alice.TokenContract.BalanceOf()
		log.Printf('Alice has {} tokens'.format(aliceTokensEnd))
		assert (aliceTokensEnd == 2), "END: Alice has incorrect number of tokens"
		bobTokensEnd = bob.TokenContract.BalanceOf()
		log.Printf('Bob has {} tokens'.format(bobTokensEnd))
		assert (bobTokensEnd == 0), "END: Bob has incorrect number of tokens"
		charlieTokensEnd = charlie.TokenContract.BalanceOf()
		log.Printf('Charlie has {} tokens'.format(charlieTokensEnd))
		assert (charlieTokensEnd == 1), "END: Charlie has incorrect number of tokens"

		log.Printf('Plasma Cash with ERC721 tokens success :)')
	*/
}
