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
	"time"
)

func main() {
	// Reading config file
	conf, err := config.NewConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running with config:", conf)

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
			Usage: "Generate N (-n flag) accounts with A (-a flag) lumens for config. Use as ./tpsloader gen-acc -n N -a A",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "n", Required: true},
				&cli.StringFlag{Name: "a", Required: true},
			},
			Action: func(c *cli.Context) {
				count := c.Int("n")
				amount := c.String("a")

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
				loadTesting(client, *conf)
			},
		},

		{
			Name:  "auto-run",
			Usage: "Run with generation of N (-n flag) creators with A1 (-a1 flag) amount and P (-p flag) payer pairs with A2 (-a2 flag) amount.",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "n", Required: true},
				&cli.StringFlag{Name: "a1", Required: true},
				&cli.IntFlag{Name: "p", Required: true},
				&cli.StringFlag{Name: "a2", Required: true},
			},
			Action: func(c *cli.Context) {
				// Getting N - creators count
				n := c.Int("n")

				// Getting A1 - creator's amount
				a1 := c.String("a1")

				// Getting P - payer pairs count
				p := c.Int("p")

				// Getting A2 - payer's amount
				a2 := c.String("a2")

				// Getting new creators and writing them to config
				creators := operations.PrepareAccounts(n, a1, *admin, client)
				conf.Creators = conf.Creators[:0]
				for _, creator := range creators {
					conf.Creators = append(conf.Creators, creator.Seed())
				}

				// Getting new payers and writing them to config
				payers := operations.PrepareAccounts(p*2, a2, *admin, client)
				conf.Payers = make(map[string]string)
				for i := 0; i < len(payers); i += 2 {
					conf.Payers[payers[i].Seed()] = payers[i+1].Seed()
				}

				// Start testing
				loadTesting(client, *conf)
			},
		},
	}

	app.Run(os.Args)
}

func loadTesting(client horizonclient.Client, conf config.Config) {
	creators := parseCreators(conf)

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
