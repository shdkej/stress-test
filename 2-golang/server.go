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
	//log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func (cs *Counter) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	cs.value += 1
	log.Printf("Connected", cs.value)
	io.WriteString(w, `hello`)
}

func (cs *Counter) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%d", cs.value-cs.prevValue)
	cs.prevValue = cs.value
}

func HTTP2TestHandler(w http.ResponseWriter, r *http.Request) {
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
