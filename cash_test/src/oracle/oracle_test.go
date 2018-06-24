package oracle

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextPlasmaBlockNum(t *testing.T) {
	interval := big.NewInt(1000)

	res := nextPlasmaBlockNum(big.NewInt(9), interval)
	assert.Equal(t, res.Cmp(big.NewInt(1000)), 0)

	res = nextPlasmaBlockNum(big.NewInt(999), interval)
	assert.Equal(t, res.Cmp(big.NewInt(1000)), 0)

	res = nextPlasmaBlockNum(big.NewInt(0), interval)
	assert.Equal(t, res.Cmp(big.NewInt(1000)), 0)

	res = nextPlasmaBlockNum(big.NewInt(1000), interval)
	assert.Equal(t, res.Cmp(big.NewInt(2000)), 0)

	res = nextPlasmaBlockNum(big.NewInt(1001), interval)
	assert.Equal(t, res.Cmp(big.NewInt(2000)), 0)

	res = nextPlasmaBlockNum(big.NewInt(1999), interval)
	assert.Equal(t, res.Cmp(big.NewInt(2000)), 0)
}
