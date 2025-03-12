package infrastructure

import "errors"

type ExchangeRateProvider struct{}

func (p *ExchangeRateProvider) GetExchangeRate(from, to string) (float32, error) {
	rates := map[string]map[string]float32{
		"USD": {"EUR": 0.91, "GBP": 0.78},
		"EUR": {"USD": 1.10, "GBP": 0.85},
		"GBP": {"USD": 1.40, "EUR": 1.20},
	}

	if rate, exists := rates[from][to]; exists {
		return rate, nil
	}

	return 0, errors.New("exchange rate not found")
}
