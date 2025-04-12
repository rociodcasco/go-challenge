package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUserBalance(c *gin.Context)
}

type TransferHandler interface {
	CreateTransfer(c *gin.Context)
	FinishTransfer(c *gin.Context)
	GetTransferByID(c *gin.Context)
}

func SetupRouter(userHandler UserHandler, transferHandler TransferHandler) *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	users := r.Group("/users")
	users.POST("/", userHandler.CreateUser)
	users.GET("/:id/balance", userHandler.GetUserBalance)

	transfers := r.Group("/transfers")
	transfers.POST("/", transferHandler.CreateTransfer)
	transfers.POST("/finish", transferHandler.FinishTransfer)
	transfers.GET("/:id", transferHandler.GetTransferByID)

	return r
}