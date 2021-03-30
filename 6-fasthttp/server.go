package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
)

type Counter struct {
	ch    chan int
	value int
}

func (cs Counter) fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	cs.ch <- 1
	fmt.Fprintf(ctx, "Hi %q", ctx.RequestURI())
}

func (cs Counter) valueHandler(ctx *fasthttp.RequestCtx) {
	log.Println(cs.value)
}

func main() {
	cs := Counter{}
	cs.ch = make(chan int, 10)
	log.Println("listening ...")
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			cs.fastHTTPHandler(ctx)
		case "/value":
			cs.valueHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	go func() {
		if err := fasthttp.ListenAndServe(":8080", m); err != nil {
			log.Fatal("Stop Server")
		}
	}()
	for {
		select {
		case <-cs.ch:
			cs.value += 1
		}
	}
}
