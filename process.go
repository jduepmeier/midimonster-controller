package midimonster

import (
	"context"
	"os"
	"os/exec"
)

type ProcessControllerProcess struct {
	cmd        *exec.Cmd
	ExecPath   string
	ConfigPath string
}

func NewProcessControllerProcess(ctx context.Context, execPath string, configPath string) (*ProcessControllerProcess, error) {
	return &ProcessControllerProcess{
		ExecPath:   execPath,
		ConfigPath: configPath,
	}, nil
}

func (pc *ProcessControllerProcess) Start(ctx context.Context) error {
	pc.cmd = exec.CommandContext(ctx, pc.ExecPath, "-c", pc.ConfigPath)
	return pc.cmd.Start()
}
func (pc *ProcessControllerProcess) Stop(ctx context.Context) error {
	if pc.cmd != nil && pc.cmd.Process != nil {
		err := pc.cmd.Process.Signal(os.Interrupt)
		if err != nil {
			return err
		}
		err = pc.cmd.Wait()
		pc.cmd = nil
		return err
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
	if pc.cmd == nil || pc.cmd.ProcessState.Exited() {
		return ProcessStatusStopped, nil
	}
	return ProcessStatusRunning, nil
}

func (pc *ProcessControllerProcess) Cleanup() {
	_ = pc.Stop(context.Background())
}
