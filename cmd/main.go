package main

import (
	"fmt"
	userHandler "go-challenge/internal/handler/user"
	"go-challenge/internal/router"
	"go-challenge/internal/storage"
	"go-challenge/pkg/user"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
  }

func main() {
	db, err := storage.SetupStorage()
	if err != nil {
		panic("failed to initializing database")
	}
	
	userManager, err := user.NewUserManager(db)
	if err != nil {
		panic("failed to initializing user manager")
	}

	userHandler, err := userHandler.NewUserHanlder(userManager)
	if err != nil {
		panic("failed to initializing user handler")
	}

	r := router.SetupRouter(userHandler)
	fmt.Println("Server is running on port 8080...")
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8000")
}