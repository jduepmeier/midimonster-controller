package midimonster

import (
	"context"

	"github.com/rs/zerolog"
)

type ProcessController interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	Status(ctx context.Context) (ProcessStatus, error)
	Cleanup()
}

type NewProcessControllerFunc = func(ctx context.Context, logger zerolog.Logger, config *Config) (ProcessController, error)

type ProcessStatus int

var (
	ProcessControllerConstructors = make(map[string]NewProcessControllerFunc)
)

const (
	ProcessStatusRunning = iota
	ProcessStatusStopped
)
