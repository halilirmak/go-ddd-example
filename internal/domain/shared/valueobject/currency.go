package valueobject

import (
	"fmt"
	"slices"

	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
)

type Currency string

const (
	Currency_EUR Currency = "EUR"
	Currency_USD Currency = "USD"
	Currency_GBP Currency = "GBP"
)

var ErrCurrency = fmt.Errorf("currency is invalid")

func NewCurrency(cur Currency) (Currency, error) {
	currency := Currency(cur)
	if !currency.isValid() {
		return "", errors.NewInvalidInputError("invalid currency", "currency-vo")
	}

	return cur, nil
}

func (c Currency) isValid() bool {
	avail := []Currency{Currency_EUR, Currency_GBP, Currency_USD}
	return slices.Contains(avail, c)
}
