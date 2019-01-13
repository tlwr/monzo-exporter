package main

func CollectAllMetrics(accessTokens []string) {
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
