package operations

import (
	"github.com/olegfomenko/tpsloader/internal/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"log"
)

func SendPayment(source *keypair.Full, destination *keypair.Full, amount string, client *horizonclient.Client) {
	log.Println("Creating payment from:", source.Address(), "to:", destination.Address())

	paymentOperation := txnbuild.Payment{
		Destination: destination.Address(),
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	result := utils.SendTransaction(source, []txnbuild.Operation{&paymentOperation}, client)
	log.Println("Payment result:", result)
}
