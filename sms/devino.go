package sms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const mockSessionID = "MOCK-SESSION-ID"
const providerDevino = "devino"

var mockMessageIDS = []string{"579700854169272358"}
var mockState = struct {
	State            int
	CreationDateUtc  interface{}
	SubmittedDateUtc interface{}
	ReportedDateUtc  interface{}
	TimeStampUtc     string
	StateDescription string
	Price            interface{}
}{
	State:            255,
	CreationDateUtc:  nil,
	SubmittedDateUtc: nil,
	ReportedDateUtc:  nil,
	TimeStampUtc:     "62135596800000",
	StateDescription: "Неизвестный",
	Price:            nil,
}

type Devino struct {
	Stor *Storage
}

func (sms *Devino) UserSessionIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q", mockSessionID)
}

func (sms *Devino) SmsSend(w http.ResponseWriter, r *http.Request) {
	srcAddr := r.FormValue("SourceAddress")
	dstAddr := r.FormValue("DestinationAddress")
	data := r.FormValue("Data")

	sms.Stor.Push(Message{
		Time:     time.Now(),
		From:     srcAddr,
		To:       dstAddr,
		Text:     data,
		Provider: providerDevino,
	})

	log.Printf("Devino send: src=%q dst=%q msg=%q", srcAddr, dstAddr, data)

	enc := json.NewEncoder(w)
	err := enc.Encode(mockMessageIDS)
	if err != nil {
		log.Printf("Devino: can send sessionID to client: %v", err)
	}
}
func (sms *Devino) SmsState(w http.ResponseWriter, r *http.Request) {
	msgID := r.FormValue("messageId")

	log.Printf("Devino state: msgId=%q", msgID)

	enc := json.NewEncoder(w)
	err := enc.Encode(mockState)
	if err != nil {
		log.Printf("Devino: state can send state to client: %v", err)
	}
}
