package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	limit   = flag.Int("l", 10, "limit of fconnection")
	port    = flag.String("p", "8080", "listen address")
	host_ip = flag.String("h", "localhost", "host ip address")
)

func main() {
	flag.Parse()
	start := time.Now()
	ch := make(chan int, *limit)

	fmt.Println("start")

	count2 := connectLongTime(*port, ch)
	<-ch

	value := GetValue(*port)
	fmt.Println("Get Value:", value)

	elapsed2 := time.Since(start)
	fmt.Println("Connect with goroutines:", count2, elapsed2)
}

func connectWithSync(port string) int {
	count := 0
	for i := 0; i < *limit; i++ {
		if load(port) == "hello" {
			count += 1
		}
	}
	return count
}

func connectAsyncWithOne(port string) int {
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		if load(port) == "hello" {
			count += 1
		}
		wg.Done()
	}()
	wg.Wait()
	return count
}

func connectAsyncWithLimit(port string) int {
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(*limit)

	for i := 0; i < *limit; i++ {
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

	for i := 0; i < *limit; i++ {
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
	for i := 0; i < *limit; i++ {
		go func() {
			if load(port) == "hello" {
				count += 1
			}
		}()
	}
	return count
}

func connectLongTime(port string, ch chan int) int {
	count := 0
	for i := 0; i < *limit; i++ {
		go func(i int) {
			loadLong(port)
			ch <- i
		}(i)
	}
	return count
}

func load(port string) string {
	res, err := http.Get("http://" + *host_ip + ":" + port)
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

func loadLong(port string) {
	res, err := http.Get("http://" + *host_ip + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)
	res.Body.Close()
}

func GetValue(port string) string {
	res, err := http.Get("http://" + *host_ip + ":" + port + "/value")
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
