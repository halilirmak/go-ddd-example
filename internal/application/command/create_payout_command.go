package command

import (
	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/application/dto"
)

type (
	Items               []uuid.UUID
	CreatePayoutCommand struct {
		Currency string
		Items    Items
	}
)

type CreatePayoutCommandResult struct {
	Currency string
	Payout   []dto.BatchPayoutDTO
}

func (i Items) IsEmpty() bool {
	return len(i) == 0
}

func (i Items) IsEqual(length int) bool {
	return len(i) == length
}
