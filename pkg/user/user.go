package user

import (
	"context"
	"fmt"
)

type Storage interface{
	CreateUser(ctx context.Context, user *User) (uint,error)
	GetUserByID(ctx context.Context, id uint) (*User, error)

	CountUsers() (int, error)
}

type UserManager struct{
	storage Storage
}

func NewUserManager(s Storage) (*UserManager, error) {
	if s == nil {
		return nil, fmt.Errorf("storage cannot be nil")
	}
	return &UserManager{storage: s}, nil
}

func (m *UserManager) CreateUser(ctx context.Context, user *User) (uint, error) {
	if user == nil {
		return 0, fmt.Errorf("user cannot be nil")
	}
	
	user.Balance = 0.0 // default balance
	id, err := m.storage.CreateUser(ctx, user)
	if err != nil {	
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (m *UserManager) GetUserBalance(ctx context.Context, id uint) (uint, error) {
	user, err := m.storage.GetUserByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}
	return user.Balance, nil
}

func (m *UserManager) GetUsersInformation() (map[string]int, error) {
	count, err := m.storage.CountUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}
	return map[string]int{"total_users": count}, nil
}