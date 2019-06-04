package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/h2non/gentleman"
	"github.com/h2non/gentleman/plugins/multipart"
)

const (
	STATE_COOKIE_NAME = "monzo_exporter_state"
	STATE_LENGTH      = 32

	START_PATH    = "/token/start"
	CALLBACK_PATH = "/token/callback"
)

func generateRandomState() string {
	randomBytes := make([]byte, STATE_LENGTH)
	_, err := rand.Read(randomBytes)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(randomBytes)
}

func (m *MonzoOAuthClient) redirectURL() string {
	return m.ExternalURL + CALLBACK_PATH
}

func (m *MonzoOAuthClient) handleJourneyStart(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()
	stateCookie := http.Cookie{Name: STATE_COOKIE_NAME, Value: state}

	http.SetCookie(w, &stateCookie)

	query := strings.Join([]string{
		"client_id=" + m.MonzoOAuthClientID,
		"redirect_uri=" + m.redirectURL(),
		"state=" + state,
		"response_type=" + "code",
	}, "&")
	monzoAuthURI := fmt.Sprintf("https://auth.monzo.com?%s", query)

	log.Printf("Redirecting user to %s\n", monzoAuthURI)
	http.Redirect(w, r, monzoAuthURI, 302)
}

func (m *MonzoOAuthClient) handleJourneyCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie(STATE_COOKIE_NAME)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - No %s cookie set", STATE_COOKIE_NAME)))
		return
	}

	requestStates, ok := r.URL.Query()["state"]

	if !ok || len(requestStates) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - state not retrievable")))
		return
	}

	requestState := requestStates[0]
	cookieState := stateCookie.Value

	if requestState != cookieState {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - cookie state and Monzo state differ")))
		return
	}

	requestCodes, ok := r.URL.Query()["code"]

	if !ok || len(requestCodes) != 1 || requestCodes[0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Monzo auth code not retrievable")))
		return
	}

	requestCode := requestCodes[0]

	fields := multipart.DataFields{
		"grant_type":    {"authorization_code"},
		"client_id":     {m.MonzoOAuthClientID},
		"client_secret": {m.MonzoOAuthClientSecret},
		"redirect_uri":  {m.redirectURL()},
		"code":          {requestCode},
	}

	authURL := "https://api.monzo.com/oauth2/token"

	client := gentleman.New()
	client.URL(authURL)
	client.Use(multipart.Fields(fields))

	log.Printf("Making POST request to %s\n", authURL)
	response, err := client.Request().Method("POST").Send()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(
			fmt.Sprintf("500 - Error making request to Monzo"),
		))
		return
	}

	log.Printf("Response to POST request to %s was %d\n", authURL, response.StatusCode)

	var authResponse MonzoAuthResponse
	err = json.Unmarshal(response.Bytes(), &authResponse)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - could not unmarshal JSON from Monzo"))
		return
	}

	expiryTime := time.Now().Add(
		time.Duration(authResponse.ExpirySeconds-300) * time.Second,
	)

	log.Println("Locking TokensBox")
	m.TokensBox.Lock.Lock()

	m.TokensBox.Tokens = append(
		m.TokensBox.Tokens,
		MonzoAccessAndRefreshTokens{
			AccessToken:  authResponse.AccessToken,
			RefreshToken: authResponse.RefreshToken,
			UserID:       authResponse.UserID,
			ExpiryTime:   expiryTime,
		},
	)
	log.Println("Appended to TokensBox")

	log.Println("Unlocking TokensBox")
	m.TokensBox.Lock.Unlock()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("201 - Tokens received and accepted"))
}

func (m *MonzoOAuthClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	log.Printf("GET %s\n", path)

	if path == START_PATH {
		m.handleJourneyStart(w, r)
		return
	}
	if path == CALLBACK_PATH {
		m.handleJourneyCallback(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not found"))
	log.Printf("Served 404 for %s\n", path)
}

func (m *MonzoOAuthClient) GetAccessTokens() ([]string, error) {
	log.Println("Getting access tokens")
	tokens := make([]string, 0)

	log.Println("Locking TokensBox")
	m.TokensBox.Lock.Lock()
	defer func() {
		log.Println("Unlocking TokensBox")
		m.TokensBox.Lock.Unlock()
	}()

	log.Printf("There are %d tokens in the box\n", len(m.TokensBox.Tokens))

	for _, accessAndRefreshToken := range m.TokensBox.Tokens {
		tokens = append(tokens, string(accessAndRefreshToken.AccessToken))
	}

	log.Println("Finished getting access tokens")
	return tokens, nil
}

func (m *MonzoOAuthClient) listen(port int) func() ([]string, error) {
	m.TokensBox = ConcurrentMonzoTokensBox{
		Lock:   sync.Mutex{},
		Tokens: make([]MonzoAccessAndRefreshTokens, 0),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m,
	}

	go server.ListenAndServe()
	return m.GetAccessTokens
}

func (m *MonzoOAuthClient) RefreshAToken() error {
	log.Println("Locking TokensBox")
	m.TokensBox.Lock.Lock()

	tokens := m.TokensBox.Tokens

	defer func() {
		log.Println("Unlocking TokensBox")
		m.TokensBox.Lock.Unlock()
	}()

	if len(tokens) == 0 {
		log.Println("No tokens to refresh. Done")
		return nil
	}

	headToken := m.TokensBox.Tokens[0]
	tailTokens := m.TokensBox.Tokens[1:]

	doWeNeedToRefresh := true // FIXME
	if doWeNeedToRefresh {
		log.Printf("Refreshing token for user %s", headToken.UserID)

		refreshedToken, err := RefreshToken(
			m.MonzoOAuthClientID, m.MonzoOAuthClientSecret,
			string(headToken.AccessToken), string(headToken.RefreshToken),
		)

		if err != nil {
			return fmt.Errorf(
				"Encountered error refreshing token for user %s => %s",
				headToken.UserID, err,
			)
		}

		headToken = refreshedToken
		log.Printf("Refreshed token for user %s", headToken.UserID)

		SetAccessTokenExpiry(headToken.UserID, headToken.ExpiryTime)
	}

	m.TokensBox.Tokens = append(tailTokens, headToken)
	return nil
}
