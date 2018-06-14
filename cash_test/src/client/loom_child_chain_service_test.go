package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

type LoomChildChainSuite struct{}

var _ = Suite(&LoomChildChainSuite{})

func LoomBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
        "jsonrpc": "2.0",
        "id": "",
        "result": {
            "response": {
            "last_block_height": 4955,
            "last_block_app_hash": "p/FQg0Muf7IUMeQAE25T1JZDwFA="
            }
        }
    }`)
}

func NewLoomRouter() (s *mux.Router) {
	s = mux.NewRouter()

	s.HandleFunc("/abci_info", LoomBlockNumberHandler).Methods("Get")
	/*	s.HandleFunc("/blocknumber", BlockNumberHandler).Methods("Get")
		s.HandleFunc("/block/{theblock}", BlockHandler).Methods("Get")
		s.HandleFunc("/proof", ProofHandler).Methods("Get")
		s.HandleFunc("/submit_block", SubmitBlockHandler).Methods("Post")
		s.HandleFunc("/send_tx", SendTransactionHandler).Methods("Post")
	*/
	return s
}

func (s *LoomChildChainSuite) TestChildChainBlockNumber(c *C) {
	router := NewLoomRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	svc := NewLoomChildChainService(localServer.URL)
	blockNumber := svc.BlockNumber()

	// Check the response body is what we expect.
	c.Assert(blockNumber, Equals, 4955)

}
