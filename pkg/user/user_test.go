package user_test

import (
	"fmt"
	"go-challenge/pkg/user"
	"testing"

	user_mocks "go-challenge/pkg/user/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -package user_mocks -source=user.go -destination=./mocks/user_mock.go


type UserTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	storageMock *user_mocks.MockStorage
	userManager *user.UserManager
}

func (suite *UserTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.storageMock = user_mocks.NewMockStorage(suite.ctrl)

	manager, err := user.NewUserManager(suite.storageMock)
	suite.NoError(err, "Failed to create user manager")
	suite.userManager = manager
}

func TestUserTestSuite(t *testing.T) {	
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestNewManager() {
	suite.Run("nil storage", func() {
		manager, err := user.NewUserManager(nil)
		suite.Nil(manager, "Expected nil manager")
		suite.Error(err, "Expected error when storage is nil")
	})
	suite.Run("valid storage", func() {
		manager, err := user.NewUserManager(suite.storageMock)
		suite.NotNil(manager, "Expected non-nil manager")
		suite.NoError(err, "Expected no error when storage is valid")
	})
}

func (suite *UserTestSuite) TestCreateUser() {
	suite.Run("nil user", func() {
		id, err := suite.userManager.CreateUser(nil)
		suite.Equal(uint(0), id, "Expected id to be 0")
		suite.Error(err, "Expected error when user is nil")
	})

	suite.Run("valid user", func() {
		expectedID := uint(1)
		user := &user.User{
			Name: "John Doe",
			DNI: "12345678",
			Email: "jdoe@gmail.com",
		}
		suite.storageMock.EXPECT().CreateUser(user).Return(expectedID, nil)

		id, err := suite.userManager.CreateUser(user)
		suite.Equal(expectedID, id, "Expected id to be 1")
		suite.NoError(err, "Expected no error when creating user")
	})
}

func (suite *UserTestSuite) TestGetUserBalance() {
	suite.Run("user not found", func() {
		userID := uint(1)
		errExpected := fmt.Errorf("user not found")
		suite.storageMock.EXPECT().GetUserByID(userID).Return(nil, errExpected)

		balance, err := suite.userManager.GetUserBalance(userID)
		suite.Equal(uint(0), balance, "Expected balance to be 0")
		suite.Error(err, "Expected error when user is not found")
		suite.Equal(err.Error(), fmt.Errorf("user not found: %w", errExpected).Error(), "Expected error to match")
	})

	suite.Run("valid user", func() {
		userID := uint(1)
		expectedBalance := uint(1000)
		user := &user.User{
			Balance: expectedBalance,
		}
		suite.storageMock.EXPECT().GetUserByID(userID).Return(user, nil)

		balance, err := suite.userManager.GetUserBalance(userID)
		suite.Equal(expectedBalance, balance, "Expected balance to be 1000")
		suite.NoError(err, "Expected no error when getting user balance")
	})
}