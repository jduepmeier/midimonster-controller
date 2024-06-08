package midimonster

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Midimonster struct {
	Path              string
	CurrentConfig     string
	LastConfig        string
	ProcessController ProcessController
	logger            zerolog.Logger
	logsChannel       chan string
	statusChannel     chan struct{}
}

func NewMidimonster(config *Config, logger zerolog.Logger, logsChannel chan string, statusChannel chan struct{}) (*Midimonster, error) {
	var err error
	ctx := context.Background()
	midi := &Midimonster{
		Path:          config.MidimonsterConfigPath,
		logger:        logger,
		logsChannel:   logsChannel,
		statusChannel: statusChannel,
	}
	constructor, ok := ProcessControllerConstructors[config.ControlType]
	if !ok {
		return nil, fmt.Errorf("cannot create process controller: unknown control type %s", config.ControlType)
	}
	logger.Info().Msgf("using process managment type %s", config.ControlType)
	midi.ProcessController, err = constructor(ctx, midi.logger, config, logsChannel, statusChannel)
	if err != nil {
		return nil, err
	}
	return midi, midi.LoadConfig()
}

func (midi *Midimonster) LoadConfig() error {
	content, err := os.ReadFile(midi.Path)
	if err != nil {
		return err
	}
	midi.CurrentConfig = string(content)
	return nil
}

func (midi *Midimonster) ReplaceConfig(ctx context.Context, content string) error {
	content = content + "\n"
	lastContent, err := os.ReadFile(midi.Path)
	if err != nil {
		return err
	}
	midi.LastConfig = string(lastContent)
	err = os.WriteFile(midi.Path, []byte(content), 0644)
	if err != nil {
		return err
	}
	midi.CurrentConfig = content
	return midi.Restart(ctx)
}

func (midi *Midimonster) Restart(ctx context.Context) error {
	return midi.ProcessController.Restart(ctx)
}
