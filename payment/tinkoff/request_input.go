package tinkoff

import (
	"encoding/json"
	"net/http"

	"github.com/vasyahuyasa/librebread/payment"
)

type initInput struct {
	TerminalKey     string `json:"TerminalKey"`
	Amount          int64  `json:"Amount"`
	OrderID         string `json:"OrderId"`
	Description     string `json:"Description"`
	Recurrent       string `json:"Recurrent"`
	NotificationURL string `json:"NotificationURL"`
	SuccessURL      string `json:"SuccessURL"`
	FailURL         string `json:"FailURL"`
	CustomerKey     string `json:"CustomerKey"`
	Token           string `json:"Token"`
}

type chargeInput struct {
	TerminalKey string `json:"TerminalKey"`
	ProcessID   int64  `json:"PaymentId"`
	RebillID    int64  `json:"RebillId"`
	Token       string `json:"Token"`
}

type getStateInput struct {
	TerminalKey string `json:"TerminalKey"`
	ProcessID   int64  `json:"PaymentId"`
	Token       string `json:"Token"`
}

func initRequest(r *http.Request) (input initInput, err error) {
	err = readJSON(r, &input)
	return input, err
}

func chargeRequest(r *http.Request) (input chargeInput, err error) {
	err = readJSON(r, &input)
	return input, err
}

func getStateRequest(r *http.Request) (input getStateInput, err error) {
	err = readJSON(r, &input)
	return input, err
}

func readJSON(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return payment.WrongJSONFormat
	}
	return nil
}
