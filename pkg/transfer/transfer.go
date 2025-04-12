package transfer

import (
	"fmt"
	"time"

	"go-challenge/pkg/user"
)

type Storage interface {
	CreateTransfer(transfer *Transfer) (uint, error)
	GetTransferByID(id uint) (*Transfer, error)
	UpdateTransfer(transfer *Transfer) error
	
	GetUserByID(id uint) (*user.User, error)
}

type TransferManager struct {
	storage Storage
}

func NewTransferManager(s Storage) (*TransferManager, error) {
	if s == nil {
		return nil, fmt.Errorf("storage cannot be nil")
	}
	return &TransferManager{storage: s}, nil
}

func (m *TransferManager) CreateTransfer(transfer *Transfer) (uint, error) {
	if transfer == nil {
		return 0, fmt.Errorf("transfer cannot be nil")
	}
	if transfer.Amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}

	_, err := m.storage.GetUserByID(transfer.FromUserID)
	if err != nil {
		return 0, fmt.Errorf("from user not found: %w", err)
	}

	_, err = m.storage.GetUserByID(transfer.ToUserID)
	if err != nil {
		return 0, fmt.Errorf("to user not found: %w", err)
	}
	if transfer.FromUserID == transfer.ToUserID {
		return 0, fmt.Errorf("from user and to user cannot be the same")
	}

	transfer.Status = Pending
	transfer.TransferDate = time.Now()
	id, err := m.storage.CreateTransfer(transfer)
	if err != nil {
		return 0, fmt.Errorf("failed to create transfer: %w", err)
	}

	return id, nil
}

func (m *TransferManager) FinishTransfer(id uint, status Status) error {
	transfer, err := m.storage.GetTransferByID(id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}
	if transfer.Status != Pending {
		return fmt.Errorf("transfer is not pending")
	}

	// update balances

	transfer.Status = status
	err = m.storage.UpdateTransfer(transfer)
	if err != nil {
		return fmt.Errorf("failed to update transfer: %w", err)
	}

	return nil
}