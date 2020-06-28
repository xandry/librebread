package mailserver

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
)

type MailMessage struct {
	Time time.Time
	From string
	To   string
	Data string
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
		messages: []MailMessage{
			{
				Time: time.Now(),
				From: "evilbunny.x@gmail.com",
				To:   "nottrack@yandex.ru",
				Data: "hello",
			},
			{
				Time: time.Now(),
				From: "evilbunny.x@gmail.com",
				To:   "nottrack@yandex.ru",
				Data: "hello",
			},
		},

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

	msgs := make([]MailMessage, 0, len(s.messages))

	for _, msg := range s.messages {
		msgs = append(msgs, msg)
	}

	return msgs
}

func (s *MailStorage) Len() int {
	return len(s.messages)
}
