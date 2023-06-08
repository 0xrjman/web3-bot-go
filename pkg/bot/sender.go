package bot

import (
	"github.com/0xrjman/web3-bot-go/pkg/db"
	"github.com/0xrjman/web3-bot-go/pkg/eth"
)

type Sender struct {
	DBClient   *db.DBClient
	EthClient  *eth.RPCClient
	EthWallet  *eth.EthKeystore
	Pending    chan db.Transaction
	StopSignal chan bool
}

func NewSender(dbClient *db.DBClient, ethClient *eth.RPCClient, ethWallet *eth.EthKeystore) *Sender {
	return &Sender{
		DBClient:   dbClient,
		EthClient:  ethClient,
		EthWallet:  ethWallet,
		Pending:    make(chan db.Transaction),
		StopSignal: make(chan bool),
	}
}

func (s *Sender) Start() {
	for {
		select {
		case redemption := <-s.Pending:
			txHash, err := s.EthWallet.SendTransaction(s.EthClient, redemption.Address, redemption.Amount)
			if err != nil {
				// handle error, e.g., log it and continue
				continue
			}
			err = s.DBClient.UpdateTransactions(redemption.ID, txHash)
			if err != nil {
				// handle error, e.g., log it and continue
				continue
			}
		case <-s.StopSignal:
			return
		}
	}
}

func (s *Sender) Stop() {
	s.StopSignal <- true
}
