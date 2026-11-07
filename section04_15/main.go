package main

import "fmt"

// const (
// 	Sunday = iota + 1
// 	Monday
// 	Tuesday
// 	Wednesday
// 	Thursday
// 	Friday
// 	Saturday
// 	upperBound
// )

type LogLevel int

const (
	Error LogLevel = iota
	Warning
	Info
	Debug
	Fatal
	upperBound
)

func main() {
	for l := range upperBound {
		fmt.Println(l)
	}
}
