package main

import (
	"fmt"
	"time"
)

func main() {
	jobIDs := make(chan int, 3)

	// Producer
	go func() {
		for i := range 10 {
			jobIDs <- i
			fmt.Println("Queued job ", i)
		}

		close(jobIDs)
	}()

	// Consumer
	for jobID := range jobIDs {
		// Add a large delay to simulate producer blocking until buffer has space.
		time.Sleep(5 * time.Second)
		fmt.Println("Processed job ", jobID)
	}
}
