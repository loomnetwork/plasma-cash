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

	dan := client.NewClient(svc, client.GetRootChain("dan"), client.GetTokenContract("dan"))
	mallory := client.NewClient(svc, client.GetRootChain("mallory"), client.GetTokenContract("mallory"))
	authority := client.NewClient(svc, client.GetRootChain("authority"),
		client.GetTokenContract("authority"))

	// Give Mallory 5 tokens
	mallory.TokenContract.Register()
	slots := []uint64{}

	danTokensStart, err := dan.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Dan has %v tokens", danTokensStart)
	if danTokensStart != 0 {
		log.Fatal("START: Dan has incorrect number of tokens")
	}
	malloryTokensStart, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %v tokens", malloryTokensStart)
	if malloryTokensStart != 5 {
		log.Fatal(fmt.Sprintf("START: Mallory has incorrect number of tokens -%d", malloryTokensStart))
	}
	currentBlock, err := authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %v", currentBlock)

	// Mallory deposits one of her coins to the plasma contract
	txHash := mallory.Deposit(6)

	depEvent, err := mallory.RootChain.DepositEventData(txHash)
	exitIfError(err)
	depositSlot1 := depEvent.Slot
	slots = append(slots, depEvent.Slot)

	txHash = mallory.Deposit(7)
	depEvent, err = mallory.RootChain.DepositEventData(txHash)
	exitIfError(err)
	//depositSlot2 := depEvent.Slot
	slots = append(slots, depEvent.Slot)

	// wait to make sure that events get fired correctly
	time.Sleep(2)

	malloryTokensPostDeposit, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %v tokens", malloryTokensPostDeposit)
	if malloryTokensPostDeposit != 3 {
		log.Fatal("POST-DEPOSIT: Mallory has incorrect number of tokens")
	}

	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %v", currentBlock)

	authority.SubmitBlock()
	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %v", currentBlock)
	block3000, err := authority.GetBlock(3000)
	exitIfError(err)
	log.Printf("block300-%v\n", block3000)

	authority.SubmitBlock()

	log.Printf("SubmitBlock\n")
	// Mallory sends her coin to Dan
	// Coin 6 was the first deposit of
	coin, err := mallory.RootChain.PlasmaCoin(depositSlot1)
	exitIfError(err)

	mallory.DebugCoinMetaData(slots)

	danaccount, err := dan.TokenContract.Account()
	exitIfError(err)
	log.Printf("account\n")

	err = mallory.SendTransaction(depositSlot1, coin.DepositBlockNum, 1, danaccount.Address) //mallory_to_dan
	exitIfError(err)
	authority.SubmitBlock()

	// Mallory attempts to exit spent coin (the one sent to Dan)
	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %v", currentBlock)

	log.Printf("Mallory trying an exit %d on block number %d\n", depositSlot1, coin.DepositBlockNum)
	mallory.StartExit(depositSlot1, 0, coin.DepositBlockNum)
	log.Printf("StartExit\n")

	// Dan"s transaction depositSlot1 included in block 5000. He challenges!
	dan.ChallengeAfter(depositSlot1, 5000)
	log.Printf("ChallengeAfter\n")
	dan.StartExit(depositSlot1, coin.DepositBlockNum, 5000)
	log.Printf("StartExit\n")

	// After 8 days pass,
	ganache, err := client.ConnectToGanache("http://localhost:8545")
	exitIfError(err)
	_, err = ganache.IncreaseTime(context.TODO(), 8*24*3600)
	exitIfError(err)
	log.Printf("increase time\n")

	authority.FinalizeExits()
	log.Printf("FinalizeExits\n")

	dan.Withdraw(depositSlot1)
	log.Printf("withdraw\n")

	//	accunt, err := dan.TokenContract.Account()
	//	exitIfError(err)
	//TODO web3
	/*
		danBalanceBefore := w3.eth.getBalance(account.address)
		dan.withdraw_bonds()
		dan_balance_after = w3.eth.getBalance(dan.TokenContract.account.address)
		if danBalanceBefore < dan_balance_after {
			log.Fatal("END: Dan did not withdraw his bonds")
		}
	*/

	malloryTokensEnd, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %v tokens", malloryTokensEnd)
	if malloryTokensEnd == 3 {
		log.Fatal("END: Mallory has incorrect number of tokens")
	}

	danTokensEnd, err := dan.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Dan has %v tokens", danTokensEnd)
	if danTokensEnd == 1 {
		log.Fatal("END: Dan has incorrect number of tokens")
	}

	log.Printf("Plasma Cash `challengeAfter` success :)")

}

// not idiomatic go, but it cleans up this sample
func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
