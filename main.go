package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
)

const (
	TLSaddr     = ":443"
	addr        = ":80"
	filename    = "messages.txt"
	helpdekFile = "helpdesk.msgp"

	tplHeader = `
	<html>
		<head>
			<style>
				ol {
					padding: 10px; 
					list-style-type: none;				
				}
				ol li {
					float: left;
					margin: 0 10px 0 0;
				}
			</style>
		</head>
		<body>
			<ol>
				<li><a href="/">sms</a></li>
				<li><a href="/helpdesk">helpdesk</a></li>
			</ol>`

	smsTableFooter = `</table>`

	helpdeskTableFooter = ``

	tplFooter = `</body></html>`
)

func helpdeskTableHeaderWithCount(feedbackCount int) string {
	const helpdeskTableHeader = `
	<table border=1>
		<caption>Helpdesk (%d)</caption>
		<thead>
			<th>Date</th>
			<th>Title</th>
			<th>Description</th>
		</thead>`

	return fmt.Sprintf(helpdeskTableHeader, feedbackCount)
}

func smsTableHeaderWithCount(messageCount int) string {
	const smsTableHeader = `
	<table border=1>
	    <caption>SMS (%d)</caption>
		<thead>
			<th>Date</th>
			<th>From</th>
			<th>Phone</th>
			<th>Msg</th>
			<th>Provider</th>
		</thead>`

	return fmt.Sprintf(smsTableHeader, messageCount)
}

func main() {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("can not open file:", err)
	}
	defer f.Close()

	stor := NewStorage(f)

	err = stor.Restore()
	if err != nil {
		log.Fatal("can not restore SMS messages:", err)
	}

	hdf, err := os.OpenFile(helpdekFile, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("can not open helpdesk file:", err)
	}
	defer hdf.Close()

	hstor := newHelpdeskStorage(hdf)

	err = hstor.Restore()
	if err != nil {
		log.Fatal("can not restore HelpDesk messages:", err)
	}

	smsru := SmsRu{stor: stor}
	devino := Devino{stor: stor}

	go func() {
		httpServer(stor, hstor, smsru)
	}()

	// devino telecom mock server
	r := chi.NewRouter()
	r.Use(caselessMatcher)

	devinoTelecomRoutes(r, devino)
	smsRuRoutes(r, smsru)
	helpdeskRoutes(r, hstor)

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

func helpdeskRoutes(mux *chi.Mux, stor *HelpdeskStorage) {
	mux.Post("/api/v2/tickets/", helpdeskEddyHandler(stor))
}

// sms.ru and stats server
func httpServer(stor *Storage, hstor *HelpdeskStorage, smsru SmsRu) {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(indexPageWrapper)
		r.Get("/", indexSmsHandler(stor))
		r.Get("/helpdesk", helpdeskIndexHandler(hstor))
	})

	smsRuRoutes(r, smsru)
	helpdeskRoutes(r, hstor)

	log.Println("start HTTP on", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Println("Web server fail:", err)
	}
}

func indexSmsHandler(stor *Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := strings.Builder{}
		b.WriteString(smsTableHeaderWithCount(stor.Len()))
		for _, msg := range stor.LastMessages(50) {
			b.WriteString("<tr>" +
				"<td>" + msg.Time.Format("2006-01-02 15:04:05") + "</td>" +
				"<td>" + msg.From + "</td>" +
				"<td>" + msg.To + "</td>" +
				"<td>" + msg.Text + "</td>" +
				"<td>" + msg.Provider + "</td>" +
				"</tr>")
		}
		b.WriteString(smsTableFooter)
		_, err := w.Write([]byte(b.String()))
		if err != nil {
			log.Printf("can not send index to client: %v", err)
		}
	}
}

func helpdeskIndexHandler(stor *HelpdeskStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b := strings.Builder{}
		b.WriteString(helpdeskTableHeaderWithCount(stor.Len()))
		for _, msg := range stor.LastMessages(50) {
			b.WriteString("<tr>" +
				"<td>" + msg.Time.Format("2006-01-02 15:04:05") + "</td>" +
				"<td>" + msg.Title + "</td>" +
				"<td>" + msg.Description + "</td>" +
				"</tr>")
		}
		b.WriteString(helpdeskTableFooter)
		_, err := w.Write([]byte(b.String()))
		if err != nil {
			log.Printf("can not send index to client: %v", err)
		}
	}
}

func indexPageWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, tplHeader)
		next.ServeHTTP(w, r)
		fmt.Fprint(w, tplFooter)
	})
}

// caselessMatcher is convert request path to lowercase
func caselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
