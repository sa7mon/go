package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	//"os/signal"
	"runtime"
	"sync"
	//"syscall"
	"time"
)

/*
	Example from: https://medium.com/hdac/producer-consumer-pattern-implementation-with-golang-6ac412cf941c

	Process a list of strings in a configurable number of simultaneous threads.
	If an error happens during processing, the loop should die as quickly as possible, and we should get the error out.

*/

type Consumer struct {
	in *chan string
	jobs chan string
}
func (c Consumer) Work(wg *sync.WaitGroup, termChan chan error) {
	defer wg.Done()
	for job := range c.jobs {
		fmt.Printf("Starting processing on: %s \n", job)
		if job == "e" {
			termChan <- errors.New("found bad item to process")
			wg.Done()
			return
		}
		d := time.Duration(rand.Intn(10)) * time.Second
		time.Sleep(d)
		fmt.Printf("Done processing on: %s \n", job)
	}
}
func (c Consumer) Consume(ctx context.Context) {
	for {
		select {
		case job := <-*c.in:
			c.jobs <- job
		case <-ctx.Done():
			close(c.jobs)
			return
		}
	}
}
type Producer struct {
	in *chan string
}
func (p Producer) Produce(items []string, termChan chan error) {
	for _, item := range items {
		*p.in <- item
	}
	termChan <- nil
}

func main() {
	itemsToProcess := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}

	const nConsumers = 2
	runtime.GOMAXPROCS(runtime.NumCPU())
	in := make(chan string, 1)
	p := Producer{&in}
	c := Consumer{&in, make(chan string, nConsumers)}

	termChan := make(chan error)
	//signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go p.Produce(itemsToProcess, termChan)

	ctx, cancelFunc := context.WithCancel(context.Background())
	go c.Consume(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(nConsumers)
	for i := 0; i < nConsumers; i++ {
		go c.Work(wg, termChan)
	}

	var errorWhileProcessing error
	errorWhileProcessing = <- termChan
	cancelFunc()
	wg.Wait()

	if errorWhileProcessing != nil {
		fmt.Println("Caught error while processing: " + errorWhileProcessing.Error())
	}
}