package tinkoff

type PaymentProcessStatus string

// Documentation on payment statuses:
// https://www.tinkoff.ru/kassa/develop/api/payments/
const (
	StatusNew             PaymentProcessStatus = "NEW"              // Payment has been created
	StatusDeadlineExpired PaymentProcessStatus = "DEADLINE_EXPIRED" // Payment time has expired
	StatusAttemptsExpired PaymentProcessStatus = "ATTEMPTS_EXPIRED" // Attempts to open the form have exhausted
	StatusFormShowed      PaymentProcessStatus = "FORM_SHOWED"      // Payment form is opened by the buyer
	StatusAuthorized      PaymentProcessStatus = "AUTHORIZED"       // Funds are reserved
	StatusRejected        PaymentProcessStatus = "REJECTED"         // Payment was canceled by the bank
	StatusConfirmed       PaymentProcessStatus = "CONFIRMED"        // Confirmed by
	StatusRefunded        PaymentProcessStatus = "REFUNDED"         // Funds returned in full
)
