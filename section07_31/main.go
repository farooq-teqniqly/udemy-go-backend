package main

import (
	"errors"
	"fmt"
	"time"
)

// custom errors
var ErrDivideByZero = errors.New("divide by zero")
var ErrNumTooLarge = errors.New("number to large")

// richer custom error
type AppError struct {
	Op      string
	Code    int
	Message string
	Time    time.Time
}

func (err AppError) Error() string {
	return err.Message
}

func NewAppError(op string, code int, message string, t time.Time) *AppError {
	return &AppError{
		Op:      op,
		Code:    code,
		Message: message,
		Time:    t,
	}
}

func DoSomething() error {
	return NewAppError("DoSomething", 100, "failed", time.Now())
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}

	if a > 1000 || b > 1000 {
		return 0, ErrNumTooLarge
	}

	return a / b, nil
}

func main() {
	// res, err := divide(1001, 1)

	// if errors.Is(err, ErrDivideByZero) {
	// 	fmt.Println("cannot divide by zero")
	// } else if errors.Is(err, ErrNumTooLarge){
	// 	fmt.Println("dividend and divisor must both be less than 1000")
	// } else {
	// 	fmt.Println(res)
	// }

	err := DoSomething()
	fmt.Println(err)
	fmt.Printf("%#v\n", err)

	if appErr, ok := errors.AsType[*AppError](err); ok {
		fmt.Printf("op=%s code=%d msg=%s time=%s\n", appErr.Op, appErr.Code, appErr.Message, appErr.Time)
	}

}
