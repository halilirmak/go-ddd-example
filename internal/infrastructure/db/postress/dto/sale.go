package postgres_dto

import "github.com/google/uuid"

type Sale struct {
	ID              uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ItemID          uuid.UUID `gorm:"not null;type:uuid;index:idx_item_id"`
	ItemPrice       float32   `gorm:"not null;type:numeric(10,2)"`
	Currency        string    `gorm:"not null;type text"`
	SellerReference string    `gorm:"not null;"`
	BatchPayoutID   uuid.UUID `gorm:"not null;type:uuid"`

	Item Item `gorm:"foreignkey:ItemID"`
}

func (Sale) TableName() string {
	return "sales"
}
