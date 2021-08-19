package midimonster

type Controller struct {
	Config      *Config
	server      *Server
	Midimonster *Midimonster
}

func NewController(config *Config) (*Controller, error) {
	var err error
	controller := &Controller{
		Config: config,
	}
	controller.Midimonster, err = NewMidimonster(config)
	controller.server = NewServer(config, controller)

	return controller, err
}

func (controller *Controller) Serve() error {
	return controller.server.Start()
}
