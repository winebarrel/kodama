package kodama

import (
	"fmt"
)

type Options struct {
	Domain   string `arg:"" required:"" env:"KODAMA_DOMAIN" help:"Accepted domain."`
	Addr     string `default:":53" env:"KODAMA_ADDR" help:"Listening address."`
	TTL      uint32 `default:"300" env:"KODAMA_TTL" help:"Dynamic Record TTL."`
	ZoneData string `env:"KODAMA_ZONE_DATA" help:"Zone file data."`
	Zone     *Zone  `kong:"-"`
}

func (options *Options) AfterApply() error {
	zone, err := NewZone(options.ZoneData)

	if err != nil {
		return fmt.Errorf("failed to parse zone file data: %w", err)
	}

	options.Zone = zone

	return nil
}
