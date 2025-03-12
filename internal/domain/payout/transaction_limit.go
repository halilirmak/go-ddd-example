package payout

import "github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"

type TransactionLimit struct {
	limit *valueobject.Money
}

func NewTransactionLimit(limit *valueobject.Money) *TransactionLimit {
	return &TransactionLimit{limit: limit}
}

func (t *TransactionLimit) Limit() *valueobject.Money {
	return t.limit
}

func (t *TransactionLimit) GetMin(amount float32) float32 {
	if amount < t.limit.Amount() {
		return amount
	}
	return t.limit.Amount()
}
