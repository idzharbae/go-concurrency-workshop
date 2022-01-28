package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func lamarKerja(ctx context.Context, perusahaan string) <-chan string {
	c := make(chan string)

	go func() {
		ticker := time.NewTicker(time.Duration(+rand.Intn(3))*time.Second + time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("Gagal melamar di %s\n", perusahaan)
				return
			case <-ticker.C:
				c <- fmt.Sprintf("Selamat anda diterima di %s!", perusahaan)
				return
			}
		}
	}()

	return c
}

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

func lamarKerjaConcurrent(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	responTokped := lamarKerja(ctx, "Tokped")
	responGojek := lamarKerja(ctx, "Gojek")

	for {
		select {
		case s := <-fanInSimple(responGojek, responTokped):
			return s, nil
		case <-ctx.Done():
			log.Println("Kelamaan, mending nganggur!")
			return "", ctx.Err()
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	result, err := lamarKerjaConcurrent(context.Background())
	if err != nil {
		log.Println(err)
	} else {
		log.Println(result)
	}

	time.Sleep(time.Second)
}
