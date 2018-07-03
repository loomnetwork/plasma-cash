package main

import (
	"client"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {

	client.InitClients("http://localhost:8545")
	client.InitTokenClient("http://localhost:8545")

	svc, err := client.NewLoomChildChainService("http://localhost:46658/rpc", "http://localhost:46658/query")
	exitIfError(err)

	alice := client.NewClient(svc, client.GetRootChain("alice"), client.GetTokenContract("alice"))
	bob := client.NewClient(svc, client.GetRootChain("bob"), client.GetTokenContract("bob"))
	eve := client.NewClient(svc, client.GetRootChain("eve"), client.GetTokenContract("eve"))
	authority := client.NewClient(svc, client.GetRootChain("authority"),
		client.GetTokenContract("authority"))

	bobTokensStart, err := bob.TokenContract.BalanceOf()
	exitIfError(err)

	// Give Eve 5 tokens
	eve.Register()

	// Eve deposits a coin
	txHash := eve.Deposit(11)
	time.Sleep(1 * time.Second)
	deposit1, err := eve.RootChain.DepositEventData(txHash)
	exitIfError(err)

	// wait to make sure that events get fired correctly
	time.Sleep(2)

	// Eve sends her plasma coin to Bob
	coin, err := eve.PlasmaCoin(deposit1.Slot)
	exitIfError(err)
	bobaccount, err := bob.TokenContract.Account()
	exitIfError(err)

	err = eve.SendTransaction(deposit1.Slot, coin.DepositBlockNum, 1, bobaccount.Address) //eve_to_bob
	exitIfError(err)

	authority.SubmitBlock()
	eve_to_bob_block, err := authority.GetBlockNumber()
	exitIfError(err)

	bob.WatchExits(deposit1.Slot)

	alicaccount, err := bob.TokenContract.Account()
	exitIfError(err)
	// Eve sends this same plasma coin to Alice
	err = eve.SendTransaction(deposit1.Slot, coin.DepositBlockNum, 1, alicaccount.Address) //eve_to_alice
	exitIfError(err)

	err = authority.SubmitBlock()
	exitIfError(err)

	eveToAliceBlock, err := authority.GetBlockNumber()
	exitIfError(err)

	// Alice attempts to exit here double-spent coin
	// Bob auto-challenges Alice's exit
	alice.StartExit(deposit1.Slot, coin.DepositBlockNum, eveToAliceBlock)

	// bob.challenge_between(deposit1.Slot, eve_to_bob_block)
	// Wait for challenge
	time.Sleep(2)
	bob.StartExit(deposit1.Slot, coin.DepositBlockNum, eve_to_bob_block)
	bob.StopWatchingExits(deposit1.Slot)

	ganache, err := client.ConnectToGanache("http://localhost:8545")
	exitIfError(err)
	_, err = ganache.IncreaseTime(context.TODO(), 8*24*3600)
	exitIfError(err)

	err = bob.Withdraw(deposit1.Slot)
	exitIfError(err)

	//TODO
	bobBalanceBefore := 0
	bobBalanceAfter := 0
	//bobBalanceBefore = w3.eth.getBalance(bob.token_contract.account.address)
	//bob.withdraw_bonds()
	//bobBalanceAfter = w3.eth.getBalance(bob.token_contract.account.address)
	if !(bobBalanceBefore < bobBalanceAfter) {
		log.Fatal("END: Bob did not withdraw his bonds")
	}

	bobTokensEnd, err := bob.TokenContract.BalanceOf()
	exitIfError(err)

	fmt.Printf("Bob has %d tokens", bobTokensEnd)
	if !(bobTokensEnd == bobTokensStart+1) {
		log.Fatal("END: Bob has incorrect number of tokens")
	}

	fmt.Printf("Plasma Cash 'challengeBetween' success :)\n")

}

// not idiomatic go, but it cleans up this sample
func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
