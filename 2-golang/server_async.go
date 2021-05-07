package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Counter struct {
	ch       chan int
	value    int
	duration time.Duration
}

func main() {
	cs := Counter{}
	cs.ch = make(chan int, 10)
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", cs.HealthCheckHandler)
	http.HandleFunc("/value", cs.GetValueHandler)
	log.Printf("listening on...")
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Stop Server")
		}
	}()
	for {
		select {
		case <-cs.ch:
			cs.value += 1
		}
	}
}

func (cs *Counter) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	cs.ch <- 1
	log.Printf("%d", cs.value)
}

func (cs *Counter) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", cs.value)
}
