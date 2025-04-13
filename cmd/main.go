package main

import (
	"fmt"
	transferHandler "go-challenge/internal/handler/transfer"
	userHandler "go-challenge/internal/handler/user"
	"go-challenge/internal/router"
	"go-challenge/internal/storage"
	"go-challenge/pkg/transfer"
	"go-challenge/pkg/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=postgres password=go-challenge dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	failOnError(err)

	storage, err := storage.SetupStorage(db)
	failOnError(err)

	
	userManager, err := user.NewUserManager(storage)
	failOnError(err)


	userHandler, err := userHandler.NewUserHandler(userManager)
	failOnError(err)


	trasferManager, err := transfer.NewTransferManager(db)
	failOnError(err)

	transferHandler, err := transferHandler.NewTransferHandler(trasferManager)
	failOnError(err)

	r := router.SetupRouter(userHandler, transferHandler)
	fmt.Println("Server is running on port 8080...")
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8000")
}

func failOnError(err error){
	if err != nil {
		panic(err)
	}
}