package router

import (
	"net/http"
	"os"

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

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/metrics", metrics.GetMetrics)

	users := r.Group("/users", AdminMiddleware)
	users.POST("/", userHandler.CreateUser)
	users.GET("/:id/balance", userHandler.GetUserBalance)

	transfers := r.Group("/transfers")
	transfers.POST("/", AdminMiddleware, transferHandler.CreateTransfer)
	transfers.POST("/finish", WebhookMiddleware, transferHandler.FinishTransfer)
	transfers.GET("/:id", AdminMiddleware, transferHandler.GetTransferByID)

	return r
}

func AdminMiddleware(c *gin.Context) {
	// Retrieve the token from the Authorization header
	token := c.GetHeader("Authorization")
	role := c.GetHeader("Role")
	if role != "admin" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	// Check if the token is valid for admin
	if token != os.Getenv("ADMIN_TOKEN") {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	// Continue to the next handler if authentication is successful
	c.Next()
}

func WebhookMiddleware(c *gin.Context) {
	// Retrieve the token from the Authorization header
	token := c.GetHeader("Authorization")
	role := c.GetHeader("Role")
	if role != "webhook" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	if token != os.Getenv("WEBHOOK_TOKEN") {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	// Continue to the next handler if authentication is successful
	c.Next()
}
