package client

import (
	"fmt"
	"net/http"

	. "gopkg.in/check.v1"
)

type LoomChildChainSuite struct{}

var _ = Suite(&LoomChildChainSuite{})

func LoomBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "123")
}

func (s *LoomChildChainSuite) TestChildChainBlockNumber(c *C) {
}
