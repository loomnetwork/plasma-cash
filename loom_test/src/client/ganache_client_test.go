package client

import (
	"context"

	. "gopkg.in/check.v1"
)

type GanacheClientTestSuite struct{}

var _ = Suite(&GanacheClientTestSuite{})

func (s *GanacheClientTestSuite) TestTimeIncrease(c *C) {
	c.Skip("Ganache must be running before this test will pass")

	ganache, err := ConnectToGanache("http://localhost:8545")
	if err != nil {
		c.Fatal(err)
	}

	timeAdj, err := ganache.IncreaseTime(context.TODO(), 60)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(timeAdj, Equals, uint32(60))

	timeAdj, err = ganache.IncreaseTime(context.TODO(), 30)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(timeAdj, Equals, uint32(90))
}
