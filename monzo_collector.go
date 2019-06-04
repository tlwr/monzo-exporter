package main

import (
	"log"
	"time"
)

type MonzoCollector struct {
	usingAccessTokens func(func([]string) error) error
	duration          time.Duration
	stop              chan bool
}

func (m *MonzoCollector) Stop() {
	log.Println("Stop: Stopping MonzoCollector")
	m.stop <- true
}

func (m *MonzoCollector) Serve() {
	log.Println("Serve: Starting MonzoCollector")
	for {
		select {
		case <-m.stop:
			log.Println("Serve: Stopping")
			m.stop <- true
			log.Println("Serve: Stopped")
			return
		default:
			log.Println("Serve: Starting metric collection")
			err := m.usingAccessTokens(CollectAllMetrics)
			log.Println("Serve: Finished metric collection")

			if err != nil {
				log.Printf("Serve: Encountered error collecting metrics => %s", err)
			}

			log.Println("Serve: Sleeping")
			time.Sleep(m.duration)
		}
	}
}

func CollectAllMetrics(accessTokens []string) error {
	log.Printf("CollectAllMetrics: Starting for %d tokens", len(accessTokens))

	for i, token := range accessTokens {
		log.Printf("CollectAllMetrics: Doing token %d of %d",
			i+1, len(accessTokens),
		)

		identity, err := GetUserIdentity(token)
		if err != nil {
			return err
		}

		err = CollectAccountMetrics(token, identity)
		if err != nil {
			return err
		}

		err = CollectPotMetrics(token, identity)
		if err != nil {
			return err
		}

		log.Printf("CollectAllMetrics: Done for user => %s", identity.UserID)
	}

	log.Printf("CollectAllMetrics: Done %d tokens", len(accessTokens))
	return nil
}

func CollectAccountMetrics(accessToken string, identity MonzoCallerIdentity) error {
	log.Printf("CollectAccountMetrics: Starting user %s", identity.UserID)

	accounts, err := ListAccounts(accessToken)

	if err != nil {
		log.Printf(
			"CollectAccountMetrics: Encountered error listing accounts for user %s => %s",
			identity.UserID, err,
		)
		return err
	}

	for _, account := range accounts {
		log.Printf(
			"CollectAccountMetrics: Getting balance for user %s", identity.UserID,
		)

		balance, err := GetBalance(accessToken, account.ID)

		if err != nil {
			log.Printf(
				"CollectAccountMetrics: Encountered error getting balance for user %s => %s",
				identity.UserID, err,
			)
			return err
		}

		SetCurrentBalance(identity.UserID, account.ID, balance.Balance)
		SetTotalBalance(identity.UserID, account.ID, balance.TotalBalance)
		SetSpendToday(identity.UserID, account.ID, balance.SpendToday)
	}

	log.Printf("CollectAccountMetrics: Done user %s", identity.UserID)
	return nil
}

func CollectPotMetrics(accessToken string, identity MonzoCallerIdentity) error {
	log.Printf("CollectPotMetrics: Starting user %s", identity.UserID)
	pots, err := ListPots(accessToken)

	if err != nil {
		log.Printf(
			"CollectPotMetrics: Encountered error listing pots for user %s => %s",
			identity.UserID, err,
		)
		return err
	}

	for _, pot := range pots {
		SetPotBalance(identity.UserID, pot.ID, pot.Name, pot.Balance)
	}

	log.Printf("CollectPotMetrics: Done user %s", identity.UserID)

	return nil
}
