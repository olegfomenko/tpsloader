package main

import (
	"github.com/olegfomenko/tpsloader/internal/config"
	"github.com/olegfomenko/tpsloader/internal/operations"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"log"
)

func main() {
	// reading config file
	conf, err := config.NewConfig("config.yml")
	if err != nil {
		panic(err)
	}

	// Getting admin keypair using secret key
	admin, _ := keypair.ParseFull(conf.AdminSeed)
	log.Println("Seed admin:", admin.Seed())
	log.Println("Address admin:", admin.Address())

	// Get horizon service API instance
	client := horizonclient.DefaultTestNetClient
	client.HorizonURL = conf.HorizonURL

	// Creating new account
	account := operations.CreateAccount(admin, client)

	// Sending 10 lumens to created account
	operations.SendPayment(admin, account, "10", client)
}
