package restapi_mapper

import (
	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/dto/response"
)

func ToPayoutResponse(batchPayouts *command.CreatePayoutCommandResult) *response.CreatePayoutResponse {
	if batchPayouts == nil {
		return nil
	}
	payoutResponse := response.CreatePayoutResponse{
		Currency: batchPayouts.Currency,
	}

	for _, batchPayout := range batchPayouts.Payout {
		batchPayoutResponse := response.BatchPayout{
			ID:              batchPayout.ID.String(),
			SellerReference: batchPayout.SellerReference,
		}
		for _, payout := range batchPayout.Payouts {
			payResponse := response.Payout{
				ID:     payout.ID.String(),
				Amount: payout.Amount,
			}
			batchPayoutResponse.Payouts = append(batchPayoutResponse.Payouts, payResponse)
		}

		for _, sale := range batchPayout.Sales {
			saleResponse := response.Sale{
				ID:               sale.ID.String(),
				ItemID:           sale.ItemID.String(),
				Price:            sale.Price,
				OriginalCurrency: sale.Currency,
			}

			batchPayoutResponse.Sales = append(batchPayoutResponse.Sales, saleResponse)
		}
		payoutResponse.BatchPayouts = append(payoutResponse.BatchPayouts, batchPayoutResponse)
	}

	return &payoutResponse
}
