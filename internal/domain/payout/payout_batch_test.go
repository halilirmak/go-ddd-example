package payout_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func createMoney(amount float32, currency valueobject.Currency) *valueobject.Money {
	m, _ := valueobject.NewMoney(amount, currency)
	return m
}

func TestPayoutBatch(t *testing.T) {
	money := createMoney(100, valueobject.Currency_GBP)

	t.Run("should add payout > 0", func(t *testing.T) {
		sellerRef := valueobject.NewSellerReference("seller-1")
		payoutBatch := payout.NewBatchPayout(uuid.New(), sellerRef)
		payout := payout.NewPayout(uuid.New(), money, sellerRef)

		assert.Len(t, payoutBatch.GetPayouts(), 0)
		payoutBatch.Add(payout)
		assert.Len(t, payoutBatch.GetPayouts(), 1)
	})

	t.Run("should split payouts depends on tx limit", func(t *testing.T) {
		limit := createMoney(20, valueobject.Currency_GBP)
		sellerRef := valueobject.NewSellerReference("seller-1")
		payoutBatch := payout.NewBatchPayout(uuid.New(), sellerRef)
		txl := payout.NewTransactionLimit(limit)

		assert.Len(t, payoutBatch.GetPayouts(), 0)
		payoutBatch.SplitPayouts(money, txl, valueobject.Currency_GBP)
		assert.Len(t, payoutBatch.GetPayouts(), 5)
	})
}
