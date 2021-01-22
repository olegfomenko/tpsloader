package pool

import (
	"context"
	"github.com/olegfomenko/tpsloader/internal/operations"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"log"
	"time"
)

type AccountTask struct {
	Creator keypair.Full
	Client  horizonclient.Client
	Name    string
}

func (task *AccountTask) Run(ctx context.Context, ready chan Task) {
	log.Println("Processing CreateAccount Task:", task.Name)

	// Making transaction
	_, err := operations.CreateAccount(task.Creator, "100", task.Client)

	// Checking results
	if err != nil {
		log.Println("Task", task.Name, "got an error:", err.(*horizonclient.Error), err.(*horizonclient.Error).Problem)
		Failed++
	} else {
		// Updating information
		log.Println("Successful finishing CreateAccount Task:", task.Name)
		Timestamps = append(Timestamps, time.Now())
		Successful++
	}

	select {
	case <-ctx.Done():
		log.Println("Context finished. Task", task.Name, "will not be ready")
	case ready <- task:
		log.Println("CreateAccount Task", task.Name, "is ready")
	}
}

func (task *AccountTask) GetName() string {
	return task.Name
}
