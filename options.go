package kodama

type Options struct {
	Domain string `short:"d" required:"" env:"KODAMA_DOMAIN" help:"Accepted domain."`
	Port   uint   `short:"p" default:"53" env:"KODAMA_PORT" help:"Listening port."`
}
