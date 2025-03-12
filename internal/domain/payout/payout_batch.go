package payout

import (
	"github.com/google/uuid"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

type BatchPayout struct {
	id        uuid.UUID
	sellerRef valueobject.SellerReference
	payouts   []*Payout
	sales     []*Sale
}

func NewBatchPayout(id uuid.UUID, ref valueobject.SellerReference) *BatchPayout {
	return &BatchPayout{
		id:        id,
		sellerRef: ref,
	}
}

func (p *BatchPayout) Add(payout *Payout) {
	p.payouts = append(p.payouts, payout)
}

func (p *BatchPayout) SplitPayouts(amount *valueobject.Money, tx *TransactionLimit, cur valueobject.Currency) error {
	remaining := amount.Amount()

	for remaining > 0 {
		payAmount := tx.GetMin(remaining)

		money, err := valueobject.NewMoney(payAmount, cur)
		if err != nil {
			return err
		}
		payout := NewPayout(uuid.New(), money, p.sellerRef)
		p.Add(payout)
		remaining -= payAmount

	}
	return nil
}

func (p *BatchPayout) RegisterSale(id, itemId, batchPayoutID uuid.UUID, itemPrice *valueobject.Money, ref valueobject.SellerReference) {
	sale := NewSale(id, itemId, batchPayoutID, itemPrice, ref)
	p.sales = append(p.sales, sale)
}

func (p *BatchPayout) GetSales() []*Sale {
	return p.sales
}

func (p *BatchPayout) GetPayouts() []*Payout {
	return p.payouts
}

func (p *BatchPayout) ID() uuid.UUID {
	return p.id
}

func (p *BatchPayout) SellerRef() valueobject.SellerReference {
	return p.sellerRef
}
