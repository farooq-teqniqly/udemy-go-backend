package main

import (
	"fmt"
)

func increaseAgeByVal(age int) {
	age *= 2
	fmt.Printf("[increaseAgeByVal]\tage: %d, address: %p\n", age, &age)
}

func increaseAgeByRef(agePtr *int) {
	if agePtr == nil {
		fmt.Println("age is nil")
		return
	}

	*agePtr = *agePtr * 3
	fmt.Printf("[increaseAgeByRef]\tage: %d, address: %p, agePtr address: %p\n", *agePtr, agePtr, &agePtr)

}

func main() {
	age := 10
	fmt.Printf("[main]\tage: %d, address: %p\n", age, &age)

	increaseAgeByVal(age)
	fmt.Printf("[main]\tage: %d, address: %p\n", age, &age)

	increaseAgeByRef(&age)
	fmt.Printf("[main]\tage: %d, address: %p\n", age, &age)
}
