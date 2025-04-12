package userHandler

import (
	"fmt"
	"go-challenge/pkg/user"

	"github.com/gin-gonic/gin"
)

type UserManager interface {
	CreateUser(user *user.User) (uint,error)
}

type UserHandler struct {
	manager UserManager
}

func NewUserHanlder(manager UserManager) (*UserHandler, error) {
	if manager == nil {
		return nil, fmt.Errorf("user manager cannot be nil")
	}
	return &UserHandler{manager: manager}, nil
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := h.manager.CreateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "user created",
		"user_id":     id,
	})
}