package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	list := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

	var wg sync.WaitGroup
	wg.Add(len(list))

	for i := 0; i < len(list); i++ {
		go func(val string) {
			fmt.Printf("Starting processing on: %s \n", val)
			defer wg.Done()
			d := time.Duration(rand.Intn(10)) * time.Second
			time.Sleep(d)
			fmt.Printf("Done processing on: %s \n", val)
		}(list[i])
	}
	wg.Wait()
	fmt.Println("Program finished executing")
}