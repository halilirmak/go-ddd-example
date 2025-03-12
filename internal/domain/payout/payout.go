package payout

import (
	"time"

	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type Payout struct {
	id              uuid.UUID
	totalAmount     *valueobject.Money
	sellerReferance valueobject.SellerReference
	createdAt       time.Time
}

func NewPayout(id uuid.UUID, amount *valueobject.Money, ref valueobject.SellerReference) *Payout {
	return &Payout{
		id:              id,
		totalAmount:     amount,
		sellerReferance: ref,
		createdAt:       time.Now(),
	}
}

func (p *Payout) ID() uuid.UUID {
	return p.id
}

func (p *Payout) TotalAmount() *valueobject.Money {
	return p.totalAmount
}

func (p *Payout) SellerReference() valueobject.SellerReference {
	return p.sellerReferance
}
