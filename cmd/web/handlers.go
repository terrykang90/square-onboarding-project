package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// GET methods

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	td := &templateData{}
	app.renderTemplate(w, r, "home", td)
}

func (app *application) loginView(w http.ResponseWriter, r *http.Request) {
	td := &templateData{}
	app.renderTemplate(w, r, "login", td)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, OAUTH_TOKEN_COOKIE_KEY)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (app *application) listCustomersView(w http.ResponseWriter, r *http.Request) {
	sort_order := r.URL.Query().Get("sort")
	limit := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")

	params := url.Values{}
	if limit != "" {
		params.Add("limit", limit)
	}
	if sort_order != "" {
		params.Add("sort_order", sort_order)
	}
	if cursor != "" {
		params.Add("cursor", cursor)
	}

	sq_api_ctrl := NewSquareApiController(r)
	customers := sq_api_ctrl.getCustomers(&params)
	app.renderTemplate(w, r, "customers", &templateData{Customers: customers})
}

func (app *application) listGiftCardsView(w http.ResponseWriter, r *http.Request) {
	sq_api_ctrl := NewSquareApiController(r)
	td := &templateData{}

	customer_id := r.URL.Query().Get("customer_id")
	sort_order := r.URL.Query().Get("sort")
	limit := r.URL.Query().Get("limit")
	cursor := r.URL.Query().Get("cursor")
	state := r.URL.Query().Get("state")

	params := url.Values{}
	if customer_id != "" {
		params.Add("customer_id", customer_id)
		customer := sq_api_ctrl.getACustomer(customer_id)
		td.Customer = customer

	}
	if limit != "" {
		params.Add("limit", limit)
	}
	if sort_order != "" {
		params.Add("sort_order", sort_order)
	}
	if cursor != "" {
		params.Add("cursor", cursor)
	}
	if state != "" {
		params.Add("state", state)
	}

	gift_cards := sq_api_ctrl.getGiftCards(&params)
	td.GiftCards = gift_cards
	td.SelectedGiftCardState = state
	app.renderTemplate(w, r, "giftcards", td)
}

func (app *application) sendAGiftCardView(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "send_gift_card_form", &templateData{})
}

func (app *application) activateAGiftCardView(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "activate_gift_card_form", &templateData{})
}

func (app *application) deactivateAGiftCardView(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "deactivate_gift_card_form", &templateData{})
}

func (app *application) loadAGiftCardView(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "load_gift_card_form", &templateData{})
}

func (app *application) squareOauthLogin(w http.ResponseWriter, r *http.Request) {
	sq_oauth_ctrl := NewSquareOauthController()
	url := sq_oauth_ctrl.getSquareOauthAuthURL(w, r)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) squareOauthCallback(w http.ResponseWriter, r *http.Request) {
	sq_oauth_ctrl := NewSquareOauthController()
	sq_oauth_ctrl.setSquareOauthToken(w, r)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (app *application) showAGiftCardView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	if id == "" {
		app.notFound(w)
		return
	}

	sq_api_ctrl := NewSquareApiController(r)
	gift_card := sq_api_ctrl.getAGiftCard(id)

	params := url.Values{}
	params.Add("gift_card_id", id)
	gift_card_activities := sq_api_ctrl.getGiftCardActivities(&params)

	app.renderTemplate(w, r, "giftcard", &templateData{GiftCard: gift_card, GiftCardActivities: gift_card_activities})
}

// POST Methods
func (app *application) createAGiftCard(w http.ResponseWriter, r *http.Request) {
	idempotency_key := uuid.New()

	params := createGiftCardPayload{
		IdempotencyKey: idempotency_key.String(),
		LocationId:     SQUARE_LOCATION_ID,
		GiftCard:       &giftCardTypePayload{Type: "DIGITAL"},
	}

	sq_api_ctrl := NewSquareApiController(r)
	gift_card := sq_api_ctrl.createAGiftCard(params)

	http.Redirect(w, r, fmt.Sprintf("/giftcards/%s", gift_card.Id), http.StatusFound)
}

func (app *application) createAGiftCardActivity(w http.ResponseWriter, r *http.Request) {
	idempotency_key := uuid.New()
	gift_card_id := r.URL.Query().Get(":id")
	if gift_card_id == "" {
		app.notFound(w)
		return
	}

	r.ParseForm()
	activity_type := r.Form.Get("type")
	currency := r.Form.Get("currency")
	reason := r.Form.Get("reason")

	params := createGiftCardActivityPayload{
		IdempotencyKey: idempotency_key.String(),
		GiftCardActivity: &giftCardActivityPayload{
			Type:       activity_type,
			LocationId: SQUARE_LOCATION_ID,
			GiftCardId: gift_card_id,
		},
	}

	if activity_type == "ACTIVATE" {
		amount, err := strconv.Atoi(r.Form.Get("amount"))
		if err != nil {
			panic(err)
		}

		params.GiftCardActivity.ActivateActivityDetails = &loadActivityDetailsPayload{
			AmountMoney: &amountMoney{
				Amount:   amount,
				Currency: currency,
			},
			BuyerInstrumentIds: []string{"test"},
		}
	}

	if activity_type == "LOAD" {
		amount, err := strconv.Atoi(r.Form.Get("amount"))
		if err != nil {
			panic(err)
		}

		params.GiftCardActivity.LoadActivityDetails = &loadActivityDetailsPayload{
			AmountMoney: &amountMoney{
				Amount:   amount,
				Currency: currency,
			},
			BuyerInstrumentIds: []string{"test"},
		}
	}

	if activity_type == "DEACTIVATE" {
		params.GiftCardActivity.DeactivateActivityDetails = &deactivateActivityDetailsPayload{
			Reason: reason,
		}
	}

	sq_api_ctrl := NewSquareApiController(r)
	sq_api_ctrl.createAGiftCardActivity(params)

	http.Redirect(w, r, fmt.Sprintf("/giftcards/%s", gift_card_id), http.StatusFound)
}

func (app *application) sendAGiftCard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gift_card_id := r.URL.Query().Get(":id")
	customer_id := r.Form.Get("customer_id")
	app.infoLog.Printf("%s", gift_card_id)

	params := sendAGiftCardPayload{
		CustomerId: customer_id,
	}

	sq_api_ctrl := NewSquareApiController(r)
	sq_api_ctrl.sendAGiftCard(gift_card_id, &params)

	http.Redirect(w, r, fmt.Sprintf("/giftcards/%s", gift_card_id), http.StatusFound)
}

func (app *application) activateAGiftCard(w http.ResponseWriter, rs *http.Request) {
	r.ParseForm()
	r.Form.Set("type", "TEST")
	app.createAGiftCardActivity(w, r)
}

func (app *application) deactivateAGiftCard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.Form.Set("type", "TEST")
	app.createAGiftCardActivity(w, r)
}

func (app *application) loadAGiftCard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	r.Form.Set("type", "TEST")
	app.createAGiftCardActivity(w, r)
}