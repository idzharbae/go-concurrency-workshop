package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func getSomething() string {
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	return "something"
}

func getSomethingWithTimeout(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	c := make(chan string)

	go func() {
		c <- getSomething()
	}()

	for {
		select {
		case s := <-c:
			return s, nil
		case <-ctx.Done():
			log.Println("Kelamaan!")
			return "", ctx.Err()
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	smth, err := getSomethingWithTimeout(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(smth)
}
