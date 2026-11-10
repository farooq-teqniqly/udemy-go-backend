package main

import (
	"fmt"
)

func shouldPanic(should bool) string {
	if should {
		panic("panicking")
	}

	return "not panicking"
}

func safeShouldPanic(should bool) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic")
		}
	}()

	return shouldPanic(should)
}

func main() {
	safeShouldPanic(true)
}
