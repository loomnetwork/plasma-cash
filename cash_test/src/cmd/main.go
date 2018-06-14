package main

import (
	"client"
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

	// Alice deposits 3 of her coins to the plasma contract and gets 3 plasma nft
	// utxos in return
	tokenID := 1
	alice.Deposit(tokenID)
	alice.Deposit(tokenID + 1)
	alice.Deposit(tokenID + 2)

	//Alice to Bob, and Alice to Charlie. We care about the Alice to Bob
	// transaction
	utxoID := 2
	blkNum := 3
	_ = alice.SendTransaction(utxoID, blkNum, 1, bob.TokenContract.Account().Address)         //aliceToBob
	_ = alice.SendTransaction(utxoID-1, blkNum-1, 1, charlie.TokenContract.Account().Address) //randomTx
	authority.SubmitBlock()

	// Bob to Charlie
	blkNum = 1000                                                                       // the prev transaction was included in block 1000
	_ = bob.SendTransaction(utxoID, blkNum, 1, charlie.TokenContract.Account().Address) //bobToCharlie
	authority.SubmitBlock()
	/*
		// Charlie should be able to submit an exit by referencing blocks 0 and 1 which
		// included his transaction.
		utxoID = 2
		prev_tx_blkNum = 1000
		exiting_tx_blkNum = 2000
		charlie.start_exit(utxoID, prev_tx_blkNum, exiting_tx_blkNum)

		// After 8 days pass, charlie's exit should be finalizable
		w3 = charlie.root_chain.w3  // get w3 instance
		increaseTime(w3, 8 * 24 * 3600)
		authority.finalize_exits()
		// Charlie should now be able to withdraw the utxo which included token 2 to his
		// wallet.

		charlie.withdraw(utxoID)

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
