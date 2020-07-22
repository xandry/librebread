package web

import (
	"context"
	"net/http"
)

type Server struct {
	//TLSServer bool
	//TLSCert   string
	//TLSKey    string

	http *http.Server
}

func NewServer(h http.Handler) *Server {

	return &Server{
		http: &http.Server{Handler: h},
	}
}

func (s *Server) ListenAndServe() error {

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
