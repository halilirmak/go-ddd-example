package request

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
)

type (
	Items               []string
	CreatePayoutRequest struct {
		Currency string `json:"currency" example:"USD"`
		Items    Items  `json:"items" example:"itemId1,itemid2"`
	}
)

func (req *CreatePayoutRequest) ToCreatePayoutCommand() (*command.CreatePayoutCommand, error) {
	var (
		itemIds             command.Items
		availableCurrencies = []string{"EUR", "USD", "GBP"}
	)
	for _, id := range req.Items.Unique() {
		itemId, err := uuid.Parse(id)
		if err != nil {
			return nil, errors.NewInvalidInputError(fmt.Sprintf("item id (%s) not valid", id), "create-payout-request")
		}
		itemIds = append(itemIds, itemId)
	}

	if itemIds.IsEmpty() {
		return nil, errors.NewInvalidInputError("item ids expected", "create-payout-request")
	}

	if req.Currency == "" {
		return nil, errors.NewInvalidInputError("currency expected", "create-payout-request")
	}

	if !slices.Contains(availableCurrencies, req.Currency) {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("currency not available (%s)", req.Currency), "create-payout-request")
	}

	return &command.CreatePayoutCommand{
		Currency: req.Currency,
		Items:    itemIds,
	}, nil
}

func (i Items) Unique() Items {
	var (
		un     = make(map[string]struct{})
		result Items
	)

	for _, item := range i {
		if _, exits := un[item]; !exits {
			un[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
