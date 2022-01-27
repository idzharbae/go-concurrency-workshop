package main

import (
	"fmt"
	"sync"
	"time"
)

func slowIO() int {
	time.Sleep(time.Second / 2)
	return 123
}

func main() {
	now := time.Now()

	// for i := 0; i < 2; i++ {
	// 	fmt.Println(slowIO())
	// }

	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(slowIO())
		}()
	}

	wg.Wait()

	fmt.Println(time.Since(now))
}
