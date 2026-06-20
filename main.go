package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)


type Log struct {
	Level       string
	ServiceName string
	Diagnostics string
	Timestamp   string
}

type Config struct {
	// Where cho should tail logs from. Since we're using fsnotify to tail
	// the logSource
	logSource string

	// Ingestor http address
	ingestorAddr string

	// interval delay for sending logs to ingestor. Server and client can neogotiate or
	// use a default
	interval time.Duration
}

type Handshake struct {
	// calatrava and cho should be able to negotiate, but default would typically be 100ms
	Interval string `json:"interval"`
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Token    string `json:"token"`
}

func main() {
	cfg := &Config{
		// used as test-logs for now
		logSource:    "./logs/test_logs.txt",
		ingestorAddr: "http://localhost:9082",
		interval:     time.Millisecond * 500,
	}

	token := "random_token"
	collector, err := cfg.createCollector(token)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := collector.tailLog(ctx); err != nil {
			log.Println(err)
			return
		}
	}()

	wg.Wait()

	// ingestor := "http://localhost:9082"
	// serviceName := "jkvs-cho"

	// initRes, err := contactIngestor(ingestor, serviceName, nil)
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	//
	// if initRes.Status != 200 {
	// 	log.Println(initRes.Message)
	// 	os.Exit(0)
	// }
	//
	// log.Printf("init with ingestor successfull. %s, %s, %+v\n", initRes.Interval, initRes.Message, initRes.Status)

	// establishEndpoint := fmt.Sprintf("%s/establish", cfg.ingestorAddr)
	// conn, err := collector.EstablishConnection(establishEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// conn.Send([]byte("i am tired of ju and jur son\n"))
	// conn.Send([]byte("spinning around in circles"))
	//
	// fmt.Printf("collector  %+v\n", collector)
}
