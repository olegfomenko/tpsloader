package pool

import (
	"context"
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
	Name   string
}

func (task *PaymentTask) Run(ctx context.Context, ready chan Task) {
	log.Println("Processing Payment Task:", task.Name)

	// Making transaction
	_, err := operations.SendPayment(task.From, task.To, task.Amount, task.Client)

	// Swapping accounts
	task.From, task.To = task.To, task.From

	if err == nil {
		// Updating information
		Timestamps = append(Timestamps, time.Now())
		Successful++
	} else {
		log.Println("Task", task.Name, "got an error:", err.(*horizonclient.Error), err.(*horizonclient.Error).Problem)
		Failed++
	}

	select {
	case <-ctx.Done():
		log.Println("Context finished. Task", task.Name, "will not be ready")
	case ready <- task:
		log.Println("Payment Task", task.Name, "is ready")
	}
}

func (task *PaymentTask) GetName() string {
	return task.Name
}
