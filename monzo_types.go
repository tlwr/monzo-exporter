package main

import (
	"sync"
	"time"
)

type MonzoAccessToken string
type MonzoAccountID string
type MonzoClientID string
type MonzoCurrency string
type MonzoMerchantID string
type MonzoPotID string
type MonzoRefreshToken string
type MonzoTransactionID string
type MonzoUserID string

type MonzoAccount struct {
	ID          MonzoAccountID `json:"id"`
	Description string         `json:"description"`
	Created     time.Time      `json:"created"`
}

type MonzoPot struct {
	ID MonzoPotID `json:"id"`

	Name string `json:"name"`

	Currency MonzoCurrency `json:"currency"`
	Balance  int64         `json:"balance"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type MonzoAPIListAccountsResponse struct {
	Accounts []MonzoAccount `json:"accounts"`
}

type MonzoAPIListPotsResponse struct {
	Pots []MonzoPot `json:"pots"`
}

type MonzoBalance struct {
	Balance      int64         `json:"balance"`
	TotalBalance int64         `json:"total_balance"`
	Currency     MonzoCurrency `json:"currency"`
	SpendToday   int64         `json:"spend_today"`
}

type MonzoCallerIdentity struct {
	Authenticated bool          `json:"authenticated"`
	ClientID      MonzoClientID `json:"client_id"`
	UserID        MonzoUserID   `json:"user_id"`
}

type MonzoAuthResponse struct {
	AccessToken   MonzoAccessToken  `json:"access_token"`
	RefreshToken  MonzoRefreshToken `json:"refresh_token"`
	UserID        MonzoUserID       `json:"user_id"`
	ExpirySeconds float64           `json:"expires_in"`
}

type MonzoAccessAndRefreshTokens struct {
	AccessToken  MonzoAccessToken
	RefreshToken MonzoRefreshToken
	UserID       MonzoUserID
	ExpiryTime   time.Time
}

type ConcurrentMonzoTokensBox struct {
	Lock   sync.Mutex
	Tokens []MonzoAccessAndRefreshTokens
}

type MonzoOAuthClient struct {
	MonzoOAuthClientID     string
	MonzoOAuthClientSecret string
	ExternalURL            string

	TokensBox ConcurrentMonzoTokensBox
}
