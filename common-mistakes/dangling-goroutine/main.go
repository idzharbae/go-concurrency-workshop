package main

import (
	"log"
	"runtime"
	"time"
)

type Publisher struct {
	c chan string
}

func NewPublisher() *Publisher {
	pub := &Publisher{c: make(chan string, 8)}
	pub.run()

	return pub
}

func (p *Publisher) run() {
	go func() {
		for data := range p.c {
			log.Println(data)
		}
	}()
}

func (p *Publisher) Publish(data string) {
	p.c <- data
}

func main() {
	data := []string{"data", "data lagi", "abcdefg", "more data"}

	for _, d := range data {
		p := NewPublisher()
		p.Publish(d)
	}

	time.Sleep(time.Second)

	log.Println("Show me da stack")

	b := make([]byte, 2048)
	runtime.Stack(b, true)

	log.Println(string(b))
}
