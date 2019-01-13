package main

import (
	"encoding/json"
	"log"

	"github.com/h2non/gentleman"
)

const (
	MONZO_API_ENDPOINT = "https://api.monzo.com"
)

func MonzoClient(accessToken string) *gentleman.Request {
	client := gentleman.New()
	client.URL("https://api.monzo.com")
	request := client.Request()
	request.SetHeader("Authorization", "Bearer "+accessToken)
	return request
}

func GetUserIdentity(accessToken string) (MonzoCallerIdentity, error) {
	var callerID MonzoCallerIdentity

	req := MonzoClient(accessToken)
	req.Path("/ping/whoami")
	log.Print("Requesting: /ping/whoami")
	resp, err := req.Send()

	if err != nil {
		log.Printf("Encountered error: /ping/whoami => %s", err)
		return callerID, err
	}
	log.Print("Finished: /ping/whoami")

	err = json.Unmarshal(resp.Bytes(), &callerID)
	if err != nil {
		return callerID, err
	}

	return callerID, nil
}

func ListAccounts(accessToken string) ([]MonzoAccount, error) {
	var accounts []MonzoAccount

	req := MonzoClient(accessToken)
	req.Path("/accounts")
	log.Print("Requesting: /accounts")
	resp, err := req.Send()

	if err != nil {
		log.Printf("Encountered error: /accounts => %s", err)
		return accounts, err
	}
	log.Printf("Finished: /accounts")

	var accountsResp MonzoAPIListAccountsResponse

	err = json.Unmarshal(resp.Bytes(), &accountsResp)
	if err != nil {
		return accounts, err
	}

	accounts = accountsResp.Accounts
	return accounts, nil
}

func ListPots(accessToken string) ([]MonzoPot, error) {
	var pots []MonzoPot

	req := MonzoClient(accessToken)
	req.Path("/pots")
	log.Print("Requesting: /pots")
	resp, err := req.Send()

	if err != nil {
		log.Printf("Encountered error: /pots => %s", err)
		return pots, err
	}
	log.Print("Finished: /pots")

	var potsResp MonzoAPIListPotsResponse

	err = json.Unmarshal(resp.Bytes(), &potsResp)
	if err != nil {
		return pots, err
	}

	pots = potsResp.Pots
	return pots, nil
}

func GetBalance(accessToken string, accountID MonzoAccountID) (MonzoBalance, error) {
	var balance MonzoBalance

	req := MonzoClient(accessToken)
	req.Path("/balance")
	req.AddQuery("account_id", string(accountID))
	log.Printf("Requesting: /balance?account_id=%s", accountID)
	resp, err := req.Send()

	if err != nil {
		log.Printf("Encountered error: /pots => %s", err)
		return balance, err
	}
	log.Printf("Finished: /balance?account_id=%s", accountID)

	err = json.Unmarshal(resp.Bytes(), &balance)
	if err != nil {
		return balance, err
	}

	return balance, nil
}
