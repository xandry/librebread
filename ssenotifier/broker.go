package ssenotifier

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/vasyahuyasa/librebread/helpdesk"
	"github.com/vasyahuyasa/librebread/mailserver"
	"github.com/vasyahuyasa/librebread/sms"
)

const (
	eventEmailRecived    = "email"
	eventHelpdeskRecived = "helpdesk"
	eventSmsRecived      = "sms"
)

type Broker struct {
	notifier       chan message
	newClients     chan chan message
	closingClients chan chan message
	clients        map[chan message]struct{}
}

type message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func NewBroker() *Broker {
	broker := &Broker{
		notifier:       make(chan message, 1),
		newClients:     make(chan chan message),
		closingClients: make(chan chan message),
		clients:        map[chan message]struct{}{},
	}

	go broker.listen()

	return broker
}

func (b *Broker) EmailNotify(msg mailserver.MailMessage) {
	b.notify(eventEmailRecived, msg)
}

func (b *Broker) HelpdeskNotify(msg helpdesk.HelpdeskMsg) {
	b.notify(eventHelpdeskRecived, msg)
}

func (b *Broker) SmsNotify(msg sms.Message) {
	b.notify(eventSmsRecived, msg)
}

func (b *Broker) ClientHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			log.Printf("streaming for client %v unsupported", r.RemoteAddr)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		messageChan := make(chan message)
		b.newClients <- messageChan

		defer func() {
			b.closingClients <- messageChan
		}()

		done := r.Context().Done()
		go func() {
			<-done
			b.closingClients <- messageChan
			close(messageChan)
		}()

		for {
			msg, ok := <-messageChan
			if !ok {
				break
			}
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg.String())
			if err != nil {
				log.Printf("can not send event to client %v: %v", r.RemoteAddr, err)
				break
			}
			flusher.Flush()
		}
	})
}

func (b *Broker) notify(event string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("can not marshal event: %v", err)
	}
	b.notifier <- message{
		Event: event,
		Data:  string(jsonData),
	}
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = struct{}{}
			log.Printf("Client added. %d registered clients", len(b.clients))

		case s := <-b.closingClients:
			delete(b.clients, s)
			log.Printf("Removed client. %d registered clients", len(b.clients))

		case event := <-b.notifier:
			for clientMessageChan := range b.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (msg *message) String() string {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("can not marshal broker message: %v", err)
		return "{}"
	}

	return string(data)
}
