package main

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"time"
)

type Message struct {
	Time time.Time
	From string
	To   string
	Text string
}

type Storage struct {
	messages []Message
	file     *os.File
}

func NewStorage(f *os.File) *Storage {
	return &Storage{
		file: f,
	}
}

func (s *Storage) Write(msg Message) error {
	return gob.NewEncoder(s.file).Encode(msg)
}

func (s *Storage) Restore() error {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(s.file)
	for {
		var msg Message
		err = dec.Decode(&msg)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func (s *Storage) LastMessages(n int) []Message {
	msgs := make([]Message, 0, n)
	max := len(s.messages) - 1

	for i := 0; max-i >= 0 && i >= n; i++ {
		msgs = append(msgs, s.messages[max-i])
	}

	return msgs
}

func (s *Storage) Push(msg Message) {
	err := s.Write(msg)
	if err != nil {
		log.Println("can not save messages:", err)
	}
	s.messages = append(s.messages, msg)

}
