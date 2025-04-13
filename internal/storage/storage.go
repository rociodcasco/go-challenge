package storage

import (
	"go-challenge/pkg/transfer"
	"go-challenge/pkg/user"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func SetupStorage(db *gorm.DB) (*Storage, error) {
	// migrations
	if err := db.AutoMigrate(&user.User{}, &transfer.Transfer{}); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}
