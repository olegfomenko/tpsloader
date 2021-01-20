package main

import (
	"github.com/olegfomenko/tpsloader/internal/account"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"log"
)

func main() {
	// Getting admin keypair using secret key
	admin, _ := keypair.ParseFull("SDR7XY33FYTDJTRF2CAXU5VGQWIQU4YOGZYPMYZ7ZAZTDGINQYMRJWZC")
	log.Println("Seed admin:", admin.Seed())
	log.Println("Address admin:", admin.Address())

	// Get horizon service API instance
	client := horizonclient.DefaultTestNetClient
	client.HorizonURL = "http://localhost:8000/"

	account.CreateAccount(admin, client)
}
