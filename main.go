package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Log struct {
	Level       string
	ServiceName string
	Diagnostics string
	Timestamp   string
}

type Config struct {
	// Where cho should tail logs from
	logSource string

	// Ingestor http address
	ingestorAddr string

	// interval delay for sending logs to ingestor. Server and client can neogotiate or
	// use a default
	interval int

	// port cho should listen on
	port int
}

type Cho interface {
	Start(
		ctx context.Context,
		addr string,
		logSource string,
		ingestorAddr string,
		interval time.Duration,
	)
}

type InitResponse struct {
	// calatrava and cho should be able to negotiate, but default would typically be 100ms
	Interval int `json:"interval"`
	Status   int `json:"status"`
}

func contactIngestor(ingestorAddr string, serviceName string, interval *time.Duration) (*InitResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ingestorAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request for ingestor. Reason: %w", err)
	}

	req.Header.Add("X-Service-Name", serviceName)
	req.Header.Add("X-Interval", fmt.Sprintf("%+v", interval))
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not contact ingestor. Reason: %w", err)
	}

	defer res.Body.Close()

	var initRes InitResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&initRes); err != nil {
		return nil, fmt.Errorf("could not decode InitResponse from calatrava. Reason: %w", err)
	}

	return &initRes, nil
}

func main() {
	ingestor := "http://localhost:9082"
	serviceName := "jkvs-cho"

	initRes, err := contactIngestor(ingestor, serviceName, nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("init with ingestor successfull. ", initRes)
}

