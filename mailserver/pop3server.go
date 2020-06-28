package mailserver

import (
	"github.com/DevelHell/popgun"
	"github.com/DevelHell/popgun/backends"
)

type PopServer struct {
	addr string
}

func NewPopServer(addr string) *PopServer {
	return &PopServer{
		addr: addr,
	}
}

func (pop *PopServer) ListenAndServe() error {
	cfg := popgun.Config{
		ListenInterface: pop.addr,
	}

	server := popgun.NewServer(cfg, backends.DummyAuthorizator{}, backends.DummyBackend{})
	err := server.Start()
	if err != nil {
		return err
	}

	var waitChan chan struct{}
	<-waitChan

	return nil
}
