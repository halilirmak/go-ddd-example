package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func TestMoneyValueObject(t *testing.T) {
	t.Run("if currency is not valid should return error", func(t *testing.T) {
		money, err := valueobject.NewMoney(100, "TRY")
		assert.Error(t, err)
		assert.Nil(t, money)
		assert.Equal(t, "invalid currency", err.Error())
	})

	t.Run("if amout is 0 should return error", func(t *testing.T) {
		money, err := valueobject.NewMoney(0, valueobject.Currency_EUR)
		assert.Error(t, err)
		assert.Nil(t, money)
		assert.Equal(t, "amount must be greater than 0 and lower than a million", err.Error())
	})

	t.Run("if amount is greater than 1_000_000 should return error", func(t *testing.T) {
		money, err := valueobject.NewMoney(1_000_000, valueobject.Currency_EUR)
		assert.Error(t, err)
		assert.Nil(t, money)
		assert.Equal(t, "amount must be greater than 0 and lower than a million", err.Error())
	})
}
