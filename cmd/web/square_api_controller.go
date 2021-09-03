package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type squareApiController struct {
	apiClient *http.Client
}

func NewSquareApiController(r *http.Request) *squareApiController {
	return &squareApiController{
		apiClient: NewSquareOauthController().getSquareHttpClient(r),
	}
}

func (controller *squareApiController) doRequest(method string, base_url string, url_params *url.Values, body_params interface{}) []byte {
	var req *http.Request
	var err error
	var paramsJson []byte
	var url string

	if body_params != nil {
		paramsJson, _ = json.Marshal(body_params)
	}

	if url_params != nil {
		url = fmt.Sprintf("%s?%s", base_url, url_params.Encode())
	} else {
		url = base_url
	}

	req, err = http.NewRequest(method, url, bytes.NewBuffer(paramsJson))
	if err != nil {
		panic(err)
	}

	resp, err := controller.apiClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(body)
	}

	return body
}

func (controller *squareApiController) getCustomers(params *url.Values) []customer {
	resp := controller.doRequest(http.MethodGet, SQUARE_CUSTOMERS_URL, params, nil)

	var customers_resp customersResponse
	json.Unmarshal(resp, &customers_resp)
	return customers_resp.Customers
}

func (controller *squareApiController) getACustomer(id string) *customer {
	url := fmt.Sprintf("%s/%s", SQUARE_CUSTOMERS_URL, id)
	resp := controller.doRequest(http.MethodGet, url, nil, nil)

	var customer_resp customerResponse
	json.Unmarshal(resp, &customer_resp)
	return customer_resp.Customer
}

func (controller *squareApiController) getAGiftCard(id string) *giftCard {
	url := fmt.Sprintf("%s/%s", SQUARE_GIFT_CARDS_URL, id)
	resp := controller.doRequest(http.MethodGet, url, nil, nil)

	var gift_card_resp giftCardResponse
	json.Unmarshal(resp, &gift_card_resp)
	return gift_card_resp.GiftCard
}

func (controller *squareApiController) getGiftCards(params *url.Values) []giftCard {
	resp := controller.doRequest(http.MethodGet, SQUARE_GIFT_CARDS_URL, params, nil)

	var gift_cards_resp giftCardsResponse
	json.Unmarshal(resp, &gift_cards_resp)
	return gift_cards_resp.GiftCards
}

func (controller *squareApiController) getGiftCardActivities(params *url.Values) []giftCardActivity {
	resp := controller.doRequest(http.MethodGet, SQUARE_GIFT_CARD_ACTIVITIES_URL, params, nil)

	var gift_card_activities_resp giftCardActivitiesResponse
	json.Unmarshal(resp, &gift_card_activities_resp)
	return gift_card_activities_resp.GiftCardActivities
}

func (controller *squareApiController) createAGiftCard(params interface{}) *giftCard {
	resp := controller.doRequest(http.MethodPost, SQUARE_GIFT_CARDS_URL, nil, params)

	var gift_card_resp giftCardResponse
	json.Unmarshal(resp, &gift_card_resp)
	return gift_card_resp.GiftCard
}

func (controller *squareApiController) createAGiftCardActivity(params interface{}) *giftCardActivity {
	resp := controller.doRequest(http.MethodPost, SQUARE_GIFT_CARD_ACTIVITIES_URL, nil, params)

	var gift_card_activity_resp giftCardActivityResponse
	json.Unmarshal(resp, &gift_card_activity_resp)
	return gift_card_activity_resp.GiftCardActivity
}

func (controller *squareApiController) sendAGiftCard(gift_card_id string, params interface{}) *giftCard {
	url := fmt.Sprintf("%s/%s/link-customer", SQUARE_GIFT_CARDS_URL, gift_card_id)
	resp := controller.doRequest(http.MethodPost, url, nil, params)

	var gift_card_resp giftCardResponse
	json.Unmarshal(resp, &gift_card_resp)
	return gift_card_resp.GiftCard
}
