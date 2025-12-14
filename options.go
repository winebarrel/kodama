package kodama

type Options struct {
	Domain string `arg:"" required:"" env:"KODAMA_DOMAIN" help:"Accepted domain."`
	NS     string `env:"KODAMA_NS" help:"NS record."`
	Addr   string `default:":53" env:"KODAMA_ADDR" help:"Listening address."`
}
