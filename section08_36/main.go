package main

import (
	"fmt"
	"time"
	"uuid"
)

type Employee struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Salary    float64
	IsActive  bool
	JoinedAt  time.Time
}

func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

func (e *Employee) Deactivate() {
	e.IsActive = false
}

func NewEmployee(firstName, lastName string, salary float64, isActive bool, joinedAt time.Time) *Employee {
	return &Employee{
		ID:        uuid.New(),
		LastName:  lastName,
		FirstName: firstName,
		Salary:    salary,
		IsActive:  isActive,
		JoinedAt:  joinedAt,
	}
}

func main() {
	employee := NewEmployee("Farooq", "Mahmud", 150000, true, time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC))

	fmt.Printf("%+v\n", *employee)
	fmt.Println(employee.FullName())
	employee.Deactivate()
	fmt.Printf("%+v\n", *employee)
}
