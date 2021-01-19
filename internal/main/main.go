package main

import (
	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"log"
)

func main() {
	admin, err := keypair.ParseFull("SDR7XY33FYTDJTRF2CAXU5VGQWIQU4YOGZYPMYZ7ZAZTDGINQYMRJWZC")
	log.Println("Seed admin:", admin.Seed())
	log.Println("Address admin:", admin.Address())

	son, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Seed son:", son.Seed())
	log.Println("Address son:", son.Address())


	client := horizon.DefaultTestNetClient
	client.HorizonURL = "http://localhost:8000/"

	accountRequest := horizon.AccountRequest{AccountID: admin.Address()}
	adminDetail, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	createAccountOp := txnbuild.CreateAccount{
		Destination: son.Address(),
		Amount:      "100",
	}

	txParams := txnbuild.TransactionParams{
		SourceAccount:        &adminDetail,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&createAccountOp},
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		BaseFee:              txnbuild.MinBaseFee,
	}

	tx, _ := txnbuild.NewTransaction(txParams)

	signedTx, _ := tx.Sign("Stellar Load Test Network", admin)
	txeBase64, _ := signedTx.Base64()
	log.Println("Transaction base64: ", txeBase64)

	resp, err := client.SubmitTransactionXDR(txeBase64)

	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:", hError.Problem)
	}

	log.Println("\nTransaction response: ", resp)
}