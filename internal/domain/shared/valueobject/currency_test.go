package valueobject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func TestCurrencyValueObject(t *testing.T) {
	t.Run("should fail to create currency if it is not in the available list", func(t *testing.T) {
		_, err := valueobject.NewCurrency("TRY")
		assert.Error(t, err)
		assert.Equal(t, "invalid currency", err.Error())
	})

	t.Run("should create new currency if it is in the available list", func(t *testing.T) {
		cur, err := valueobject.NewCurrency("EUR")
		assert.NoError(t, err)
		assert.Equal(t, valueobject.Currency_EUR, cur)
	})
}
