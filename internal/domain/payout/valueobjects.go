package payout

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type BatchPayouts []*BatchPayout

func (p *BatchPayouts) Add(batch *BatchPayout) {
	*p = append(*p, batch)
}

func (p BatchPayouts) BatchPayoutIDs() []uuid.UUID {
	var result []uuid.UUID
	for _, batch := range p {
		result = append(result, batch.ID())
	}
	return result
}

var ErrTotalAmount = fmt.Errorf("payment amaount not found  in map for seller")

type TotalAmountsBySeller map[valueobject.SellerReference]*valueobject.Money

func (s TotalAmountsBySeller) Get(seller valueobject.SellerReference) (*valueobject.Money, error) {
	money, ok := s[seller]

	if !ok {
		return nil, errors.NewContextualError(fmt.Sprintf("no payout amount found for seller: [%s]", seller), "seller-payouts-valueobject")
	}
	return money, nil
}
