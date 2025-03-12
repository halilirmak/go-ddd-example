package item

import (
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	GetNotSoldItemsByID(ctx context.Context, ids []uuid.UUID) (Items, error)
}
