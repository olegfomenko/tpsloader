package pool

import (
	"github.com/olegfomenko/tpsloader/internal/operations"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"log"
	"time"
)

var (
	CreatedAccounts []keypair.Full
)

type AccountTask struct {
	Creator keypair.Full
	Client  horizonclient.Client
}

func (task *AccountTask) Run(ch chan struct{}) {
	for len(ch) == 0 {
		log.Println("Starting create operation in task")

		// Creating new transaction timestamp
		timestamp := TransactionTimestamp{
			Start:  time.Now(),
			Status: false,
		}

		// Making transaction
		kp, err := operations.CreateAccount(task.Creator, "100", task.Client)

		// Checking results
		if err != nil {
			log.Println("Task got an error:", err.(*horizonclient.Error))
			Failed++
		} else {
			// Updating timestamp
			timestamp.Finish = time.Now()
			timestamp.Status = true

			Timestamps = append(Timestamps, timestamp)
			CreatedAccounts = append(CreatedAccounts, kp)

			Successful++
		}
	}
}
