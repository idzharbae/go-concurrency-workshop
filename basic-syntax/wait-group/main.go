package main

import (
	"fmt"
	"time"
)

func slowIO() int {
	time.Sleep(time.Second / 2)
	return 123
}

func main() {
	now := time.Now()

	for i := 0; i < 2; i++ {
		fmt.Println(slowIO())
	}

	fmt.Println(time.Since(now))
}
