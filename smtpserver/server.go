package smtpserver

import (
	"errors"
	"io"
	"io/ioutil"
	"log"

	"github.com/emersion/go-smtp"
)

type Server struct {
}

type session struct{}

func (svr *Server) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username != "username" || password != "password" {
		return nil, errors.New("Invalid username or password")
	}
	return &session{}, nil
}

func (s *session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
	}
	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
