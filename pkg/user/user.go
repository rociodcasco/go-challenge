package user

import (
	"fmt"
)

type Storage interface{
	CreateUser(user *User) (uint,error)
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