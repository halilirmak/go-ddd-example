package valueobject

import (
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
)

const (
	maxAmount float32 = 1_000_000.00
)

type Money struct {
	currency Currency
	amount   float32
}

func NewMoney(amount float32, currency Currency) (*Money, error) {
	cur, err := NewCurrency(currency)
	if err != nil {
		return nil, err
	}
	money := Money{amount: amount, currency: cur}
	if !money.isValid() {
		return nil, errors.NewInvalidInputError("amount must be greater than 0 and lower than a million", "money-vo")
	}
	return &money, nil
}

func (m *Money) isValid() bool {
	return m.amount > 0 && m.amount < maxAmount
}

func (m *Money) Amount() float32 {
	return m.amount
}

func (m *Money) Currency() Currency {
	return m.currency
}

func (m *Money) GreaterThan(t float32) bool {
	return m.amount > t
}

func (c Currency) String() string {
	return string(c)
}
