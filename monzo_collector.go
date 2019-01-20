package main

import (
	"log"
	"time"
)

type MonzoCollector struct {
	accessTokenGetter func() ([]string, error)
	duration          time.Duration
	stop              chan bool
}

func (m *MonzoCollector) Stop() {
	log.Println("Stopping MonzoCollector")
	m.stop <- true
}

func (m *MonzoCollector) Serve() {
	log.Println("Starting MonzoCollector")
	for {
		select {
		case <-m.stop:
			m.stop <- true
			return
		default:
			CollectAllMetrics(m.accessTokenGetter)
			time.Sleep(m.duration)
		}
	}
}

func CollectAllMetrics(
	accessTokensGetter func() ([]string, error),
) {

	accessTokens, err := accessTokensGetter()

	if err != nil {
		panic(err)
	}

	for _, token := range accessTokens {

		identity, err := GetUserIdentity(token)

		if err != nil {
			panic(err)
		}

		CollectAccountMetrics(token, identity)
		CollectPotMetrics(token, identity)
	}
}

func CollectAccountMetrics(accessToken string, identity MonzoCallerIdentity) {
	accounts, err := ListAccounts(accessToken)

	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		balance, err := GetBalance(accessToken, account.ID)

		if err != nil {
			panic(err)
		}

		SetCurrentBalance(identity.UserID, account.ID, balance.Balance)
		SetTotalBalance(identity.UserID, account.ID, balance.TotalBalance)
		SetSpendToday(identity.UserID, account.ID, balance.SpendToday)
	}
}

func CollectPotMetrics(accessToken string, identity MonzoCallerIdentity) {
	pots, err := ListPots(accessToken)

	if err != nil {
		panic(err)
	}

	for _, pot := range pots {
		SetPotBalance(identity.UserID, pot.ID, pot.Name, pot.Balance)
	}
}
