package operations

import (
	"github.com/olegfomenko/tpsloader/internal/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"log"
)

const maxOperations = 50

func CreateAccount(creator keypair.Full, amount string, client horizonclient.Client) (keypair.Full, error) {
	creation := generateKeypair()

	// Creating operation
	createAccountOp := txnbuild.CreateAccount{
		Destination: creation.Address(),
		Amount:      amount,
	}

	_, err := utils.SendTransaction(creator, []txnbuild.Operation{&createAccountOp}, client)
	if err != nil {
		log.Println("Create operation successful")
	}

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

func PrepareAccounts(count int, amount string, admin keypair.Full, client horizonclient.Client) []keypair.Full {
	var prepared []keypair.Full

	for count > 0 {
		var operations []txnbuild.Operation
		var currentPreparation []keypair.Full

		currentCnt := maxOperations
		if currentCnt > count {
			currentCnt = count
		}

		// Creating operation pack
		for len(operations) < currentCnt {
			creation := generateKeypair()
			currentPreparation = append(currentPreparation, creation)

			// Creating operation
			createAccountOp := txnbuild.CreateAccount{
				Destination: creation.Address(),
				Amount:      amount,
			}
			operations = append(operations, &createAccountOp)
		}

		_, err := utils.SendTransaction(admin, operations, client)

		if err != nil {
			log.Println("Error while preparing creator accounts:", err)
		} else {
			log.Println("Successful created", currentCnt, "accounts")
			// Adding created keypairs to result
			prepared = append(prepared, currentPreparation...)
			count -= currentCnt
		}
	}

	return prepared
}
