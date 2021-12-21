package domain

import (
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
)

var (
	zeroValue int64 = 0
	oneValue  int64 = 1
)

type UserSuite struct {
	suite.Suite
	User User
}

func (suite *UserSuite) SetupTest() {
	suite.User = User{
		ID:     0,
		Amount: 0,
	}
}

func (suite UserSuite) TestUser_Deposit() {
	// negative value case
	suite.ErrorIs(suite.User.Deposit(-1), ErrNegativeAmount)
	suite.Equal(zeroValue, suite.User.Amount)
	// zero value case
	suite.ErrorIs(suite.User.Deposit(0), ErrZeroAmount)
	suite.Equal(zeroValue, suite.User.Amount)
	// overflow case
	suite.User.Amount = 1
	suite.ErrorIs(suite.User.Deposit(math.MaxInt64), ErrOverflow)
	suite.Equal(oneValue, suite.User.Amount)
	suite.User.Amount = 0
	// normal cases
	suite.NoError(suite.User.Deposit(1))
	suite.Equal(oneValue, suite.User.Amount)
	suite.User.Amount = 0

	suite.NoError(suite.User.Deposit(math.MaxInt64))
	suite.Equal(int64(math.MaxInt64), suite.User.Amount)
}

func (suite UserSuite) TestUser_Withdraw() {
	// negative value case
	suite.ErrorIs(suite.User.Withdraw(-1), ErrNegativeAmount)
	suite.Equal(zeroValue, suite.User.Amount)
	// zero value case
	suite.ErrorIs(suite.User.Withdraw(0), ErrZeroAmount)
	suite.Equal(zeroValue, suite.User.Amount)
	// insufficient case
	suite.ErrorIs(suite.User.Withdraw(suite.User.Amount+1), ErrInsufficientFunds)
	suite.Equal(zeroValue, suite.User.Amount)
	// normal cases
	suite.User.Amount = math.MaxInt64
	suite.NoError(suite.User.Withdraw(1))
	suite.Equal(int64(math.MaxInt64-1), suite.User.Amount)

	suite.NoError(suite.User.Withdraw(math.MaxInt64 - 1))
	suite.Equal(zeroValue, suite.User.Amount)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
