package main

import (
	"fmt"
)

func checkType(i any) {
	switch v := i.(type) {
	case int:
		fmt.Printf("%d is an integer\n", v)
	case string:
		fmt.Printf("%s is a string\n", v)
	case bool:
		fmt.Printf("%t is a boolean\n", v)
	default:
		fmt.Printf("I don't know what %T is\n", v)
	}
}

func simpleSwitchExample() {
	minTemp := 30
	maxTemp := 100
	currentTemp := 20

	for currentTemp <= 120 {
		switch {
		case currentTemp < minTemp:
			fmt.Printf("%d is too cold!\n", currentTemp)
		case currentTemp > maxTemp:
			fmt.Printf("%d is too hot!\n", currentTemp)
		default:
			fmt.Printf("%d is within bounds.\n", currentTemp)
		}

		currentTemp++
	}
}

func main() {
	// simpleSwitchExample()
	checkType(10)
	checkType("Farooq")
	checkType(true)
	checkType(5.33333333)
}
