package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)


type Handshake struct {
	// calatrava and cho should be able to negotiate, but default would typically be 100ms
	Interval string `json:"interval"`
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Token    string `json:"token"`
}

func main() {
	var err error
	var tomlConfig string

	flag.StringVar(&tomlConfig, "config", "default", "path to toml config file. Uses default configs otherwise")
	flag.Parse()

	cfg := &Config{}

	if tomlConfig == "default" {
		log.Println("using default config")
		cfg = defaultConfig()
	} else {
		cfg, err = parseConfig(tomlConfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	token := "random_token"
	collector, err := cfg.createCollector(token)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("config: ", cfg.String())
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Go(func() {
		if err := collector.tailLog(ctx); err != nil {
			log.Println(err)
			return
		}
	})

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
