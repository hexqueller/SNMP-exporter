package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AgentConfig struct {
	Name      string `yaml:"name"`
	Target    string `yaml:"target"`
	Port      uint16 `yaml:"port"`
	Community string `yaml:"community"`
	Version   string `yaml:"version"`
	OID       string `yaml:"oid"`
	DataPort  uint16 `yaml:"data_port"`
}

type Config struct {
	Agents []AgentConfig `yaml:"agents"`
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
