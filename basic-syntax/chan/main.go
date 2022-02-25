package main

import "fmt"

func main() {
	ch := make(chan string)

	go func() {
		defer close(ch)

		for i := 0; i < 5; i++ {
			ch <- "ping"
		}
	}()

	for data := range ch {
		fmt.Println(data, "pong")
	}
}
