package domain

import (
	"errors"
	"fmt"
	"time"
)

const (
	Deposit     OperationType = "DEPOSIT"
	Withdraw    OperationType = "WITHDRAW"
	TransferOut OperationType = "TRANSFER OUT"
	TransferIn  OperationType = "TRANSFER IN"
)

var (
	ErrNonTransferOperation     = errors.New("it isn't transfer operation")
	ErrIncorrectOperationParams = errors.New("this operation is incorrect")
)

// OperationType describes type of Operation.
type OperationType string

// Operation represents a transaction event.
type Operation struct {
	Initiator *User         `json:"initiator"`
	Type      OperationType `json:"type"`
	Amount    float64       `json:"amount"`
	Timestamp time.Time     `json:"timestamp"`
	Receiver  *User         `json:"receiver,omitempty"`
}

// RepositoryOperation is restricted type of Operation for Repository aims.
type RepositoryOperation struct {
	InitiatorID int64         `json:"initiator_id"`
	Type        OperationType `json:"type"`
	Amount      float64       `json:"amount"`
	Timestamp   time.Time     `json:"timestamp"`
	ReceiverID  int64         `json:"receiver_id,omitempty"`
}

// IsTransfer returns true if Operation type is Transfer and false otherwise.
func (operation Operation) IsTransfer() bool {
	return operation.Type == TransferIn || operation.Type == TransferOut
}

// Validate is necessary in order to correlate field values and type value.
func (operation Operation) Validate() error {
	// check type
	if operation.Type != Deposit && operation.Type != Withdraw &&
		operation.Type != TransferIn && operation.Type != TransferOut {
		return fmt.Errorf("incorrect operation type: <%w>", ErrIncorrectOperationParams)
	}
	// check correlation of type and users' quantity
	if operation.Initiator == nil {
		return fmt.Errorf("initiator can't be nil: <%w>", ErrIncorrectOperationParams)
	}
	if operation.IsTransfer() {
		if operation.Receiver == nil {
			return fmt.Errorf("receiver can't be nil in transfer operation: <%w>",
				ErrIncorrectOperationParams)
		}
	} else {
		if operation.Receiver != nil {
			return fmt.Errorf("receiver can't be non nil in "+
				"non transfer operation: <%w>", ErrIncorrectOperationParams)
		}
	}
	return nil
}

// Reverse changes Operation type on the opposite if it is transfer and switch users.
func (operation Operation) Reverse() (*Operation, error) {
	if err := operation.Validate(); err != nil {
		return nil, fmt.Errorf("operation's validation is failed: <%w>", err)
	}
	if !operation.IsTransfer() {
		return nil, fmt.Errorf("can't reverse non duplex operation: <%w>", ErrNonTransferOperation)
	}
	reversed := Operation{
		Amount:    operation.Amount,
		Timestamp: operation.Timestamp,
	}
	reversed.Initiator = operation.Receiver
	reversed.Receiver = operation.Initiator
	switch operation.Type {
	case TransferIn:
		reversed.Type = TransferOut
	case TransferOut:
		reversed.Type = TransferIn
	default:
		return nil, fmt.Errorf("unsupported operation type: <%s>", operation.Type)
	}
	return &reversed, nil
}
