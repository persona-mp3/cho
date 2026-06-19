package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

var mu = sync.Mutex{}
var services = make(map[string]Service)

type Service struct {
	id       string
	name     string
	interval string
}

type HandshakeResponse struct {
	Token    string `json:"token"`
	Interval string `json:"interval"`
	Status   int    `json:"status"`
	Message  string `json:"message"`
}

func initialHandshake(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")

	serviceName := req.Header.Get("X-Service-Name")
	interval := req.Header.Get("X-Interval")
	encoder := json.NewEncoder(res)

	var response HandshakeResponse

	if len(strings.ReplaceAll(serviceName, " ", "")) == 0 {
		response.Message = "Unknown service"
		response.Status = 400
		response.Interval = ""
		res.WriteHeader(400)
		if err := encoder.Encode(&response); err != nil {
			log.Printf("could not encode response to client. Reason: %+v\n", err)
			return
		}
		return
	}

	id := generateHash(serviceName)

	mu.Lock()
	services[id] = Service{name: serviceName, interval: interval}
	mu.Unlock()

	response.Message = "OK"
	response.Status = 200
	response.Token = id

	log.Println("token::", id)

	if len(interval) == 0 {
		response.Interval = "10ms"
	}

	res.WriteHeader(200)
	encoder.Encode(&response)
	log.Printf("added new service %s\n", serviceName)

}

func main() {
	http.HandleFunc("/", initialHandshake)

	addr := "localhost:9082"
	log.Println("starting server at http://", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("could not start server. Reason: %s", err)
	}
}
