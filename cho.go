package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	LogThreshold    = 100
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
	sourceFile, err := os.OpenFile(cfg.LogSource, os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not open logSource. Reason: %w", err)
	}

	return &Cho{
		token:        token,
		tailedLogs:   []Log{},
		mockLogs:     []string{},
		ingestorAddr: cfg.IngestorAddr,
		source:       sourceFile,
		interval:     cfg.Interval,
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


// readLastLog reads the latest appended log by evaluating the number of bytes written
// It returns the newFileSize, the contents written and an error if [Stat] call failed, or
// [ReadAt] failed
func (c *Cho) readLastLog(originalFileSize int64) (int64, []byte, error) {
	fileInfo, err := c.source.Stat()
	if err != nil {
		return 0, nil, err
	}

	newFileSize := fileInfo.Size()

	bytesWritten := newFileSize - originalFileSize

	buffer := make([]byte, bytesWritten)
	n, err := c.source.ReadAt(buffer, originalFileSize)
	_ = n
	if err != nil && !errors.Is(err, io.EOF){
		return 0, nil, err
	}

	return newFileSize, buffer, nil
}

func (cho *Cho) publishLogs() {
	for _, log := range cho.mockLogs {
		fmt.Println(log)
	}
}


func (cho *Cho) cleanUp() {
	cho.source.Close()
}

