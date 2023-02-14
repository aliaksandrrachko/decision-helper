package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DATASOURCE_URL_ENV_NAME = "DATASOURCE_URL"
)

type Config struct {
	DataSource struct {
		Url        string `yaml:"url"`
		DriverName string `yaml:"driver-name"`
	} `yaml:"datasource"`

	Server struct {
		Host           string   `yaml:"host"`
		Port           string   `yaml:"port"`
		TrustedProxies []string `yaml:"trusted-proxies"`
	} `yaml:"server"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	datasourceEnv := os.Getenv(DATASOURCE_URL_ENV_NAME)
	if datasourceEnv != "" {
		config.DataSource.Url = datasourceEnv
	}

	return config, nil
}
