package transferHandler

import (
	"fmt"
	"go-challenge/pkg/transfer"

	"github.com/gin-gonic/gin"
)

type TransferManager interface {
	CreateTransfer(transfer *transfer.Transfer) (uint, error)
}

type TransferHandler struct {
	transferManager TransferManager
}

func NewTransferHandler(transferManager TransferManager) (*TransferHandler, error) {
	if transferManager == nil {
		return nil, fmt.Errorf("transfer manager cannot be nil")
	}
	return &TransferHandler{transferManager: transferManager}, nil
}

func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var transfer transfer.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		errasd := fmt.Errorf("failed to bind JSON: %w", err)
		c.JSON(400, gin.H{"error": errasd.Error()})
		return
	}

	id, err := h.transferManager.CreateTransfer(&transfer)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "transfer created",
		"transfer_id": id,
		"status": transfer.Status,
	})
}