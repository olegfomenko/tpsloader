package utils

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

func GetAccountDetail(keypair *keypair.Full, client *horizonclient.Client) *horizon.Account {
	accountRequest := horizonclient.AccountRequest{AccountID: keypair.Address()}
	adminDetail, err := client.AccountDetail(accountRequest)

	if err != nil {
		panic(err)
	}

	return &adminDetail
}

func SendTransaction(kp *keypair.Full, operations []txnbuild.Operation, client *horizonclient.Client) *horizon.Transaction {
	signer := GetAccountDetail(kp, client)

	// Creating transaction that holds create-account-operation
	txParams := txnbuild.TransactionParams{
		SourceAccount:        signer,
		IncrementSequenceNum: true,
		Operations:           operations,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		BaseFee:              txnbuild.MinBaseFee,
	}

	tx, _ := txnbuild.NewTransaction(txParams)

	// Signing and encoding transaction
	signedTx, _ := tx.Sign("Stellar Load Test Network", kp)

	// Submitting transaction ans print response
	resp, err := client.SubmitTransaction(signedTx)

	if err != nil {
		panic(err)
	}

	return &resp
}
