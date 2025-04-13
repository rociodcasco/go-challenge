package transferexpirer

import (
	"fmt"
	"go-challenge/pkg/transfer"
	"time"
)

type TransferManager interface {
	ExpireTransfer(transferID uint) error
	GetAllPendingTransfers() ([]transfer.Transfer, error)
}

type TransferExpirer struct {
	transferManager TransferManager
}

func NewTransferExpirer(tm TransferManager) (*TransferExpirer, error) {
	if tm == nil {
		return nil, fmt.Errorf("transfer manager cannot be nil")
	}
	return &TransferExpirer{transferManager: tm}, nil
}

func (te *TransferExpirer) Start() chan struct{} {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
		select {
			case <- ticker.C:
				te.ExpireTransfers()
			case <- quit:
				ticker.Stop()
				return
			}
		}
 	}()
	return quit
}


func (te *TransferExpirer) ExpireTransfers() {
	pendingTransfers, err := te.transferManager.GetAllPendingTransfers()
	if err != nil {
		return
	}

	for _, transfer := range pendingTransfers {
		if time.Since(transfer.TransferDate) > 1*time.Minute {
			err := te.transferManager.ExpireTransfer(transfer.ID)
			if err != nil {
				// Handle error
			}
		}
	}
}