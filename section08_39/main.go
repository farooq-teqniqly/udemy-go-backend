package main

import (
	"fmt"
)

type Number interface {
	int | float32 | float64
}

func Product[T Number](numbers ...T) T {
	var total T = 1

	for _, num := range numbers {
		total *= num
	}

	return total
}

func Sum[T int | float32 | float64](numbers ...T) T {
	var total T

	for _, num := range numbers {
		total += num
	}

	return total
}

func main() {
	fmt.Println(Sum(1, 2, 3))
	fmt.Printf("%.2f\n", Sum(1.2, 2.2, 3.2))
	fmt.Printf("%.2f\n", Sum[float32](1.2, 2.2, 3.2)) //specify generic type
	// fmt.Printf("%.2f\n", Sum(true, false, true))	// compile time error

	fmt.Println(Product(1, 4, 5))
	fmt.Println(Product(1.2, 3.4, 5.7))
}
