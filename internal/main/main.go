package main

import (
	"context"
	"fmt"
	"github.com/olegfomenko/tpsloader/internal/config"
	"github.com/olegfomenko/tpsloader/internal/operations"
	"github.com/olegfomenko/tpsloader/internal/pool"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*
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

	// Creating new account
	account := operations.CreateAccount(admin, client)

	// Sending 10 lumens to created account
	operations.SendPayment(admin, account, "10", client)

	client := horizonclient.Client{
		HorizonURL: conf.HorizonURL,
		HTTP:       http.DefaultClient,
	}

	/*	paymentPool := pool.TaskPool{
		ThreadCount: 1,
		Task:
	}

	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Minute)
	pl := pool.PoolImpl{}

	/*for _, creator := range creators {
		pl.Submit(ctx, &pool.AccountTask{
			Creator: creator,
			Client: client,
		})
	}

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
*/

func main() {
	// Reading config file
	conf, err := config.NewConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Run9ing with config:", conf)

	admin, err := keypair.ParseFull(conf.AdminSeed)
	if err != nil {
		log.Fatal(err)
	}

	client := horizonclient.Client{
		HorizonURL: conf.HorizonURL,
		HTTP:       http.DefaultClient,
	}

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:  "gen-acc",
			Usage: "Generate N accounts with A lumens for config. Use as ./tpsloader gen-acc N A",
			Action: func(c *cli.Context) {
				if c.NArg() < 2 {
					log.Fatal("No creators count woa provided")
				}

				count, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					log.Fatal(err)
				}

				amount := c.Args().Get(1)

				creators := operations.PrepareAccounts(count, amount, *admin, client)

				for _, creator := range creators {
					fmt.Println("Seed:", creator.Seed(), "Address:", creator.Address())
				}
			},
		},

		{
			Name:  "run",
			Usage: "Run load testing based on configured seeds",
			Action: func(c *cli.Context) {
				creators := parseCreators(*conf)
				log.Println("Creating context with timeout", time.Duration(conf.Duration)*time.Millisecond)
				ctx, _ := context.WithTimeout(context.TODO(), time.Duration(conf.Duration)*time.Millisecond)
				pl := pool.PoolImpl{}

				// Starting creators
				for _, creator := range creators {
					pl.Submit(ctx, &pool.AccountTask{
						Creator: creator,
						Client:  client,
					})
				}

				// Starting Payers
				for k, v := range conf.Payers {
					from, _ := keypair.ParseFull(k)
					to, _ := keypair.ParseFull(v)

					pl.Submit(ctx, &pool.PaymentTask{
						From:   *from,
						To:     *to,
						Amount: conf.Amount,
						Client: client,
					})
				}

				<-ctx.Done()

				for _, tx := range pool.Timestamps {
					fmt.Println(tx)
				}

				log.Println("MAC TPS IS", getMaxTPS(time.Minute))
			},
		},
	}

	app.Run(os.Args)
}

func parseCreators(conf config.Config) []keypair.Full {
	var creators []keypair.Full

	for _, c := range conf.Creators {
		kp, err := keypair.ParseFull(c)
		if err != nil {
			log.Fatal(err)
		}

		creators = append(creators, *kp)
	}

	return creators
}

func getMaxTPS(delta time.Duration) float64 {
	var maxTPS float64 = 0

	for l, r := 0, 0; l < len(pool.Timestamps); l++ {
		for ; r+1 < len(pool.Timestamps); r++ {
			if delta < pool.Timestamps[r+1].Finish.Sub(pool.Timestamps[l].Finish) {
				break
			}
		}

		log.Println("Checking", l, "...", r)
		var curTPS float64 = float64(time.Second) * float64(r-l+1) / float64(delta)
		log.Println("[", l, "...", r, "] TPS is", curTPS)

		if curTPS > maxTPS {
			maxTPS = curTPS
		}
	}

	return maxTPS
}
