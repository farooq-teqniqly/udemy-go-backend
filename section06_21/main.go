package main

import (
	"fmt"
	"slices"
)

func printArrayExample() {
	var numbers [2]int
	fmt.Printf("%+v\n", numbers)

	numbers[0] = 1
	numbers[1] = 2
	fmt.Printf("%+v\n", numbers)

	numbers[1] = 3
	fmt.Printf("%+v\n", numbers)

	fmt.Print("[")
	for i := range numbers {
		fmt.Printf("%d ", numbers[i])
	}
	fmt.Print("]")
}

func printArray[T any](arr []T) {
	fmt.Printf("Items: %+v, Length: %d, Capacity: %d\n", arr, len(arr), cap(arr))
}

func sliceExample() {
	names := []string{"Farooq", "Morgana", "Lucia", "Noor", "Yasin"}
	printArray(names)

	items := make([]int, 3, 5)
	printArray(items)

	items = append(items, 1)
	items = append(items, 2)
	items = append(items, 3)
	items = append(items, 4)
	printArray(items)

	items = slices.Delete(items, 0, 3)
	items = slices.Clip(items)
	printArray(items)

	items = slices.Insert(items, len(items), 10)
	items = slices.Insert(items, len(items), 20)
	printArray(items)

	if slices.Contains(items, 20) {
		fmt.Println("Found 20")
	}
}

func mapExample() {
	grades := map[string]int{
		"Bob":   90,
		"Jane":  92,
		"Alice": 85,
	}

	fmt.Printf("%+v\n", grades)

	alice := grades["Alice"]
	fmt.Printf("%d\n", alice)

	foo := grades["Foo"]
	fmt.Printf("%d\n", foo)

	dan, ok := grades["Dan"]

	if !ok {
		fmt.Println("Dan does not exist.")
	} else {
		fmt.Printf("%d\n", dan)
	}
}

func main() {
	// printArrayExample()
	// sliceExample()
	// mapExample()
}
