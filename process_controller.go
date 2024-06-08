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
	Logs(ctx context.Context, oldest uint64) ([]string, uint64, error)
	Cleanup()
}

type NewProcessControllerFunc = func(ctx context.Context, logger zerolog.Logger, config *Config, logsChannel chan string, statusChannel chan struct{}) (ProcessController, error)

type ProcessStatus int

func (status ProcessStatus) Text() string {
	switch status {
	case ProcessStatusRunning:
		return "running"
	case ProcessStatusStopped:
		return "stopped"
	}
	return ""
}

var (
	ProcessControllerConstructors = make(map[string]NewProcessControllerFunc)
)

const (
	ProcessStatusRunning = iota
	ProcessStatusStopped
)
