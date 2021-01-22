package utils

import (
	"github.com/olegfomenko/tpsloader/internal/config"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
)

func GetAccountDetail(keypair keypair.Full, client horizonclient.Client) horizon.Account {
	accountRequest := horizonclient.AccountRequest{AccountID: keypair.Address()}
	adminDetail, err := client.AccountDetail(accountRequest)

	if err != nil {
		panic(err)
	}

	return adminDetail
}

func SendTransaction(kp keypair.Full, operations []txnbuild.Operation, client horizonclient.Client) (horizon.Transaction, error) {
	signer := GetAccountDetail(kp, client)
	conf := config.GetConfig()

	// Creating transaction that holds create-operations-operation
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &signer,
		IncrementSequenceNum: true,
		Operations:           operations,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		BaseFee:              txnbuild.MinBaseFee,
	}

	tx, _ := txnbuild.NewTransaction(txParams)

	// Signing and encoding transaction
	signedTx, _ := tx.Sign(conf.Passphrase, &kp)

	// Encoding transaction
	txeBase64, _ := signedTx.Base64()
	// log.Println("Transaction base64: ", txeBase64)

	// Submitting transaction ans print response
	resp, err := client.SubmitTransactionXDR(txeBase64)
	return resp, err
}
