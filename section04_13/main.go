package main

import "fmt"

const host string = "localhost"

func main() {

	fmt.Printf("%s", host+"\n")

	const pi float32 = 3.14

	fmt.Printf("%v\n", pi)
}
