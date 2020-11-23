package main

import (
	"errors"
	"os"
	"strings"

	"github.com/dispatch-simulator/internal/config"
	"github.com/dispatch-simulator/internal/runner"
	"go.melnyk.org/mlog"
	"go.melnyk.org/mlog/console"
	"go.melnyk.org/mlog/nolog"
	"gopkg.in/yaml.v2"
)

const (
	currentConfigVersion = 1
)

type appconfig struct {
	Version int           `yaml:"version"`
	Runner  runner.Config `yaml:"runner"`
	Logger  struct {
		Type  string `yaml:"type"`
		Level string `yaml:"level"`
	} `yaml:"logger"`

	// Internal part
	logbook mlog.Logbook
}

func getDefaultConfig() *appconfig {
	cfg := &appconfig{Version: currentConfigVersion}
	cfg.Logger.Type = "console"
	cfg.Logger.Level = "info"
	return cfg
}

func getConfig() (*appconfig, error) {
	cfg := &appconfig{}

	cfgcont, err := config.GetConfig()
	if err != nil {
		return cfg, err
	}

	err = yaml.UnmarshalStrict(cfgcont, &cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.Version > currentConfigVersion {
		return cfg, errors.New("Config version is newer than supported")
	}

	// Do all required configuration stuff
	switch cfg.Logger.Type {
	case "console":
		cfg.logbook = console.NewLogbook(os.Stdout)
	case "nolog":
		cfg.logbook = nolog.NewLogbook()
	default:
		return cfg, errors.New("Only console or nolog is supported as possible logger type")
	}

	if cfg.Logger.Level != "" {
		if level, err := leveTolevel(cfg.Logger.Level); err == nil {
			cfg.logbook.SetLevel(mlog.Default, level)
		} else {
			return cfg, err
		}
	}

	return cfg, err
}

func (cfg *appconfig) cleanup() {
	// Do config cleanup here
}

func leveTolevel(level string) (mlog.Level, error) {
	switch strings.ToLower(level) {
	case "verbose":
		return mlog.Verbose, nil
	case "info":
		return mlog.Info, nil
	case "warning":
		return mlog.Warning, nil
	case "error":
		return mlog.Error, nil
	case "fatal":
		return mlog.Fatal, nil
	default:
		return mlog.Fatal, errors.New("Unknown log level:" + level)
	}
}
