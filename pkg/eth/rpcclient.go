package eth

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
)

type RPCClient struct {
	clients []*ethclient.Client
	index   int
}

func NewRPCClient(urls []string) (*RPCClient, error) {
	clients := make([]*ethclient.Client, len(urls))
	for i, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			return nil, err
		}
		clients[i] = client
	}
	return &RPCClient{clients: clients, index: 0}, nil
}

func (c *RPCClient) CallRPCMethod(method func(*ethclient.Client) (interface{}, error)) (interface{}, error) {
	for i := 0; i < len(c.clients); i++ {
		result, err := method(c.clients[c.index])
		if err == nil {
			return result, nil
		}
		// Rotate to the next client if an error occurs.
		c.index = (c.index + 1) % len(c.clients)
	}
	return nil, errors.New("all RPC clients failed")
}
