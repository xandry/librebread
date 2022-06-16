package tinkoff

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/vasyahuyasa/librebread/payment"
)

func SetStatusHandler(p *payment.Payment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		process, provider, err := getPaymentProcessAndProvider(r, p)
		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		PaymentProcessStatus := chi.URLParam(r, "status")
		redirectUrl := fmt.Sprintf("/payment/%d", process.ProcessID)

		switch PaymentProcessStatus {
		default:
			log.Printf("Tinkoff: %v; Status: %s", payment.ErrUnknownPaymentProcessStatus, PaymentProcessStatus)

		case string(StatusDeadlineExpired):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

		case string(StatusAttemptsExpired):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

		case string(StatusFormShowed):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

		case string(StatusAuthorized):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

		case string(StatusRefunded):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

		case string(StatusConfirmed):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

			redirectUrl = strings.Replace(process.SuccessURL, "${Success}", "true", -1)
			redirectUrl = strings.Replace(redirectUrl, "${ErrorCode}", "0", -1)

		case string(StatusRejected):
			process.Status = PaymentProcessStatus
			p.UpdateProcess(process)

			redirectUrl = strings.Replace(process.FailURL, "${Success}", "false", -1)
			redirectUrl = strings.Replace(redirectUrl, "${ErrorCode}", "9999", -1)

		}

		if process.NotificationURL != "" {
			needToSendNotification := false

			switch PaymentProcessStatus {
			case string(StatusAuthorized):
				needToSendNotification = true
			case string(StatusConfirmed):
				needToSendNotification = true
			case string(StatusRejected):
				needToSendNotification = true
			case string(StatusRefunded):
				needToSendNotification = true
			}

			if needToSendNotification {
				client, _ := p.GetClientByID(process.ClientID)

				ok, err := sendNotification(process, client, provider)

				process.NotificationResponseOkReceived = false

				if err != nil {
					log.Printf("Tinkoff: %v", err)
				} else if ok {
					process.NotificationResponseOkReceived = true
				} else {
					log.Println("Tinkoff: incorrect response to the notification was received")
				}

				p.UpdateProcess(process)
			}
		}

		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
	}
}

func SendNotificationHandler(p *payment.Payment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		process, provider, err := getPaymentProcessAndProvider(r, p)
		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		process.NotificationResponseOkReceived = false

		client, _ := p.GetClientByID(process.ClientID)

		ok, err := sendNotification(process, client, provider)
		if err != nil {
			log.Printf("Tinkoff: %v", err)
		} else if ok {
			process.NotificationResponseOkReceived = true
		} else {
			log.Println("Tinkoff: incorrect response to the notification was received")
		}

		p.UpdateProcess(process)

		redirectUrl := fmt.Sprintf("/payment/%d", process.ProcessID)
		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
	}
}

func getPaymentProcessAndProvider(r *http.Request, p *payment.Payment) (payment.PaymentProcess, payment.Provider, error) {
	processID, err := payment.GetProcessIDFromURL(r)
	if err != nil {
		return payment.PaymentProcess{}, payment.Provider{}, err
	}

	process, err := p.GetProcessByID(processID)
	if err != nil {
		return payment.PaymentProcess{}, payment.Provider{}, err
	}

	provider, err := p.GetProviderByID(process.ProviderID)
	if err != nil {
		return payment.PaymentProcess{}, payment.Provider{}, err
	}

	if provider.Type != payment.TinkoffProvider {
		return payment.PaymentProcess{}, payment.Provider{}, payment.ErrIncorrectPaymentProvider
	}

	return process, provider, nil
}
