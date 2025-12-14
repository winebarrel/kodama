package kodama

type Options struct {
	Domain string `short:"d" required:"" env:"KODAMA_DOMAIN" help:"Accepted domain."`
	Addr   string `short:"a" default:":53" env:"KODAMA_ADDR" help:"Listening address."`
}
