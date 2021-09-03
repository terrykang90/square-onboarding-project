package main

type customer struct {
	Id           string `json:"id"`
	GivenName    string `json:"given_name"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
	CreatedAt    string `json:"created_at"`
}

type customerResponse struct {
	Customer *customer `json:"customer"`
}

type customersResponse struct {
	Customers []customer `json:"customers"`
}

type giftCard struct {
	Id           string       `json:"id"`
	Type         string       `json:"type"`
	GanSource    string       `json:"gan_source"`
	State        string       `json:"state"`
	CreatedAt    string       `json:"created_at"`
	BalanceMoney *amountMoney `json:"balance_money"`
	CustomerIds  []string     `json:"customer_ids"`
}

type giftCardResponse struct {
	GiftCard *giftCard `json:"gift_card"`
}

type giftCardsResponse struct {
	GiftCards []giftCard `json:"gift_cards"`
}

type amountMoney struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type giftCardActivity struct {
	Id                        string `json:"id"`
	Type                      string `json:"type"`
	GiftCardId                string `json:"gift_card_id"`
	GiftCardGan               string `json:"gift_card_gan"`
	GiftCardBalanceMoney      *amountMoney
	CreatedAt                 string                            `json:"created_at"`
	ActivateActivityDetails   *loadActivityDetailsPayload       `json:"activate_activity_details"`
	LoadActivityDetails       *loadActivityDetailsPayload       `json:"load_activity_details"`
	DeactivateActivityDetails *deactivateActivityDetailsPayload `json:"deactivate_activity_details"`
}

type giftCardActivityResponse struct {
	GiftCardActivity *giftCardActivity `json:"gift_card_activity"`
}

type giftCardActivitiesResponse struct {
	GiftCardActivities []giftCardActivity `json:"gift_card_activities"`
}

type createGiftCardPayload struct {
	IdempotencyKey string               `json:"idempotency_key"`
	LocationId     string               `json:"location_id"`
	GiftCard       *giftCardTypePayload `json:"gift_card"`
}

type loadActivityDetailsPayload struct {
	AmountMoney        *amountMoney `json:"amount_money"`
	BuyerInstrumentIds []string     `json:"buyer_payment_instrument_ids"`
}

type deactivateActivityDetailsPayload struct {
	Reason string `json:"reason"`
}

type giftCardActivityPayload struct {
	LocationId                string                            `json:"location_id"`
	GiftCardId                string                            `json:"gift_card_id"`
	Type                      string                            `json:"type"`
	ActivateActivityDetails   *loadActivityDetailsPayload       `json:"activate_activity_details,omitempty"`
	LoadActivityDetails       *loadActivityDetailsPayload       `json:"load_activity_details,omitempty"`
	DeactivateActivityDetails *deactivateActivityDetailsPayload `json:"deactivate_activity_details,omitempty"`
}

type createGiftCardActivityPayload struct {
	IdempotencyKey   string                   `json:"idempotency_key"`
	GiftCardActivity *giftCardActivityPayload `json:"gift_card_activity"`
}

type giftCardTypePayload struct {
	Type string `json:"type"`
	Gan  string `json:"gan"`
}

type sendAGiftCardPayload struct {
	CustomerId string `json:"customer_id"`
}

type templateData struct {
	CSRFToken             string
	PageName              string
	IsAuthenticated       bool
	Customers             []customer
	GiftCards             []giftCard
	GiftCardActivities    []giftCardActivity
	Customer              *customer
	GiftCard              *giftCard
	SelectedGiftCardState string
}
