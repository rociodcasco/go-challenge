package main

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	transferHandler "go-challenge/internal/handler/transfer"
	userHandler "go-challenge/internal/handler/user"
	"go-challenge/internal/metrics"
	"go-challenge/internal/router"
	"go-challenge/internal/storage"
	transferexpirer "go-challenge/internal/transferExpirer"
	"go-challenge/pkg/transfer"
	"go-challenge/pkg/user"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	failOnError(err)

	storage, err := storage.SetupStorage(db)
	failOnError(err)

	userManager, err := user.NewUserManager(storage)
	failOnError(err)

	userHandler, err := userHandler.NewUserHandler(userManager, logger)
	failOnError(err)

	trasferManager, err := transfer.NewTransferManager(storage)
	failOnError(err)

	transferHandler, err := transferHandler.NewTransferHandler(trasferManager, logger)
	failOnError(err)

	// Initialize the transfer expirer
	transferExpirer, err := transferexpirer.NewTransferExpirer(trasferManager, logger)
	failOnError(err)
	transferExpirer.Start()

	metrics, err := metrics.NewMetricsManager(trasferManager, userManager)
	failOnError(err)

	r := router.SetupRouter(userHandler, transferHandler, metrics)
	fmt.Println("Server is running on port 8000...")
	// Listen and Server in 0.0.0.0:8000
	r.Run(":8000")
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}
