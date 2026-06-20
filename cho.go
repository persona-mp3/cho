package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Cho struct {
	token        string
	tailedLogs   []Log
	ingestorAddr string
	// File is opened for readOnly access
	source   *os.File
	interval time.Duration
}

// A new collector is returned that will can be used to tail
// the log file and eventually send to the ingestor.
// The token provided is from the Handshake response from the ingestor
func (cfg *Config) createCollector(token string) (*Cho, error) {
	// check if src file exists
	source, err := os.Open(cfg.logSource)
	if err != nil {
		return nil, fmt.Errorf("could not open logSource. Reason: %w", err)
	}

	return &Cho{
		token:        token,
		tailedLogs:   []Log{},
		ingestorAddr: cfg.ingestorAddr,
		source:       source,
		interval:     cfg.interval,
	}, nil
}

func (cho *Cho) tailLog(ctx context.Context, send chan any) {
}

