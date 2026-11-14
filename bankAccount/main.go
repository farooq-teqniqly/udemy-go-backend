// Package main is a small exercise that models a bank account and a savings account
// that embeds it.
package main

import (
	"errors"
	"fmt"
	"strings"
)

// Account is a bank account held by one owner. The balance is unexported, so it can
// only change through Deposit and Withdraw (and only code in this package can touch
// the field directly).
type Account struct {
	AccountNumber string
	balance       float64
	OwnerName     string
}

// NewAccount returns an Account with the given opening balance. It trims surrounding
// whitespace from accountNumber and ownerName, and returns an error if either is
// blank or if balance is negative.
func NewAccount(accountNumber, ownerName string, balance float64) (*Account, error) {
	trimmedAccountNumber := strings.TrimSpace(accountNumber)

	if trimmedAccountNumber == "" {
		return nil, errors.New("account number not specified")
	}

	trimmedOwnerName := strings.TrimSpace(ownerName)

	if trimmedOwnerName == "" {
		return nil, errors.New("owner name not specified")
	}

	if balance < 0 {
		return nil, errors.New("opening balance cannot be negative")
	}

	return &Account{
		AccountNumber: trimmedAccountNumber,
		OwnerName:     trimmedOwnerName,
		balance:       balance,
	}, nil
}

// String returns the owner name, account number, and balance.
func (a *Account) String() string {
	return fmt.Sprintf("%s (%s, %.2f)", a.OwnerName, a.AccountNumber, a.Balance())
}

// Deposit adds amount to the balance. It returns an error, and leaves the balance
// unchanged, if amount is negative.
func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("deposit amount must be positive")
	}

	a.balance += amount
	return nil
}

// Withdraw subtracts amount from the balance. It returns an error, and leaves the
// balance unchanged, if amount is negative or larger than the balance.
func (a *Account) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("withdrawal amount must be positive")
	}

	remaining := a.balance - amount

	if remaining < 0 {
		return errors.New("insufficient balance to complete withdrawal - choose a smaller amount")
	}

	a.balance = remaining
	return nil
}

// Balance returns the current balance.
func (a *Account) Balance() float64 {
	return a.balance
}

// SavingsAccount is an Account that earns interest. It embeds Account, so Deposit,
// Withdraw, and Balance are promoted to it.
type SavingsAccount struct {
	Account
	// InterestRate is applied as a fraction, so 0.025 means 2.5%.
	InterestRate float64
}

// NewSavingsAccount returns a SavingsAccount. It returns an error if rate is negative
// or if the underlying Account cannot be created; in that case the error from
// NewAccount is wrapped and can be inspected with errors.Is or errors.As.
func NewSavingsAccount(accountNumber, ownerName string, balance float64, rate float64) (*SavingsAccount, error) {
	if rate < 0 {
		return nil, errors.New("interest rate must be positive")
	}

	a, err := NewAccount(accountNumber, ownerName, balance)

	if err != nil {
		return nil, fmt.Errorf("could not create savings account: %w", err)
	}

	return &SavingsAccount{
		Account:      *a,
		InterestRate: rate,
	}, nil
}

// String returns the owner name, account number, interest rate, and balance. It hides
// the String method promoted from Account.
func (sa *SavingsAccount) String() string {
	return fmt.Sprintf("%s (%s, %.2f, %.2f)", sa.OwnerName, sa.AccountNumber, sa.InterestRate, sa.Balance())
}

// AddInterest deposits the current balance multiplied by InterestRate. The error from
// Deposit, if any, is wrapped.
func (sa *SavingsAccount) AddInterest() error {
	interest := sa.balance * sa.InterestRate
	err := sa.Deposit(interest)

	if err != nil {
		return fmt.Errorf("error depositing interest: %w", err)
	}

	return nil
}

func main() {
	sa, err := NewSavingsAccount("12345", "Farooq", 20000.00, 0.025)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sa)

	a, err := NewAccount("45678", "Alice", 1000.00)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(a)
}
