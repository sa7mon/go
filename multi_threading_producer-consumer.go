package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

/*
	Process a list of strings in a configurable number of simultaneous threads.
	If an error happens during processing, the loop should die as quickly as possible, and we should get the error out.

	PROBLEMS: Items get skipped if the number of threads is not cleanly divisible by the number of items.
*/

func produce(doneChan chan bool, tasksChan chan string) {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

	for _, val := range list {
		tasksChan <- val
	}
	doneChan <- true
}

func consume(tasksChan chan string, errorChan chan error, doneChan chan bool) {
	for {
		select {
			case x, ok := <-errorChan:
				if ok {
					fmt.Printf("Value %v was read.\n", x)
					return
				} else {
					fmt.Println("Channel closed!")
					return
				}
			default:
				fmt.Println("No value ready, moving on.")
		}

		msg := <-tasksChan // Pull item out of channel to consume (process in some way)
		fmt.Printf("Starting processing on: %s \n", msg)

		if msg == "e" {
			errorChan <- errors.New("raising an error")
			doneChan <- true
			return
		}

		d := time.Duration(rand.Intn(10)) * time.Second
		time.Sleep(d)
		fmt.Printf("Done processing on: %s \n", msg)
	}
}

func main() {
	var done = make(chan bool)
	var tasks = make(chan string)
	var errorChan = make(chan error)

	numThreads := 2

	go produce(done, tasks)

	for i := 0; i < numThreads; i++ {
		go consume(tasks, errorChan, done)
	}
	<-done
	fmt.Println("Got to here")
}