package db

import (
	"time"
)

type GasRedemptionUser struct {
	// tableName         struct{} `pg:"gas_redemption_users"`
	Id                int     `pg:",pk"`
	Address           string  `pg:",notnull"`
	Amount            float64 `pg:",notnull"`
	Signature         string
	SignatureVerified bool
	TransactionHash   string
	Sent              bool
	Status            string

	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
}
