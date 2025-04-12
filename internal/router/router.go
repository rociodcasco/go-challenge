package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
}

func SetupRouter(userHandler UserHandler) *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	users := r.Group("/users")

	users.POST("/", userHandler.CreateUser)

	return r
}