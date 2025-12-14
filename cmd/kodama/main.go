package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/kodama"
)

var (
	version string
)

func init() {
	if levelStr := os.Getenv("KODAMA_LOG_LEVEL"); levelStr != "" {
		var level slog.Level
		if err := level.UnmarshalText([]byte(levelStr)); err == nil {
			slog.SetLogLoggerLevel(level)
		}
	}
}

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
		slog.Error(err.Error())
		os.Exit(1)
	}
}
