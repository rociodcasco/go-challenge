package userHandler

import (
	"fmt"
	"go-challenge/pkg/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserManager interface {
	CreateUser(user *user.User) (uint,error)
	GetUserBalance(id uint) (uint, error)
}

type UserHandler struct {
	manager UserManager
}

func NewUserHandler(manager UserManager) (*UserHandler, error) {
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
	fmt.Println("USER:", user)
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

func (h *UserHandler) GetUserBalance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	balance, err := h.manager.GetUserBalance(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"id": id,
		"balance": balance,
	})
}