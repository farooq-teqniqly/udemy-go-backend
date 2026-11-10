package main

import (
	"errors"
	"fmt"
	"strings"
)

// variadic function example
func sum(numbers ...int) int {
	total := 0

	for _, number := range numbers {
		total += number
	}

	return total
}

// returning an error
func divide(a, b int) (int, error) {
	if (b == 0) {
		return 0, errors.New("divide by zero")
	}

	return a / b, nil
} 

// named return values
func splitName(fullName string) (firstName, lastName string) {
	parts := strings.Split(fullName, " ")
	firstName = parts[0]
	lastName = parts[1]
	return
}



func main() {
	// res := sum(1, 2, 3, 4, 5)
	// fmt.Println(res)

	// res, err := divide(10, 0)

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(res)
	// }

	firstName, lastName := splitName("Farooq Mahmud")
	fmt.Println(firstName)
	fmt.Println(lastName)
}
