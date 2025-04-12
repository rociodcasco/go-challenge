package storage

import (
	"go-challenge/pkg/user"

	"gorm.io/gorm"
)

func (s *Storage) CreateUser(user *user.User) (uint, error) {
	result := s.db.Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (s *Storage) GetUserByID(id uint) (*user.User, error) {
	var user user.User
	result := s.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *Storage) UpdateBalances(fromUserID, toUserID uint, amount float64) error {
	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update from user's balance
	if err := tx.Model(&user.User{}).Where("id = ?", fromUserID).UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update to user's balance
	if err := tx.Model(&user.User{}).Where("id = ?", toUserID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}