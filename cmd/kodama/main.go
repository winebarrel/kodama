package main

import (
	"context"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/kodama"
)

var (
	version string
)

func parseArgs() *kodama.Options {
	var cli struct {
		kodama.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	options.Version = version
	server := kodama.NewServer(options)
	err := server.Start(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
