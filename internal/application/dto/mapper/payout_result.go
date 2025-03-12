package application_mapper

import (
	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/application/dto"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func NewPayoutResultFromEntity(batchPayouts payout.BatchPayouts, cur valueobject.Currency) *command.CreatePayoutCommandResult {
	batchPayoutsDto := []dto.BatchPayoutDTO{}
	for _, batchPayout := range batchPayouts {
		batchPayoutDto := dto.BatchPayoutDTO{
			ID:              batchPayout.ID(),
			SellerReference: batchPayout.SellerRef().String(),
		}

		for _, pay := range batchPayout.GetPayouts() {
			payout := dto.PayoutDTO{
				ID:       pay.ID(),
				Amount:   pay.TotalAmount().Amount(),
				Currency: pay.TotalAmount().Currency().String(),
			}

			batchPayoutDto.Payouts = append(batchPayoutDto.Payouts, payout)
		}

		for _, sale := range batchPayout.GetSales() {
			sale := dto.Sale{
				ID:            sale.ID(),
				ItemID:        sale.ItemID(),
				BatchPayoutID: sale.BatchPayoutID(),
				Price:         sale.ItemPrice().Amount(),
				Currency:      sale.ItemPrice().Currency().String(),
			}

			batchPayoutDto.Sales = append(batchPayoutDto.Sales, sale)
		}
		batchPayoutsDto = append(batchPayoutsDto, batchPayoutDto)
	}

	result := command.CreatePayoutCommandResult{
		Currency: cur.String(),
		Payout:   batchPayoutsDto,
	}

	return &result
}
