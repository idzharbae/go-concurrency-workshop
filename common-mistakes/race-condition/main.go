package main

import (
	"log"
	"math/rand"
	"time"
)

var m = make(map[string]int)

func main() {
	rand.Seed(time.Now().Unix())

	hewans := []string{"monyet", "gajah", "babi", "anjing", "kucing"}

	for i := 0; i < 16; i++ {
		go func() {
			key := hewans[rand.Intn(len(hewans))]
			m[key] += 1
		}()
	}

	time.Sleep(time.Second)

	for key, val := range m {
		log.Println(key, val)
	}
}
