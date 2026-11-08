package main

import "fmt"

func ifElseExample() {
	maxTemp := 100
	minTemp := 30
	currentTemp := 0

	for currentTemp < 120 {
		if currentTemp > maxTemp {
			fmt.Println("Too hot!")
		} else if currentTemp < minTemp {
			fmt.Println("Too cold")
		} else {
			fmt.Println("Temp within bounds")
		}

		currentTemp++
	}
}

func mapExample() {
	accessMap := map[string]bool{
		"Farooq": true,
		"Noor":   false,
	}

	userToCheck := "Bubba"
	hasAccess, keyExists := accessMap[userToCheck]

	if !hasAccess || !keyExists {
		fmt.Println(userToCheck + " does not have access.")
		return
	}

	if hasAccess && keyExists {
		fmt.Println(userToCheck + " has access.")
		return
	}

}

func main() {
	// ifElseExample()
	mapExample()
}
