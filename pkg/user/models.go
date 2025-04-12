package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name	 string `json:"name"`
	DNI string `json:"dni"`
	Email    string `json:"email"`
	Balance  uint   `json:"balance" gorm:"check:balance>=0"`
}