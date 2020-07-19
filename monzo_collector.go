package main

import (
	"fmt"
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

		SetUserLatestCollect(identity.UserID)

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

		log.Printf(
			"CollectAccountMetrics: Getting transactions for user %s", identity.UserID,
		)

		transactions, err := GetTransactionsSinceDay(
			accessToken, account.ID, time.Now(),
		)

		if err != nil {
			log.Printf(
				"CollectAccountMetrics: Encountered error getting transactions for user %s => %s",
				identity.UserID, err,
			)
			return err
		}

		summaries := make(map[string]MonzoTransactionsSummary, 0)
		for _, transaction := range transactions {
			summaryKey := fmt.Sprintf(
				"%s/%s",
				transaction.Category, transaction.Description,
			)

			if _, ok := summaries[summaryKey]; !ok {
				summaries[summaryKey] = MonzoTransactionsSummary{
					Amount:      transaction.Amount,
					Category:    transaction.Category,
					Description: transaction.Description,
				}
			} else {
				summaries[summaryKey] = MonzoTransactionsSummary{
					Amount:      summaries[summaryKey].Amount + transaction.Amount,
					Category:    transaction.Category,
					Description: transaction.Description,
				}
			}
		}

		for _, summary := range summaries {
			SetTransactionsAmountToday(identity.UserID, account.ID, summary)
		}
	}

	log.Printf("CollectAccountMetrics: Done user %s", identity.UserID)
	return nil
}

func CollectPotMetrics(accessToken string, identity MonzoCallerIdentity) error {
	log.Printf("CollectPotMetrics: Starting user %s", identity.UserID)

	accounts, err := ListAccounts(accessToken)

	if err != nil {
		log.Printf(
			"CollectPotMetrics: Encountered error listing accounts for user %s => %s",
			identity.UserID, err,
		)
		return err
	}

	for _, account := range accounts {
		pots, err := ListPots(accessToken, account.ID)

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
	}

	return nil
}
