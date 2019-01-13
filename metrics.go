package main

import (
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

	potBalanceMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "monzo_pot_balance",
			Help: "Shows the individual pot balance",
		},
		[]string{"user_id", "pot_id", "pot_name"},
	)
)

func RegisterCustomMetrics() {
	prometheus.MustRegister(currentBalanceMetric)
	prometheus.MustRegister(totalBalanceMetric)
	prometheus.MustRegister(spendTodayMetric)
	prometheus.MustRegister(potBalanceMetric)
}

func SetCurrentBalance(
	userID MonzoUserID,
	accountID MonzoAccountID,
	balance int64,
) {
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
	potBalanceMetric.With(
		prometheus.Labels{
			"user_id":  string(userID),
			"pot_id":   string(potID),
			"pot_name": potName,
		},
	).Set(float64(balance))
}
