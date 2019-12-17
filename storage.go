package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
)

type Message struct {
	Time     time.Time
	From     string
	To       string
	Text     string
	Provider string
}

type HelpdeskMsg struct {
	Time         time.Time
	TypeId       int
	PriorityId   int
	DepartmentId int
	Title        string
	Description  string
}

type Storage struct {
	messages []Message
	rw       io.ReadWriter
}

type HelpdeskStorage struct {
	messagesMu sync.Mutex
	messages   []HelpdeskMsg
	rw         io.ReadWriter
	enc        *codec.Encoder
}

func NewStorage(rw io.ReadWriter) *Storage {
	return &Storage{
		rw: rw,
	}
}

func newHelpdeskStorage(rw io.ReadWriter) *HelpdeskStorage {
	return &HelpdeskStorage{
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

func (s *HelpdeskStorage) Write(msg HelpdeskMsg) error {
	if s.enc == nil {
		s.enc = codec.NewEncoder(s.rw, &codec.MsgpackHandle{})
	}

	return s.enc.Encode(msg)
}

func (s *HelpdeskStorage) Restore() error {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	dec := codec.NewDecoder(s.rw, &codec.MsgpackHandle{})

	for {
		var msg HelpdeskMsg
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

func (s *HelpdeskStorage) LastMessages(n int) []HelpdeskMsg {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	msgs := make([]HelpdeskMsg, 0, n)
	max := len(s.messages) - 1

	for i := 0; max-i >= 0 && i < n; i++ {
		msgs = append(msgs, s.messages[max-i])
	}

	return msgs
}

func (s *HelpdeskStorage) Push(msg HelpdeskMsg) {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	err := s.Write(msg)
	if err != nil {
		log.Println("can not save message to log:", err)
		return
	}

	s.messages = append(s.messages, msg)

}

func (s *HelpdeskStorage) Len() int {
	return len(s.messages)
}
