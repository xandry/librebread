package mailserver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/textproto"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
)

type MailMessage struct {
	ID        int
	MessageID string
	SentOn    time.Time
	RecivedOn time.Time
	From      string
	To        string
	Subject   string
	Body      string
	Data      string
}

type MailStorage struct {
	messagesMu *sync.RWMutex
	messages   []MailMessage

	rw  io.ReadWriter
	enc *codec.Encoder
}

func NewStorage(rw io.ReadWriter) *MailStorage {
	return &MailStorage{
		messagesMu: &sync.RWMutex{},
		messages:   []MailMessage{},

		rw: rw,
	}
}

func (s *MailStorage) Restore() error {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	dec := codec.NewDecoder(s.rw, &codec.MsgpackHandle{})

	for {
		var msg MailMessage

		err := dec.Decode(&msg)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		s.messages = append(s.messages, msg)
	}
}

func (s *MailStorage) Push(msg MailMessage) error {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	err := s.write(msg)
	if err != nil {
		log.Println("can not write message to stream:", err)
		return err
	}

	s.messages = append(s.messages, msg)

	return nil
}

func (s *MailStorage) write(msg MailMessage) error {
	if s.enc == nil {
		s.enc = codec.NewEncoder(s.rw, &codec.MsgpackHandle{})
	}

	return s.enc.Encode(msg)
}

func (s *MailStorage) LastMessages() []MailMessage {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	msgs := append([]MailMessage{}, s.messages...)

	return msgs
}

func (s *MailStorage) Len() int {
	return len(s.messages)
}

func messageFromReader(r io.Reader) (MailMessage, error) {
	tr := textproto.NewReader(bufio.NewReader(r))

	headers, err := tr.ReadMIMEHeader()
	if err != nil {
		return MailMessage{}, err
	}

	msg := MailMessage{
		SentOn:    time.Now(),
		RecivedOn: time.Now(),
		MessageID: fmt.Sprintf("%d.%d@%d", time.Now().Unix(), rand.Int(), rand.Int()),
	}

	for header, v := range headers {
		switch header {
		case "Subject":
			msg.Subject = v[0]
		case "Message-Id":
			msg.MessageID = v[0]
		case "Date":
			t, err := time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", v[0])
			if err != nil {
				return MailMessage{}, err
			}

			msg.SentOn = t
		case "From":
			msg.From = v[0]
		case "To":
			msg.To = v[0]
		}
	}

	for {
		l, err := tr.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return MailMessage{}, err
		}

		msg.Body = msg.Body + "\n" + l
	}

	return msg, nil
}
