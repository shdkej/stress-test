package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const HOST_IP = "localhost"
const HOST_PORT = "8080"

//const HOST_IP = "193.123.251.27"
const LIMIT = 1

func main() {
	start := time.Now()

	fmt.Println("start")
	start = time.Now()
	count2 := connectLongTime(HOST_PORT)
	elapsed2 := time.Since(start)
	fmt.Println("Connect with multiple goroutines:", count2, elapsed2)

	time.Sleep(10 * time.Second)

	value := GetValue(HOST_PORT)
	fmt.Println("Get Value:", value)
	fmt.Println("Connect with 1 goroutine:", elapsed2)
}

func connectAsyncWithLimit(port string) int {
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(LIMIT)

	for i := 0; i < LIMIT; i++ {
		go func() {
			if load(port) == "hello" {
				count += 1
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return count
}

func connectAsyncDifferentAdd(port string) int {
	count := 0
	wg := sync.WaitGroup{}

	for i := 0; i < LIMIT; i++ {
		wg.Add(1)
		go func() {
			if load(port) == "hello" {
				count += 1
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return count
}

func connectAsyncWithoutWG(count int, port string) int {
	for i := 0; i < LIMIT; i++ {
		go func() {
			if load(port) == "hello" {
				count += 1
			}
		}()
	}
	return count
}

func connectLongTime(port string) int {
	count := 0
	for i := 0; i < LIMIT; i++ {
		go func() {
			if load(port) == "hello" {
				count += 1
			}
		}()
	}
	return count
}

func load(port string) string {
	client := &http.Client{Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	},
	}
	res, err := client.Get("https://" + HOST_IP + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func loadLong(port string) string {
	res, err := http.Get("http://" + HOST_IP + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)
	return string(body)
}

func GetValue(port string) string {
	res, err := http.Get("http://" + HOST_IP + ":" + port + "/value")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
