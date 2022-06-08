package tinkoff

import (
	"encoding/json"
	"net/http"

	"github.com/vasyahuyasa/librebread/payment"
)

type InitInput struct {
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

type GetStateInput struct {
	TerminalKey string `json:"TerminalKey"`
	PaymentID   int64  `json:"PaymentId"`
	Token       string `json:"Token"`
}

func readJSON(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return payment.WrongJSONFormat
	}
	return nil
}
