package payment

import "time"

type TemplatePayment struct {
	PaymentID                      int64
	Provider                       ProviderType
	CreatedOn                      time.Time
	PaymentURL                     string
	SuccessURL                     string
	FailURL                        string
	NotificationURL                string
	NotificationResponseOkReceived bool
	Status                         string
	Amount                         int64
	Description                    string
	Recurrent                      bool
	ClientID                       string
	OrderID                        string
}
