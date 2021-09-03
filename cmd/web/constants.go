package main

const SQUARE_BASE_URL = "https://connect.squareupsandbox.com"
const SQUARE_OAUTH_AUTH_URL = SQUARE_BASE_URL + "/oauth2/authorize"
const SQUARE_OAUTH_TOKEN_URL = SQUARE_BASE_URL + "/oauth2/token"
const SQUARE_CUSTOMERS_URL = SQUARE_BASE_URL + "/v2/customers"
const SQUARE_GIFT_CARDS_URL = SQUARE_BASE_URL + "/v2/gift-cards"
const SQUARE_GIFT_CARD_ACTIVITIES_URL = SQUARE_GIFT_CARDS_URL + "/activities"
const SQUARE_LOCATION_ID = "LGW01TE9MP793"

const OAUTH_STATE_COOKIE_KEY = "oauthstate"
const OAUTH_TOKEN_COOKIE_KEY = "oauthtoken"
