package userHandler

import (
	"context"
	"fmt"
	"go-challenge/pkg/user"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserManager interface {
	CreateUser(ctx context.Context, user *user.User) (uint, error)
	GetUserBalance(ctx context.Context, id uint) (uint, error)
}

type UserHandler struct {
	logger  *slog.Logger
	manager UserManager
}

func NewUserHandler(manager UserManager, logger *slog.Logger) (*UserHandler, error) {
	if manager == nil {
		return nil, fmt.Errorf("user manager cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	return &UserHandler{
		manager: manager,
		logger:  logger,
	}, nil
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	h.logger.InfoContext(c.Request.Context(), "Creating user...")
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id, err := h.manager.CreateUser(c.Request.Context(), &user)
	if err != nil {
		// Lo correcto seria separar el error en base a su tipo, pero lo dejo así por simplicidad
		h.logger.ErrorContext(c.Request.Context(), "Failed to create user", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"message": "user created",
		"user_id": id,
	})
}

func (h *UserHandler) GetUserBalance(c *gin.Context) {
	h.logger.InfoContext(c.Request.Context(), "Getting user balance...")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	balance, err := h.manager.GetUserBalance(c.Request.Context(), uint(id))
	if err != nil {
		// Lo correcto seria separar el error en base a su tipo, pero lo dejo así por simplicidad
		h.logger.ErrorContext(c.Request.Context(), "Failed to get user balance", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"id":      id,
		"balance": balance,
	})
}
