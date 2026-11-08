package main

import "fmt"

func countdown(start int) {
	i := start

	//while style loop
	for i >= 0 {
		fmt.Println(i)
		i--
	}
}

func printEvenNumbers(start int, stop int) {
	i := start

	// "infinite" loop
	for {
		if i > stop {
			break
		}

		if i%2 == 0 {
			fmt.Println(i)
		}

		i++
	}
}

func printArray() {
	items := [4]string{"Go", "C#", "Python", "Rust"}

	// if index (or item) not needed, use discard "_"
	for index, item := range items {
		fmt.Printf("%d:%s\n", index, item)
	}
}

func main() {
	// countdown(10)
	// printEvenNumbers(10, 20)
	printArray()
}
