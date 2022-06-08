package tinkoff

import "github.com/vasyahuyasa/librebread/payment"

type initOutput struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	PaymentID   int64  `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
	PaymentURL  string `json:"PaymentURL"`
}

type getStateOutput struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	PaymentID   int64  `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
}

func initResponse(payment payment.Payment, terminalKey string) initOutput {
	return initOutput{
		TerminalKey: terminalKey,
		Amount:      payment.Amount,
		OrderID:     payment.OrderID,
		Success:     true,
		Status:      payment.Status,
		PaymentID:   payment.PaymentID,
		ErrorCode:   "0",
		PaymentURL:  payment.PaymentURL,
	}
}

func getStateResponse(payment payment.Payment, terminalKey string) getStateOutput {
	return getStateOutput{
		TerminalKey: terminalKey,
		Amount:      payment.Amount,
		OrderID:     payment.OrderID,
		Success:     true,
		Status:      payment.Status,
		PaymentID:   payment.PaymentID,
		ErrorCode:   "0",
	}
}
