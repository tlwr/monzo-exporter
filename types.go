package main

import (
	"time"
)

type MonzoAccountID string
type MonzoCategory string
type MonzoClientID string
type MonzoCurrency string
type MonzoMerchantID string
type MonzoPotID string
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
