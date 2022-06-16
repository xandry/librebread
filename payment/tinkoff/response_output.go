package tinkoff

import "github.com/vasyahuyasa/librebread/payment"

type initOutput struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	ProcessID   int64  `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
	PaymentURL  string `json:"PaymentURL"`
}

type chargeOutput struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	ProcessID   int64  `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
}

type getStateOutput struct {
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	ProcessID   int64  `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
}

func initResponse(paymentProcess payment.PaymentProcess, terminalKey string) initOutput {
	return initOutput{
		TerminalKey: terminalKey,
		Amount:      paymentProcess.Amount,
		OrderID:     paymentProcess.OrderID,
		Success:     true,
		Status:      paymentProcess.Status,
		ProcessID:   paymentProcess.ProcessID,
		ErrorCode:   "0",
		PaymentURL:  paymentProcess.PaymentURL,
	}
}

func chargeResponse(paymentProcess payment.PaymentProcess, terminalKey string) chargeOutput {
	return chargeOutput{
		TerminalKey: terminalKey,
		Amount:      paymentProcess.Amount,
		OrderID:     paymentProcess.OrderID,
		Success:     true,
		Status:      paymentProcess.Status,
		ProcessID:   paymentProcess.ProcessID,
		ErrorCode:   "0",
	}
}

func getStateResponse(paymentProcess payment.PaymentProcess, terminalKey string) getStateOutput {
	return getStateOutput{
		TerminalKey: terminalKey,
		Amount:      paymentProcess.Amount,
		OrderID:     paymentProcess.OrderID,
		Success:     true,
		Status:      paymentProcess.Status,
		ProcessID:   paymentProcess.ProcessID,
		ErrorCode:   "0",
	}
}
