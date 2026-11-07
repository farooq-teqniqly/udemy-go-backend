package main

import "fmt"

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	levelCount
)

var levels = [levelCount]string{LevelTrace: "Trace", LevelDebug: "Debug", LevelInfo: "Info", LevelWarn: "Warning", LevelError: "Error"}

func (l LogLevel) String() string {
	if l < LevelTrace || l >= levelCount {
		return "Unknown"
	}

	return levels[l]
}

func printLevel(level LogLevel) {
	fmt.Printf("[%s:%d]\n", level.String(), level)
}

func main() {
	for l := range levelCount {
		printLevel(l)
	}

	printLevel(100)
	printLevel(-100)
}
