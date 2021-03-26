package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const HOST_IP = "172.17.0.2"
const HOST_PORT = "30080"

//const HOST_IP = "193.123.251.27"
const LIMIT = 1000

func main() {
	start := time.Now()

	fmt.Println("start")
	//count1 := connectWithSync("8080")
	//elapsed1 := time.Since(start)
	//fmt.Println("Connect Synconisly:", count1, elapsed1)

	//count := 0
	//count5 := connectAsyncWithoutWG(count, "8080")
	//elapsed5 := time.Since(start)
	//fmt.Println("Connect without wait group:", count5, elapsed5)

	//start = time.Now()
	//connectAsyncWithOne()

	fmt.Println("start2")
	start = time.Now()
	count2 := connectLongTime(HOST_PORT)
	//count2 := connectAsyncWithLimit("8080")
	//count2 = connectAsyncWithLimit("8081")
	//count2 = connectAsyncWithLimit("8082")
	elapsed2 := time.Since(start)
	time.Sleep(20 * time.Second)
	fmt.Println("Connect with multiple goroutines:", count2, elapsed2)

	value := GetValue(HOST_PORT)
	fmt.Println("Get Value:", value)
	fmt.Println("Connect with 1 goroutine:", elapsed2)
}

func connectWithSync(port string) int {
	count := 0
	for i := 0; i < LIMIT; i++ {
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
			if loadLong(port) == "hello" {
				count += 1
			}
		}()
	}
	return count
}

func load(port string) string {
	res, err := http.Get("http://" + HOST_IP + ":" + port)
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
