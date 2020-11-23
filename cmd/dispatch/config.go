package main

import (
	"fmt"

	"github.com/dispatch-simulator/internal/config"
	"github.com/dispatch-simulator/internal/runner"
	"gopkg.in/yaml.v2"
)

type appconfig struct {
	Runner runner.Config `yaml:"runner"`
}

func getDefaultConfig() *appconfig {
	cfg := &appconfig{}
	return cfg
}

func getConfig() (*appconfig, error) {
	cfg := &appconfig{}

	cfgcont, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return cfg, err
	}

	err = yaml.UnmarshalStrict(cfgcont, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
