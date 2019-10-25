package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	msgID    = "000000-10000000"
	mockSend = `100
` + msgID + `
balance=4122.56`

	mockStats = `100
103`

	providerSmsRu = "smsru"
)

type SmsRu struct {
	stor *Storage
}

type smsInfo struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	SmsID      string `json:"sms_id"`
}

type smsRuJsonResponse struct {
	Status     string             `json:"status"`
	StatusCode int                `json:"status_code"`
	StatusText string             `json:"status_text"`
	Sms        map[string]smsInfo `json:"sms"`
	Balance    float64            `json:"balance"`
}

type smsRuJsonStatusResponse struct {
	Status     string  `json:"status"`
	StatusCode int     `json:"status_code"`
	Balance    float64 `json:"balance"`
}

func (sms *SmsRu) Send(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	text := r.FormValue("text")
	isJson := r.FormValue("json")

	sms.stor.Push(Message{
		Time:     time.Now(),
		From:     from,
		To:       to,
		Text:     text,
		Provider: providerSmsRu,
	})

	log.Printf("SmsRu send: from=%q to=%q text=%q", from, to, text)

	var err error
	if isJson == "1" {
		enc := json.NewEncoder(w)
		err = enc.Encode(smsRuJsonResponse{
			Status:     "OK",
			StatusCode: 100,
			StatusText: "OK",
			Sms: map[string]smsInfo{
				to: smsInfo{
					Status:     "OK",
					StatusCode: 100,
					SmsID:      msgID,
				},
			},
			Balance: 4122.56,
		})

	} else {
		_, err = fmt.Fprint(w, mockSend)
	}

	if err != nil {
		log.Println("SmsRu: can not send response to client:", err)
	}
}
func (sms *SmsRu) Status(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("id")
	isJson := r.FormValue("json")

	log.Printf("SmsRu status: id=%q ", ID)

	var err error
	if isJson == "1" {
		enc := json.NewEncoder(w)
		err = enc.Encode(smsRuJsonStatusResponse{
			Status:     "OK",
			StatusCode: 100,
			Balance:    4122.56,
		})
	} else {
		_, err = fmt.Fprint(w, mockStats)
	}

	if err != nil {
		log.Println("SmsRu: can not send status to client:", err)
	}
}
