package user

import (
	"fmt"
)

type Storage interface{
	CreateUser(user *User) (uint,error)
	GetUserByID(id uint) (*User, error)
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

func (m *UserManager) CreateUser(user *User) (uint, error) {
	if user == nil {
		return 0, fmt.Errorf("user cannot be nil")
	}
	
	user.Balance = 0.0 // default balance
	id, err := m.storage.CreateUser(user)
	if err != nil {	
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (m *UserManager) GetUserBalance(id uint) (uint, error) {
	user, err := m.storage.GetUserByID(id)
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}
	return user.Balance, nil
}