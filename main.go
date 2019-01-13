package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "0.0.1"

	monzoClientId     = kingpin.Flag("monzo-client-id", "Monzo client id").Default("").OverrideDefaultFromEnvar("MONZO_CLIENT_ID").String()
	monzoClientSecret = kingpin.Flag("monzo-client-secret", "Monzo client secret").Default("").OverrideDefaultFromEnvar("MONZO_CLIENT_SECRET").String()
	monzoAccessTokens = kingpin.Flag("monzo-access-tokens", "Monzo access tokens comma separated").Default("").OverrideDefaultFromEnvar("MONZO_ACCESS_TOKENS").String()

	metricsScrapeInterval = kingpin.Flag("scrape-interval", "Time in seconds between scrapes").Default("60").OverrideDefaultFromEnvar("METRICS_SCRAPE_INTERVAL").Int64()
	metricsPort           = kingpin.Flag("metrics-port", "The port to bind to for serving metrics").Default("8080").OverrideDefaultFromEnvar("METRICS_PORT").Int()
)

func main() {
	kingpin.Parse()

	if *monzoAccessTokens == "" {
		fmt.Println("-monzo-access-tokens is required")
		os.Exit(1)
	}

	monzoAccessTokensList := strings.Split(*monzoAccessTokens, ",")

	go func() {
		for true {
			for _, token := range monzoAccessTokensList {

				identity, err := GetUserIdentity(token)

				if err != nil {
					panic(err)
				}

				accounts, err := ListAccounts(token)

				if err != nil {
					panic(err)
				}

				for _, account := range accounts {
					balance, err := GetBalance(token, account.ID)

					if err != nil {
						panic(err)
					}

					SetCurrentBalance(identity.UserID, account.ID, balance.Balance)
					SetTotalBalance(identity.UserID, account.ID, balance.TotalBalance)
					SetSpendToday(identity.UserID, account.ID, balance.SpendToday)
				}

				pots, err := ListPots(token)

				if err != nil {
					panic(err)
				}

				for _, pot := range pots {
					SetPotBalance(identity.UserID, pot.ID, pot.Name, pot.Balance)
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()

	RegisterCustomMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9036", nil)
}
