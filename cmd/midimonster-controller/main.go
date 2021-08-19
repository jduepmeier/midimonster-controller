package main

import (
	"midimonster"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

type opts struct {
	ConfigPath string `short:"c" long:"config" description:"path to config path"`
}

func main() {
	opts := opts{}
	logger := zerolog.New(os.Stderr)
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	config, err := midimonster.ReadConfig(opts.ConfigPath)
	if err != nil {
		logger.Err(err).Msg("cannot read config")
	}

	controller, err := midimonster.NewController(config)
	if err != nil {
		logger.Err(err).Msg("cannot create controller")
	}
	err = controller.Serve()
	if err != nil {
		logger.Err(err).Msg("cannot serve")
	}
}
