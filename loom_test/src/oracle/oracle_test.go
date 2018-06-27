package oracle

import (
	"math/big"
	"testing"

	pctypes "github.com/loomnetwork/go-loom/builtin/types/plasma_cash"
	"github.com/pkg/errors"
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

type fakeEthPlasmaClient struct {
	curPlasmaBlockNum int64
	plasmaChain       []int64
	allowSubmit       bool
}

func (c *fakeEthPlasmaClient) Init() error {
	return nil
}

func (c *fakeEthPlasmaClient) CurrentPlasmaBlockNum() (*big.Int, error) {
	return big.NewInt(c.curPlasmaBlockNum), nil
}

func (c *fakeEthPlasmaClient) LatestEthBlockNum() (uint64, error) {
	return 0, nil
}

func (c *fakeEthPlasmaClient) SubmitPlasmaBlock(blockNum *big.Int, merkleRoot [32]byte) error {
	if !c.allowSubmit {
		return errors.New("EthPlasmaClient.SubmitPlasmaBlock shouldn't have been called")
	}
	c.plasmaChain = append(c.plasmaChain)
	return nil
}

func (c *fakeEthPlasmaClient) FetchDeposits(startBlock, endBlock uint64) ([]*pctypes.DepositRequest, error) {
	return nil, nil
}

type fakeDAppChainPlasmaClient struct {
	curPlasmaBlockNum int64
	plasmaChain       []int64
}

func (c *fakeDAppChainPlasmaClient) Init() error {
	return nil
}

func (c *fakeDAppChainPlasmaClient) CurrentPlasmaBlockNum() (*big.Int, error) {
	return big.NewInt(c.curPlasmaBlockNum), nil
}

func (c *fakeDAppChainPlasmaClient) PlasmaBlockAt(blockNum *big.Int) (*pctypes.PlasmaBlock, error) {
	bn := blockNum.Int64()
	for _, b := range c.plasmaChain {
		if b == bn {
			return &pctypes.PlasmaBlock{
				MerkleHash: make([]byte, 32, 32),
			}, nil
		}
	}
	return nil, errors.New("block not found")
}

func (c *fakeDAppChainPlasmaClient) FinalizeCurrentPlasmaBlock() error {
	return nil
}

func (c *fakeDAppChainPlasmaClient) Deposit(deposit *pctypes.DepositRequest) error {
	return nil
}

func createTestFakes() (*fakeEthPlasmaClient, *fakeDAppChainPlasmaClient, *PlasmaBlockWorker) {
	ethPlasmaClient := &fakeEthPlasmaClient{}
	dappPlasmaClient := &fakeDAppChainPlasmaClient{}
	return ethPlasmaClient, dappPlasmaClient,
		&PlasmaBlockWorker{
			ethPlasmaClient:     ethPlasmaClient,
			dappPlasmaClient:    dappPlasmaClient,
			plasmaBlockInterval: 1000,
		}
}
func TestSyncPlasmaBlocksWithEthereumWithNewChain(t *testing.T) {
	ethPlasmaClient, dappPlasmaClient, w := createTestFakes()
	w.Init()

	// No blocks should be sent to Ethereum plasma chain
	ethPlasmaClient.curPlasmaBlockNum = 0
	dappPlasmaClient.curPlasmaBlockNum = 0
	if err := w.syncPlasmaBlocksWithEthereum(); err != nil {
		t.Fatal(err)
	}
}

func TestSyncPlasmaBlocksWithEthereum(t *testing.T) {
	ethPlasmaClient, dappPlasmaClient, w := createTestFakes()
	w.Init()

	// TODO: setup ethPlasmaClient.plasmaChain & dappPlasmaClient.plasmaChain
	ethPlasmaClient.curPlasmaBlockNum = 0
	dappPlasmaClient.curPlasmaBlockNum = 0

	if err := w.syncPlasmaBlocksWithEthereum(); err != nil {
		t.Fatal(err)
	}
}
