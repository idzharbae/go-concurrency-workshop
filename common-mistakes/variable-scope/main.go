package main

import (
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			log.Printf("Saya worker nomer %d\n", i)
		}()
	}

	wg.Wait()
}
