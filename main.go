package main

import (
	"fmt"
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

	monzoOAuthClientID     = kingpin.Flag("monzo-oauth-client-id", "Monzo OAuth client id").Default("").OverrideDefaultFromEnvar("MONZO_OAUTH_CLIENT_ID").String()
	monzoOAuthClientSecret = kingpin.Flag("monzo-oauth-client-secret", "Monzo OAuth client secret").Default("").OverrideDefaultFromEnvar("MONZO_OAUTH_CLIENT_SECRET").String()
	monzoOAuthPort         = kingpin.Flag("monzo-oauth-port", "The port to bind to for serving OAuth").Default("8080").OverrideDefaultFromEnvar("MONZO_OAUTH_PORT").Int()
	monzoOAuthExternalURL  = kingpin.Flag("monzo-oauth-external-url", "The URL on which the exporter will be reachable").Default("").OverrideDefaultFromEnvar("MONZO_OAUTH_EXTERNAL_URL").String()

	monzoAccessTokens = kingpin.Flag("monzo-access-tokens", "Monzo access tokens comma separated").Default("").OverrideDefaultFromEnvar("MONZO_ACCESS_TOKENS").String()

	metricsScrapeInterval = kingpin.Flag("scrape-interval", "Time in seconds between scrapes").Default("30").OverrideDefaultFromEnvar("METRICS_SCRAPE_INTERVAL").Int64()
	metricsPort           = kingpin.Flag("metrics-port", "The port to bind to for serving metrics").Default("9036").OverrideDefaultFromEnvar("METRICS_PORT").Int()
)

func main() {
	kingpin.Parse()

	var getMonzoAccessTokens func() ([]string, error)

	if *monzoAccessTokens != "" {
		getMonzoAccessTokens = func() ([]string, error) {
			return strings.Split(*monzoAccessTokens, ","), nil
		}
	} else if *monzoOAuthClientID != "" && *monzoOAuthClientSecret != "" && *monzoOAuthExternalURL != "" {

		getMonzoAccessTokens = (&MonzoOAuthClient{
			MonzoOAuthClientID:     *monzoOAuthClientID,
			MonzoOAuthClientSecret: *monzoOAuthClientSecret,
			ExternalURL:            *monzoOAuthExternalURL,
		}).listen(*monzoOAuthPort)
	} else {
		fmt.Println("One of the following options is required:")
		fmt.Println("  - ONLY   --monzo-access-tokens")
		fmt.Println("  - ALL OF --monzo-oauth-client-id AND --monzo-oauth-client-secret AND --monzo-oauth-external-url")
		os.Exit(1)
	}

	duration := time.Duration(*metricsScrapeInterval) * time.Second

	RegisterCustomMetrics()

	supervisor := suture.NewSimple("MonzoCollector")
	supervisor.Add(&MonzoCollector{
		getMonzoAccessTokens,
		duration,
		make(chan bool),
	})
	defer supervisor.Stop()
	supervisor.ServeBackground()

	http.ListenAndServe(fmt.Sprintf(":%d", *metricsPort), promhttp.Handler())
}
