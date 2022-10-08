package librepayment

import (
	"fmt"
	"os"
)

type PaymentUrlGenerator struct {
	baseURL string
}

func NewPaymentUrlGeneratorFromENV() *PaymentUrlGenerator {
	baseURL := "http://localhost"

	url, ok := os.LookupEnv("LIBREPAYMENT_BASE_URL")
	if ok {
		baseURL = url
	}

	return &PaymentUrlGenerator{
		baseURL: baseURL,
	}
}

func (g *PaymentUrlGenerator) Generate(id string) string {
	return fmt.Sprintf("%s/librepayments/%s", g.baseURL, id)
}
