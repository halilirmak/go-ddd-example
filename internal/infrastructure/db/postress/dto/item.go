package postgres_dto

import (
	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/item"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type Item struct {
	ID              uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ItemName        string    `gorm:"not null"`
	Price           float32   `gorm:"not null;type:numeric(10,2)"`
	Currency        string    `gorm:"not null"`
	SellerReference string    `gorm:"not null"`
}

type Items []Item

func (Item) TableName() string {
	return "items"
}

func (i *Item) ToEntity() (*item.Item, error) {
	currency, err := valueobject.NewCurrency(valueobject.Currency(i.Currency))
	if err != nil {
		return nil, err
	}
	money, err := valueobject.NewMoney(i.Price, currency)
	if err != nil {
		return nil, err
	}

	return item.NewItem(i.ID, item.ProductName(i.ItemName), money, valueobject.NewSellerReference(i.SellerReference)), nil
}

func (is Items) ToEntity() (item.Items, error) {
	var entities item.Items
	for _, item := range is {
		entitiy, err := item.ToEntity()
		if err != nil {
			return nil, err
		}
		entities = append(entities, entitiy)
	}
	return entities, nil
}

func (is Items) IsEmpty() bool {
	return len(is) == 0
}
