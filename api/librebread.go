//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config oapi-codgen.yml openapi_librebread.yml

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type smser interface {
	LastMessages(limit int64) (SMSList, error)
	Create(from, to, text, provider string) (string, error)
}

type ticketer interface {
	LastTickets(limit int64) (HelpdeskEddyTicketList, error)
	Create(title, description string, typeID, priorityID, departmentID int) error
}

type LibreBread struct {
	sms    smser
	ticket ticketer
	re     *renderer
}

var _ ServerInterface = (*LibreBread)(nil)

func NewLibrebread(sms smser, ticket ticketer) *LibreBread {
	return &LibreBread{
		sms:    sms,
		ticket: ticket,
		re:     newRenderer(),
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

	smses := make(SMSList, len(messages))

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

func (lb *LibreBread) PostLibreSend(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	if from == "" {
		http.Error(w, "from param is required", http.StatusBadRequest)
		return
	}

	to := r.FormValue("to")
	if to == "" {
		http.Error(w, "to param is required", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "text param is required", http.StatusBadRequest)
		return
	}

	id, err := lb.sms.Create(from, to, text, "LibreSMS")
	if err != nil {
		log.Printf("can not create libre SMS: %v", err)
		http.Error(w, fmt.Sprintf("can not create libre SMS: %v", err), http.StatusInternalServerError)
		return
	}

	response := LibreBreadSMSIds{id}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("can send response to client: %v", err)
	}
}

func (lb *LibreBread) PostLibreCheck(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}

func (lb *LibreBread) GetHelpdeskEddyTicket(w http.ResponseWriter, r *http.Request, params GetHelpdeskEddyTicketParams) {
	var limit int64 = 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	tickets, err := lb.ticket.LastTickets(limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("can not get helpdeskeddy tickets: %v", err), http.StatusInternalServerError)
		log.Printf("can not get helpdeskeddy messages: %v", err)
		return
	}

	if params.Json != nil && *params.Json {
		enc := json.NewEncoder(w)

		err = enc.Encode(tickets)
		if err != nil {
			http.Error(w, fmt.Sprintf("can not encode messages: %v", err), http.StatusInternalServerError)
			log.Printf("can not encode messages: %v", err)
		}

		return
	}

	err = lb.re.renderHelpdeskeddy(w, tickets)
	if err != nil {
		http.Error(w, fmt.Sprintf("can not render messages: %v", err), http.StatusInternalServerError)
		log.Printf("can not render messages: %v", err)
	}
}

func (lb *LibreBread) PostHelpdeskEddyTicket(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" || description == "" {
		http.Error(w, "title and description required", http.StatusBadRequest)
		return
	}

	typeID := atoiOrZero(r.FormValue("type_id"))
	priorityID := atoiOrZero(r.FormValue("priority_id"))
	departmentID := atoiOrZero(r.FormValue("department_id"))

	err := lb.ticket.Create(title, description, typeID, priorityID, departmentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("can not create HelpdeskEddy ticket: %v", err), http.StatusInternalServerError)
		log.Printf("can not create HelpdeskEddy ticket: %v", err)
		return
	}

	log.Printf("HelpdeskEddy ticket created")
}

func atoiOrZero(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}
