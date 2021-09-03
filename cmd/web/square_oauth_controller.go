package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/justinas/nosurf"
	"golang.org/x/oauth2"
)

type squareOauthController struct {
	oauthConfig *oauth2.Config
}

func NewSquareOauthController() *squareOauthController {
	clientId := os.Getenv("SQUARE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("SQUARE_OAUTH_CLIENT_SECRET")
	return &squareOauthController{
		oauthConfig: &oauth2.Config{
			RedirectURL:  "http://localhost:8000/oauth/square/callback",
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       []string{"CUSTOMERS_WRITE", "CUSTOMERS_READ", "GIFTCARDS_READ", "GIFTCARDS_WRITE"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  SQUARE_OAUTH_AUTH_URL,
				TokenURL: SQUARE_OAUTH_TOKEN_URL,
			},
		},
	}
}

func (controller *squareOauthController) getSquareOauthAuthURL(w http.ResponseWriter, r *http.Request) string {
	state := nosurf.Token(r)

	controller.setOauthStateCookie(w, state)

	return controller.oauthConfig.AuthCodeURL(state)
}

func (controller *squareOauthController) setSquareOauthToken(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie(OAUTH_STATE_COOKIE_KEY)

	if r.FormValue("state") != oauthState.Value {
		panic("invalid square oauth state")
	}

	oauthCode := r.FormValue("code")
	token, err := controller.oauthConfig.Exchange(context.Background(), oauthCode)
	if err != nil {
		panic(err.Error())
	}

	controller.setOauthTokenCookie(w, token)
}

func (controller *squareOauthController) getSquareOAuthToken(r *http.Request) *oauth2.Token {
	tokenCookie, err := r.Cookie(OAUTH_TOKEN_COOKIE_KEY)
	if err != nil || tokenCookie == nil {
		return nil
	}

	var token oauth2.Token
	data, err := base64.StdEncoding.DecodeString(tokenCookie.Value)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &token)
	if err != nil {
		return nil
	}

	return &token
}

func (controller *squareOauthController) setOauthStateCookie(w http.ResponseWriter, state string) {
	expire := time.Now().Add(time.Minute)

	cookie := http.Cookie{
		Name:     OAUTH_STATE_COOKIE_KEY,
		Value:    state,
		HttpOnly: true,
		Expires:  expire,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func (controller *squareOauthController) setOauthTokenCookie(w http.ResponseWriter, token *oauth2.Token) string {
	data, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	encodedToken := base64.StdEncoding.EncodeToString(data)
	setCookie(w, OAUTH_TOKEN_COOKIE_KEY, encodedToken)
	return encodedToken
}

func (controller *squareOauthController) getSquareHttpClient(r *http.Request) *http.Client {
	token := controller.getSquareOAuthToken(r)

	if token == nil {
		panic("Square Oauth token is not set.")
	}

	return controller.oauthConfig.Client(context.Background(), token)
}
