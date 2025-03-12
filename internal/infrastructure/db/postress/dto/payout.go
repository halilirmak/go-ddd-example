package postgres_dto

import "github.com/google/uuid"

type Payout struct {
	ID              uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Amount          float32   `gorm:"not null;type:numeric(10,2)"`
	Currency        string    `gorm:"not null"`
	SellerReference string    `gorm:"not null"`
}

func (Payout) TableName() string {
	return "payouts"
}
