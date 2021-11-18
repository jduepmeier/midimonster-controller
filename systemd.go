// systemd is only useful on linux
//go:build linux && !nosystemd

package midimonster

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/rs/zerolog"
)

type ProcessControllerSystemd struct {
	conn            *dbus.Conn
	unitName        string
	logger          zerolog.Logger
	logs            *RingBuffer
	journalReader   *sdjournal.JournalReader
	journalWait     sync.WaitGroup
	journalStopChan chan time.Time
}

func NewProcessControllerSystemd(ctx context.Context, logger zerolog.Logger, config *Config) (ProcessController, error) {
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return nil, err
	}
	newLogger := logger.With().Str("module", "systemd").Logger()

	journalConfig := sdjournal.JournalReaderConfig{
		NumFromTail: 0,
		Matches: []sdjournal.Match{
			{
				Field: sdjournal.SD_JOURNAL_FIELD_SYSTEMD_UNIT,
				Value: config.Systemd.UnitName,
			},
		},
	}
	journal, err := sdjournal.NewJournalReader(journalConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create journal reader: %w", err)
	}

	pc := &ProcessControllerSystemd{
		conn:            conn,
		unitName:        config.Systemd.UnitName,
		logger:          newLogger,
		logs:            NewRingBuffer(1024),
		journalReader:   journal,
		journalStopChan: make(chan time.Time, 1),
	}
	pc.journalWait.Add(1)
	go func() {
		defer pc.journalWait.Done()
		pc.startJournalWatcher()
	}()
	return pc, nil
}

func (pc *ProcessControllerSystemd) startJournalWatcher() {
	var buf bytes.Buffer
	go pc.journalReader.Follow(pc.journalStopChan, &buf)
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		text := scanner.Text()
		pc.logs.Append(text)
	}
	if scanner.Err() != nil {
		pc.logger.Err(scanner.Err()).Msgf("journal scanner finished with error")
	}
	pc.logger.Debug().Msg("finished scanning journal")
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
	prop, err := pc.conn.GetUnitPropertyContext(ctx, pc.unitName, "ActiveState")
	if err != nil {
		return ProcessStatusStopped, err
	}
	pc.logger.Debug().Msgf("unit is %s", prop.Value.String())
	if prop.Value.String() == "\"active\"" {
		return ProcessStatusRunning, nil
	} else {
		return ProcessStatusStopped, nil
	}
}

func (pc *ProcessControllerSystemd) Cleanup() {
	pc.conn.Close()
	pc.journalStopChan <- time.Now()
	pc.journalWait.Done()
}

func (pc *ProcessControllerSystemd) Logs(ctx context.Context, oldest uint64) ([]string, uint64, error) {
	return pc.logs.GetFromOldest(oldest), pc.logs.Newest(), nil
}
