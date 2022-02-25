package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func work(input int) int {
	time.Sleep(time.Second) // Working... Zzzzz
	result := rand.Intn(1000) * input
	return result
}

func jobConsumer(jobCh <-chan int, resultCh chan<- int) {
	for job := range jobCh {
		resultCh <- work(job)
	}
}

func main() {
	workersCount := 4

	var wg sync.WaitGroup
	wg.Add(workersCount)

	jobCh := make(chan int)

	for i := 0; i < 32; i++ {
		jobCh <- i
	}

	resultCh := make(chan int)

	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()
			jobConsumer(jobCh, resultCh)
		}()
	}

	go func() {
		defer close(resultCh)
		wg.Wait()
	}()

	for r := range resultCh {
		log.Println(r)
	}
}
