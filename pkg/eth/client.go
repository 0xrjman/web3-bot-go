package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/0xrjman/web3-bot-go/config"
)

type ETHClient struct {
	Client *ethclient.Client
}

// func NewETHClient(cfg *config.Config) (*ETHClient, error) {
// 	client, err := ethclient.Dial(cfg.EthereumURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ETHClient{Client: client}, nil
// }
