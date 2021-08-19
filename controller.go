package midimonster

import "github.com/rs/zerolog"

type Controller struct {
	Config      *Config
	server      *Server
	Midimonster *Midimonster
	logger      zerolog.Logger
}

func NewController(config *Config, logger zerolog.Logger) (*Controller, error) {
	var err error
	controller := &Controller{
		Config: config,
	}
	controller.Midimonster, err = NewMidimonster(config, logger)
	controller.server = NewServer(config, controller)

	return controller, err
}

func (controller *Controller) Serve() error {
	controller.logger.Info().Msgf("start server")
	return controller.server.Start()
}
