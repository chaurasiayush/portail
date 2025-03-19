package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Forwards []ForwardRule `yaml:"forwards"`
}

type ForwardRule struct {
	Listen   string      `yaml:"listen"`
	Forward  string      `yaml:"forward"`
	Protocol string      `yaml:"protocol"` // tcp or udp
	TLS      *TLSOptions `yaml:"tls,omitempty"`
}

type TLSOptions struct {
	Enabled    bool `yaml:"enabled"`
	SkipVerify bool `yaml:"skip_verify"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
