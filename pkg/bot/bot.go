package bot

import (
	"fmt"
	"time"

	"github.com/0xrjman/web3-bot-go/config"
	"github.com/0xrjman/web3-bot-go/pkg/db"
	"github.com/0xrjman/web3-bot-go/pkg/eth"
)

type Bot struct {
	dbClient  *db.DBClient
	rpcClient *eth.RPCClient
	ethWallet *eth.EthKeystore
	gasPrice  int
}

func NewBot(cfg *config.Config) (*Bot, error) {
	dbClient, err := db.NewClient(cfg)
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return nil, err
	}
	fmt.Println("Database connection established")

	rpcClient, err := eth.NewRPCClient([]string{cfg.EthereumURL})
	if err != nil {
		fmt.Println("Failed to connect to the Ethereum node")
		return nil, err
	}
	fmt.Println("Ethereum node connection established")

	ethWallet := eth.LoadKeyStore(cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println("Ethereum wallet connection established")

	return &Bot{
		dbClient:  dbClient,
		rpcClient: rpcClient,
		ethWallet: ethWallet,
	}, nil
}

func (b *Bot) Run() error {
	fmt.Println("Bot configuration:")
	// fmt.Printf("Ethereum URL: %s\n", b.rpcClient.CallRPCMethod())
	fmt.Printf("Ethereum wallet address: %s\n", b.ethWallet.Key.Accounts())
	fmt.Printf("Gas price: %d\n", b.gasPrice)
	// TODO: implement the main loop of the bot
	fmt.Println("Bot is running...")

	for {
		// Get pending transactions from the database
		fmt.Println("Checking for pending transactions...")
		transactions, err := b.dbClient.GetPendingTransactions()
		if err != nil {
			return err
		}

		// TODO: Send transactions and update the database
		fmt.Printf("There are %d pending transactions\n", len(transactions))

		// Sleep for a while before the next loop
		time.Sleep(10 * time.Second)
	}
}
