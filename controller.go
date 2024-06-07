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
		logger: logger.With().Str("component", "controller").Logger(),
	}
	controller.server = NewServer(config, controller, logger)
	controller.Midimonster, err = NewMidimonster(config, logger, controller.server.logsChannel)

	return controller, err
}

func (controller *Controller) Serve() error {
	controller.logger.Info().Msgf("start server")
	return controller.server.Start()
}
