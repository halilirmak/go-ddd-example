package dto

import "github.com/google/uuid"

type PayoutDTO struct {
	ID       uuid.UUID
	Amount   float32
	Currency string
}

type Sale struct {
	ID            uuid.UUID
	ItemID        uuid.UUID
	BatchPayoutID uuid.UUID
	Price         float32
	Currency      string
}

type BatchPayoutDTO struct {
	ID              uuid.UUID
	SellerReference string
	Payouts         []PayoutDTO
	Sales           []Sale
}
