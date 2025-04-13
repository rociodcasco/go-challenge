package storage

import (
	"context"
	"go-challenge/pkg/user"

	"gorm.io/gorm"
)

func (s *Storage) CreateUser(ctx context.Context, user *user.User) (uint, error) {
	result := s.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	var user user.User
	result := s.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *Storage) UpdateBalances(ctx context.Context, fromUserID, toUserID uint, amount uint) error {
	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
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

func (s *Storage) CountUsers() (int, error) {
	var count int64
	result := s.db.Model(&user.User{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(count), nil
}