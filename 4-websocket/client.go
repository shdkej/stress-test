package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

const LIMIT = 35000
const HOST = "192.168.10.165"

func main() {
	u := url.URL{Scheme: "ws", Host: HOST + ":8080", Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	start := time.Now()

	fmt.Println("start")
	start = time.Now()
	count := connectAsyncWithLimit(u.String())
	elapsed := time.Since(start)
	fmt.Println("Connect with goroutines:", count, elapsed)

	res, err := http.Get("http://" + HOST + ":8080/check")
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
	wg.Add(LIMIT)

	for i := 0; i < LIMIT; i++ {
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
	err = c.WriteMessage(websocket.TextMessage, []byte("write"))

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

	timer := time.NewTimer(120 * time.Second)
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
