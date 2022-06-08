package payment

import (
	"errors"
)

var (
	ErrPaymentNotFound             = errors.New("Payment not found")
	ErrPaymentNotProvidedInContext = errors.New("Payment not provided in context")
	ErrPaymentAlreadyExists        = errors.New("Payment already exists")
	ErrUnknownPaymentStatus        = errors.New("Unknown payment status")

	ErrProviderNotFound           = errors.New("Provider not found")
	ErrProviderIDAlreadyExists    = errors.New("Provider with this ID already exists")
	ErrProviderLoginAlreadyExists = errors.New("Provider with this login already exists")
	ErrIncorrectPaymentProvider   = errors.New("Incorrect payment provider")

	ErrClientNotFound        = errors.New("Client not found")
	ErrClientAlreadyExists   = errors.New("Client already exists")
	ErrRebillIDAlreadyExists = errors.New("RebillID already exists")

	ErrRequestNotPaymentID = errors.New("PaymentID param is required and must be a number")

	WrongJSONFormat = errors.New("Wrong json format")
)
