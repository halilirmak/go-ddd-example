package payout

import (
	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type Sale struct {
	id              uuid.UUID
	itemID          uuid.UUID
	batchPayoutID   uuid.UUID
	itemPrice       *valueobject.Money
	sellerReferance valueobject.SellerReference
}

func NewSale(id, itemId, batchPayoutID uuid.UUID, price *valueobject.Money, ref valueobject.SellerReference) *Sale {
	return &Sale{
		id:              id,
		itemID:          itemId,
		itemPrice:       price,
		sellerReferance: ref,
		batchPayoutID:   batchPayoutID,
	}
}

func (s *Sale) ID() uuid.UUID {
	return s.id
}

func (s *Sale) ItemID() uuid.UUID {
	return s.itemID
}

func (s *Sale) ItemPrice() *valueobject.Money {
	return s.itemPrice
}

func (s *Sale) SellerReference() valueobject.SellerReference {
	return s.sellerReferance
}

func (s *Sale) BatchPayoutID() uuid.UUID {
	return s.batchPayoutID
}
