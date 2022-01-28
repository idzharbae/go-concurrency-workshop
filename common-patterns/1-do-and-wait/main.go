package main

import (
	"log"
	"sync"
	"time"
)

func doSomething() {
	time.Sleep(time.Second)
	log.Println("Done!")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	now := time.Now()
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			doSomething()
		}()
	}

	wg.Wait()

	log.Println(time.Since(now))
}
