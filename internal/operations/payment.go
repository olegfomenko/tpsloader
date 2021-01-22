package operations

import (
	"github.com/olegfomenko/tpsloader/internal/utils"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"log"
)

func SendPayment(source keypair.Full, destination keypair.Full, amount string, client horizonclient.Client) (horizon.Transaction, error) {
	log.Println("Creating payment from:", source.Address(), "to:", destination.Address())

	paymentOperation := txnbuild.Payment{
		Destination: destination.Address(),
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	result, err := utils.SendTransaction(source, []txnbuild.Operation{&paymentOperation}, client)
	if err != nil {
		log.Println("Gon an error while payment operation:", err)
	} else {
		log.Println("Payment successful!")
	}

	return result, err
}
