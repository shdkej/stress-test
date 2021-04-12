package main

import (
	"net"
	"net/http"
	"time"
)

func main() {

	keepAliveTimeout := 600 * time.Second
	timeout := 2 * time.Second
	defaultTansport := &http.Transport{
		Dial: (&net.Dialer{
			KeepAlive: keepAliveTimeout}).Dial,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{
		Transport: defaultTransport,
		Timeout:   timeout,
	}
}
