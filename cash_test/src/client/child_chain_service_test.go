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

// same here
func SubmitBlockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blknum := vars["blknum"]
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, blknum)
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
}

func (s *ChildChainSuite) TestSendTransaction(c *C) {

}
