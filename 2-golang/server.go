package main

import (
	"fmt"
	"io"
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
	log.Printf("listening on...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (cs *Counter) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	cs.value += 1
	log.Printf("Connected %d", cs.value)
	io.WriteString(w, `2`)
}

func (cs *Counter) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", cs.value-cs.prevValue)
	cs.prevValue = cs.value
}
