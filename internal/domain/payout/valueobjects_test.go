package payout_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func TestPayoutValueObjects(t *testing.T) {
	ref1 := valueobject.SellerReference("rick")
	ref2 := valueobject.SellerReference("morty")
	t.Run("BatchPayoutIDs should return ids of batchpayouts", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		money, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)

		payout1 := payout.NewPayout(uuid.New(), money, ref1)
		payout2 := payout.NewPayout(uuid.New(), money, ref1)
		b1 := payout.NewBatchPayout(id1, ref1)
		b2 := payout.NewBatchPayout(id2, ref1)
		b1.Add(payout1)
		b2.Add(payout2)

		var batchPayouts payout.BatchPayouts

		batchPayouts.Add(b1)
		batchPayouts.Add(b2)

		ids := batchPayouts.BatchPayoutIDs()

		assert.Len(t, ids, 2)
		assert.ElementsMatch(t, ids, []uuid.UUID{id1, id2})
	})

	t.Run("SellerPayoutAmounts GetAmountBySeller should return correct amount", func(t *testing.T) {
		money1, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)
		sellerPayoutAmounts := payout.TotalAmountsBySeller{
			ref1: money1,
		}

		value, err := sellerPayoutAmounts.Get(ref1)
		assert.NoError(t, err)
		assert.Equal(t, float32(100), value.Amount())
	})

	t.Run("SellerPayoutAmounts GetAmountBySeller should fail seller ref not found", func(t *testing.T) {
		money1, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)
		sellerPayoutAmounts := payout.TotalAmountsBySeller{
			ref1: money1,
		}

		value, err := sellerPayoutAmounts.Get(ref2)
		assert.Error(t, err)
		assert.Nil(t, value)
		assert.Equal(t, "no payout amount found for seller: [morty]", err.Error())
	})
}
