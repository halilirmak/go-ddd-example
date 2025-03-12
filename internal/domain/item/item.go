package item

import (
	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type Item struct {
	id              uuid.UUID
	name            ProductName
	price           *valueobject.Money
	sellerReferance valueobject.SellerReference
}

func NewItem(id uuid.UUID, name ProductName, price *valueobject.Money, ref valueobject.SellerReference) *Item {
	return &Item{
		id:              id,
		name:            name,
		price:           price,
		sellerReferance: ref,
	}
}

func (i *Item) ID() uuid.UUID {
	return i.id
}

func (i *Item) Name() ProductName {
	return i.name
}

func (i *Item) Price() *valueobject.Money {
	return i.price
}

func (i *Item) SellerReference() valueobject.SellerReference {
	return i.sellerReferance
}
