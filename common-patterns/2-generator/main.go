package main

import (
	"log"
	"math/rand"
	"time"
)

func generateNumber() int {
	time.Sleep(time.Second) // Generating...
	return rand.Intn(1000)
}

func numberGenerator(n int) <-chan int {
	resultChan := make(chan int)
	go func() {
		defer close(resultChan)
		for i := 0; i < n; i++ {
			resultChan <- generateNumber()
		}
	}()

	return resultChan
}

func main() {
	resultChan := numberGenerator(10)

	for result := range resultChan {
		log.Println(result)
	}
}
