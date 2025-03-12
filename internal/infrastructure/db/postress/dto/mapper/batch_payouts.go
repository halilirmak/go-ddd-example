package postgres_mapper

import (
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	postgres_dto "github.com/cryptoPickle/go-ddd-example/internal/infrastructure/db/postress/dto"
)

type BatchPayout struct {
	BatchPayouts []postgres_dto.BatchPayout
	Sales        []postgres_dto.Sale
}

func ToDatabaseBatchPayout(data payout.BatchPayouts) *BatchPayout {
	var (
		batchs []postgres_dto.BatchPayout
		sales  []postgres_dto.Sale
	)
	for _, b := range data {
		for _, p := range b.GetPayouts() {
			batch := postgres_dto.BatchPayout{
				BatchPayoutID: b.ID(),
				PayoutID:      p.ID(),

				Payout: postgres_dto.Payout{
					ID:              p.ID(),
					Amount:          p.TotalAmount().Amount(),
					Currency:        p.TotalAmount().Currency().String(),
					SellerReference: p.SellerReference().String(),
				},
			}

			batchs = append(batchs, batch)
		}

		for _, s := range b.GetSales() {
			sale := postgres_dto.Sale{
				ItemID:          s.ItemID(),
				ItemPrice:       s.ItemPrice().Amount(),
				Currency:        s.ItemPrice().Currency().String(),
				SellerReference: s.SellerReference().String(),
				BatchPayoutID:   s.BatchPayoutID(),
			}
			sales = append(sales, sale)
		}
	}

	result := &BatchPayout{
		BatchPayouts: batchs,
		Sales:        sales,
	}
	return result
}
