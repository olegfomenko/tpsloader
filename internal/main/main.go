package main

import (
	"context"
	"fmt"
	"github.com/olegfomenko/tpsloader/internal/config"
	"github.com/olegfomenko/tpsloader/internal/pool"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"log"
	"net/http"
	"time"
)

func main() {
	// Reading config file
	conf, err := config.NewConfig("config.yml")
	if err != nil {
		panic(err)
	}

	// Getting admin keypair using secret key
	admin, _ := keypair.ParseFull(conf.AdminSeed)
	log.Println("Seed admin:", admin.Seed())
	log.Println("Address admin:", admin.Address())

	/*// Creating new account
	account := operations.CreateAccount(admin, client)

	// Sending 10 lumens to created account
	operations.SendPayment(admin, account, "10", client)*/

	client := horizonclient.Client{
		HorizonURL: conf.HorizonURL,
		HTTP:       http.DefaultClient,
	}

	/*	paymentPool := pool.TaskPool{
		ThreadCount: 1,
		Task:
	}*/

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Minute)
	pl := pool.PoolImpl{}

	/*for _, creator := range creators {
		pl.Submit(ctx, &pool.AccountTask{
			Creator: creator,
			Client: client,
		})
	}*/

	for k, v := range conf.Payers {
		from, _ := keypair.ParseFull(k)
		to, _ := keypair.ParseFull(v)

		pl.Submit(ctx, &pool.PaymentTask{
			From:   *from,
			To:     *to,
			Amount: "10",
			Client: client,
		})
	}

	<-ctx.Done()

	for _, tx := range pool.Timestamps {
		fmt.Println(tx)
	}
}
