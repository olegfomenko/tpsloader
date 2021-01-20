package operations

import (
	"github.com/olegfomenko/tpsloader/internal/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"log"
)

func CreateAccount(creator *keypair.Full, client *horizonclient.Client) *keypair.Full {
	creation := generateKeypair()

	// Creating operation
	createAccountOp := txnbuild.CreateAccount{
		Destination: creation.Address(),
		Amount:      "100",
	}

	result := utils.SendTransaction(creator, []txnbuild.Operation{&createAccountOp}, client)
	log.Println("Create operations transaction result:", result)

	return creation
}

func generateKeypair() *keypair.Full {
	kp, err := keypair.Random()
	if err != nil {
		panic(err)
	}

	log.Println("Seed son:", kp.Seed())
	log.Println("Address son:", kp.Address())

	return kp
}
