package userHandler

import (
	"fmt"
	"go-challenge/pkg/user"

	"github.com/gin-gonic/gin"
)

type UserManager interface {
	CreateUser(user *user.User) error
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
	err := h.manager.CreateUser(&user.User{
		ID: 123,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user created"})
}