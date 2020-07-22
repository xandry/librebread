//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=api --generate types,chi-server -o librebread.gen.go librebread.yml

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type smser interface {
	LastMessages(limit int64) (SMSes, error)
}

type LibreBread struct {
	sms smser
	re  *renderer
}

func NewLibrebread(sms smser) *LibreBread {
	return &LibreBread{
		sms: sms,
		re:  newRenderer(),
	}
}

func (lb *LibreBread) Get(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/sms", http.StatusTemporaryRedirect)
}

func (lb *LibreBread) GetSms(w http.ResponseWriter, r *http.Request, params GetSmsParams) {
	var limit int64 = 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	messages, err := lb.sms.LastMessages(limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("can not get messages: %v", err), http.StatusInternalServerError)
		log.Printf("can not get messages: %v", err)
		return
	}

	smses := make(SMSes, len(messages))

	for i, m := range messages {
		smses[i] = SMS{
			From:     m.From,
			ID:       m.ID,
			Provider: m.Provider,
			Text:     m.Text,
			Time:     m.Time,
			To:       m.To,
		}
	}

	if params.Json != nil && *params.Json {
		enc := json.NewEncoder(w)

		err = enc.Encode(smses)
		if err != nil {
			http.Error(w, fmt.Sprintf("can not encode messages: %v", err), http.StatusInternalServerError)
			log.Printf("can not encode messages: %v", err)
		}

		return
	}

	err = lb.re.renderSms(w, smses)
	if err != nil {
		http.Error(w, fmt.Sprintf("can not render messages: %v", err), http.StatusInternalServerError)
		log.Printf("can not render messages: %v", err)
	}
}
