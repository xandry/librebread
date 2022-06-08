package payment

import (
	"log"
	"net/http"
	"text/template"
)

const listPaymentsLimit = 100

func (p *LibrePayment) IndexPaymentHandler(w http.ResponseWriter, r *http.Request) {
	listPayments := p.LastPayments(listPaymentsLimit)
	templatePayments, err := p.ConvertToTemplatePayments(listPayments)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	payments := struct {
		Payments []TemplatePayment
	}{
		Payments: templatePayments,
	}

	indexTemplate := template.Must(template.New("index").Parse(tplIndex))

	if err := indexTemplate.Execute(w, payments); err != nil {
		log.Printf("Payment: %v", err)
		return
	}
}

func (p *LibrePayment) ViewPaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentID, err := GetPaymentIDFromURL(r)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	payment, err := p.GetPaymentByID(paymentID)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	templatePayment, err := p.ConvertToTemplatePayment(payment)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	provider, err := p.GetProviderByID(payment.ProviderID)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	switch provider.Type {
	default:
		log.Printf("Payment: %v", ErrIncorrectPaymentProvider)
		return
	case TinkoffProvider:
		viewTemplate := template.Must(template.New("view").Parse(tplTinkoffView))

		if err := viewTemplate.Execute(w, templatePayment); err != nil {
			log.Printf("Payment: %v", err)
			return
		}
	}
}
