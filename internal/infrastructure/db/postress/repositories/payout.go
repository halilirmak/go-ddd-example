package repositories

import (
	"context"

	"gorm.io/gorm"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
	postgres_mapper "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress/dto/mapper"
)

type PostgresPayoutRepository struct {
	db *gorm.DB
}

func NewPostgresPayoutRepository(db *gorm.DB) payout.PayoutRepository {
	return &PostgresPayoutRepository{
		db: db,
	}
}

func (pr *PostgresPayoutRepository) GetTransactionLimit(ctx context.Context) (*payout.TransactionLimit, error) {
	money, _ := valueobject.NewMoney(100, valueobject.Currency_GBP)
	items := payout.NewTransactionLimit(money)
	return items, nil
}

func (pr *PostgresPayoutRepository) TxCreatePayouts(ctx context.Context, payouts payout.BatchPayouts) error {
	params := postgres_mapper.ToDatabaseBatchPayout(payouts)
	return pr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&params.BatchPayouts).Error; err != nil {
			return errors.NewContextualError("batch payout transaction failed", "payout-repo").Wrap(err)
		}

		if err := tx.Create(&params.Sales).Error; err != nil {
			return errors.NewContextualError("sale transaction failed", "payout-repo").Wrap(err)
		}
		return nil
	})
}
