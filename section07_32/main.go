package main

import (
	"fmt"
)

func closeResource() {
	fmt.Println("closing resource")
}

func main() {
	fmt.Println("opening resource")
	defer closeResource()

	for range 10 {
		fmt.Println("using resource")
	}
}
