package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func printRandomNumber(delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(delay)
	fmt.Printf("Random number %d\n", rand.IntN(20))
}

func main() {
	workerCount := 10
	var wg sync.WaitGroup

	wg.Add(workerCount)

	fmt.Println("\nStarting...")

	for i := range workerCount {
		go printRandomNumber(time.Duration(i)*time.Second, &wg)
	}

	fmt.Println("\nWaiting...")

	wg.Wait()
	fmt.Println("\nExiting...")
}
