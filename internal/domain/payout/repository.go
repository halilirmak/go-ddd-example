package payout

import (
	"context"
)

type PayoutRepository interface {
	GetTransactionLimit(ctx context.Context) (*TransactionLimit, error)
	TxCreatePayouts(ctx context.Context, payouts BatchPayouts) error
}
