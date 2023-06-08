package db

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/0xrjman/web3-bot-go/config"
)

type DBClient struct {
	Conn *pg.DB
}

func NewClient(cfg *config.Config) (*DBClient, error) {
	fmt.Println("Connecting to database...")
	fmt.Println("Database configuration:")
	fmt.Printf("URL: %s\n", cfg)
	fmt.Println(fmt.Sprint(cfg.DBHost, ":", cfg.DBPort))
	psqlDB := pg.Connect(&pg.Options{
		Addr:     fmt.Sprint(cfg.DBHost, ":", cfg.DBPort),
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	})

	err := createSchema(psqlDB)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
		return nil, err
	}

	return &DBClient{Conn: psqlDB}, nil
}

func createSchema(
	db *pg.DB,
) error {
	models := []interface{}{(*GasRedemptionUser)(nil)}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type Transaction struct {
	ID        int
	Address   string
	Amount    int
	Signature string
	Verified  bool
	Sent      bool
	TxHash    string
}

func (c *DBClient) GetPendingTransactions() ([]Transaction, error) {
	var users []GasRedemptionUser
	if err := c.Conn.Model(&users).Where("verified = TRUE AND sent = FALSE").Select(); err != nil {
		return nil, err
	}

	// map to Transaction
	var transactions []Transaction
	for _, user := range users {
		transactions = append(transactions, Transaction{
			ID:        user.Id,
			Address:   user.Address,
			Amount:    int(user.Amount),
			Signature: user.Signature,
			Verified:  user.SignatureVerified,
			Sent:      user.Sent,
			TxHash:    user.TransactionHash,
		})
	}
	return transactions, nil
}

func (c *DBClient) UpdateTransactions(
	id int,
	hash string,
) error {
	fmt.Printf("Update transaction %d with hash %s\n", id, hash)
	return nil
}
