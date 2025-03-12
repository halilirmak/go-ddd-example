package usecase

import (
	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/item"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/shared_services"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type PayoutCalculator interface {
	// This method generates payouts for the sellers lets say seller sold 3 items with
	// 100GBP 300GBP 70GBP and limit is 100 GBP, it would create for the seller
	// 5 payouts that would have look like [100GBP, 100GBP, 100GBP, 100GBP, 70GBP]
	GeneratePayouts(items item.Items, txl *payout.TransactionLimit, cur valueobject.Currency, amounts payout.TotalAmountsBySeller) (payout.BatchPayouts, error)

	// This method is for converting limit to the desired exchange rate. Lets say you have 100GBP limit
	// the value is not same for USD or EUR. If someone requires diffrent currency, transaction  limit should
	// reflect that
	ConvertTransactionLimit(limit *valueobject.Money, cur valueobject.Currency) (*payout.TransactionLimit, error)

	// This method would calculate the total amounts of the pays by seller
	// It calls exchange service multiple times which is not desired but this can be improved
	// by adding a short lived cache (if application does not require to have live rate)
	CalculateTotalAmountBySeller(items item.Items, cur valueobject.Currency) (payout.TotalAmountsBySeller, error)
}

type Payout struct {
	exchangeRateService shared_services.ExchangeRateService
}

func NewPayoutCalculator(service shared_services.ExchangeRateService) PayoutCalculator {
	return &Payout{exchangeRateService: service}
}

func (pc *Payout) GeneratePayouts(items item.Items, txl *payout.TransactionLimit, cur valueobject.Currency, amounts payout.TotalAmountsBySeller) (payout.BatchPayouts, error) {
	var payouts payout.BatchPayouts

	for ref, sellerItems := range items.AggregateByReferance() {
		total, err := amounts.Get(ref)
		if err != nil {
			return nil, err
		}

		batchPayout := payout.NewBatchPayout(uuid.New(), ref)
		if err := batchPayout.SplitPayouts(total, txl, cur); err != nil {
			return nil, err
		}

		for _, item := range sellerItems {
			batchPayout.RegisterSale(uuid.New(), item.ID(), batchPayout.ID(), item.Price(), item.SellerReference())
		}

		payouts.Add(batchPayout)
	}

	return payouts, nil
}

func (p *Payout) ConvertTransactionLimit(limit *valueobject.Money, cur valueobject.Currency) (*payout.TransactionLimit, error) {
	txLimit, err := p.exchangeRateService.Convert(limit, cur)
	if err != nil {
		return nil, err
	}
	return payout.NewTransactionLimit(txLimit), nil
}

func (pc *Payout) CalculateTotalAmountBySeller(items item.Items, cur valueobject.Currency) (payout.TotalAmountsBySeller, error) {
	sellerPayouts := make(payout.TotalAmountsBySeller)

	for ref, sellerItems := range items.AggregateByReferance() {
		var total float32

		for _, item := range sellerItems {
			exchanged, err := pc.exchangeRateService.Convert(item.Price(), cur)
			if err != nil {
				return nil, err
			}
			total += exchanged.Amount()
		}

		money, err := valueobject.NewMoney(total, cur)
		if err != nil {
			return nil, err
		}

		sellerPayouts[ref] = money

	}
	return sellerPayouts, nil
}
