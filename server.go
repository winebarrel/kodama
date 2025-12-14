package kodama

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)

	{
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stop
			cancel()
		}()
	}

	for _, net := range []string{"udp", "tcp"} {
		eg.Go(func() error {
			s := &dns.Server{Addr: addr, Net: net, Handler: handler}
			context.AfterFunc(ctx, func() {
				s.Shutdown()
			})
			return s.ListenAndServe()
		})
	}

	return eg.Wait()
}
