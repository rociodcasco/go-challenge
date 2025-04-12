package storage

import (
	"go-challenge/pkg/transfer"
)

func (s *Storage) CreateTransfer(transfer *transfer.Transfer) (uint, error) {
	result := s.db.Create(transfer)
	if result.Error != nil {
		return 0, result.Error
	}
	return transfer.ID, nil
}

func (s *Storage) GetTransferByID(id uint) (*transfer.Transfer, error) {
	var transfer transfer.Transfer
	result := s.db.First(&transfer, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transfer, nil
}

func (s *Storage) UpdateTransfer(transfer *transfer.Transfer) error {
	result := s.db.Save(transfer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}