package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
)

const (
	addr      = ":8099"
	filename  = "messages.gob"
	tplHeader = `<html><body><table border=1><thead><th>Date</th><th>From</th><th>Phone</th><th>Msg</th></thead>`
	tplFooter = `</table></body></html>`
)

func main() {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("can not open file:", err)
	}

	stor := NewStorage(f)
	err = stor.Restore()
	if err != nil {
		log.Fatal("can not restore messages:", err)
	}

	smsru := SmsRu{stor: stor}
	devino := Devino{stor: stor}

	r := chi.NewRouter()
	r.Get("/", indexHandler(stor))
	r.Route("/rest", func(r chi.Router) {
		r.Post("/user/sessionid", devino.UserSessionIdHandler)
		r.Post("/sms/send", devino.SmsSend)
		r.Post("/sms/state", devino.SmsState)
	})
	r.Route("/sms", func(r chi.Router) {
		r.Post("/user/send", smsru.Send)
		r.Post("/user/status", smsru.Status)
	})

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("Web server fail:", err)
	}
}

func indexHandler(stor *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := strings.Builder{}
		b.WriteString(tplHeader)
		for _, msg := range stor.LastMessages(50) {
			b.WriteString("<tr>" + "<td>" + msg.Time.Format("2006-2006-01-02 15:04:05") + "</td>" + "<td>" + msg.From + "</td>" + "<td>" + msg.To + "</td>" + "<td>" + msg.Text + "</td>" + "</tr>")
		}
		b.WriteString(tplFooter)
		_, err := w.Write([]byte(b.String()))
		if err != nil {
			log.Printf("can not send index to client: %v", err)
		}
	}
}
