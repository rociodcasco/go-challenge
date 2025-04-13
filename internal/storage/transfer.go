package storage

import (
	"context"
	"go-challenge/pkg/transfer"
)

func (s *Storage) CreateTransfer(ctx context.Context, transfer *transfer.Transfer) (uint, error) {
	result := s.db.WithContext(ctx).Create(transfer)
	if result.Error != nil {
		return 0, result.Error
	}
	return transfer.ID, nil
}

func (s *Storage) GetTransferByID(ctx context.Context, id uint) (*transfer.Transfer, error) {
	var transfer transfer.Transfer
	result := s.db.WithContext(ctx).First(&transfer, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transfer, nil
}

func (s *Storage) UpdateTransfer(ctx context.Context, transfer *transfer.Transfer) error {
	result := s.db.WithContext(ctx).Save(transfer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAllPendingTransfers(ctx context.Context) ([]transfer.Transfer, error) {
	var transfers []transfer.Transfer
	result := s.db.WithContext(ctx).Where("status = ?", transfer.Pending).Find(&transfers)
	if result.Error != nil {
		return nil, result.Error
	}
	return transfers, nil
}

func (s *Storage) GetTransfersInformation() (map[string]int, error) {
	var pending int64
	pendingResult := s.db.Model(&transfer.Transfer{}).Where("status = ?", "PENDING").Count(&pending)
	if pendingResult.Error != nil {
		return nil, pendingResult.Error
	}
	var completed int64
	completedResult := s.db.Model(&transfer.Transfer{}).Where("status = ?", "COMPLETED").Count(&completed)
	if completedResult.Error != nil {
		return nil, completedResult.Error
	}
	var failed int64
	failedResult := s.db.Model(&transfer.Transfer{}).Where("status = ?", "FAILED").Count(&failed)
	if failedResult.Error != nil {
		return nil, failedResult.Error
	}
	return map[string]int{
		"total":     int(pending + completed + failed),
		"pending":   int(pending),
		"completed": int(completed),
		"failed":    int(failed),
	}, nil
}