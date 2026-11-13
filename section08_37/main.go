package main

import (
	"fmt"
	"time"
	"uuid"
)

type Person interface {
	FullName() string
}

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

func (e *Employee) String() string {
	return fmt.Sprintf("%s (ID: %s, Salary: %.2f, Active: %t, Joined: %s)", e.FullName(), e.ID, e.Salary, e.IsActive, e.JoinedAt.Format("2006-01-02"))
}

func (e *Employee) Equal(other *Employee) bool {
	return e.ID == other.ID
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
	farooq1 := NewEmployee("Farooq", "Mahmud", 150000, true, time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC))
	farooq2 := NewEmployee("Farooq", "Mahmud", 150000, true, time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC))
	fmt.Println(farooq1 == farooq2) // false because pointers to two different objects in memory

	farooq2.ID = farooq1.ID
	fmt.Println(farooq1 == farooq2)     // false because pointers to two different objects in memory
	fmt.Println(farooq1.Equal(farooq2)) // true because same ID's.

	fmt.Println(farooq1)
	fmt.Println(farooq2)
}
