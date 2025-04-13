package metrics

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type TransferInformer interface {
	GetTransfersInformation() (map[string]int, error)
}

type UserInformer interface {
	GetUsersInformation() (map[string]int, error)
}

type MetricsManager struct {
	transferInformer TransferInformer
	userInformer     UserInformer
}

func NewMetricsManager(transferInformer TransferInformer, userInformer UserInformer) (*MetricsManager, error) {
	if transferInformer == nil {
		return nil, fmt.Errorf("transfer informer cannot be nil")
	}
	return &MetricsManager{
		transferInformer: transferInformer,
		userInformer:     userInformer,
	}, nil
}

func (m *MetricsManager) GetMetrics(c *gin.Context) {
	transferCount, err := m.transferInformer.GetTransfersInformation()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	userCount, err := m.userInformer.GetUsersInformation()
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"transfers": transferCount,
		"users":     userCount,
	})
}
