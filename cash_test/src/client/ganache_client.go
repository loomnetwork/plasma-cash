package client

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// GanacheClient defines typed wrappers for the Ethereum RPC API provided by Ganache.
type GanacheClient struct {
	*ethclient.Client
	rpcClient *rpc.Client
}

// ConnectToGanache connects a client to the given URL.
func ConnectToGanache(url string) (*GanacheClient, error) {
	rpcClient, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &GanacheClient{ethclient.NewClient(rpcClient), rpcClient}, nil
}

// IncreaseTime will advance the blockchain forward in time.
// The secs parameter specifies how many seconds to advance time by.
// Returns the total time adjustment in seconds.
func (c *GanacheClient) IncreaseTime(ctx context.Context, secs uint32) (uint32, error) {
	var result uint32
	if err := c.rpcClient.CallContext(ctx, &result, "evm_increaseTime", secs); err != nil {
		return 0, err
	}
	return result, nil
}
