package mailserver

import (
	"io"

	"github.com/emersion/go-smtp"
)

type backend struct {
	store *MailStorage
}

type session struct {
	store *MailStorage
}

func (bkd *backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &session{
		store: bkd.store,
	}, nil
}

func (bkd *backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &session{
		store: bkd.store,
	}, nil
}

func (s *session) Mail(from string, opts smtp.MailOptions) error {
	return nil
}

func (s *session) Rcpt(to string) error {
	return nil
}

func (s *session) Data(r io.Reader) error {
	msg, err := messageFromReader(r)
	if err != nil {
		return err
	}

	err = s.store.Push(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
