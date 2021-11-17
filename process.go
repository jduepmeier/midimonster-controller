package midimonster

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

type ProcessControllerProcess struct {
	cmd          *exec.Cmd
	ExecPath     string
	ConfigPath   string
	WorkDir      string
	stderr       io.ReadCloser
	stdout       io.ReadCloser
	logger       zerolog.Logger
	runningMutex sync.Mutex
}

func NewProcessControllerProcess(ctx context.Context, logger zerolog.Logger, config *Config) (ProcessController, error) {
	newLogger := logger.With().Str("process-controller", "process").Logger()
	newLogger.Info().Msg("init")
	if config.Process.WorkDir == "" {
		config.Process.WorkDir = path.Dir(config.Process.BinPath)
	}
	configPath, err := filepath.Abs(config.MidimonsterConfigPath)
	if err != nil {
		newLogger.Err(err).Msgf("cannot make path %s absolute", config.MidimonsterConfigPath)
		configPath = config.MidimonsterConfigPath
	}
	return &ProcessControllerProcess{
		ExecPath:   config.Process.BinPath,
		ConfigPath: configPath,
		WorkDir:    config.Process.WorkDir,
		logger:     newLogger,
	}, nil
}

func init() {
	ProcessControllerConstructors["process"] = NewProcessControllerProcess
}

func (pc *ProcessControllerProcess) Start(ctx context.Context) (err error) {
	pc.logger.Info().Msgf("start midimonster (%s  %s)", pc.ExecPath, pc.ConfigPath)
	backgroundCtx := context.Background()
	pc.cmd = exec.CommandContext(backgroundCtx, pc.ExecPath, pc.ConfigPath)
	pc.cmd.Dir = pc.WorkDir
	pc.stderr, err = pc.cmd.StderrPipe()
	if err != nil {
		return err
	}
	pc.stdout, err = pc.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go pc.startReader(pc.stdout, "stdout", &wg)
	go pc.startReader(pc.stderr, "stderr", &wg)
	err = pc.cmd.Start()
	if err != nil {
		pc.logger.Err(err).Msg("could not start midimonster")
	}
	go pc.waitForExit(&wg)
	return err
}

func (pc *ProcessControllerProcess) waitForExit(wg *sync.WaitGroup) {
	pc.runningMutex.Lock()
	defer pc.runningMutex.Unlock()
	pc.logger.Debug().Msgf("midimonster pid: %d", pc.cmd.Process.Pid)
	wg.Wait()
	pc.cmd.Wait()
	pc.logger.Info().Msgf("midimonster exit code %d", pc.cmd.ProcessState.ExitCode())
	pc.cmd = nil
}

func (pc *ProcessControllerProcess) startReader(reader io.ReadCloser, id string, wg *sync.WaitGroup) {
	scanner := bufio.NewReader(reader)
	var err error
	var line []byte
	for {
		line, _, err = scanner.ReadLine()
		if err != nil {
			break
		}
		pc.logger.Debug().Msgf("midimonster (%s): %s", id, line)
	}
	if err != nil {
		pc.logger.Err(err).Msgf("error scanning %s", id)
	}
	pc.logger.Debug().Msgf("finished reading %s", id)
	reader.Close()
	wg.Done()
}

func (pc *ProcessControllerProcess) Stop(ctx context.Context) error {
	if pc.cmd != nil && pc.cmd.Process != nil {
		pc.logger.Info().Msg("stop midimonster")
		err := pc.cmd.Process.Signal(os.Interrupt)
		if err != nil {
			return err
		}
		pc.runningMutex.Lock()
		defer pc.runningMutex.Unlock()
		return err
	} else {
		pc.logger.Info().Msg("midimonster is not running")
	}
	return nil
}

func (pc *ProcessControllerProcess) Restart(ctx context.Context) error {
	err := pc.Stop(ctx)
	if err != nil {
		return err
	}
	return pc.Start(ctx)
}
func (pc *ProcessControllerProcess) Status(ctx context.Context) (ProcessStatus, error) {
	if pc.cmd == nil {
		return ProcessStatusStopped, nil
	}
	if pc.cmd.Process != nil {
		pc.logger.Debug().Msgf("midimonster pid is %d", pc.cmd.Process.Pid)
	}
	if pc.cmd.ProcessState != nil && pc.cmd.ProcessState.Exited() {
		return ProcessStatusStopped, nil
	}
	return ProcessStatusRunning, nil
}

func (pc *ProcessControllerProcess) Cleanup() {
	_ = pc.Stop(context.Background())
}
