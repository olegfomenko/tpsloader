package operations

import (
	"github.com/olegfomenko/tpsloader/internal/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"log"
)

func CreateAccount(creator keypair.Full, amount string, client horizonclient.Client) (keypair.Full, error) {
	creation := generateKeypair()

	// Creating operation
	createAccountOp := txnbuild.CreateAccount{
		Destination: creation.Address(),
		Amount:      amount,
	}

	result, err := utils.SendTransaction(creator, []txnbuild.Operation{&createAccountOp}, client)
	log.Println("Create operations transaction result:", result)

	return creation, err
}

func generateKeypair() keypair.Full {
	kp, err := keypair.Random()
	if err != nil {
		panic(err)
	}

	log.Println("Seed son:", kp.Seed())
	log.Println("Address son:", kp.Address())

	return *kp
}

func PrepareCreators(count int, admin keypair.Full, client horizonclient.Client) []keypair.Full {
	var prepared []keypair.Full

	for len(prepared) != count {
		kp, err := CreateAccount(admin, "1000000", client)

		if err == nil {
			prepared = append(prepared, kp)
		} else {
			log.Println("Error while preparing creator accounts:", err)
		}
	}

	return prepared
}
