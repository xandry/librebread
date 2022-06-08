package payment

import (
	"sync"
	"time"
)

const (
	minPaymentID = 1000
	maxPaymentID = 10000000
	minRebillID  = 1000
	maxRebillID  = 10000000
)

type ProviderType string

const (
	TinkoffProvider  ProviderType = "Tinkoff"
	SberbankProvider ProviderType = "Sberbank"
)

type URLScheme string

const (
	httpScheme  URLScheme = "http"
	httpsScheme URLScheme = "https"
)

type MemoryStorage struct {
	storageMu                  sync.Mutex
	providers                  map[int64]Provider
	payments                   map[int64]Payment
	listPaymentIDsInOrder      []int64
	clients                    map[string]Client
	clientIDsIndexesByRebillID map[int64]string
}

type Payment struct {
	PaymentID                      int64
	ProviderID                     int64
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

type Client struct {
	ClientID   string
	RebillID   int64
	CardNumber string
}

type Provider struct {
	ProviderID       int64
	Type             ProviderType
	Login            string
	Password         string
	PaymentURLScheme URLScheme
}

var mockProviders = map[int64]Provider{
	1: {
		ProviderID:       1,
		Type:             TinkoffProvider,
		Login:            "tinkoff-root",
		Password:         "root",
		PaymentURLScheme: httpScheme,
	},
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		providers:                  mockProviders,
		payments:                   make(map[int64]Payment),
		clients:                    make(map[string]Client),
		clientIDsIndexesByRebillID: make(map[int64]string),
	}
}

// Payment
func (s *MemoryStorage) HasPaymentByID(paymentID int64) bool {
	_, ok := s.payments[paymentID]
	return ok
}

func (s *MemoryStorage) GetPaymentByID(paymentID int64) (Payment, error) {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if !s.HasPaymentByID(paymentID) {
		return Payment{}, ErrPaymentNotFound
	}
	return s.payments[paymentID], nil
}

func (s *MemoryStorage) AddPayment(p Payment) error {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if s.HasPaymentByID(p.PaymentID) {
		return ErrPaymentAlreadyExists
	}

	s.payments[p.PaymentID] = p
	s.listPaymentIDsInOrder = append(s.listPaymentIDsInOrder, p.PaymentID)
	return nil
}

func (s *MemoryStorage) UpdatePayment(p Payment) error {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if !s.HasPaymentByID(p.PaymentID) {
		return ErrPaymentNotFound
	}

	s.payments[p.PaymentID] = p
	return nil
}

func (s *MemoryStorage) LastPayments(n int) []Payment {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	payments := make([]Payment, 0, n)
	length := len(s.payments) - 1

	for i := 0; length-i >= 0 && i < n; i++ {
		payments = append(payments, s.payments[s.listPaymentIDsInOrder[length-i]])
	}
	return payments
}

// Client
func (s *MemoryStorage) HasClientByID(clientID string) bool {
	_, ok := s.clients[clientID]
	return ok
}

func (s *MemoryStorage) HasRebillId(rebillID int64) bool {
	_, ok := s.clientIDsIndexesByRebillID[rebillID]
	return ok
}

func (s *MemoryStorage) GetClientByID(clientID string) (Client, error) {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if !s.HasClientByID(clientID) {
		return Client{}, ErrClientNotFound
	}
	return s.clients[clientID], nil
}

func (s *MemoryStorage) AddClient(c Client) error {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if s.HasClientByID(c.ClientID) {
		return ErrClientAlreadyExists
	}

	if s.HasRebillId(c.RebillID) {
		return ErrRebillIDAlreadyExists
	}

	s.clients[c.ClientID] = c
	s.clientIDsIndexesByRebillID[c.RebillID] = c.ClientID
	return nil
}

// Provider
func (s *MemoryStorage) HasProviderByID(providerID int64) bool {
	_, ok := s.providers[providerID]
	return ok
}

func (s *MemoryStorage) HasProviderByLogin(login string) bool {
	for id := range s.providers {
		if s.providers[id].Login == login {
			return true
		}
	}
	return false
}

func (s *MemoryStorage) GetProviderByID(providerID int64) (Provider, error) {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if !s.HasProviderByID(providerID) {
		return Provider{}, ErrProviderNotFound
	}
	return s.providers[providerID], nil
}

func (s *MemoryStorage) GetProviderByLogin(login string) (Provider, error) {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	for id := range s.providers {
		if s.providers[id].Login == login {
			return s.providers[id], nil
		}
	}
	return Provider{}, ErrProviderNotFound
}

func (s *MemoryStorage) AddProvider(p Provider) error {
	s.storageMu.Lock()
	defer s.storageMu.Unlock()

	if s.HasProviderByID(p.ProviderID) {
		return ErrProviderIDAlreadyExists
	}

	if s.HasProviderByLogin(p.Login) {
		return ErrProviderLoginAlreadyExists
	}

	s.providers[p.ProviderID] = p
	return nil
}
