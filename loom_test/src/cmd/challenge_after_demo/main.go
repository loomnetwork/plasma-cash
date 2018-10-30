package main

import (
	"client"
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func main() {

	maxIteration := 30
	sleepPerIteration := 2000 * time.Millisecond

	var hostile bool
	flag.BoolVar(&hostile, "hostile", false, "run the demo with a hostile Plasma Cash operator")
	flag.Parse()

	if hostile {
		log.Println("Testing with a hostile Plasma Cash operator")
	}

	client.InitClients("http://localhost:8545")
	client.InitTokenClient("http://localhost:8545")
	ganache, err := client.ConnectToGanache("http://localhost:8545")
	exitIfError(err)

	svc, err := client.NewLoomChildChainService(hostile, "http://localhost:46658/rpc", "http://localhost:46658/query")
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
	if notEquals(danTokensStart, 0) {
		log.Fatal("START: Dan has incorrect number of tokens")
	}
	malloryTokensStart, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %v tokens", malloryTokensStart)
	if notEquals(malloryTokensStart, 5) {
		log.Fatal(fmt.Sprintf("START: Mallory has incorrect number of tokens -%d", malloryTokensStart))
	}
	currentBlock, err := authority.GetBlockNumber()
	exitIfError(err)

	_, err = ganache.HeaderByNumber(context.TODO(), nil)
	exitIfError(err)

	// Mallory deposits one of her coins to the plasma contract
	txHash := mallory.Deposit(big.NewInt(6))
	currentBlock, err = client.PollForBlockChange(authority, currentBlock, maxIteration, sleepPerIteration)
	if err != nil {
		panic(err)
	}

	depEvent, err := mallory.RootChain.DepositEventData(txHash)
	exitIfError(err)
	depositSlot1 := depEvent.Slot
	slots = append(slots, depEvent.Slot)

	txHash = mallory.Deposit(big.NewInt(7))

	currentBlock, err = client.PollForBlockChange(authority, currentBlock, maxIteration, sleepPerIteration)
	if err != nil {
		panic(err)
	}

	depEvent, err = mallory.RootChain.DepositEventData(txHash)
	exitIfError(err)
	slots = append(slots, depEvent.Slot)

	malloryTokensPostDeposit, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %v tokens", malloryTokensPostDeposit)
	if notEquals(malloryTokensPostDeposit, 3) {
		log.Fatal("POST-DEPOSIT: Mallory has incorrect number of tokens")
	}

	exitIfError(err)
	currentBlock, err = authority.GetBlockNumber()
	exitIfError(err)

	// Mallory sends her coin to Dan
	// Coin 6 was the first deposit of
	coin, err := mallory.RootChain.PlasmaCoin(depositSlot1)
	exitIfError(err)

	mallory.DebugCoinMetaData(slots)

	danAccount, err := dan.TokenContract.Account()
	exitIfError(err)

	err = mallory.SendTransaction(depositSlot1, coin.DepositBlockNum, big.NewInt(1), danAccount.Address) //mallory_to_dan
	exitIfError(err)

	currentBlock, err = client.PollForBlockChange(authority, currentBlock, maxIteration, sleepPerIteration)
	if err != nil {
		panic(err)
	}

	exitIfError(err)
	plasmaBlock3, err := authority.GetBlockNumber()
	exitIfError(err)

	// Mallory attempts to exit spent coin (the one sent to Dan)
	log.Printf("Mallory trying an exit %d on block number %d\n", depositSlot1, coin.DepositBlockNum)
	mallory.StartExit(depositSlot1, big.NewInt(0), coin.DepositBlockNum)

	// Dan's transaction depositSlot1 included in plasmaBlock3. He challenges!
	dan.ChallengeAfter(depositSlot1, plasmaBlock3)
	dan.StartExit(depositSlot1, coin.DepositBlockNum, plasmaBlock3)

	// After 8 days pass,
	_, err = ganache.IncreaseTime(context.TODO(), 8*24*3600)
	exitIfError(err)

	authority.FinalizeExits()

	dan.Withdraw(depositSlot1)

	danBalanceBefore, err := ganache.BalanceAt(context.TODO(), common.HexToAddress(danAccount.Address), nil)
	exitIfError(err)
	err = dan.WithdrawBonds()

	exitIfError(err)
	danBalanceAfter, err := ganache.BalanceAt(context.TODO(), common.HexToAddress(danAccount.Address), nil)
	exitIfError(err)

	if !(danBalanceBefore.Cmp(danBalanceAfter) < 0) {
		log.Fatal("END: Dan did not withdraw his bonds")
	}

	malloryTokensEnd, err := mallory.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Mallory has %v tokens", malloryTokensEnd)
	if notEquals(malloryTokensEnd, 3) {
		log.Fatal("END: Mallory has incorrect number of tokens")
	}

	danTokensEnd, err := dan.TokenContract.BalanceOf()
	exitIfError(err)
	log.Printf("Dan has %v tokens", danTokensEnd)
	if notEquals(danTokensEnd, 1) {
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

func notEquals(x *big.Int, y int64) bool {
	if x.Cmp(big.NewInt(y)) != 0 {
		return true
	} else {
		return false
	}
}
