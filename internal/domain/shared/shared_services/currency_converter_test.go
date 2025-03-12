package shared_services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/shared_services"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type MockExchangeRateProvider struct {
	mock.Mock
}

func (m *MockExchangeRateProvider) GetExchangeRate(from string, to string) (float32, error) {
	args := m.Called(from, to)
	money, ok := args.Get(0).(float32)
	if !ok {
		return 0, args.Error(1)
	}

	return money, args.Error(1)
}

func createMoney(amount float32, currency valueobject.Currency) *valueobject.Money {
	m, _ := valueobject.NewMoney(amount, currency)
	return m
}

func NewCurrencyConverter(t *testing.T) {
	mockExchange := new(MockExchangeRateProvider)
	service := shared_services.NewCurrencyConverter(mockExchange)

	t.Run("successfull conversation", func(t *testing.T) {
		money := createMoney(100, valueobject.Currency_EUR)

		mockExchange.On("GetExchangeRate", string(valueobject.Currency_EUR), string(valueobject.Currency_USD)).Return(float32(1.20), nil).Once()

		convertedMoney, err := service.Convert(money, valueobject.Currency_USD)

		assert.NoError(t, err)
		assert.NotNil(t, convertedMoney)
		assert.InDelta(t, float32(120), convertedMoney.Amount(), 0.0001)
		mockExchange.AssertExpectations(t)
	})

	t.Run("failed conversation", func(t *testing.T) {
		money := createMoney(100, valueobject.Currency_EUR)

		mockExchange.On("GetExchangeRate", string(valueobject.Currency_EUR), string(valueobject.Currency_USD)).Return(0, assert.AnError).Once()

		convertedMoney, err := service.Convert(money, valueobject.Currency_USD)

		assert.Error(t, err)
		assert.Nil(t, convertedMoney)
		mockExchange.AssertExpectations(t)
	})
}
