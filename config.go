package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	LogThreshold    = 100
	DefaultInterval = 500 * time.Millisecond
)

type Config struct {
	Name string
	// Where cho should tail logs from. Since we're using fsnotify to tail
	// the logSource
	LogSource string

	// Ingestor http address
	IngestorAddr string

	// interval delay for sending logs to ingestor. Server and client can neogotiate or
	// use a default
	Interval time.Duration

	LogThreshold int
}

func (cfg *Config) String() string {
	return fmt.Sprintf("Config: {name: %s, logSource: %s ingestorAddr: %s, interval: %+v, threshold: %d}",
		cfg.Name, cfg.LogSource, cfg.IngestorAddr, cfg.Interval.String(), cfg.LogThreshold,
	)
}

func parseConfig(configFile string) (*Config, error) {
	cfg := defaultConfig()

	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	data, err := toml.Decode(string(content), cfg)
	if err != nil {
		return nil, fmt.Errorf("could not decode toml file. Reason: %w", err)
	}

	_ = data
	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		// used as test-logs for now
		Name:         "cho",
		LogSource:    "./logs/test_logs.txt",
		IngestorAddr: "http://localhost:9082",
		Interval:     DefaultInterval,
		LogThreshold: LogThreshold,
	}
}
