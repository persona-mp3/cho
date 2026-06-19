package main

import (
	"fmt"
	"log"
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

// Used to persist long-lived sessions with calatrava. This allows the underlying tcp-connection
// passed around in a streaming way. The connection is tied to the lifetime to the request via
// the context. The flusher, flushes
type Conn struct {
}

// EstablishConnection creates a new httpRequest to the server provided
// in the establishEndpoint. The connection returned should be able to be
// written to on demand. It's not expected that the server sends messages
// except from closing the connection or minor event changes.
func (c *Cho) EstablishConnection(establishEndpoint string) (*Conn, error) {
	log.Println("establishing keep-alive connection with server")
	return &Conn{}, nil
}

func (c *Conn) Send(data []byte) error {
	return nil
}


