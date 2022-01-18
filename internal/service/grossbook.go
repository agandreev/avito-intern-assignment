package service

import (
	"github.com/agandreev/avito-intern-assignment/internal/domain"
	"github.com/agandreev/avito-intern-assignment/internal/repository"
	"context"
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
	User(id int64) (*domain.User, error)
}

// OperationRepository describes UserStorage methods.
type OperationRepository interface {
	AddOperation(ctx context.Context, operation domain.Operation) error
	Operations(id, offset int64, mode domain.SortingMode) ([]domain.RepositoryOperation, error)
}

// Converter converts amount of money from one currency to RUB.
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

// DepositMoney increases user balance by id and updates db.
func (grossBook *GrossBook) DepositMoney(id int64, amount float64) (
	*domain.Operation, error) {
	grossBook.log.Printf("DEPOSIT: <%f>RUB to <%d> processing...", amount, id)
	// get user or create it
	user, err := grossBook.Users.User(id)
	if err != nil {
		switch err {
		// create empty raw in db
		case repository.ErrNoSuchUser:
			if err = grossBook.Users.AddUser(id); err != nil {
				return nil, fmt.Errorf("grossbook get user error: <%w>", err)
			}
		default:
			return nil, fmt.Errorf("grossbook get user error: <%w>", err)
		}
	}
	// increase User's amount
	if err = user.Deposit(amount); err != nil {
		return nil, fmt.Errorf("grossbook deposit error: <%w>", err)
	}
	operation := domain.Operation{
		Initiator: user,
		Type:      domain.Deposit,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	// update db
	if err = grossBook.Users.AddOperation(context.Background(), operation); err != nil {
		return nil, fmt.Errorf("grossbook update error: <%w>", err)
	}
	grossBook.log.Printf("DEPOSIT: <%f>RUB from <%d> was processed successful",
		amount, id)
	return &operation, nil
}

// WithdrawMoney decreases domain.User's balance and updates db.
func (grossBook *GrossBook) WithdrawMoney(id int64, amount float64, currency string) (
	*domain.Operation, error) {
	grossBook.log.Printf("WITHDRAW: <%f> from <%d> processing...", amount, id)
	// get user
	user, err := grossBook.Users.User(id)
	if err != nil {
		return nil, fmt.Errorf("grossbook get user error: <%w>", err)
	}
	// convert amount to RUB
	if len(currency) != 0 {
		convertedAmount, err := grossBook.Exchange.Convert(currency, amount)
		if err != nil {
			return nil, fmt.Errorf("gorssbook withdraw conversion error: <%w>", err)
		}
		amount = convertedAmount
	}
	// decrease user's balance
	if err = user.Withdraw(amount); err != nil {
		return nil, fmt.Errorf("grossbook withdraw error: <%w>", err)
	}
	operation := domain.Operation{
		Initiator: user,
		Type:      domain.Withdraw,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	// update db
	if err = grossBook.Users.AddOperation(context.Background(), operation); err != nil {
		return nil, fmt.Errorf("grossbook update error: <%w>", err)
	}
	grossBook.log.Printf("WITHDRAW: <%f>RUB from <%d> was processed successful",
		amount, id)
	return &operation, nil
}

// TransferMoney transfers money from one domain.User to another and updates db.
func (grossBook *GrossBook) TransferMoney(ownerID, receiverID int64, amount float64) (
	*domain.Operation, error) {
	grossBook.log.Printf("TRANSFER: <%f>RUB from <%d> to <%d> processing...",
		amount, ownerID, receiverID)
	// get users
	owner, err := grossBook.Users.User(ownerID)
	if err != nil {
		return nil, fmt.Errorf("grossbook get owner error: <%w>", err)
	}
	receiver, err := grossBook.Users.User(receiverID)
	if err != nil {
		return nil, fmt.Errorf("grossbook get receiver error: <%w>", err)
	}
	if owner.ID == receiverID {
		return nil, fmt.Errorf("grossbook can't transfer money for the same user")
	}
	// decrease and increase balances
	if err = owner.Withdraw(amount); err != nil {
		return nil, fmt.Errorf("grossbook owner withdraw error: <%w>", err)
	}
	if err = receiver.Deposit(amount); err != nil {
		return nil, fmt.Errorf("grossbook receiver deposit error: <%w>", err)
	}
	operation := domain.Operation{
		Initiator: owner,
		Type:      domain.TransferOut,
		Amount:    amount,
		Timestamp: time.Now(),
		Receiver:  receiver,
	}
	// update db
	if err = grossBook.Users.AddOperation(context.Background(), operation); err != nil {
		return nil, fmt.Errorf("grossbook transfer update error: <%w>", err)
	}
	grossBook.log.Printf("TRANSFER: <%f>RUB from <%d> to <%d> was processed successful",
		amount, ownerID, receiverID)
	// hide second side amount for safety
	operation.Receiver = &domain.User{ID: receiverID}
	return &operation, nil
}

// Balance returns domain.User's balance from db.
func (grossBook GrossBook) Balance(id int64) (*domain.User, error) {
	grossBook.log.Printf("BALANCE: by <%d> processing...", id)
	user, err := grossBook.Users.User(id)
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
	if _, err := grossBook.Users.User(id); err != nil {
		return nil, fmt.Errorf("can't load history: <%w>", err)
	}
	operations, err := grossBook.Users.Operations(id, offset, mode)
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
