package main

import (
	"log"
	"math/rand"
	"time"
)

func fanInSimple(cs ...<-chan string) <-chan string {
	c := make(chan string)
	for _, ci := range cs { // spawn channel based on the number of input channel

		go func(cv <-chan string) { // cv is a channel value
			for {
				c <- <-cv
			}
		}(ci) // send each channel to

	}
	return c
}

func hewanGenerator() <-chan string {
	ch := make(chan string)
	hewans := []string{"monyet", "gajah", "babi", "anjing", "kucing"}

	go func() {
		for {
			time.Sleep(time.Second)
			i := rand.Intn(len(hewans))
			ch <- hewans[i]
		}
	}()

	return ch
}

func tumbuhanGenerator() <-chan string {
	ch := make(chan string)
	tumbuhans := []string{"jamur", "terong", "kacang", "semangka", "melon"}

	go func() {
		for {
			time.Sleep(time.Second)
			i := rand.Intn(len(tumbuhans))
			ch <- tumbuhans[i]
		}
	}()

	return ch
}

func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10; i++ {
		log.Println(<-fanInSimple(hewanGenerator(), tumbuhanGenerator()))
	}
}
