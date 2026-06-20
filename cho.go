package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	LogThreshold    = 10
	DefaultInterval = 500 * time.Millisecond
)

type Cho struct {
	// Provided by calatrava after initialising handshake. This will
	// be used for subsequent requests instead of the service name
	token string

	tailedLogs   []Log
	mockLogs     []string
	ingestorAddr string

	// File is opened for readOnly access
	source *os.File

	// interval for sending out logs to calatrava
	interval time.Duration
}

// A new collector is returned that will can be used to tail
// the log file and eventually send to the ingestor.
// The token provided is from the Handshake response from the ingestor
func (cfg *Config) createCollector(token string) (*Cho, error) {
	sourceFile, err := os.OpenFile(cfg.logSource, os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not open logSource. Reason: %w", err)
	}

	return &Cho{
		token:        token,
		tailedLogs:   []Log{},
		mockLogs:     []string{},
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
	notification := make(chan struct{}, 100)

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

	fileInfo, err := cho.source.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	log.Println("originalFileSize: ", fileSize)

	ticker := time.NewTicker(cho.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if len(cho.mockLogs) < LogThreshold {
				continue
			}
			cho.publishLogs()
			cho.mockLogs = []string{}
		default:
		}

		select {
		case <-parentCtx.Done():
			return nil
		case <-notification:
			newSize, content, err := cho.readLastLog(fileSize)
			if err != nil {
				return fmt.Errorf("could not read last log. %w", err)
			}

			fileSize = newSize
			cho.mockLogs = append(cho.mockLogs, string(content))
		default:
		}
	}

}

func (cho *Cho) cleanUp() {
	cho.source.Close()
}

func (cho *Cho) publishLogs() {
	for _, log := range cho.mockLogs {
		fmt.Println(log)
	}
}
