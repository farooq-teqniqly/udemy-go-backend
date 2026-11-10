package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	divisionOperation = "division"
	divideByZeroErrorMessage = "cannot divide by zero"
)

type MathError struct {
	Operation string
	InputX int
	InputY int
	Message string
}

func (e MathError) Error() string {
	var inputs []string

	if (e.Operation == divisionOperation) {
		inputs = append(inputs, fmt.Sprintf("x=%d", e.InputX))
		inputs = append(inputs, fmt.Sprintf("x=%d", e.InputY))
	}

	return fmt.Sprintf("Math error in %s (%s): %s", e.Operation, strings.Join(inputs, ","), e.Message)

}

func sum (numbers ...int) int {
	defer fmt.Println("sum finished")

	total := 0

	for _, n := range numbers {
		total += n
	}

	return total
}

func divide(x, y int) (int, error) {
	if (y != 0) {
		return x / y, nil
	}

	return 0, &MathError{
			Operation: divisionOperation,
			InputX: x,
			InputY: y,
			Message: divideByZeroErrorMessage,
		}
}

func main() {
	// fmt.Println(sum(1, 2, 3, 4))

	res, err := divide(10, 0)

	if mathErr, ok := errors.AsType[*MathError](err); ok {
		fmt.Println(mathErr)
	} else {
		fmt.Println(res)
	}
}
