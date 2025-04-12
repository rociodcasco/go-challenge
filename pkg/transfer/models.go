package transfer

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const(
	Pending Status = "PENDING"
	Completed Status = "COMPLETED"
	Failed Status = "FAILED"
)

type Transfer struct {
	gorm.Model
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	Amount        float64   `json:"amount"`
	Description   string `json:"description"`
	TransferDate  time.Time `json:"transfer_date"`
	Status        Status `json:"status"`
}