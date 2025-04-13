package transferHandler

import (
	"fmt"
	"go-challenge/pkg/transfer"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransferManager interface {
	CreateTransfer(transfer *transfer.Transfer) (uint, error)
	FinishTransfer(id uint, status transfer.Status) error
	GetTransferByID(id uint) (*transfer.Transfer, error)
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
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to bind JSON: %w", err).Error()})
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
	})
}

func (h *TransferHandler) FinishTransfer(c *gin.Context) {
	var transfer transfer.Transfer
	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("failed to bind JSON: %w", err).Error()})
		return
	}

	err := h.transferManager.FinishTransfer(transfer.ID, transfer.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("transfer %d finished with status %s", transfer.ID, transfer.Status),
	})
}

func (h *TransferHandler) GetTransferByID(c *gin.Context) {
	id := c.Param("id")
	transferID, err := strconv.Atoi(id)
	fmt.Println("transferID", transferID)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Errorf("invalid transfer ID: %w", err).Error()})
		return
	}
	fmt.Println("transferID", uint(transferID))
	transfer, err := h.transferManager.GetTransferByID(uint(transferID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"transfer": transfer,
	})
}