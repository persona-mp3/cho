package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Handshake struct {
	Interval string `json:"interval"`
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Token    string `json:"token"`
}

func contactIngestor(ingestorAddr string, serviceName string, interval *time.Duration) (*Handshake, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ingestorAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request for ingestor. Reason: %w", err)
	}

	req.Header.Add("X-Service-Name", serviceName)

	if interval != nil {
		req.Header.Add("X-Interval", fmt.Sprintf("%s", interval.String()))
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not contact ingestor. Reason: %w", err)
	}

	defer res.Body.Close()

	var initRes Handshake
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&initRes); err != nil {
		return nil, fmt.Errorf("could not decode Handshake from calatrava. Reason: %w", err)
	}

	return &initRes, nil
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
