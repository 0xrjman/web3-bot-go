package bot

import (
	"time"

	"github.com/0xrjman/web3-bot-go/pkg/db"
)

type Fetcher struct {
	DBClient   *db.DBClient
	Interval   time.Duration
	Pending    chan db.Transaction
	StopSignal chan bool
}

func NewFetcher(dbClient *db.DBClient, interval time.Duration) *Fetcher {
	return &Fetcher{
		DBClient:   dbClient,
		Interval:   interval,
		Pending:    make(chan db.Transaction),
		StopSignal: make(chan bool),
	}
}

func (f *Fetcher) Start() {
	ticker := time.NewTicker(f.Interval)

	for {
		select {
		case <-ticker.C:
			redemptions, err := f.DBClient.GetPendingTransactions()
			if err != nil {
				// handle error, e.g., log it
				continue
			}

			for _, redemption := range redemptions {
				f.Pending <- redemption
			}
		case <-f.StopSignal:
			ticker.Stop()
			return
		}
	}
}

func (f *Fetcher) Stop() {
	f.StopSignal <- true
}
