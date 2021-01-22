package pool

import (
	"github.com/olegfomenko/tpsloader/internal/operations"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"log"
	"time"
)

type PaymentTask struct {
	From   keypair.Full
	To     keypair.Full
	Amount string
	Client horizonclient.Client
}

func (task *PaymentTask) Run(ch chan struct{}) {
	for len(ch) == 0 {
		log.Println("Starting payment operation in task")

		// Creating new transaction timestamp
		timestamp := TransactionTimestamp{
			Start:  time.Now(),
			Status: false,
		}

		// Making transaction
		_, err := operations.SendPayment(task.From, task.To, task.Amount, task.Client)
		task.From, task.To = task.To, task.From

		if err == nil {
			// Updating timestamp
			timestamp.Finish = time.Now()
			timestamp.Status = true
			Timestamps = append(Timestamps, timestamp)
		}
	}
}
