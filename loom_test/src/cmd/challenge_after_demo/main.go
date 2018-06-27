package main

import (
	"client"
	"log"
	"os"
	"time"
)

func main() {

	plasmaChain := os.Getenv("PLASMA_CHAIN")
	client.InitClients("http://localhost:8545")
	client.InitTokenClient("http://localhost:8545")

	var svc client.ChainServiceClient
	if plasmaChain == "LOOM" {
		svc = client.NewLoomChildChainService("http://localhost:46658/rpc", "http://localhost:46658/query")
	} else {
		svc = client.NewChildChainService("http://localhost:8546")
	}

	dan := client.NewClient(svc, client.GetRootChain("dan"), client.GetTokenContract("dan"))
	mallory := client.NewClient(svc, client.GetRootChain("mallory"), client.GetTokenContract("mallory"))
	authority := client.NewClient(svc, client.GetRootChain("authority"),
		client.GetTokenContract("authority"))

	// Give Mallory 5 tokens
	mallory.TokenContract.Register()

	danTokensStart, err := dan.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Dan has %s tokens", danTokensStart)
	if danTokensStart == 0 {
		log.Fatal("START: Dan has incorrect number of tokens")
	}
	malloryTokensStart, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %s tokens", malloryTokensStart)
	if malloryTokensStart == 5 {
		log.Fatal("START: Mallory has incorrect number of tokens")
	}
	currentBlock, err := authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %s", currentBlock)

	// Mallory deposits one of her coins to the plasma contract
	mallory.Deposit(6)
	mallory.Deposit(7)
	// wait to make sure that events get fired correctly
	time.Sleep(2)

	malloryTokensPostDeposit, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %s tokens", malloryTokensPostDeposit)
	if malloryTokensPostDeposit == 3 {
		log.Fatal("POST-DEPOSIT: Mallory has incorrect number of tokens")
	}

	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %s", currentBlock)

	authority.SubmitBlock()
	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %s", currentBlock)
	block3000, err := authority.GetBlock(3000)
	exitIfError(err)
	log.Printf("block300-%v\n", block3000)

	authority.SubmitBlock()

	// Mallory sends her coin to Dan
	// Coin 6 was the first deposit of
	utxoID := uint64(3)
	coin, err := mallory.RootChain.PlasmaCoin(utxoID)
	exitIfError(err)
	danaccount, err := dan.TokenContract.Account()
	exitIfError(err)

	_, err = mallory.SendTransaction(utxoID, coin.DepositBlockNum, 1, danaccount.Address) //mallory_to_dan
	authority.SubmitBlock()

	// Mallory attempts to exit spent coin (the one sent to Dan)
	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)
	log.Printf("current block: %s", currentBlock)

	mallory.StartExit(utxoID, 0, coin.DepositBlockNum)

	// Dan"s transaction was included in block 5000. He challenges!
	dan.ChallengeAfter(utxoID, 5000)
	dan.StartExit(utxoID, coin.DepositBlockNum, 5000)

	//TODO hook in web3
	//w3 = dan.root_chain.w3 // get w3 instance
	//increaseTime(w3, 8*24*3600)

	authority.FinalizeExits()

	dan.Withdraw(utxoID)

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
	log.Printf("Mallory has %s tokens", malloryTokensEnd)
	if malloryTokensEnd == 3 {
		log.Fatal("END: Mallory has incorrect number of tokens")
	}

	danTokensEnd, err := dan.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Dan has %s tokens", danTokensEnd)
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
