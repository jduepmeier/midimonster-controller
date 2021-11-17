package midimonster

import (
	"context"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/rs/zerolog"
)

type ProcessControllerSystemd struct {
	conn     *dbus.Conn
	unitName string
	logger   zerolog.Logger
}

func NewProcessControllerSystemd(ctx context.Context, logger zerolog.Logger, config *Config) (ProcessController, error) {
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return nil, err
	}
	newLogger := logger.With().Str("module", "systemd").Logger()
	return &ProcessControllerSystemd{
		conn:     conn,
		unitName: config.Systemd.UnitName,
		logger:   newLogger,
	}, nil
}

func init() {
	ProcessControllerConstructors["systemd"] = NewProcessControllerSystemd
}

func (pc *ProcessControllerSystemd) Start(ctx context.Context) error {
	result := make(chan string, 1)
	defer close(result)
	jobId, err := pc.conn.StartUnitContext(ctx, pc.unitName, "replace", result)
	if err != nil {
		return err
	}
	resultString := <-result
	if resultString != "" {
		pc.logger.Info().Msgf("%d: %s", jobId, resultString)
	}
	return nil
}
func (pc *ProcessControllerSystemd) Stop(ctx context.Context) error {
	result := make(chan string, 1)
	defer close(result)
	jobId, err := pc.conn.StopUnitContext(ctx, pc.unitName, "replace", result)
	if err != nil {
		return err
	}
	resultString := <-result
	if resultString != "" {
		pc.logger.Info().Msgf("%d: %s", jobId, resultString)
	}
	return nil
}
func (pc *ProcessControllerSystemd) Restart(ctx context.Context) error {
	result := make(chan string, 1)
	defer close(result)
	jobId, err := pc.conn.ReloadOrRestartUnitContext(ctx, pc.unitName, "replace", result)
	if err != nil {
		return err
	}
	resultString := <-result
	if resultString != "" {
		pc.logger.Info().Msgf("%d: %s", jobId, resultString)
	}
	return nil
}
func (pc *ProcessControllerSystemd) Status(ctx context.Context) (ProcessStatus, error) {
	prop, err := pc.conn.GetServicePropertyContext(ctx, pc.unitName, "ActiveState")
	if err != nil {
		return ProcessStatusStopped, err
	}
	if prop.Value.String() == "running" {
		return ProcessStatusRunning, nil
	} else {
		return ProcessStatusStopped, nil
	}
}

func (pc *ProcessControllerSystemd) Cleanup() {
	pc.conn.Close()
}
