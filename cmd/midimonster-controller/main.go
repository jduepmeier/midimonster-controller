package main

import (
	"fmt"
	"midimonster"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

var build = "dev"

type opts struct {
	ConfigPath  string `short:"c" long:"config" description:"path to config path"`
	LogLevel    string `short:"l" long:"loglevel" description:"loglevel"`
	Development bool   `short:"d" long:"dev" description:"start in dev mode"`
	Version     bool   `short:"v" long:"version" description:"show version"`
}

func main() {
	opts := opts{
		LogLevel: "warn",
	}
	logger := zerolog.New(os.Stderr)
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	if opts.Version {
		fmt.Printf("%s - %s\n", os.Args[0], build)
		return
	}

	level, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil {
		logger.Err(err).Msg("cannot parse loglevel")
		return
	}
	zerolog.SetGlobalLevel(level)

	config, err := midimonster.ReadConfig(opts.ConfigPath)
	if err != nil {
		logger.Err(err).Msg("cannot read config")
		return
	}
	config.Development = opts.Development

	controller, err := midimonster.NewController(config, logger)
	if err != nil {
		logger.Err(err).Msg("cannot create controller")
		return
	}
	err = controller.Serve()
	if err != nil {
		logger.Err(err).Msg("cannot serve")
	}
}
