package shared_services

import (
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type ExchangeRateProvider interface {
	GetExchangeRate(from, to string) (float32, error)
}

type ExchangeRateService interface {
	Convert(from *valueobject.Money, to valueobject.Currency) (*valueobject.Money, error)
}

type CurrencyConverter struct {
	service ExchangeRateProvider
}

func NewCurrencyConverter(provider ExchangeRateProvider) ExchangeRateService {
	return &CurrencyConverter{service: provider}
}

func (c *CurrencyConverter) Convert(from *valueobject.Money, to valueobject.Currency) (*valueobject.Money, error) {
	if from.Currency() == to {
		return from, nil
	}

	rate, err := c.service.GetExchangeRate(string(from.Currency()), string(to))
	if err != nil {
		return nil, err
	}

	convertedAmount := from.Amount() * rate
	return valueobject.NewMoney(convertedAmount, to)
}
