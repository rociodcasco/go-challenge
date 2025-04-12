package storage

import (
	"go-challenge/pkg/transfer"
	"go-challenge/pkg/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct{
	db *gorm.DB
}

func SetupStorage() (*Storage, error) {
	dsn := "host=db user=postgres password=go-challenge dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}	
	// migrations 
	if err := db.AutoMigrate(&user.User{}, &transfer.Transfer{}); err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}