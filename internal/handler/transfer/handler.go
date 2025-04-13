package transferHandler

import (
	"context"
	"fmt"
	"go-challenge/pkg/transfer"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransferManager interface {
	CreateTransfer(ctx context.Context, transfer *transfer.Transfer) (uint, error)
	FinishTransfer(ctx context.Context, id uint, status transfer.Status) error
	GetTransferByID(ctx context.Context, id uint) (*transfer.Transfer, error)
}

type TransferHandler struct {
	logger *slog.Logger
	transferManager TransferManager
}

func NewTransferHandler(transferManager TransferManager, logger *slog.Logger) (*TransferHandler, error) {
	if transferManager == nil {
		return nil, fmt.Errorf("transfer manager cannot be nil")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	return &TransferHandler{
		transferManager: transferManager,
		logger: logger,
	}, nil
}

func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	h.logger.InfoContext(c.Request.Context(), "Creating transfer...")
	var transfer transfer.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to bind JSON: %w", err).Error()})
		return
	}

	id, err := h.transferManager.CreateTransfer(c.Request.Context(), &transfer)
	if err != nil {
		// Lo correcto seria separar el error en base a su tipo, pero lo dejo así por simplicidad
		h.logger.ErrorContext(c.Request.Context(), "Failed to create transfer", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"message": "transfer created",
		"transfer_id": id,
	})
}

func (h *TransferHandler) FinishTransfer(c *gin.Context) {
	h.logger.InfoContext(c.Request.Context(), "Finishing transfer...")
	var transfer transfer.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to bind JSON: %w", err).Error()})
		return
	}

	err := h.transferManager.FinishTransfer(c.Request.Context(), transfer.ID, transfer.Status)
	if err != nil {
		// Lo correcto seria separar el error en base a su tipo, pero lo dejo así por simplicidad
		h.logger.ErrorContext(c.Request.Context(), "Failed to finish transfer", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("transfer %d finished with status %s", transfer.ID, transfer.Status),
	})
}

func (h *TransferHandler) GetTransferByID(c *gin.Context) {
	h.logger.InfoContext(c.Request.Context(), "Getting transfer by ID...")
	id := c.Param("id")
	transferID, err := strconv.Atoi(id)
	fmt.Println("transferID", transferID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("invalid transfer ID: %w", err).Error()})
		return
	}

	transfer, err := h.transferManager.GetTransferByID(c.Request.Context(), uint(transferID))
	if err != nil {
		// Lo correcto seria separar el error en base a su tipo, pero lo dejo así por simplicidad
		h.logger.ErrorContext(c.Request.Context(), "Failed to get transfer", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"transfer": transfer,
	})
}