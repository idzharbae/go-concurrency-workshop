package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

func beliTepung(ctx context.Context) error {
	r := rand.Intn(3)

	if r > 0 {
		return errors.New("Tepung nya habis")
	}

	// Nakar tepung dulu
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Berhasil beli tepung")
			return nil
		case <-ctx.Done():
			log.Println("Gajadi beli tepung")
			return ctx.Err()
		}
	}
}

func beliTelor(ctx context.Context) error {
	r := rand.Intn(3)

	if r > 0 {
		return errors.New("Telor nya habis")
	}

	// Nimbang telor dulu
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Berhasil beli telor")
			return nil
		case <-ctx.Done():
			log.Println("Gajadi beli telor")
			return ctx.Err()
		}
	}
}

func beliGula(ctx context.Context) error {
	r := rand.Intn(3)

	if r > 0 {
		return errors.New("Gula nya habis")
	}

	// Nimbang telor dulu
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Berhasil beli Gula")
			return nil
		case <-ctx.Done():
			log.Println("Gajadi beli Gula")
			return ctx.Err()
		}
	}
}

func bikinKue(ctx context.Context) error {
	errChan := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan bool)

	var wg sync.WaitGroup
	wg.Add(2)

	// Mesen tepung
	go func(ctx context.Context) {
		defer wg.Done()

		err := beliTepung(ctx)
		if err != nil {
			errChan <- err
		}
	}(ctx)

	// Mesen telor
	go func(ctx context.Context) {
		defer wg.Done()

		err := beliTelor(ctx)
		if err != nil {
			errChan <- err
		}
	}(ctx)

	go func() {
		defer func() {
			done <- true
			close(errChan)
		}()

		wg.Wait()
	}()

	for {
		select {
		case <-done:
			return nil
		case err := <-errChan:
			go func() {
				for err := range errChan {
					log.Println(err)
				}
			}()
			return err
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())

	err := bikinKue(context.Background())
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Berhasil bikin kue")
	}

	time.Sleep(time.Second)
}
