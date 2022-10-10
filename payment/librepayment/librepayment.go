package librepayment

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type LibrePayment struct {
	stor *MemoryStorage
}

type Payment struct {
	CreatedAt time.Time
	ID        string
	Amount    float64
	Merchant  string
	Status    string

	Payload map[string]string
}

func NewDefaultLibrePyament() *LibrePayment {
	return &LibrePayment{
		stor: NewMemoryStorage(),
	}
}

func (p *LibrePayment) Register(merchant string, amount float64, payload map[string]string) (string, error) {
	internalPaymentID := generateID()

	err := p.stor.Add(time.Now(), internalPaymentID, amount, merchant, payload)
	if err != nil {
		return "", err
	}

	return internalPaymentID, nil
}

func (p *LibrePayment) Status(id string) (Payment, error) {
	storagePayment, err := p.stor.Get(id)
	if err != nil {
		return Payment{}, err
	}

	return storagePaymentToEntity(storagePayment), nil
}

func (p *LibrePayment) Confirm(id string) error {
	return p.stor.SetPaymentStatus(id, StatusConfirmed)
}

func (p *LibrePayment) AllPaymentsDescOrder() ([]Payment, error) {
	storagePayments, err := p.stor.GetAllDesc()
	if err != nil {
		return nil, err
	}

	payments := make([]Payment, len(storagePayments))

	for i, p := range storagePayments {
		payments[i] = storagePaymentToEntity(p)
	}

	return payments, nil
}

func generateID() string {
	return uuid.NewV4().String()
}

func storagePaymentToEntity(p StoragePayment) Payment {
	payment := Payment{
		CreatedAt: p.CreatedAt,
		ID:        p.ID,
		Amount:    p.Amount,
		Merchant:  p.Merchant,
		Status:    p.Status.String(),
		Payload:   map[string]string{},
	}

	for k, v := range p.Payload {
		payment.Payload[k] = v
	}

	return payment
}
