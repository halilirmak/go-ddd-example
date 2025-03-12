package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/item"
	postgres_dto "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress/dto"
)

type PostgresItemRepository struct {
	db *gorm.DB
}

func NewPostgresItemRepository(db *gorm.DB) item.ItemRepository {
	return &PostgresItemRepository{
		db: db,
	}
}

func (ir *PostgresItemRepository) GetNotSoldItemsByID(ctx context.Context, ids []uuid.UUID) (item.Items, error) {
	var items postgres_dto.Items

	// if err := ir.db.Joins("LEFT JOIN sales ON items.id = sales.item_id").
	// 	Where("items.id IN ? AND sales.item_id IS NULL", ids).
	// 	Find(&items).Error; err != nil {
	// 	return nil, err
	// }

	if err := ir.db.Where("id IN ? AND NOT EXISTS (SELECT 1 FROM sales WHERE sales.item_id = items.id)", ids).
		Find(&items).Error; err != nil {
		return nil, errors.NewContextualError("fetching sold items failed", "items-repo").Wrap(err)
	}

	if items.IsEmpty() {
		return nil, errors.NewNotFoundError("no unsold item found for provided ids", "item-repository")
	}

	entities, err := items.ToEntity()
	if err != nil {
		return nil, err
	}
	return entities, nil
}
