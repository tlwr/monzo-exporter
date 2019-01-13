package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	MONZO_API_ENDPOINT = "https://api.monzo.com"
)

func GetUserIdentity(accessToken string) (MonzoCallerIdentity, error) {
	var callerID MonzoCallerIdentity

	client := &http.Client{}
	req, err := http.NewRequest("GET", MONZO_API_ENDPOINT+"/ping/whoami", nil)

	if err != nil {
		return callerID, err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	log.Print("Requesting: Monzo ping/whoami")
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Encountered error: Monzo /ping/whoami => %s", err)
		return callerID, err
	}

	err = json.Unmarshal(body, &callerID)
	if err != nil {
		return callerID, err
	}

	return callerID, nil
}

func ListAccounts(accessToken string) ([]MonzoAccount, error) {
	var accounts []MonzoAccount

	client := &http.Client{}
	req, err := http.NewRequest("GET", MONZO_API_ENDPOINT+"/accounts", nil)

	if err != nil {
		return accounts, err
	}

	log.Print("Requesting: Monzo /accounts")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Encountered error: Monzo /accounts request => %s", err)
		return accounts, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accounts, err
	}

	var accountsResp MonzoAPIListAccountsResponse

	err = json.Unmarshal(body, &accountsResp)
	if err != nil {
		return accounts, err
	}

	accounts = accountsResp.Accounts
	return accounts, nil
}

func ListPots(accessToken string) ([]MonzoPot, error) {
	var pots []MonzoPot

	client := &http.Client{}
	req, err := http.NewRequest("GET", MONZO_API_ENDPOINT+"/pots", nil)

	if err != nil {
		return pots, err
	}

	log.Print("Requesting: Monzo /pots")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Encountered error: Monzo /pots request => %s", err)
		return pots, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return pots, err
	}

	var potsResp MonzoAPIListPotsResponse

	err = json.Unmarshal(body, &potsResp)
	if err != nil {
		return pots, err
	}

	pots = potsResp.Pots
	return pots, nil
}

func GetBalance(accessToken string, accountID MonzoAccountID) (MonzoBalance, error) {
	var balance MonzoBalance

	url := fmt.Sprintf("%s/balance?account_id=%s", MONZO_API_ENDPOINT, accountID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return balance, err
	}

	log.Printf("Requesting: Monzo /balance/account_id=%s", accountID)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Encountered error: Monzo /pots request => %s", err)
		return balance, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return balance, err
	}

	err = json.Unmarshal(body, &balance)
	if err != nil {
		return balance, err
	}

	return balance, nil
}
