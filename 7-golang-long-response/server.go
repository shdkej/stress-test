package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Counter struct {
	value     int
	prevValue int
	duration  time.Duration
}

func main() {
	cs := Counter{}
	http.HandleFunc("/", cs.HealthCheckHandler)
	http.HandleFunc("/value", cs.GetValueHandler)
	log.Printf("listening on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (cs *Counter) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	cs.value += 1
	s := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-s.C:
			return
		}
	}
}

func (cs *Counter) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", cs.value-cs.prevValue)
	cs.prevValue = cs.value
}
