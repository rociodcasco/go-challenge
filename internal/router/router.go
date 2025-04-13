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

type MetricsManager interface {
	GetMetrics(c *gin.Context)
}

func SetupRouter(userHandler UserHandler, transferHandler TransferHandler, metrics MetricsManager) *gin.Engine {
	r := gin.Default()

	auth := gin.BasicAuth(gin.Accounts{
		"rocio":  "casco", 
		"pepe": "123",
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/metrics", metrics.GetMetrics)

	users := r.Group("/users")
	users.POST("/", userHandler.CreateUser)
	users.GET("/:id/balance",auth, userHandler.GetUserBalance)

	transfers := r.Group("/transfers")
	transfers.POST("/", auth, transferHandler.CreateTransfer)
	transfers.POST("/finish",auth, transferHandler.FinishTransfer)
	transfers.GET("/:id", transferHandler.GetTransferByID)

	return r
}