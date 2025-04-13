package transfer_test

import (
	"fmt"
	"go-challenge/pkg/transfer"
	transfer_mocks "go-challenge/pkg/transfer/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -package transfer_mocks -source=transfer.go -destination=./mocks/transfer_mock.go

type TransferTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	storageMock *transfer_mocks.MockStorage
	manager *transfer.TransferManager
}

func (suite *TransferTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.storageMock = transfer_mocks.NewMockStorage(suite.ctrl)

	manager, err := transfer.NewTransferManager(suite.storageMock)
	suite.NoError(err, "Failed to create transfer manager")
	suite.manager = manager
}

func (suite *TransferTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestTransferTestSuite(t *testing.T) {
    suite.Run(t, new(TransferTestSuite))
}

func (suite *TransferTestSuite) TestNewManager() {
	suite.Run("nil storage", func() {
		manager, err := transfer.NewTransferManager(nil)
		suite.Nil(manager, "Expected nil manager")
		suite.Error(err, "Expected error when storage is nil")
	})
	suite.Run("valid storage", func() {
		manager, err := transfer.NewTransferManager(suite.storageMock)
		suite.NotNil(manager, "Expected non-nil manager")
		suite.NoError(err, "Expected no error when storage is valid")
	})
}

func (suite *TransferTestSuite) TestCreateTransfer() {
	suite.Run("nil transfer", func() {
		id, err := suite.manager.CreateTransfer(nil)
		suite.Equal(uint(0), id, "Expected id to be 0")
		suite.Error(err, "Expected error when transfer is nil")
	})

	suite.Run("zero amount", func() {
		transfer := &transfer.Transfer{Amount: 0}
		id, err := suite.manager.CreateTransfer(transfer)
		suite.Equal(uint(0), id, "Expected id to be 0")
		suite.Error(err, "Expected error when amount is zero")
	})

	suite.Run("same from and to user", func() {
		transfer := &transfer.Transfer{FromUserID: 1, ToUserID: 1}
		id, err := suite.manager.CreateTransfer(transfer)
		suite.Equal(uint(0), id, "Expected id to be 0")
		suite.Error(err, "Expected error when from and to user are the same")
	})

	suite.Run("storage error", func() {
		transfer := &transfer.Transfer{
			FromUserID: 1,
			ToUserID: 2,
			Amount: 100,
			Description: "Test transfer",
		}
		errExpected := fmt.Errorf("storage error")
		suite.storageMock.EXPECT().CreateTransfer(transfer).Return(uint(0), errExpected)
		id, err := suite.manager.CreateTransfer(transfer)
		suite.Equal(uint(0), id, "Expected id to be 0")
		suite.Error(err, "Expected error when storage returns an error")
		suite.Equal(err.Error(), fmt.Errorf("failed to create transfer: %w", errExpected).Error(), "Expected error to match")
	})

	suite.Run("valid transfer", func() {
		expectedID := uint(1)
		expectedTransfer := &transfer.Transfer{
			Status: transfer.Pending,
			TransferDate: time.Now(),
			FromUserID: 1,
			ToUserID: 2,
			Amount: 100,
			Description: "Test transfer",
		}
		
		suite.storageMock.EXPECT().CreateTransfer(expectedTransfer).Return(expectedID, nil)

		id, err := suite.manager.CreateTransfer(expectedTransfer)
		suite.NoError(err, "Expected no error when creating valid transfer")
		suite.Equal(expectedID, id, "Expected id to match")
	})
}

func (suite *TransferTestSuite) TestFinishTransfer() {
	suite.Run("transfer not found", func() {
		id := uint(1)
		suite.storageMock.EXPECT().GetTransferByID(id).Return(nil, fmt.Errorf("not found"))
		err := suite.manager.FinishTransfer(id, transfer.Completed)
		suite.Error(err, "Expected error when transfer is not found")
		suite.Equal(err.Error(), fmt.Errorf("transfer not found: %w", fmt.Errorf("not found")).Error(), "Expected error to match")
	})

	suite.Run("transfer not pending", func() {
		id := uint(1)
		existingTransfer := &transfer.Transfer{Status: transfer.Completed}
		suite.storageMock.EXPECT().GetTransferByID(id).Return(existingTransfer, nil)
		err := suite.manager.FinishTransfer(id, transfer.Completed)
		suite.Error(err, "Expected error when transfer is not pending")
		suite.Equal(err.Error(), fmt.Errorf("transfer is not pending").Error(), "Expected error to match")
	})

	suite.Run("storage update error", func() {
		id := uint(1)
		existingTransfer := &transfer.Transfer{Status: transfer.Pending}
		suite.storageMock.EXPECT().GetTransferByID(id).Return(existingTransfer, nil)
		errExpected := fmt.Errorf("storage update error")
		suite.storageMock.EXPECT().UpdateBalances(existingTransfer.FromUserID, existingTransfer.ToUserID, existingTransfer.Amount).Return(errExpected)
		err := suite.manager.FinishTransfer(id, transfer.Completed)
		suite.Error(err, "Expected error when storage update fails")
		suite.Equal(err.Error(), fmt.Errorf("failed to update balances: %w", errExpected).Error(), "Expected error to match")
	})

	suite.Run("valid finish transfer completed", func() {
		id := uint(1)
		expectedFromUserID := uint(1)
		expectedToUserID := uint(2)
		existingTransfer := &transfer.Transfer{
			FromUserID: expectedFromUserID,
			ToUserID: expectedToUserID,
			Amount: 100,
			Description: "Test transfer",
			Status: transfer.Pending,
		}
		
		suite.storageMock.EXPECT().GetTransferByID(id).Return(existingTransfer, nil)

		suite.storageMock.EXPECT().UpdateBalances(expectedFromUserID, expectedToUserID, existingTransfer.Amount).Return(nil)

		suite.storageMock.EXPECT().UpdateTransfer(existingTransfer).DoAndReturn(
			func(t *transfer.Transfer) error {
				existingTransfer.Status = transfer.Completed
				suite.Require().Equal(t, existingTransfer, "Expected transfer to be updated")
				return nil
			},
		)
		err := suite.manager.FinishTransfer(id, transfer.Completed)
		suite.NoError(err, "Expected no error when finishing transfer")
		suite.Equal(transfer.Completed, existingTransfer.Status, "Expected transfer status to be updated")
	})

	suite.Run("valid finish transfer failed", func() {
		id := uint(1)
		expectedFromUserID := uint(1)
		expectedToUserID := uint(2)
		existingTransfer := &transfer.Transfer{
			FromUserID: expectedFromUserID,
			ToUserID: expectedToUserID,
			Amount: 100,
			Description: "Test transfer",
			Status: transfer.Pending,
		}
		
		suite.storageMock.EXPECT().GetTransferByID(id).Return(existingTransfer, nil)

		suite.storageMock.EXPECT().UpdateTransfer(existingTransfer).DoAndReturn(
			func(t *transfer.Transfer) error {
				existingTransfer.Status = transfer.Failed
				suite.Require().Equal(t, existingTransfer, "Expected transfer to be updated")
				return nil
			},
		)
		err := suite.manager.FinishTransfer(id, transfer.Failed)
		suite.NoError(err, "Expected no error when finishing transfer")
		suite.Equal(transfer.Failed, existingTransfer.Status, "Expected transfer status to be updated")
	})
}

func (suite *TransferTestSuite) TestGetTransferByID() {
	suite.Run("transfer not found", func() {
		id := uint(1)
		suite.storageMock.EXPECT().GetTransferByID(id).Return(nil, fmt.Errorf("not found"))
		transfer, err := suite.manager.GetTransferByID(id)
		suite.Nil(transfer, "Expected transfer to be nil")
		suite.Error(err, "Expected error when transfer is not found")
		suite.Equal(err.Error(), fmt.Errorf("transfer not found: %w", fmt.Errorf("not found")).Error(), "Expected error to match")
	})

	suite.Run("valid transfer", func() {
		id := uint(1)
		expectedTransfer := &transfer.Transfer{
			FromUserID: 1,
			ToUserID: 2,
			Amount: 100,
		}
		suite.storageMock.EXPECT().GetTransferByID(id).Return(expectedTransfer, nil)
		transfer, err := suite.manager.GetTransferByID(id)
		suite.NoError(err, "Expected no error when getting valid transfer")
		suite.Equal(expectedTransfer, transfer, "Expected transfer to match")
	})
}