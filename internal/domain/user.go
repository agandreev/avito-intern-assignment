package domain

import (
	"errors"
	"fmt"
	"math"
)

var (
	eps = math.Nextafter(1, 2) - 1

	ErrZeroAmount        = errors.New("can't operate with zero values")
	ErrNegativeAmount    = errors.New("doesn't work with negative amounts")
	ErrOverflow          = errors.New("can't hold so big amount of money")
	ErrInsufficientFunds = errors.New("user hasn't enough money")
)

// User represents user entity in our service.
type User struct {
	ID     int64   `json:"id"`
	Amount float64 `json:"amount,omitempty"`
}

// Deposit increases User's amount.
func (user *User) Deposit(amount float64) error {
	// check for machine zero or minimal value
	if amount >= 0 && amount < eps {
		return fmt.Errorf("deposit error: <%w>", ErrZeroAmount)
	}
	// check for negative case
	if amount < 0 {
		return fmt.Errorf("deposit error: <%w>", ErrNegativeAmount)
	}
	// check for overflow
	if user.Amount > math.MaxFloat64-amount {
		return fmt.Errorf("deposit error: <%w>", ErrOverflow)
	}
	user.Amount += amount
	return nil
}

// Withdraw decreases User's amount.
func (user *User) Withdraw(amount float64) error {
	// check for machine zero or minimal value
	if amount >= 0 && amount < eps {
		return fmt.Errorf("withdraw error: <%w>", ErrZeroAmount)
	}
	// check for negative case
	if amount < 0 {
		return fmt.Errorf("withdraw error: <%w>", ErrNegativeAmount)
	}
	// check for overflow
	if user.Amount < amount {
		return fmt.Errorf("withdraw error: <%w>", ErrInsufficientFunds)
	}
	user.Amount -= amount
	return nil
}
