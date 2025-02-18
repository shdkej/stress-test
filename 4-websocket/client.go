package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	limit   = flag.Int("l", 10, "limit of fconnection")
	port    = flag.String("p", "8080", "listen address")
	host_ip = flag.String("h", "localhost", "host ip address")
)

func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *host_ip + ":" + *port, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	start := time.Now()

	fmt.Println("start")
	start = time.Now()
	count := connectAsyncWithLimit(u.String())
	elapsed := time.Since(start)
	fmt.Println("Connect with goroutines:", count, elapsed)

	res, err := http.Get("http://" + *host_ip + ":" + *port + "/check")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

func connectAsyncWithLimit(url string) int {
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(*limit)

	for i := 0; i < *limit; i++ {
		go func() {
			loadTimer(url)
			wg.Done()
		}()
	}
	wg.Wait()
	return count
}

func load(url string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("recv err:", err)
			return
		}
		log.Printf("recv: %s", message)
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte("client"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}

}

func loadTimer(url string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	err = c.WriteMessage(websocket.TextMessage, []byte(`1`))

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	timer := time.NewTimer(300 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-done:
			return
		case <-timer.C:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func loadTicker(url string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
