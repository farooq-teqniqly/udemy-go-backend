package main

import (
	"errors"
	"fmt"
	"strings"
)

type Account struct {
	AccountNumber string
	balance float64
	OwnerName string
}

func NewAccount(accountNumber, ownerName string, balance float64) (*Account, error) {
	trimmedAccountNumber := strings.TrimSpace(accountNumber)

	if (trimmedAccountNumber == "") {
		return nil, errors.New("account number not specified")
	}

	trimmedOwnerName := strings.TrimSpace(ownerName)

	if (trimmedOwnerName == "") {
		return nil, errors.New("owner name not specified")
	}

	if (balance < 0) {
		return nil, errors.New("opening balance cannot be negative")
	}

	return &Account{
		AccountNumber: trimmedAccountNumber,
		OwnerName: trimmedOwnerName,
		balance: balance,
	}, nil
}

func (a *Account) String() string {
	return fmt.Sprintf("%s (%s, %.2f)", a.OwnerName, a.AccountNumber, a.Balance())
}

func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("deposit amount must be positive")
	}

	a.balance += amount
	return nil
}

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

func (a *Account) Balance() float64 {
	return a.balance
}

type SavingsAccount struct {
	Account
	InterestRate float64
}

func NewSavingsAccount(accountNumber, ownerName string, balance float64, rate float64) (*SavingsAccount, error) {
	if rate < 0 {
		return nil, errors.New("interest rate must be positive")
	}

	a, err := NewAccount(accountNumber, ownerName, balance)

	if err != nil {
		return nil, fmt.Errorf("could not create savings account: %w", err)
	}

	return &SavingsAccount{
		Account: *a,
		InterestRate: rate,
	}, nil
}

func (sa *SavingsAccount) String() string {
	return fmt.Sprintf("%s (%s, %.2f, %.2f)", sa.OwnerName, sa.AccountNumber, sa.InterestRate, sa.Balance())
}

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
