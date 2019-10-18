package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	mockSend = `100
000000-10000000
balance=4122.56`

	mockStats = `100
103`
)

type SmsRu struct {
	stor *Storage
}

func (sms *SmsRu) Send(w http.ResponseWriter, r *http.Request) {
	apiID := r.FormValue("api_id")
	from := r.FormValue("from")
	to := r.FormValue("to")
	text := r.FormValue("text")

	sms.stor.Push(Message{
		Time: time.Now(),
		From: from,
		To:   to,
		Text: text,
	})

	log.Printf("SmsRu: send api_id=%q from=%q to=%q text=%q", apiID, from, to, text)

	_, err := fmt.Fprint(w, mockSend)
	if err != nil {
		log.Println("SmsRu: can not send response to client:", err)
	}
}
func (sms *SmsRu) Status(w http.ResponseWriter, r *http.Request) {
	apiID := r.FormValue("api_id")
	ID := r.FormValue("id")

	log.Printf("SmsRu: status api_id=%q id=%q ", apiID, ID)

	_, err := fmt.Fprint(w, mockStats)
	if err != nil {
		log.Println("SmsRu: can not send status to client:", err)
	}
}
