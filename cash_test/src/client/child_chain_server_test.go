package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	fmt.Fprintf(w, "theBlock")
}

func BlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	resp := ` {
        "jsonrpc": "2.0",
        "id": "",
        "result": {
            "response": {
            "last_block_height": 4955,
            "last_block_app_hash": "p/FQg0Muf7IUMeQAE25T1JZDwFA="
            }
        }
	}`

	type jResponse struct {
		Last_block_height   int
		Last_block_app_hash string
	}

	type jResult struct {
		Response jResponse
	}

	type jBlock struct {
		Jsonrpc string
		Id      string
		Result  jResult
	}

	var jblock jBlock

	err := json.Unmarshal([]byte(resp), &jblock)
	if err != nil {
		fmt.Print(err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strconv.Itoa(jblock.Result.Response.Last_block_height))
}

// ??? why no return blknum???
func BlockHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// blknum := vars["theblock"]
	blknum := `10`
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, blknum)
}

// not even sure if can extract those vars
func ProofHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blknum := vars["blknum"]
	uid := vars["uid"]
	fmt.Print(blknum)
	fmt.Print(uid)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "proof")
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

	req, err := http.NewRequest("GET", "/block", nil)
	if err != nil {
		c.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CurrentBlockHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	c.Assert(rr.Code, Equals, http.StatusOK)

	// Check the response body is what we expect.
	expected := `theBlock`
	c.Assert(rr.Body.String(), Equals, expected)

}

func (s *ChildChainSuite) TestBlockNumber(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	req, err := http.NewRequest("GET", "/blocknumber", nil)
	if err != nil {
		c.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(BlockNumberHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	c.Assert(rr.Code, Equals, http.StatusOK)

	// Check the response body is what we expect.
	expected := `4955`
	c.Assert(rr.Body.String(), Equals, expected)

}

func (s *ChildChainSuite) TestBlock(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	req, err := http.NewRequest("GET", "/block/10", nil)
	if err != nil {
		c.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(BlockHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	c.Assert(rr.Code, Equals, http.StatusOK)

	// Check the response body is what we expect.
	expected := `10`
	c.Assert(rr.Body.String(), Equals, expected)

}

func (s *ChildChainSuite) TestProof(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	var jsonStr = []byte(`{"blknum":"1234", "uid":"what_is_id?"}`)

	req, err := http.NewRequest("GET", "/proof", bytes.NewBuffer(jsonStr))
	if err != nil {
		c.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ProofHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	c.Assert(rr.Code, Equals, http.StatusOK)

	// Check the response body is what we expect.
	expected := `proof`
	c.Assert(rr.Body.String(), Equals, expected)

}

func (s *ChildChainSuite) TestSubmit(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	var jsonStr = []byte(`{"blknum":"1234"}`)

	req, err := http.NewRequest("POST", "/submit_block", bytes.NewBuffer(jsonStr))
	if err != nil {
		c.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SubmitBlockHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	c.Assert(rr.Code, Equals, http.StatusOK)

	// Check the response body is what we expect.
	expected := ``
	c.Assert(rr.Body.String(), Equals, expected)

}

func (s *ChildChainSuite) TestSendTransaction(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	var jsonStr = []byte(`{"tx":"tx"}`)

	req, err := http.NewRequest("POST", "/send_tx", bytes.NewBuffer(jsonStr))
	if err != nil {
		c.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SendTransactionHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	c.Assert(rr.Code, Equals, http.StatusOK)

	// Check the response body is what we expect.
	expected := `Transaction`
	c.Assert(rr.Body.String(), Equals, expected)

}
