package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const mockSessionID = "MOCK-SESSION-ID"

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
	stor *Storage
}

func (sms *Devino) UserSessionIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%q", mockSessionID)
}

func (sms *Devino) SmsSend(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("SessionID")
	srcAddr := r.FormValue("SourceAddress")
	dstAddr := r.FormValue("DestinationAddress")
	data := r.FormValue("Data")

	sms.stor.Push(Message{
		Time: time.Now(),
		From: srcAddr,
		To:   dstAddr,
		Text: data,
	})

	log.Printf("Devino: send session=%q src=%q dst=%q msg=%q", sessionID, srcAddr, dstAddr, data)

	enc := json.NewEncoder(w)
	err := enc.Encode(mockMessageIDS)
	if err != nil {
		log.Printf("Devino: can send sessionID to client: %v", err)
	}
}
func (sms *Devino) SmsState(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("SessionID")
	msgID := r.FormValue("messageId")

	log.Printf("Devino: session=%q msgId=%q", sessionID, msgID)

	enc := json.NewEncoder(w)
	err := enc.Encode(mockState)
	if err != nil {
		log.Printf("Devino: state can send state to client: %v", err)
	}
}
