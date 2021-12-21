package domain

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrZeroAmount        = errors.New("can't operate with zero values")
	ErrNegativeAmount    = errors.New("doesn't work with negative amounts")
	ErrOverflow          = errors.New("can't hold so big amount of money")
	ErrInsufficientFunds = errors.New("user hasn't enough money")
)

//todo: thick about interface
type User struct {
	ID     int64
	Amount int64
}

func (user *User) Deposit(amount int64) error {
	if amount == 0 {
		return fmt.Errorf("deposit error: <%w>", ErrZeroAmount)
	}
	if amount < 0 {
		return fmt.Errorf("deposit error: <%w>", ErrNegativeAmount)
	}
	if user.Amount > math.MaxInt64-amount {
		return fmt.Errorf("deposit error: <%w>", ErrOverflow)
	}
	user.Amount += amount
	return nil
}

func (user *User) Withdraw(amount int64) error {
	if amount == 0 {
		return fmt.Errorf("withdraw error: <%w>", ErrZeroAmount)
	}
	if amount < 0 {
		return fmt.Errorf("withdraw error: <%w>", ErrNegativeAmount)
	}
	if user.Amount < amount {
		return fmt.Errorf("withdraw error: <%w>", ErrInsufficientFunds)
	}
	user.Amount -= amount
	return nil
}
