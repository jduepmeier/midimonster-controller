package midimonster

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	DefaultConfigPaths = []string{
		"config.yaml",
		"~/.config/bpm/config.yaml",
	}
)

type Config struct {
	MidimonsterConfigPath string `yaml:"configPath"`
	UnitName              string `yaml:"unitName"`
	BindAddr              string `yaml:"bind"`
	Port                  uint16 `yaml:"port"`
}

func DefaultConfig() *Config {
	return &Config{
		MidimonsterConfigPath: "/etc/midimonster/midimonster.cfg",
		BindAddr:              "0.0.0.0",
		Port:                  8080,
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

func dumpYaml(path string, obj interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(obj)
}

func expandPath(path string) string {
	if path == "~" {
		path = "$HOME"
	} else if strings.HasPrefix(path, "~/") {
		path = strings.Replace(path, "~/", "$HOME/", 1)
	}
	return os.ExpandEnv(path)
}
