package kodama

import (
	"fmt"
	"net"
)

type Options struct {
	Domain string            `arg:"" required:"" env:"KODAMA_DOMAIN" help:"Accepted domain."`
	NS     map[string]string `env:"KODAMA_NS" help:"NS record. (e.g., ns.example.com=203.0.113.0)"`
	Addr   string            `default:":53" env:"KODAMA_ADDR" help:"Listening address."`
	TTL    uint32            `default:"300" env:"KODAMA_TTL" help:"Record TTL."`
}

func (options *Options) AfterApply() error {
	for _, ipstr := range options.NS {
		ip := net.ParseIP(ipstr)

		if ip == nil {
			return fmt.Errorf("invalid IP address: %s", ipstr)
		}
	}
	return nil
}
