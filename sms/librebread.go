package sms

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

const providerLibreBread = "LibreBread"

type LibreBread struct {
	Stor     *Storage
	Notifier SmsNotifier
	lastId   int64
}

func (sms *LibreBread) Send(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	text := r.FormValue("text")

	msg := Message{
		Time:     time.Now(),
		From:     from,
		To:       to,
		Text:     text,
		Provider: providerLibreBread,
	}
	sms.Stor.Push(msg)

	sms.Notifier.SmsNotify(msg)

	log.Printf("LibreBread send: from=%q to=%q text=%q", from, to, text)

	// number of sms with UTF-16 encoding
	chars := utf8.RuneCountInString(text)
	numSMS := 1

	if chars > 70 {
		numSMS = int(math.Ceil(float64(chars) / float64(67)))
	}

	var ids []int64
	for i := 0; i < numSMS; i++ {
		ids = append(ids, atomic.AddInt64(&sms.lastId, 1))
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(ids)

	if err != nil {
		log.Println("SmsRu: can not send response to client:", err)
	}
}
func (sms *LibreBread) Check(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
