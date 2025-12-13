package kodama

import (
	"fmt"

	"github.com/miekg/dns"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	*Options
}

func NewServer(options *Options) *Server {
	server := &Server{
		Options: options,
	}

	return server
}

func (svr *Server) Start() error {
	addr := fmt.Sprintf(":%d", svr.Port)
	handler := &Handler{}
	eg := &errgroup.Group{}

	for _, net := range []string{"udp", "tcp"} {
		eg.Go(func() error {
			s := &dns.Server{Addr: addr, Net: net, Handler: handler}
			defer s.Shutdown()
			return s.ListenAndServe()
		})
	}

	return eg.Wait()
}
