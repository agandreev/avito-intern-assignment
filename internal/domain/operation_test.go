package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type OperationSuite struct {
	suite.Suite
	Operation Operation
}

func (suite *OperationSuite) SetupTest() {
	suite.Operation = Operation{}
}

func (suite OperationSuite) TestOperation_IsTransfer() {
	suite.Operation.Type = TransferIn
	suite.True(suite.Operation.IsTransfer())
	suite.Operation.Type = TransferOut
	suite.True(suite.Operation.IsTransfer())
	suite.Operation.Type = Deposit
	suite.False(suite.Operation.IsTransfer())
	suite.Operation.Type = Withdraw
	suite.False(suite.Operation.IsTransfer())
}

func (suite OperationSuite) TestOperation_Reverse() {
	_, err := suite.Operation.Reverse()
	suite.Error(err)

	suite.Operation = Operation{
		Initiator: &User{},
		Type:      TransferIn,
		Amount:    0,
		Timestamp: time.Time{},
		Receiver:  &User{},
	}
	_, err = suite.Operation.Reverse()
	suite.NoError(err)
}

func (suite OperationSuite) TestOperation_Validate() {
	suite.Operation.Type = ""
	suite.ErrorIs(suite.Operation.Validate(), ErrIncorrectOperationParams)

	suite.Operation.Type = Deposit
	suite.ErrorIs(suite.Operation.Validate(), ErrIncorrectOperationParams)

	suite.Operation.Initiator = &User{}
	suite.NoError(suite.Operation.Validate())
	suite.Operation.Receiver = &User{}
	suite.ErrorIs(suite.Operation.Validate(), ErrIncorrectOperationParams)

	suite.Operation.Type = TransferIn
	suite.Operation.Receiver = nil
	suite.ErrorIs(suite.Operation.Validate(), ErrIncorrectOperationParams)

	suite.Operation.Receiver = &User{}
	suite.NoError(suite.Operation.Validate())
}

func TestOperationSuite(t *testing.T) {
	suite.Run(t, new(OperationSuite))
}
