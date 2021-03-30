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
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", cs.HTTP2TestHandler)
	http.HandleFunc("/value", cs.GetValueHandler)
	log.Printf("listening on...")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func (cs *Counter) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", cs.value-cs.prevValue)
	cs.prevValue = cs.value
}

func (cs *Counter) HTTP2TestHandler(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); ok {
		options := &http.PushOptions{
			Header: http.Header{
				"Accept-Encoding": r.Header["Accept-Encoding"],
			},
		}
		if err := pusher.Push("/styles.css", options); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	w.Write([]byte("HTTP2 Test"))
}
