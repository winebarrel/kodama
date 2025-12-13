package kodama

type Options struct {
	Domain string
	Port   uint `short:"p" default:"53" env:"KODAMA_PORT" help:"Listening port."`
}
