package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ChildChainSuite struct{}

var _ = Suite(&ChildChainSuite{})

func CurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	resp := `f844c0b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`
	fmt.Fprintf(w, resp)
}

func BlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "4955")
}

func BlockHandler(w http.ResponseWriter, r *http.Request) {
	blknum := `f8a2f85df85b808001945194b63f10691e46635b27925100cfc0a5ceca62b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, blknum)
}

func ProofHandler(w http.ResponseWriter, r *http.Request) {
	proofstring := "AAAAAAAAAAKdBSjS0ahpIEal4CqiRzU5/MeqvWl59n9KfDhLFvI8AQ=="
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, proofstring)
}

func SubmitBlockHandler(w http.ResponseWriter, r *http.Request) {
	submitReturn := "0x0384d7e5062fc799fef4aaec1404d753f170e235fc8796a35c23732141de25b3"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, submitReturn)
}

func SendTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Transaction")
}

func NewRouter() (s *mux.Router) {
	s = mux.NewRouter()

	s.HandleFunc("/block", CurrentBlockHandler).Methods("Get")
	s.HandleFunc("/blocknumber", BlockNumberHandler).Methods("Get")
	s.HandleFunc("/block/{theblock}", BlockHandler).Methods("Get")
	s.HandleFunc("/proof", ProofHandler).Methods("Get")
	s.HandleFunc("/submit_block", SubmitBlockHandler).Methods("Post")
	s.HandleFunc("/send_tx", SendTransactionHandler).Methods("Post")

	return s
}

func (s *ChildChainSuite) TestCurrentBlock(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	svc := NewChildChainService(localServer.URL)
	_, currentBlock := svc.CurrentBlock()
	respString := "f844c0b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	c.Assert(currentBlock.blockId, Equals, respString)

}

func (s *ChildChainSuite) TestBlockNumber(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	svc := NewChildChainService(localServer.URL)
	blockNumber := svc.BlockNumber()

	// Check the response body is what we expect.
	c.Assert(blockNumber, Equals, 4955)

}

func (s *ChildChainSuite) TestBlock(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	svc := NewChildChainService(localServer.URL)
	blknum := 1
	_, block := svc.Block(blknum)
	respString := "f8a2f85df85b808001945194b63f10691e46635b27925100cfc0a5ceca62b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	c.Assert(block.blockId, Equals, respString)
}

func (s *ChildChainSuite) TestProof(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	svc := NewChildChainService(localServer.URL)
	blknum := 1000
	uid := 2
	_, proof := svc.Proof(blknum, uid)
	respString := "AAAAAAAAAAKdBSjS0ahpIEal4CqiRzU5/MeqvWl59n9KfDhLFvI8AQ=="
	c.Assert(proof.proofstring, Equals, respString)
}

func (s *ChildChainSuite) TestSubmit(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	svc := NewChildChainService(localServer.URL)
	block := Block{blockId: "f8a6f860f85e028203e8019450bce46ff7f6b92e4d383e4ada3ecba9e86d1292b842009f6f4a7be02d290320d0ba7370719f711edc2516908ef36728f18eaeea0bc90c2f806f64f7a914a6eb9933b4a9c8a47aafa597fdce883d7a41e8f2e18bb9e98f1bb8420036a21d66b53418e0b8e2e5f1c0d5ec3bc623620af3f43b95ef3fea199275833d23f25e335daddd5124682eae9591a99628ace8cbfae384e0ebeb5faf32d516bb1b"}
	err := svc.SubmitBlock(&block)
	fmt.Print(err)
	c.Assert(err, Equals, nil)
}

func (s *ChildChainSuite) TestSendTransaction(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

}
