package midimonster

import (
	"context"
	"io/ioutil"

	"github.com/rs/zerolog"
)

type Midimonster struct {
	Path              string
	CurrentConfig     string
	LastConfig        string
	ProcessController ProcessController
	logger            zerolog.Logger
}

func NewMidimonster(config *Config, logger zerolog.Logger) (*Midimonster, error) {
	var err error
	ctx := context.Background()
	midi := &Midimonster{
		Path: config.MidimonsterConfigPath,
	}
	midi.ProcessController, err = NewProcessControllerSystemd(ctx, midi.logger, config.UnitName)
	if err != nil {
		return nil, err
	}
	return midi, midi.LoadConfig()
}

func (midi *Midimonster) LoadConfig() error {
	content, err := ioutil.ReadFile(midi.Path)
	if err != nil {
		return err
	}
	midi.CurrentConfig = string(content)
	return nil
}

func (midi *Midimonster) ReplaceConfig(ctx context.Context, content string) error {
	lastContent, err := ioutil.ReadFile(midi.Path)
	if err != nil {
		return err
	}
	midi.LastConfig = string(lastContent)
	err = ioutil.WriteFile(midi.Path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return midi.Restart(ctx)
}

func (midi *Midimonster) Restart(ctx context.Context) error {
	return midi.ProcessController.Restart(ctx)
}
