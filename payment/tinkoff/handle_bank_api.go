package tinkoff

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	paymentpkg "github.com/vasyahuyasa/librebread/payment"
)

// Инициализирует платеж
func InitHandler(p *paymentpkg.LibrePayment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input InitInput

		if err := readJSON(r, &input); err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		provider, err := p.GetProviderByLogin(input.TerminalKey)
		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		if provider.Type != paymentpkg.TinkoffProvider {
			log.Printf("Tinkoff: %v", paymentpkg.ErrIncorrectPaymentProvider)
			return
		}

		var client paymentpkg.Client
		if input.CustomerKey != "" {
			if p.HasClientByID(input.CustomerKey) {
				client, err = p.GetClientByID(input.CustomerKey)
				if err != nil {
					log.Printf("Tinkoff: %v", err)
					return
				}
			} else {
				client, err = p.CreateClient(input.CustomerKey)
				if err != nil {
					log.Printf("Tinkoff: %v", err)
					return
				}
			}
		}

		isRecurrent := false
		if input.Recurrent == "Y" {
			isRecurrent = true
		}

		scheme := provider.PaymentURLScheme
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		paymentID := p.GeneratePaymentID()

		paymentURL := fmt.Sprintf("%s://%s/payment/%d", scheme, host, paymentID)

		payment := paymentpkg.Payment{
			PaymentID:       paymentID,
			ProviderID:      provider.ProviderID,
			CreatedOn:       time.Now(),
			PaymentURL:      paymentURL,
			SuccessURL:      input.SuccessURL,
			FailURL:         input.FailURL,
			NotificationURL: input.NotificationURL,
			Status:          string(StatusNew),
			Amount:          input.Amount,
			Description:     input.Description,
			Recurrent:       isRecurrent,
			ClientID:        client.ClientID,
			OrderID:         input.OrderID,
		}

		p.AddPayment(payment)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(initResponse(payment, input.TerminalKey))
	}
}

// Выполняет автоплатеж
func ChargeHandler(p *paymentpkg.LibrePayment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: доделать
		w.Write([]byte("OK"))
	}
}

func GetStateHandler(p *paymentpkg.LibrePayment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input GetStateInput

		if err := readJSON(r, &input); err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		if input.PaymentID == 0 {
			log.Printf("Tinkoff: %v", paymentpkg.ErrRequestNotPaymentID)
			return
		}

		payment, err := p.GetPaymentByID(input.PaymentID)

		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(getStateResponse(payment, input.TerminalKey))
	}
}
