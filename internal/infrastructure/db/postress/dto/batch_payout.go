package postgres_dto

import (
	"github.com/google/uuid"
)

type BatchPayout struct {
	PayoutID      uuid.UUID `gorm:"primary_key;type:uuid"`
	BatchPayoutID uuid.UUID `gorm:"type:uuid;index:idx_batch_payout_id"`

	Payout Payout `gorm:"foreignkey:PayoutID"`
}

func (BatchPayout) TableName() string {
	return "batch_payouts"
}
