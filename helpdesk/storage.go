package helpdesk

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
)

type HelpdeskMsg struct {
	Time         time.Time
	TypeId       int
	PriorityId   int
	DepartmentId int
	Title        string
	Description  string
}

type HelpdeskStorage struct {
	messagesMu sync.Mutex
	messages   []HelpdeskMsg
	rw         io.ReadWriter
	enc        *codec.Encoder
}

func NewStorage(rw io.ReadWriter) *HelpdeskStorage {
	return &HelpdeskStorage{
		rw: rw,
	}
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
