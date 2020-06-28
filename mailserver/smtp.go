package mailserver

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/emersion/go-smtp"
)

type backend struct {
}

type session struct{}

func (bkd *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &session{}, nil
}

func (bkd *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
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
