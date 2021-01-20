package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type State struct {
	mutex sync.Mutex
	queue Queue
}

/*
	FIFO queue
 */
type Queue struct {
	items []string
}

func (q Queue) String() string {
	return fmt.Sprintf("%v", q.items)
}

/*
	Constructor
 */
func NewQueue(initial []string) Queue {
	return Queue{items: initial}
}

func (q *Queue) Add(item string) {
	q.items = append(q.items, item)
}

func (q *Queue) Get() string {
	if len(q.items) == 0 {
		return ""
	}
	popped := q.items[0]  // Get top element
	q.items = q.items[1:] // Remove top element
	return popped
}

func (q *Queue) Length() int {
	return len(q.items)
}

func Work(state *State, wg *sync.WaitGroup, workerID int) {
	for {
		// Get work to do
		state.mutex.Lock()
		workToDo := state.queue.Get()
		state.mutex.Unlock()
		if workToDo	== "" {
			fmt.Printf("[worker%v] No work to do. Exiting\n", workerID)
			wg.Done()
			return
		}

		fmt.Printf("[worker%v] Received work to do:%v \n", workerID, workToDo)
		d := time.Duration(rand.Intn(5)) * time.Second
		time.Sleep(d)
		fmt.Printf("[worker%v] Done working on: %s \n", workerID, workToDo)
	}
}

func main() {
	const numWorkers = 8

	globalState := &State{queue: NewQueue([]string{"a", "b"})}

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go Work(globalState, wg, i)
	}

	wg.Wait()
	fmt.Println("[main] All workers have completed")
}