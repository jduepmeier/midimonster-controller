package midimonster

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	DefaultConfigPaths = []string{
		"config.yaml",
		"~/.config/midimonster-controller/config.yaml",
	}
)

type WebsocketConfig struct {
	LoopDuration time.Duration `yaml:"loopDuration"`
}

type Config struct {
	MidimonsterConfigPath string        `yaml:"configPath"`
	BindAddr              string        `yaml:"bind"`
	Port                  uint16        `yaml:"port"`
	Systemd               ConfigSystemd `yaml:"systemd"`
	Process               ConfigProcess `yaml:"process"`
	ControlType           string        `yaml:"controlType"`
	Development           bool
	Websocket             WebsocketConfig `yaml:"websocket"`
}

type ConfigSystemd struct {
	UnitName string `yaml:"unitName"`
}

type ConfigProcess struct {
	BinPath string   `yaml:"binPath"`
	WorkDir string   `yaml:"workDir"`
	Args    []string `yaml:"args"`
}

func DefaultConfig() *Config {
	return &Config{
		MidimonsterConfigPath: "/etc/midimonster/midimonster.cfg",
		BindAddr:              "0.0.0.0",
		Port:                  8080,
		ControlType:           "systemd",
		Systemd: ConfigSystemd{
			UnitName: "midimonster.service",
		},
		Process: ConfigProcess{
			Args: []string{},
		},
		Websocket: WebsocketConfig{
			LoopDuration: 5 * time.Second,
		},
	}
}

func ReadConfig(path string) (*Config, error) {
	config := DefaultConfig()
	if path != "" {
		err := loadYaml(path, &config)
		if err != nil {
			return config, fmt.Errorf("cannot %s", err)
		}
	} else {
		for _, path := range DefaultConfigPaths {
			err := loadYaml(expandPath(path), &config)
			if err == nil {
				break
			}
		}
	}

	config.MidimonsterConfigPath = expandPath(config.MidimonsterConfigPath)
	config.Process.BinPath = expandPath(config.Process.BinPath)

	return config, nil
}

func loadYaml(path string, obj interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(obj)
}

func expandPath(path string) string {
	if path == "~" {
		path = "$HOME"
	} else if strings.HasPrefix(path, "~/") {
		path = strings.Replace(path, "~/", "$HOME/", 1)
	}
	return os.ExpandEnv(path)
}
