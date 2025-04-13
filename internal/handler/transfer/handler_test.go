package transferHandler_test

import (
	"encoding/json"
	"fmt"
	transferHandler "go-challenge/internal/handler/transfer"
	transfer_mocks "go-challenge/internal/handler/transfer/mocks"
	userHandler "go-challenge/internal/handler/user"
	"go-challenge/internal/router"
	"go-challenge/pkg/transfer"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -package transferHandler_mock -source=handler.go -destination=./mocks/transferHandler_mock.go


type TransferHandlerTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	transferManagerMock *transfer_mocks.MockTransferManager
	transferHandler *transferHandler.TransferHandler
	logger *slog.Logger
	router *gin.Engine
}

func (suite *TransferHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.transferManagerMock = transfer_mocks.NewMockTransferManager(suite.ctrl)
	suite.logger = slog.Default()
	handler, err := transferHandler.NewTransferHandler(suite.transferManagerMock, suite.logger)
	suite.NoError(err, "Failed to create transfer handler")
	suite.transferHandler = handler

	suite.router = router.SetupRouter(&userHandler.UserHandler{}, suite.transferHandler)
}

func (suite *TransferHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestTransferHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(TransferHandlerTestSuite))
}

func (suite *TransferHandlerTestSuite) TestNewTransferHandler() {
	suite.Run("nil transfer manager", func() {
		handler, err := transferHandler.NewTransferHandler(nil, suite.logger)
		suite.Nil(handler, "Expected nil handler")
		suite.Error(err, "Expected error when transfer manager is nil")
	})
	suite.Run("nil logger", func() {
		handler, err := transferHandler.NewTransferHandler(suite.transferManagerMock, nil)
		suite.Nil(handler, "Expected nil handler")
		suite.Error(err, "Expected error when logger is nil")
	})
	suite.Run("valid transfer manager", func() {
		handler, err := transferHandler.NewTransferHandler(suite.transferManagerMock, suite.logger)
		suite.NotNil(handler, "Expected non-nil handler")
		suite.NoError(err, "Expected no error when transfer manager is valid")
	})
}

func (suite *TransferHandlerTestSuite) TestCreateTransfer() {
	expectedTransfer := &transfer.Transfer{
		FromUserID: 1,
		ToUserID:   2,
		Amount:     100,
	}
	suite.Run("invalid JSON", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/transfers/", strings.NewReader("{invalid json}"))
		suite.router.ServeHTTP(w, req)

		suite.Equal(400, w.Code, "Expected status code 400")
	})
	suite.Run("valid transfer", func() {

		expectedResponse := map[string]any{
			"message":     "transfer created",
			"transfer_id": uint(1),
		}
	
		transferJson, _ := json.Marshal(expectedTransfer)
		w := httptest.NewRecorder()
	
		suite.transferManagerMock.EXPECT().CreateTransfer(gomock.Any(), expectedTransfer).Return(uint(1), nil).Times(1)
	
		req, _ := http.NewRequest("POST", "/transfers/", strings.NewReader(string(transferJson)))
		suite.router.ServeHTTP(w, req)
	
		m := make(map[string]any)
		err := json.Unmarshal(w.Body.Bytes(), &m)
		suite.Require().NoError(err, "Failed to unmarshal response body")
		// Check the response
		suite.Equal(200, w.Code, "Expected status code 200")
		suite.Equal(expectedResponse["message"], m["message"], "Expected response body to match")
		suite.Equal(expectedResponse["transfer_id"], uint(m["transfer_id"].(float64)), "Expected response body to match")
	})
}

func (suite *TransferHandlerTestSuite) TestFinishTransfer() {
	expectedTransfer := &transfer.Transfer{
		Status: transfer.Completed,
	}
	expectedTransfer.ID = 1
	suite.Run("invalid JSON", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/transfers/finish", strings.NewReader("{invalid json}"))
		suite.router.ServeHTTP(w, req)

		suite.Equal(400, w.Code, "Expected status code 400")
	})
	suite.Run("valid transfer", func() {

		expectedResponse := map[string]any{
			"message": fmt.Sprintf("transfer %d finished with status %s", expectedTransfer.ID, expectedTransfer.Status),
		}
	
		transferJson, _ := json.Marshal(expectedTransfer)
		w := httptest.NewRecorder()
	
		suite.transferManagerMock.EXPECT().FinishTransfer(gomock.Any(), expectedTransfer.ID, expectedTransfer.Status).Return(nil).Times(1)
	
		req, _ := http.NewRequest("POST", "/transfers/finish", strings.NewReader(string(transferJson)))
		suite.router.ServeHTTP(w, req)
	
		m := make(map[string]any)
		err := json.Unmarshal(w.Body.Bytes(), &m)
		suite.Require().NoError(err, "Failed to unmarshal response body")
		// Check the response
		suite.Equal(200, w.Code, "Expected status code 200")
		suite.Equal(expectedResponse["message"], m["message"], "Expected response body to match")
	})
}

func (suite *TransferHandlerTestSuite) TestGetTransferByID() {
	expectedTransfer := &transfer.Transfer{
		FromUserID: 1,
		ToUserID:   2,
		Amount:     100,
		Description: "Test transfer",
		Status: transfer.Pending,
	}
	expectedTransfer.ID = uint(1)
	suite.Run("transfer not found", func() {
		w := httptest.NewRecorder()
		errExpected := fmt.Errorf("transfer not found")
		suite.transferManagerMock.EXPECT().GetTransferByID(gomock.Any(), uint(1)).Return(nil, errExpected)

		req, _ := http.NewRequest("GET", "/transfers/1", nil)
		suite.router.ServeHTTP(w, req)

		suite.Equal(500, w.Code, "Expected status code 500")
	})

	suite.Run("valid transfer", func() {
	
		expectedResponse := map[string]any{
			"from_user_id": expectedTransfer.FromUserID,
			"to_user_id":   expectedTransfer.ToUserID,
			"amount":       expectedTransfer.Amount,
			"description":  expectedTransfer.Description,
			"status":       "PENDING",
		}
	
		w := httptest.NewRecorder()
	
		suite.transferManagerMock.EXPECT().GetTransferByID(gomock.Any(), uint(1)).Return(expectedTransfer, nil)
	
		req, _ := http.NewRequest("GET", "/transfers/1", nil)
		suite.router.ServeHTTP(w, req)
	
		m := make(map[string]any)
		err := json.Unmarshal(w.Body.Bytes(), &m)
		suite.Require().NoError(err, "Failed to unmarshal response body")
	
		suite.Equal(200, w.Code, "Expected status code 200")
		
		m = m["transfer"].(map[string]any)
		suite.Equal(expectedResponse["from_user_id"], uint(m["from_user_id"].(float64)), "Expected response body to match")
		suite.Equal(expectedResponse["to_user_id"], uint(m["to_user_id"].(float64)), "Expected response body to match")	
		suite.Equal(expectedResponse["amount"], uint(m["amount"].(float64)), "Expected response body to match")
		suite.Equal(expectedResponse["description"], m["description"], "Expected response body to match")
		suite.Equal(expectedResponse["status"], m["status"], "Expected response body to match")
	})
}