package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
	Process a list of strings in a configurable number of simultaneous threads.
	If an error happens during processing, the loop should die as quickly as possible, and we should get the error out.

	PROBLEMS: Waits until each chunk is completed fully before grabbing a new chunk. This means most of the time
				we won't be running the specified number of threads.
 */

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func main() {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"}
	listChunks := chunkSlice(list, 2)

	var wg sync.WaitGroup

	processingError := ""

	beforeLoop:
	for i, chunk := range listChunks {
		if processingError != "" {
			break beforeLoop
		}
		wg.Add(len(chunk))
		for i := 0; i < len(chunk); i++ {
			go func(val string) {
				fmt.Printf("Starting processing on: %s \n", val)
				defer wg.Done()
				if val == "c" {
					processingError = "found the c"
					return
				}
				d := time.Duration(rand.Intn(10)) * time.Second
				time.Sleep(d)
				fmt.Printf("Done processing on: %s \n", val)
			}(chunk[i])
		}
		wg.Wait()
		fmt.Printf("Done with chunk %v \n", i)
	}
	fmt.Println("All done with job")
}