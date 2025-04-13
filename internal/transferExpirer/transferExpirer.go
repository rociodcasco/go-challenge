package transferexpirer

import (
	"context"
	"fmt"
	"go-challenge/pkg/transfer"
	"log/slog"
	"time"
)

type TransferManager interface {
	ExpireTransfer(ctx context.Context, transferID uint) error
	GetAllPendingTransfers(ctx context.Context) ([]transfer.Transfer, error)
}

type TransferExpirer struct {
	logger          *slog.Logger
	transferManager TransferManager
}

func NewTransferExpirer(tm TransferManager, logger *slog.Logger) (*TransferExpirer, error) {
	if tm == nil {
		return nil, fmt.Errorf("transfer manager cannot be nil")
	}
	return &TransferExpirer{
		transferManager: tm,
		logger:          logger,
	}, nil
}

func (te *TransferExpirer) Start() chan struct{} {
	ctx := context.Background()
	te.logger.InfoContext(ctx, "Starting transfer expirer...")
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				te.logger.InfoContext(ctx, "Expiring transfers...")
				te.ExpireTransfers(ctx)
			case <-quit:
				te.logger.InfoContext(ctx, "Stopping transfer expirer...")
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

func (te *TransferExpirer) ExpireTransfers(ctx context.Context) {
	pendingTransfers, err := te.transferManager.GetAllPendingTransfers(ctx)
	if err != nil {
		return
	}

	te.logger.InfoContext(ctx, "Found pending transfers", "count", len(pendingTransfers))

	for _, transfer := range pendingTransfers {
		if time.Since(transfer.TransferDate) > 1*time.Minute {
			err := te.transferManager.ExpireTransfer(ctx, transfer.ID)
			if err != nil {
				te.logger.ErrorContext(ctx, "Failed to expire transfer", "transferID", transfer.ID, "error", err)
				continue
			}
		}
	}
}
