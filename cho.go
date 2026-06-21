package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Cho struct {
	// Provided by calatrava after initialising handshake. This will
	// be used for subsequent requests instead of the service name
	token string

	tailedLogs []Log
	rawLogs    []string

	// HTTP address of calatrava instance running
	ingestorAddr string

	// source is the file opened for read only access for cho to tail
	source *os.File

	// interval for sending out logs to calatrava. If [logThreshold] has been set, cho
	// will no longer send logs during this [interval] but will batch them when the [tailedLogs]
	// have reach the [logThreshold]
	interval time.Duration

	// logThreshold is the amount of logs to be held before logs are sent over to calatrava. If
	logThreshold int
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
		rawLogs:      []string{},
		ingestorAddr: cfg.IngestorAddr,
		source:       sourceFile,
		interval:     cfg.Interval,
		logThreshold: cfg.LogThreshold,
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
			if len(cho.rawLogs) < cho.logThreshold {
				continue
			}

			logs, err := parseLogs(cho.rawLogs)
			if err != nil {
				log.Println("failed to parse logs: ", err)
				continue
			}
			_ = logs

			cho.publishLogs()
			cho.rawLogs = []string{}
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

			if content != nil {
				cho.rawLogs = append(cho.rawLogs, string(content))
			}

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

	if newFileSize == 0 {
		log.Printf("warn: logFile has been truncated. Returned size of %d\n", newFileSize)
		return 0, nil, nil
	}

	bytesWritten := newFileSize - originalFileSize

	buffer := make([]byte, bytesWritten)
	n, err := c.source.ReadAt(buffer, originalFileSize)
	_ = n
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, nil, err
	}

	return newFileSize, buffer, nil
}

func parseLogs(logs []string) ([]*Log, error) {
	parsedLogs := []*Log{}

	for _, line := range logs {
		for entry := range strings.SplitSeq(line, "\n") {
			if len(strings.TrimSpace(entry)) == 0 {
				continue
			}

			parsedLog, err := JSONParser(entry)
			if err != nil {
				log.Printf("could not parse log: %s. Reason: %s\n", entry, err)
				continue
			}

			log.Println("logEntry:", entry)
			log.Println(parsedLog.String())
			parsedLogs = append(parsedLogs, parsedLog)

		}
	}

	return parsedLogs, nil
}

func (cho *Cho) publishLogs() {
}

func (cho *Cho) cleanUp() {
	cho.source.Close()
	log.Println("closed source file")
}
