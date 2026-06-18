package main

type Log struct {
	Level       string
	ServiceName string
	Information string
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
