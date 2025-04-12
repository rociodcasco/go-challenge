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