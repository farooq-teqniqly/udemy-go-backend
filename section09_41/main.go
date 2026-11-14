package main

import (
	"errors"
	"fmt"
	"strings"
)

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

func NewAddress(street, city, state, zipCode string) (Address, error) {
	street = strings.TrimSpace(street)
	if street == "" {
		return Address{}, errors.New("address not provided")
	}

	city = strings.TrimSpace(city)
	if city == "" {
		return Address{}, errors.New("city not provided")
	}

	state = strings.TrimSpace(state)
	if len(state) != 2 {
		return Address{}, errors.New("state must be two letter abbreviation")
	}

	zipCode = strings.TrimSpace(zipCode)
	if len(zipCode) != 5 {
		return Address{}, errors.New("zip code must be five characters")
	}

	return Address{
		Street:  street,
		City:    city,
		State:   state,
		ZipCode: zipCode,
	}, nil
}

func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s %s", a.Street, a.City, a.State, a.ZipCode)
}

type Customer struct {
	CustomerID      int
	Name            string
	Email           string
	BillingAddress  Address
	ShippingAddress Address
}

func NewCustomer(id int, name, email string, billingAddress, shippingAddress Address) (Customer, error) {
	if id < 1 {
		return Customer{}, fmt.Errorf("invalid id value: '%d'; id should be a positive integer", id)
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Customer{}, errors.New("name not provided")
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return Customer{}, errors.New("email not provided")
	}

	if !strings.Contains(email, "@") {
		return Customer{}, errors.New("email does not appear to be valid")
	}

	return Customer{
		CustomerID:      id,
		Name:            name,
		Email:           email,
		BillingAddress:  billingAddress,
		ShippingAddress: shippingAddress,
	}, nil
}

func (c Customer) String() string {
	return fmt.Sprintf("%s (Email: %s, Billing address: %s, Shipping address: %s)", c.Name, c.Email, c.BillingAddress, c.ShippingAddress)
}

func main() {
	address, err := NewAddress("123 Main St", "Bothell", "WA", "98029")

	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := NewCustomer(1001, "Farooq Mahmud", "farooq@foo.com", address, address)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(c)
}
