package main

import "fmt"

type Payable interface {
	MonthlyPay() float64
	fmt.Stringer
}

type SalariedEmployee struct {
	Name         string
	AnnualSalary float64
}

func (e *SalariedEmployee) MonthlyPay() float64 {
	return e.AnnualSalary / 12.0
}

func (e *SalariedEmployee) String() string {
	return fmt.Sprintf("%s (Salary: %.2f, MonthlyPay: %.2f)", e.Name, e.AnnualSalary, e.MonthlyPay())
}

type HourlyEmployee struct {
	Name        string
	Rate        float64
	HoursWorked float64
}

func (e *HourlyEmployee) MonthlyPay() float64 {
	return e.Rate * e.HoursWorked
}

func (e *HourlyEmployee) String() string {
	return fmt.Sprintf("%s (Rate: %.2f, HoursWorked: %.2f, MonthlyPay: %.2f)", e.Name, e.Rate, e.HoursWorked, e.MonthlyPay())
}

type CommissionEmployee struct {
	Name             string
	BaseAnnualSalary float64
	Rate             float64
}

func (e *CommissionEmployee) MonthlyPay() float64 {
	return (e.BaseAnnualSalary / 12.0) * (1 + e.Rate)
}

func (e *CommissionEmployee) String() string {
	return fmt.Sprintf("%s (Salary: %.2f, Rate: %.2f, MonthlyPay: %.2f)", e.Name, e.BaseAnnualSalary, e.Rate, e.MonthlyPay())
}

func PrintEmployee(employee Payable) {
	fmt.Println(employee)
}

func processPayroll(employees []Payable) {
	fmt.Println("\n--- Processing Payroll ---")
	total := 0.0

	for _, emp := range employees {
		PrintEmployee(emp)
		total += emp.MonthlyPay()
	}

	fmt.Println("\n--- Total payroll for month ---")
	fmt.Printf("%.2f", total)
	fmt.Println("----------------------")
}

func main() {
	se := &SalariedEmployee{
		Name:         "Farooq",
		AnnualSalary: 200000,
	}

	ce := &CommissionEmployee{
		Name:             "Alice",
		BaseAnnualSalary: 120000,
		Rate:             0.10,
	}

	he := &HourlyEmployee{
		Name:        "Colleen",
		HoursWorked: 160,
		Rate:        25.00,
	}

	employees := []Payable{se, ce, he}

	processPayroll(employees)

}
