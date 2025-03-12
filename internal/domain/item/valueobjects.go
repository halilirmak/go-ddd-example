package item

import (
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type ProductName string

func NewProductName(name string) ProductName {
	// some validation for product name
	return ProductName(name)
}

func (p ProductName) String() string {
	return string(p)
}

type (
	Items           []*Item
	SellerAggregate map[valueobject.SellerReference]Items
)

func (i Items) AggregateByReferance() SellerAggregate {
	aggregate := make(SellerAggregate)
	for _, item := range i {
		aggregate[item.SellerReference()] = append(aggregate[item.SellerReference()], item)
	}
	return aggregate
}

func (i Items) IsEqual(length int) bool {
	return len(i) == length
}

func (i Items) IsEmpty() bool {
	return len(i) == 0
}
