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

func BlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "123")
}

func NewRouter() (s *mux.Router) {
	s = mux.NewRouter()

	s.HandleFunc("/BlockNumber", BlockNumberHandler).Methods("Get")

	return s
}

func (s *ChildChainSuite) TestChildChainBlockNumber(c *C) {
	router := NewRouter()
	localServer := httptest.NewServer(router)
	defer localServer.Close()

	req, err := http.NewRequest("GET", "/BlockNumber", nil)
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
	expected := `123`
	c.Assert(rr.Body.String(), Equals, expected)
}
