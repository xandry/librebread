package tinkoff

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	paymentpkg "github.com/vasyahuyasa/librebread/payment"
)

func SetStatusHandler(p *paymentpkg.LibrePayment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payment, provider, err := getPaymentAndProvider(r, p)

		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		paymentStatus := chi.URLParam(r, "status")
		redirectUrl := fmt.Sprintf("/payment/%d", payment.PaymentID)

		switch paymentStatus {
		default:
			log.Printf("Tinkoff: %v; Status: %s", paymentpkg.ErrUnknownPaymentStatus, paymentStatus)

		case string(StatusDeadlineExpired):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

		case string(StatusAttemptsExpired):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

		case string(StatusFormShowed):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

		case string(StatusAuthorized):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

		case string(StatusRefunded):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

		case string(StatusConfirmed):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

			redirectUrl = strings.Replace(payment.SuccessURL, "${Success}", "true", -1)
			redirectUrl = strings.Replace(redirectUrl, "${ErrorCode}", "0", -1)

		case string(StatusRejected):
			payment.Status = paymentStatus
			p.UpdatePayment(payment)

			redirectUrl = strings.Replace(payment.FailURL, "${Success}", "false", -1)
			redirectUrl = strings.Replace(redirectUrl, "${ErrorCode}", "9999", -1)

		}

		if payment.NotificationURL != "" {
			needToSendNotification := false

			switch paymentStatus {
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
				client, _ := p.GetClientByID(payment.ClientID)

				ok, err := sendNotification(payment, client, provider)

				payment.NotificationResponseOkReceived = false

				if err != nil {
					log.Printf("Tinkoff: %v", err)
				} else if ok {
					payment.NotificationResponseOkReceived = true
				} else {
					log.Println("Tinkoff: incorrect response to the notification was received")
				}

				p.UpdatePayment(payment)
			}
		}

		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
	}
}

func SendNotificationHandler(p *paymentpkg.LibrePayment) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payment, provider, err := getPaymentAndProvider(r, p)

		if err != nil {
			log.Printf("Tinkoff: %v", err)
			return
		}

		client, _ := p.GetClientByID(payment.ClientID)

		ok, err := sendNotification(payment, client, provider)

		payment.NotificationResponseOkReceived = false

		if err != nil {
			log.Printf("Tinkoff: %v", err)
		} else if ok {
			payment.NotificationResponseOkReceived = true
		} else {
			log.Println("Tinkoff: incorrect response to the notification was received")
		}

		p.UpdatePayment(payment)

		redirectUrl := fmt.Sprintf("/payment/%d", payment.PaymentID)
		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
	}
}

func getPaymentAndProvider(r *http.Request, p *paymentpkg.LibrePayment) (paymentpkg.Payment, paymentpkg.Provider, error) {
	paymentID, err := paymentpkg.GetPaymentIDFromURL(r)

	if err != nil {
		return paymentpkg.Payment{}, paymentpkg.Provider{}, err
	}

	payment, err := p.GetPaymentByID(paymentID)

	if err != nil {
		return paymentpkg.Payment{}, paymentpkg.Provider{}, err
	}

	provider, err := p.GetProviderByID(payment.ProviderID)

	if err != nil {
		return paymentpkg.Payment{}, paymentpkg.Provider{}, err
	}

	if provider.Type != paymentpkg.TinkoffProvider {
		return paymentpkg.Payment{}, paymentpkg.Provider{}, paymentpkg.ErrIncorrectPaymentProvider
	}

	return payment, provider, nil
}
