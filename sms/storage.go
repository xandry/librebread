package sms

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

type Message struct {
	Time     time.Time
	From     string
	To       string
	Text     string
	Provider string
}

type Storage struct {
	messages []Message
	rw       io.ReadWriter
}

func NewStorage(rw io.ReadWriter) *Storage {
	return &Storage{
		rw: rw,
	}
}

func (s *Storage) Write(msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(s.rw, string(b))
	return err
}

func (s *Storage) Restore() error {
	scanner := bufio.NewScanner(s.rw)

	for scanner.Scan() {
		ln := scanner.Text()
		if ln == "" {
			return nil
		}
		var msg Message
		err := json.Unmarshal([]byte(ln), &msg)
		if err != nil {
			return err
		}
		s.messages = append(s.messages, msg)
	}

	return nil
}

func (s *Storage) LastMessages(n int) []Message {
	msgs := make([]Message, 0, n)
	max := len(s.messages) - 1

	for i := 0; max-i >= 0 && i < n; i++ {
		msgs = append(msgs, s.messages[max-i])
	}

	return msgs
}

func (s *Storage) Push(msg Message) {
	err := s.Write(msg)
	if err != nil {
		log.Println("can not save message to log:", err)
		return
	}

	s.messages = append(s.messages, msg)

}

func (s *Storage) Len() int {
	return len(s.messages)
}
