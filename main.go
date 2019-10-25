package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
)

const (
	TLSaddr   = ":443"
	addr      = ":80"
	filename  = "messages.txt"
	tplHeader = `<html><body><table border=1><thead><th>Date</th><th>From</th><th>Phone</th><th>Msg</th><th>Provider</th></thead>`
	tplFooter = `</table></body></html>`
)

func main() {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("can not open file:", err)
	}
	defer f.Close()

	stor := NewStorage(f)

	err = stor.Restore()
	if err != nil {
		log.Fatal("can not restore messages:", err)
	}

	smsru := SmsRu{stor: stor}
	devino := Devino{stor: stor}

	go func() {
		httpServer(stor, smsru)
	}()

	// devino telecom mock server
	r := chi.NewRouter()
	r.Use(caselessMatcher)

	devinoTelecomRoutes(r, devino)
	smsRuRoutes(r, smsru)

	log.Println("start HTTPS on", TLSaddr)
	err = http.ListenAndServeTLS(TLSaddr, "cert/server.crt", "cert/server.key", r)
	if err != nil {
		log.Println("TLS Web server fail:", err)
	}
}

func devinoTelecomRoutes(r chi.Router, devino Devino) {
	r.Route("/rest", func(r chi.Router) {
		r.Post("/user/sessionid", devino.UserSessionIdHandler)
		r.Post("/sms/send", devino.SmsSend)
		r.Post("/sms/state", devino.SmsState)
	})

	r.Route("/rest/v2", func(r chi.Router) {
		r.Post("/sms/send", devino.SmsSend)
		r.Post("/sms/state", devino.SmsState)
	})
}

func smsRuRoutes(mux *chi.Mux, smsru SmsRu) {
	mux.Route("/sms", func(r chi.Router) {
		r.Post("/send", smsru.Send)
		r.Post("/status", smsru.Status)
	})
}

// sms.ru and stats server
func httpServer(stor *Storage, smsru SmsRu) {
	r := chi.NewRouter()
	r.Get("/", indexHandler(stor))
	smsRuRoutes(r, smsru)
	log.Println("start HTTP on", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Println("Web server fail:", err)
	}
}

func indexHandler(stor *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := strings.Builder{}
		b.WriteString(tplHeader)
		for _, msg := range stor.LastMessages(50) {
			b.WriteString("<tr>" +
				"<td>" + msg.Time.Format("2006-2006-01-02 15:04:05") + "</td>" +
				"<td>" + msg.From + "</td>" +
				"<td>" + msg.To + "</td>" +
				"<td>" + msg.Text + "</td>" +
				"<td>" + msg.Provider + "</td>" +
				"</tr>")
		}
		b.WriteString(tplFooter)
		_, err := w.Write([]byte(b.String()))
		if err != nil {
			log.Printf("can not send index to client: %v", err)
		}
	}
}

// caselessMatcher is convert request path to lowercase
func caselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
