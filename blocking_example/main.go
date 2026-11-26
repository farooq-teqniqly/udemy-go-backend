package main

import "fmt"

func main() {
	messages := make(chan string)
	defer close(messages)

	go func() {
		fmt.Println("waiting to send a message")
		// Commenting out the line below induces a deadlock. Go runtime detects no senders, forward progress not
		// possible so kills the process.
		messages <- "hello"
		fmt.Println("message sent")
	}()

	fmt.Println("main waiting to receive message")
	message := <-messages
	fmt.Printf("main received message: %s", message)
}
