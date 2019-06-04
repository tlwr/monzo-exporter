package main

import (
	"encoding/json"
	"log"
	"time"

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
	log.Print("GetUserIdentity: Requesting: /ping/whoami")
	resp, err := req.Send()

	if err != nil {
		log.Printf("GetUserIdentity: Encountered error: /ping/whoami => %s", err)
		return callerID, err
	}
	log.Println("GetUserIdentity: Finished: /ping/whoami")

	err = json.Unmarshal(resp.Bytes(), &callerID)
	if err != nil {
		log.Printf("GetUserIdentity: Encountered error unmarshalling => %s", err)
		return callerID, err
	}

	return callerID, nil
}

func ListAccounts(accessToken string) ([]MonzoAccount, error) {
	var accounts []MonzoAccount

	req := MonzoClient(accessToken)
	req.Path("/accounts")
	log.Print("ListAccounts: Requesting: /accounts")
	resp, err := req.Send()

	if err != nil {
		log.Printf("ListAccounts: Encountered error: /accounts => %s", err)
		return accounts, err
	}
	log.Printf("ListAccounts: Finished: /accounts")

	var accountsResp MonzoAPIListAccountsResponse

	err = json.Unmarshal(resp.Bytes(), &accountsResp)
	if err != nil {
		log.Printf("ListAccounts: Encountered error unmarshalling => %s", err)
		return accounts, err
	}

	accounts = accountsResp.Accounts
	log.Printf("ListAccounts: Done")
	return accounts, nil
}

func ListPots(accessToken string) ([]MonzoPot, error) {
	var pots []MonzoPot

	req := MonzoClient(accessToken)
	req.Path("/pots")
	log.Print("ListPots: Requesting: /pots")
	resp, err := req.Send()

	if err != nil {
		log.Printf("ListPots: Encountered error: /pots => %s", err)
		return pots, err
	}
	log.Print("ListPots Finished: /pots")

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
	log.Printf("GetBalance: Requesting: /balance?account_id=%s", accountID)
	resp, err := req.Send()

	if err != nil {
		log.Printf("GetBalance: Encountered error: /pots => %s", err)
		return balance, err
	}
	log.Printf("GetBalance: Finished: /balance?account_id=%s", accountID)

	err = json.Unmarshal(resp.Bytes(), &balance)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func RefreshToken(clientId string, clientSecret string, accessToken string, refreshToken string) (MonzoAccessAndRefreshTokens, error) {
	var returnTokens MonzoAccessAndRefreshTokens

	req := MonzoClient(accessToken)
	req.Path("/oauth2/token?grant_type=refresh_token")
	req.AddQuery("grant_type", "refresh_token")
	req.AddQuery("client_id", clientId)
	req.AddQuery("client_secret", clientSecret)
	log.Printf("RefreshToken: Requesting: /oauth2/token")
	resp, err := req.Send()

	if err != nil {
		log.Printf(
			"RefreshToken: Encountered error: /oauth2/token?grant_type=refresh_token => %s",
			err,
		)
		return returnTokens, err
	}

	var authResponse MonzoAuthResponse
	err = json.Unmarshal(resp.Bytes(), &authResponse)

	if err != nil {
		log.Printf(
			"RefreshToken: Encountered error unmarshalling refresh token response => %s",
			err,
		)
		return returnTokens, err
	}

	expiryTime := time.Now().Add(
		time.Duration(authResponse.ExpirySeconds-300) * time.Second,
	)

	log.Printf(
		"RefreshTokeN: Refreshed access token for %s", authResponse.UserID,
	)
	return MonzoAccessAndRefreshTokens{
		AccessToken:  authResponse.AccessToken,
		RefreshToken: authResponse.RefreshToken,
		UserID:       authResponse.UserID,
		ExpiryTime:   expiryTime,
	}, nil
}
