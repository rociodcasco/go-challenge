package user

import (
	"fmt"
)

type Storage interface{
	CreateUser(user *User) error
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

func (m *UserManager) CreateUser(user *User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	err := m.storage.CreateUser(user)
	if err != nil {	
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}