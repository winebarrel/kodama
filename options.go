package kodama

import (
	"fmt"
	"net"
)

type Options struct {
	Domain  string            `arg:"" required:"" env:"KODAMA_DOMAIN" help:"Accepted domain."`
	NS      map[string]string `env:"KODAMA_NS" help:"NS records. (e.g., ns.example.com=203.0.113.0)"`
	TXT     map[string]string `env:"KODAMA_TXT" help:"TXT records. (e.g., spf.example.com=...)"`
	CNAME   map[string]string `env:"KODAMA_CNAME" help:"CNAME records. (e.g., alias.example.com=...)"`
	Addr    string            `default:":53" env:"KODAMA_ADDR" help:"Listening address."`
	TTL     uint32            `default:"300" env:"KODAMA_TTL" help:"Dynamic Record TTLL."`
	Version string            `kong:"-"`
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
