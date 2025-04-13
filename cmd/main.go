package main

import (
	"fmt"
	transferHandler "go-challenge/internal/handler/transfer"
	userHandler "go-challenge/internal/handler/user"
	"go-challenge/internal/router"
	"go-challenge/internal/storage"
	"go-challenge/pkg/transfer"
	"go-challenge/pkg/user"
)

func main() {
	db, err := storage.SetupStorage()
	failOnError(err)

	
	userManager, err := user.NewUserManager(db)
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