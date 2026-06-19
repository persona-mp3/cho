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
	defer mu.Unlock()
	services[id] = Service{name: serviceName, interval: interval}

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

type Message struct {
	Message string `json:"message"`
}

func establish(res http.ResponseWriter, req *http.Request) {
	flusher, ok := res.(http.Flusher)
	if !ok {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(Message{Message: "streaming not supported"})
		return
	}

	token := req.Header.Get("X-Token")
	if token == "" {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(Message{Message: "missing token"})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(res)
	enc.Encode(Message{Message: "established"})
	flusher.Flush()

	// read what the client sends, respond to each message
	dec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	for dec.More() {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			log.Printf("error: failed to decode client message: %v", err)
			return
		}
		log.Printf("client [+] %s", msg.Message)

		// respond to each message
		if err := enc.Encode(Message{Message: "ack: " + msg.Message}); err != nil {
			log.Printf("error: failed to write response: %v", err)
			return
		}
		flusher.Flush()
	}

	log.Println("info: client stream ended")
}

// func establish(res http.ResponseWriter, req *http.Request) {
// 	log.Println("do you know who else want's to establish a connection w calatrava?")
// 	m := Message{Message: "testing"}
// 	res.WriteHeader(300)
// 	enc := json.NewEncoder(res)
// 	enc.Encode(m)
// 	log.Println("wooo")
// }

func main() {
	http.HandleFunc("/", initialHandshake)
	http.HandleFunc("/establish", establish)

	addr := "localhost:9082"
	log.Println("starting server at http://", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("could not start server. Reason: %s", err)
	}
}
