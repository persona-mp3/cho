package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Cho struct {
	// Provided by calatrava after initialising handshake. This will 
	// be used for subsequent requests instead of the service name
	token        string

	tailedLogs   []Log
	ingestorAddr string

	// File is opened for readOnly access
	source   *os.File

	// interval for sending out logs to calatrava
	interval time.Duration
}

// A new collector is returned that will can be used to tail
// the log file and eventually send to the ingestor.
// The token provided is from the Handshake response from the ingestor
func (cfg *Config) createCollector(token string) (*Cho, error) {
	sourceFile, err := os.Open(cfg.logSource)
	if err != nil {
		return nil, fmt.Errorf("could not open logSource. Reason: %w", err)
	}

	return &Cho{
		token:        token,
		tailedLogs:   []Log{},
		ingestorAddr: cfg.ingestorAddr,
		source:       sourceFile,
		interval:     cfg.interval,
	}, nil
}

func (cho *Cho) tailLog(parentCtx context.Context) error {
	abs, err := filepath.Abs(cho.source.Name())
	if err != nil {
		return err
	}
	defer cho.cleanUp()

	logDir := filepath.Dir(abs)
	notification := make(chan struct{})

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	go func() {
		defer close(notification)
		if err := watch(ctx, logDir, abs, notification); err != nil {
			log.Fatal(err)
			return
		}
	}()

	log.Println("tailer started....")
	log.Println()

	for {
		select {
		case <-parentCtx.Done():
			return nil
		case <-notification:
			log.Println(" > write event occured")
		}
	}

}

func (cho *Cho) cleanUp() {
	cho.source.Close()
}
