package main

import (
	"log"
	"sync"
	"time"
)

type job struct {
	id        int
	result    int
	processed bool
}

func processJob(unprocessed <-chan job, processed chan<- job) {
	for j := range unprocessed {
		log.Printf("processing job: %v\n", j)
		time.Sleep(5 * time.Second)

		j.result *= 2
		j.processed = true
		processed <- j

		log.Printf("processed job: %v\n", j)
	}
}

func main() {
	const maxJobs = 3
	unprocessed := make(chan job, maxJobs)
	processed := make(chan job, maxJobs)

	var wg sync.WaitGroup

	for _ = range maxJobs {
		wg.Go(func() {
			processJob(unprocessed, processed)
		})
	}

	for id := range maxJobs {
		j := job{
			id:        id,
			result:    id + 1,
			processed: false,
		}

		unprocessed <- j

		log.Printf("sent job: %v\n", j)
	}
	close(unprocessed)

	go func() {
		wg.Wait()
		close(processed)
	}()

	for j := range processed {
		log.Printf("received processed job: %v\n", j)
	}
}
