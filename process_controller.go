package midimonster

import "context"

type ProcessController interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	Status(ctx context.Context) (ProcessStatus, error)
	Cleanup()
}

type ProcessStatus int

const (
	ProcessStatusRunning = iota
	ProcessStatusStopped
)
