package payment

import (
	"log"
	"net/http"
	"text/template"
)

const listPaymentProcessesLimit = 100

func (p *Payment) IndexPaymentHandler(w http.ResponseWriter, r *http.Request) {
	listProcesses := p.LastProcesses(listPaymentProcessesLimit)
	templateProcesses, err := p.ConvertToTemplatePaymentProcesses(listProcesses)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	paymentProcessesData := struct {
		NumberOfProcesses int
		Processes         []TemplatePaymentProcess
	}{
		NumberOfProcesses: p.ProcessesLen(),
		Processes:         templateProcesses,
	}

	indexTemplate := template.Must(template.New("index").Parse(tplIndex))

	if err := indexTemplate.Execute(w, paymentProcessesData); err != nil {
		log.Printf("Payment: %v", err)
		return
	}
}

func (p *Payment) ViewPaymentHandler(w http.ResponseWriter, r *http.Request) {
	processID, err := GetProcessIDFromURL(r)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	paymentProcess, err := p.GetProcessByID(processID)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	templatePaymentProcess, err := p.ConvertToTemplatePaymentProcess(paymentProcess)

	if err != nil {
		log.Printf("Payment: %v", err)
		return
	}

	provider, err := p.GetProviderByID(paymentProcess.ProviderID)

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

		if err := viewTemplate.Execute(w, templatePaymentProcess); err != nil {
			log.Printf("Payment: %v", err)
			return
		}
	}
}
