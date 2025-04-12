package storage

import "go-challenge/pkg/user"

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