package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/semaphore"
)

func work(input int) int {
	time.Sleep(time.Second) // Working... Zzzzz
	result := rand.Intn(1000) * input
	return result
}

func jobGenerator() <-chan int {
	jobCh := make(chan int)

	go func() {
		defer close(jobCh)
		for i := 0; i < 32; i++ {
			jobCh <- i
		}
	}()

	return jobCh
}

func main() {
	workersCount := 4
	var sem = semaphore.NewWeighted(int64(workersCount))

	jobCh := jobGenerator()

	resultCh := make(chan int)
	done := make(chan bool)

	ctx := context.Background()

	go func() {
		for r := range resultCh {
			log.Println(r)
		}

		done <- true
	}()

	for job := range jobCh {
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Fatalf("Failed to acquire semaphore: %v", err)
		}

		go func() {
			defer sem.Release(1)
			resultCh <- work(job)
		}()
	}

	go func() {
		defer close(resultCh)
		sem.Acquire(ctx, int64(workersCount))
	}()

	<-done
}
