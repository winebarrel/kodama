package kodama

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	TCP *dns.Server
	UDP *dns.Server
}

func NewServer(options *Options) *Server {
	dns.HandleFunc(options.Domain, ServeDNS)

	server := &Server{
		TCP: &dns.Server{Addr: options.Addr, Net: "tcp"},
		UDP: &dns.Server{Addr: options.Addr, Net: "udp"},
	}

	return server
}

func (svr *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)

	{
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stop
			cancel()
		}()
	}

	for _, s := range []*dns.Server{svr.TCP, svr.UDP} {
		eg.Go(func() error {
			context.AfterFunc(ctx, func() {
				s.Shutdown()
			})
			return s.ListenAndServe()
		})
	}

	return eg.Wait()
}
