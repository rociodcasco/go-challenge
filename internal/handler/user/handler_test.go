package userHandler_test

import (
	"encoding/json"
	transferHandler "go-challenge/internal/handler/transfer"
	userHandler "go-challenge/internal/handler/user"
	userHandler_mock "go-challenge/internal/handler/user/mocks"
	"go-challenge/internal/router"
	"go-challenge/pkg/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -package userHandler_mock -source=handler.go -destination=./mocks/userHandler_mock.go


type UserHandlerTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	managerMock *userHandler_mock.MockUserManager
	handler *userHandler.UserHandler
	router *gin.Engine
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.managerMock = userHandler_mock.NewMockUserManager(suite.ctrl)

	handler, err := userHandler.NewUserHandler(suite.managerMock)
	suite.NoError(err, "Failed to create user handler")
	suite.handler = handler

	suite.router = router.SetupRouter(suite.handler, &transferHandler.TransferHandler{})
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) TestNewUserHandler() {
	suite.Run("nil user manager", func() {
		handler, err := userHandler.NewUserHandler(nil)
		suite.Nil(handler, "Expected nil handler")
		suite.Error(err, "Expected error when user manager is nil")
	})

	suite.Run("valid user manager", func() {
		handler, err := userHandler.NewUserHandler(suite.managerMock)
		suite.NotNil(handler, "Expected non-nil handler")
		suite.NoError(err, "Expected no error when user manager is valid")
	})
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
	suite.Run("invalid user", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/", nil)
		suite.router.ServeHTTP(w, req)
		suite.Equal(400, w.Code)
		suite.Contains(w.Body.String(), `{"error":"invalid request"}`)
	})
	
	suite.Run("valid user", func() {
		user := &user.User{
			Name: "John Doe",
			DNI: "12345678",
			Email: "jdoe@gmail.com",
		}
		suite.managerMock.EXPECT().CreateUser(user).Return(uint(1), nil).Times(1)

		userJson, _ := json.Marshal(user)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/", strings.NewReader(string(userJson)))
		suite.router.ServeHTTP(w, req)
		suite.Equal(200, w.Code)
	})
}

func (suite *UserHandlerTestSuite) TestGetUserBalance() {
	suite.Run("invalid user id", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/invalid_id/balance", nil)
		suite.router.ServeHTTP(w, req)
		suite.Equal(400, w.Code)
		suite.Contains(w.Body.String(), `{"error":"invalid user id"}`)
	})

	suite.Run("valid user id", func() {
		userIdStr := "1"
		userId := uint(1)
		balance := uint(1000)
		suite.managerMock.EXPECT().GetUserBalance(userId).Return(balance, nil).Times(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/"+userIdStr+"/balance", nil)
		suite.router.ServeHTTP(w, req)
		suite.Equal(200, w.Code)
	})
}