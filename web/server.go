package web

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vasyahuyasa/librebread/sms"
)

type Server struct {
	//TLSServer bool
	//TLSCert   string
	//TLSKey    string

	smsStore *sms.SqliteStorage
	http     *http.Server
	handler  http.Handler
	re       *renderer
}

func NewServer(smsStore *sms.SqliteStorage) *Server {

	s := &Server{
		smsStore: smsStore,
		re:       newRenderer(),
	}

	s.initRoutes()

	s.http = &http.Server{Handler: s.handler}

	return s
}

func (s *Server) ListenAndServe() error {

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) initRoutes() {
	r := chi.NewMux()

	r.Get("/", redirect("/sms", http.StatusTemporaryRedirect))
	r.Get("/sms", smsIndexHandler(s.smsStore, s.re))

	s.handler = r
}
