package main

import (
	"fmt"
)

type Counter struct {
	value chan int
}

func main() {
	fmt.Println("vim-go")

	done1 := make(chan int)
	i := 0
	go func() {
		for i = 0; i < 10; i++ {
			done1 <- i
		}
	}()

	var c int
	for {
		select {
		case c = <-done1:
			println(c)
		}
	}

}

func (c Counter) increment() {
}
