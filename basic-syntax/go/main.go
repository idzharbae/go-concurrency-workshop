package main

import (
	"fmt"
	"time"
)

func helloworld(i int) {
	fmt.Println("Hello world!", i)
}

func main() {
	for i := 0; i < 10; i++ {
		go helloworld(i)
	}

	fmt.Println("ayaya")

	time.Sleep(time.Second)
}
