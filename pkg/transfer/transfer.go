package transfer

import (
	"context"
	"fmt"
	"time"

	"go-challenge/pkg/user"
)

type Storage interface {
	CreateTransfer(ctx context.Context, transfer *Transfer) (uint, error)
	GetTransferByID(ctx context.Context, id uint) (*Transfer, error)
	UpdateTransfer(ctx context.Context, transfer *Transfer) error
	GetAllPendingTransfers(ctx context.Context) ([]Transfer, error)

	GetTransfersInformation()(map[string]int, error)

	GetUserByID(ctx context.Context, id uint) (*user.User, error)

	UpdateBalances(ctx context.Context, fromUserID, toUserID uint, amount uint) error
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

func (m *TransferManager) CreateTransfer(ctx context.Context, transfer *Transfer) (uint, error) {
	if transfer == nil {
		return 0, fmt.Errorf("transfer cannot be nil")
	}
	if transfer.Amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}

	if transfer.FromUserID == transfer.ToUserID {
		return 0, fmt.Errorf("from user and to user cannot be the same")
	}

	transfer.Status = Pending
	transfer.TransferDate = time.Now()
	id, err := m.storage.CreateTransfer(ctx, transfer)
	if err != nil {
		return 0, fmt.Errorf("failed to create transfer: %w", err)
	}

	return id, nil
}

func (m *TransferManager) FinishTransfer(ctx context.Context, id uint, status Status) error {
	transfer, err := m.storage.GetTransferByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}
	if transfer.Status != Pending {
		return fmt.Errorf("transfer is not pending")
	}

	if status == Completed {
		err = m.storage.UpdateBalances(ctx, transfer.FromUserID, transfer.ToUserID, transfer.Amount)
		if err != nil {
			return fmt.Errorf("failed to update balances: %w", err)
		}
	}

	transfer.Status = status
	err = m.storage.UpdateTransfer(ctx, transfer)
	if err != nil {
		return fmt.Errorf("failed to update transfer: %w", err)
	}

	return nil
}

func (m *TransferManager) GetTransferByID(ctx context.Context, id uint) (*Transfer, error) {
	transfer, err := m.storage.GetTransferByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transfer not found: %w", err)
	}
	return transfer, nil
}

func (m *TransferManager) ExpireTransfer(ctx context.Context, id uint) error {
	transfer, err := m.storage.GetTransferByID(ctx, id)
	if err != nil {
		return fmt.Errorf("transfer not found: %w", err)
	}
	if transfer.Status != Pending {
		return fmt.Errorf("transfer is not pending")
	}

	transfer.Status = Failed
	err = m.storage.UpdateTransfer(ctx, transfer)
	if err != nil {
		return fmt.Errorf("failed to update transfer: %w", err)
	}

	return nil
}

func (m *TransferManager) GetAllPendingTransfers(ctx context.Context) ([]Transfer, error) {
	transfers, err := m.storage.GetAllPendingTransfers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending transfers: %w", err)
	}
	return transfers, nil
}

func (m *TransferManager) GetTransfersInformation() (map[string]int, error) {
	return m.storage.GetTransfersInformation()
}