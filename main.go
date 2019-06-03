package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thejerf/suture"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "0.0.1"

	monzoOAuthClientID        = kingpin.Flag("monzo-oauth-client-id", "Monzo OAuth client id").Default("").OverrideDefaultFromEnvar("MONZO_OAUTH_CLIENT_ID").String()
	monzoOAuthClientSecret    = kingpin.Flag("monzo-oauth-client-secret", "Monzo OAuth client secret").Default("").OverrideDefaultFromEnvar("MONZO_OAUTH_CLIENT_SECRET").String()
	monzoOAuthPort            = kingpin.Flag("monzo-oauth-port", "The port to bind to for serving OAuth").Default("8080").OverrideDefaultFromEnvar("MONZO_OAUTH_PORT").Int()
	monzoOAuthExternalURL     = kingpin.Flag("monzo-oauth-external-url", "The URL on which the exporter will be reachable").Default("").OverrideDefaultFromEnvar("MONZO_OAUTH_EXTERNAL_URL").String()
	monzoOAuthRefreshInterval = kingpin.Flag("monzo-oauth-refresh-interval", "Time in seconds between OAuth token refreshes").Default("10").OverrideDefaultFromEnvar("MONZO_OAUTH_REFRESH_INTERVAL").Int64()

	monzoAccessTokens = kingpin.Flag("monzo-access-tokens", "Monzo access tokens comma separated").Default("").OverrideDefaultFromEnvar("MONZO_ACCESS_TOKENS").String()

	metricsScrapeInterval = kingpin.Flag("scrape-interval", "Time in seconds between scrapes").Default("30").OverrideDefaultFromEnvar("METRICS_SCRAPE_INTERVAL").Int64()
	metricsPort           = kingpin.Flag("metrics-port", "The port to bind to for serving metrics").Default("9036").OverrideDefaultFromEnvar("METRICS_PORT").Int()
)

func main() {
	kingpin.Parse()

	var getMonzoAccessTokens func() ([]string, error)
	var monzoOAuthClient MonzoOAuthClient

	if *monzoAccessTokens != "" {
		getMonzoAccessTokens = func() ([]string, error) {
			return strings.Split(*monzoAccessTokens, ","), nil
		}
	} else if *monzoOAuthClientID != "" && *monzoOAuthClientSecret != "" && *monzoOAuthExternalURL != "" {

		monzoOAuthClient.MonzoOAuthClientID = *monzoOAuthClientID
		monzoOAuthClient.MonzoOAuthClientSecret = *monzoOAuthClientSecret
		monzoOAuthClient.ExternalURL = *monzoOAuthExternalURL

		getMonzoAccessTokens = monzoOAuthClient.listen(*monzoOAuthPort)
	} else {
		fmt.Println("One of the following options is required:")
		fmt.Println("  - ONLY   --monzo-access-tokens")
		fmt.Println("  - ALL OF --monzo-oauth-client-id AND --monzo-oauth-client-secret AND --monzo-oauth-external-url")
		os.Exit(1)
	}

	RegisterCustomMetrics()

	supervisor := suture.NewSimple("MonzoCollector")
	supervisor.Add(&MonzoCollector{
		getMonzoAccessTokens,
		time.Duration(*metricsScrapeInterval) * time.Second,
		make(chan bool),
	})
	defer supervisor.Stop()
	supervisor.ServeBackground()

	tickerOAuthInterval := time.NewTicker(
		time.Duration(*monzoOAuthRefreshInterval) * time.Second,
	)

	if *monzoAccessTokens == "" {
		go func() {
			for _ = range tickerOAuthInterval.C {
				log.Println("Refreshing OAuth tokens")
			}
		}()
	} else {
		log.Println("Skipping starting OAuth token refresher")
	}

	http.ListenAndServe(fmt.Sprintf(":%d", *metricsPort), promhttp.Handler())
}
