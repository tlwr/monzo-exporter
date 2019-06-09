package main

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	currentBalanceMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_current_balance",
			Help: "Shows the currently spendable account balance",
		},
		[]string{"user_id", "account_id"},
	)

	totalBalanceMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_total_balance",
			Help: "Shows the total account balance including pots",
		},
		[]string{"user_id", "account_id"},
	)

	spendTodayMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_spend_today",
			Help: "Shows the spend amount spent today",
		},
		[]string{"user_id", "account_id"},
	)

	transactionsAmountToday = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_transactions_amount_today",
			Help: "Shows the amount transacted today for a transaction description",
		},
		[]string{
			"user_id", "account_id",
			"description", "category",
		},
	)

	potBalanceMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_pot_balance",
			Help: "Shows the individual pot balance",
		},
		[]string{"user_id", "pot_id", "pot_name"},
	)

	userLatestCollectMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_user_latest_collect",
			Help: "Shows the unix timestamp expiry for most recent data collection",
		},
		[]string{"user_id"},
	)

	accessTokenExpiryMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_access_token_expiry",
			Help: "Shows the unix timestamp expiry for the access token",
		},
		[]string{"user_id"},
	)

	monzoAPIResponseCodeMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "monzo_api_response_code",
			Help: "Shows the response codes per endpoint from the Monzo API",
		},
		[]string{"response_code", "endpoint"},
	)
)

func RegisterCustomMetrics() {
	prometheus.MustRegister(currentBalanceMetric)
	prometheus.MustRegister(totalBalanceMetric)
	prometheus.MustRegister(spendTodayMetric)
	prometheus.MustRegister(transactionsAmountToday)
	prometheus.MustRegister(potBalanceMetric)
	prometheus.MustRegister(userLatestCollectMetric)
	prometheus.MustRegister(accessTokenExpiryMetric)
	prometheus.MustRegister(monzoAPIResponseCodeMetric)
}

func SetCurrentBalance(
	userID MonzoUserID,
	accountID MonzoAccountID,
	balance int64,
) {
	log.Printf(
		"Setting monzo_current_balance for user %s for account %s to %d",
		userID, accountID, balance,
	)

	currentBalanceMetric.With(
		prometheus.Labels{
			"user_id":    string(userID),
			"account_id": string(accountID),
		},
	).Set(float64(balance))
}

func SetTotalBalance(
	userID MonzoUserID,
	accountID MonzoAccountID,
	balance int64,
) {
	log.Printf(
		"Setting monzo_total_balance for user %s for account %s to %d",
		userID, accountID, balance,
	)

	totalBalanceMetric.With(
		prometheus.Labels{
			"user_id":    string(userID),
			"account_id": string(accountID),
		},
	).Set(float64(balance))
}

func SetSpendToday(
	userID MonzoUserID,
	accountID MonzoAccountID,
	spend int64,
) {
	log.Printf(
		"Setting monzo_spend_today for user %s for account %s to %d",
		userID, accountID, spend,
	)

	spendTodayMetric.With(
		prometheus.Labels{
			"user_id":    string(userID),
			"account_id": string(accountID),
		},
	).Set(float64(spend))
}

func SetPotBalance(
	userID MonzoUserID,
	potID MonzoPotID,
	potName string,
	balance int64,
) {
	log.Printf(
		"Setting monzo_pot_balance for user %s for pot %s to %d",
		userID, potID, balance,
	)

	potBalanceMetric.With(
		prometheus.Labels{
			"user_id":  string(userID),
			"pot_id":   string(potID),
			"pot_name": potName,
		},
	).Set(float64(balance))
}

func SetUserLatestCollect(userID MonzoUserID) {
	timestamp := time.Now().Unix()

	log.Printf(
		"Setting monzo_user_latest_collect for user %s to %d",
		userID, timestamp,
	)

	userLatestCollectMetric.With(
		prometheus.Labels{
			"user_id": string(userID),
		},
	).Set(float64(timestamp))
}

func SetAccessTokenExpiry(
	userID MonzoUserID,
	expiryTime time.Time,
) {
	log.Printf(
		"Setting monzo_access_token_expiry for user %s to %d",
		userID, expiryTime.Unix(),
	)

	accessTokenExpiryMetric.With(
		prometheus.Labels{
			"user_id": string(userID),
		},
	).Set(float64(expiryTime.Unix()))
}

func IncMonzoAPIResponseCode(
	endpoint string,
	responseCode int,
) {
	log.Printf(
		"Incrementing monzo_api_response_code %d for endpoint %s",
		responseCode, endpoint,
	)

	monzoAPIResponseCodeMetric.With(
		prometheus.Labels{
			"endpoint":      endpoint,
			"response_code": fmt.Sprintf("%d", responseCode),
		},
	).Inc()
}

func SetTransactionsAmountToday(
	userID MonzoUserID, accountID MonzoAccountID,
	transactionsSummary MonzoTransactionsSummary,
) {
	log.Printf(
		"Setting monzo_transactions_amount_today for User %s for Account %s for Description %s for Category %s to %d",
		userID, accountID,
		transactionsSummary.Description, transactionsSummary.Category,
		transactionsSummary.Amount,
	)

	transactionsAmountToday.With(
		prometheus.Labels{
			"user_id":     string(userID),
			"account_id":  string(accountID),
			"description": transactionsSummary.Description,
			"category":    transactionsSummary.Category,
		},
	).Set(float64(transactionsSummary.Amount))
}

func ResetTransactionsAmountToday() {
	log.Println("Resetting monzo_transactions_amount_today")
	transactionsAmountToday.Reset()
	log.Println("Reset monzo_transactions_amount_today")
}
