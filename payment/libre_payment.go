package payment

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type LibrePayment struct {
	stor *MemoryStorage
}

func NewLibrePayment() *LibrePayment {
	return &LibrePayment{
		stor: NewMemoryStorage(),
	}
}

// Payment
func GetPaymentIDFromURL(r *http.Request) (paymentID int64, err error) {
	paymentID, err = strconv.ParseInt(chi.URLParam(r, "paymentID"), 10, 64)

	if err != nil {
		return paymentID, err
	}

	if paymentID == 0 {
		return paymentID, ErrRequestNotPaymentID
	}

	return paymentID, nil
}

func (p *LibrePayment) HasPaymentByID(paymentID int64) bool {
	return p.stor.HasPaymentByID(paymentID)
}

func (p *LibrePayment) GetPaymentByID(paymentID int64) (Payment, error) {
	return p.stor.GetPaymentByID(paymentID)
}

func (p *LibrePayment) AddPayment(payment Payment) error {
	return p.stor.AddPayment(payment)
}

func (p *LibrePayment) UpdatePayment(payment Payment) error {
	return p.stor.UpdatePayment(payment)
}

func (p *LibrePayment) LastPayments(n int) []Payment {
	return p.stor.LastPayments(n)
}

func (p *LibrePayment) ConvertToTemplatePayment(payment Payment) (tplPayment TemplatePayment, err error) {
	provider, err := p.GetProviderByID(payment.ProviderID)
	if err != nil {
		return TemplatePayment{}, err
	}

	tplPayment = TemplatePayment{
		PaymentID:                      payment.PaymentID,
		Provider:                       provider.Type,
		CreatedOn:                      payment.CreatedOn,
		PaymentURL:                     payment.PaymentURL,
		SuccessURL:                     payment.SuccessURL,
		FailURL:                        payment.FailURL,
		NotificationURL:                payment.NotificationURL,
		NotificationResponseOkReceived: payment.NotificationResponseOkReceived,
		Status:                         payment.Status,
		Amount:                         payment.Amount,
		Description:                    payment.Description,
		Recurrent:                      payment.Recurrent,
		ClientID:                       payment.ClientID,
		OrderID:                        payment.OrderID,
	}

	return tplPayment, nil
}

func (p *LibrePayment) ConvertToTemplatePayments(payments []Payment) (tplPayments []TemplatePayment, err error) {
	for _, payment := range payments {
		tplPayment, err := p.ConvertToTemplatePayment(payment)
		if err != nil {
			return []TemplatePayment{}, err
		}

		tplPayments = append(tplPayments, tplPayment)
	}

	return tplPayments, nil
}

// Client
func (p *LibrePayment) HasClientByID(clientID string) bool {
	return p.stor.HasClientByID(clientID)
}

func (p *LibrePayment) HasRebillId(rebillID int64) bool {
	return p.stor.HasRebillId(rebillID)
}

func (p *LibrePayment) GetClientByID(clientID string) (Client, error) {
	return p.stor.GetClientByID(clientID)
}

func (p *LibrePayment) CreateClient(clientID string) (Client, error) {
	client := Client{
		ClientID:   clientID,
		RebillID:   p.generateRebillID(),
		CardNumber: "430000******" + p.generate4Digits(),
	}

	return client, p.stor.AddClient(client)
}

func (p *LibrePayment) generateRebillID() int64 {
	rebillID := rand.Int63n(maxRebillID-minRebillID) + minRebillID

	if p.stor.HasRebillId(rebillID) {
		return p.generateRebillID()
	}

	return rebillID
}

func (p *LibrePayment) generate4Digits() (str4digits string) {
	for i := 0; i < 4; i++ {
		str4digits += strconv.Itoa(rand.Intn(10))
	}

	return str4digits
}

// Provider
func (p *LibrePayment) HasProviderByID(providerID int64) bool {
	return p.stor.HasProviderByID(providerID)
}

func (p *LibrePayment) HasProviderByLogin(login string) bool {
	return p.stor.HasProviderByLogin(login)
}

func (p *LibrePayment) GetProviderByID(providerID int64) (Provider, error) {
	return p.stor.GetProviderByID(providerID)
}

func (p *LibrePayment) GetProviderByLogin(login string) (Provider, error) {
	return p.stor.GetProviderByLogin(login)
}

func (p *LibrePayment) AddProvider(provider Provider) error {
	return p.stor.AddProvider(provider)
}

func (p *LibrePayment) GeneratePaymentID() int64 {
	paymentID := rand.Int63n(maxPaymentID-minPaymentID) + minPaymentID

	if p.stor.HasPaymentByID(paymentID) {
		return p.GeneratePaymentID()
	}

	return paymentID
}
