package service

import (
	"avitoInternAssignment/internal/domain"
	"avitoInternAssignment/internal/repository"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// GrossBookRepository combines UserRepository and OperationRepository.
type GrossBookRepository interface {
	UserRepository
	OperationRepository
	Shutdown()
}

// UserRepository describes UserStorage methods.
type UserRepository interface {
	AddUser(id int64) error
	GetUser(id int64) (*domain.User, error)
}

// OperationRepository describes UserStorage methods.
type OperationRepository interface {
	AddOperation(operation domain.Operation) error
	GetOperations(id, offset int64, mode domain.SortingMode) ([]domain.RepositoryOperation, error)
}

// Converter converts amount of money from one currency to another
type Converter interface {
	Convert(from string, amount float64) (float64, error)
}

// GrossBook represents this service logic.
type GrossBook struct {
	Users    GrossBookRepository
	Exchange Converter
	log      *logrus.Logger
}

// NewGrossBook sets GrossBook fields and returns pointer.
func NewGrossBook(users GrossBookRepository, exchange Converter,
	log *logrus.Logger) *GrossBook {
	return &GrossBook{
		Users:    users,
		Exchange: exchange,
		log:      log,
	}
}

// DepositMoney deposits money and updates db.
func (grossBook *GrossBook) DepositMoney(id int64, amount float64) (
	*domain.Operation, error) {
	grossBook.log.Printf("DEPOSIT: <%f>RUB to <%d> processing...", amount, id)
	user, err := grossBook.Users.GetUser(id)
	if err != nil {
		switch err {
		case repository.ErrNoSuchUser:
			if err = grossBook.Users.AddUser(id); err != nil {
				return nil, fmt.Errorf("grossbook get user error: <%w>", err)
			}
		default:
			return nil, fmt.Errorf("grossbook get user error: <%w>", err)
		}
	}
	operation := domain.Operation{
		Initiator: user,
		Type:      domain.Deposit,
		Amount:    amount,
	}
	if err = user.Deposit(amount); err != nil {
		return nil, fmt.Errorf("grossbook deposit error: <%w>", err)
	}
	operation.Timestamp = time.Now()
	if err = grossBook.Users.AddOperation(operation); err != nil {
		return nil, fmt.Errorf("grossbook update error: <%w>", err)
	}
	grossBook.log.Printf("DEPOSIT: <%f>RUB from <%d> was processed successful",
		amount, id)
	return &operation, nil
}

// WithdrawMoney withdraws money and updates db.
func (grossBook *GrossBook) WithdrawMoney(id int64, amount float64, currency string) (
	*domain.Operation, error) {
	grossBook.log.Printf("WITHDRAW: <%f> from <%d> processing...", amount, id)
	user, err := grossBook.Users.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("grossbook get user error: <%w>", err)
	}
	operation := domain.Operation{
		Initiator: user,
		Type:      domain.Withdraw,
		Amount:    amount,
	}
	if len(currency) != 0 {
		if amount, err = grossBook.Exchange.Convert(currency, amount); err != nil {
			return nil, fmt.Errorf("gorssbook withdraw conversion error: <%w>", err)
		}
		operation.Amount = amount
	}
	if err = user.Withdraw(amount); err != nil {
		return nil, fmt.Errorf("grossbook withdraw error: <%w>", err)
	}
	operation.Timestamp = time.Now()
	if err = grossBook.Users.AddOperation(operation); err != nil {
		return nil, fmt.Errorf("grossbook update error: <%w>", err)
	}
	grossBook.log.Printf("WITHDRAW: <%f>RUB from <%d> was processed successful",
		amount, id)
	return &operation, nil
}

// TransferMoney transfers money and updates db.
func (grossBook *GrossBook) TransferMoney(ownerID, receiverID int64, amount float64) (
	*domain.Operation, error) {
	grossBook.log.Printf("TRANSFER: <%f>RUB from <%d> to <%d> processing...",
		amount, ownerID, receiverID)
	owner, err := grossBook.Users.GetUser(ownerID)
	if err != nil {
		return nil, fmt.Errorf("grossbook get owner error: <%w>", err)
	}
	receiver, err := grossBook.Users.GetUser(receiverID)
	if err != nil {
		return nil, fmt.Errorf("grossbook get receiver error: <%w>", err)
	}
	operation := domain.Operation{Initiator: owner,
		Type:   domain.TransferOut,
		Amount: amount,
	}
	if owner.ID == receiverID {
		return nil, fmt.Errorf("grossbook can't transfer money for the same user")
	}
	if err = owner.Withdraw(amount); err != nil {
		return nil, fmt.Errorf("grossbook owner withdraw error: <%w>", err)
	}
	if err = receiver.Deposit(amount); err != nil {
		return nil, fmt.Errorf("grossbook receiver deposit error: <%w>", err)
	}
	operation.Timestamp = time.Now()
	operation.Receiver = receiver
	if err = grossBook.Users.AddOperation(operation); err != nil {
		return nil, fmt.Errorf("grossbook transfer update error: <%w>", err)
	}
	grossBook.log.Printf("TRANSFER: <%f>RUB from <%d> to <%d> was processed successful",
		amount, ownerID, receiverID)
	// hide amount for safety
	operation.Receiver = &domain.User{ID: receiverID}
	return &operation, nil
}

// Balance returns domain.User's balance from db.
func (grossBook GrossBook) Balance(id int64) (*domain.User, error) {
	grossBook.log.Printf("BALANCE: by <%d> processing...", id)
	user, err := grossBook.Users.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("grossbook get owner error: <%w>", err)
	}
	grossBook.log.Printf("BALANCE: by <%d> was processed successful", id)
	return user, nil
}

// History returns slice of domain.RepositoryOperation from db.
func (grossBook GrossBook) History(id, offset int64, mode domain.SortingMode) (
	[]domain.RepositoryOperation, error) {
	grossBook.log.Printf("HISTORY: by <%d> processing...", id)
	if _, err := grossBook.Users.GetUser(id); err != nil {
		return nil, fmt.Errorf("can't load history: <%w>", err)
	}
	operations, err := grossBook.Users.GetOperations(id, offset, mode)
	if err != nil {
		return nil, fmt.Errorf("can't load history: <%w>", err)
	}
	grossBook.log.Printf("HISTORY: by <%d> was processed successful", id)
	return operations, nil
}

// Shutdown gracefully shuts this service down.
func (grossBook GrossBook) Shutdown() {
	grossBook.Users.Shutdown()
}
