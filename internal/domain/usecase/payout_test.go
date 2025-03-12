package usecase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/item"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/usecase"
)

type MockExchangeRateCalculator struct {
	mock.Mock
}

func (m *MockExchangeRateCalculator) Convert(from *valueobject.Money, to valueobject.Currency) (*valueobject.Money, error) {
	args := m.Called(from, to)
	money, ok := args.Get(0).(*valueobject.Money)
	if !ok {
		return nil, args.Error(1)
	}

	return money, args.Error(1)
}

func createMoney(amount float32, currency valueobject.Currency) *valueobject.Money {
	m, _ := valueobject.NewMoney(amount, currency)
	return m
}

func TestPayout(t *testing.T) {
	var (
		mockExchange = new(MockExchangeRateCalculator)
		service      = usecase.NewPayoutCalculator(mockExchange)
		money        = createMoney(100, valueobject.Currency_GBP)
		limit        = createMoney(50, valueobject.Currency_GBP)
		seller1      = valueobject.SellerReference("a")
		seller2      = valueobject.SellerReference("b")
		seller3      = valueobject.SellerReference("c")
		item1        = item.NewItem(uuid.New(), item.NewProductName("rick"), money, seller1)
		item2        = item.NewItem(uuid.New(), item.NewProductName("morty"), money, seller1)
		item3        = item.NewItem(uuid.New(), item.NewProductName("some shirt"), money, seller2)
		item4        = item.NewItem(uuid.New(), item.NewProductName("tshirt"), money, seller3)
	)
	t.Run("successful payout for one seller", func(t *testing.T) {
		var (
			id1   = uuid.New()
			id2   = uuid.New()
			item1 = item.NewItem(id1, item.NewProductName("rick"), money, seller1)
			item2 = item.NewItem(id2, item.NewProductName("morty"), money, seller1)
			items = item.Items{item1, item2}
			txl   = payout.NewTransactionLimit(limit)
		)

		// asuming 1:1
		mockExchange.On("Convert", mock.Anything, valueobject.Currency_GBP).Return(money, nil).Twice()

		total, err := service.CalculateTotalAmountBySeller(items, valueobject.Currency_GBP)
		assert.NoError(t, err)

		po, err := service.GeneratePayouts(items, txl, valueobject.Currency_GBP, total)

		assert.NoError(t, err)
		// one seller gets 1 payout batch
		assert.Len(t, po, 1)
		mockExchange.AssertExpectations(t)

		for _, p := range po {
			// 2 items sold there shoild be 2 sales
			assert.Len(t, p.GetSales(), 2)
			sales1 := p.GetSales()[0]
			sales2 := p.GetSales()[1]
			assert.Equal(t, id1, sales1.ItemID())
			assert.Equal(t, id2, sales2.ItemID())

			assert.Equal(t, item1.Price().Amount(), sales1.ItemPrice().Amount())
			assert.Equal(t, item1.Price().Currency(), sales1.ItemPrice().Currency())
		}
	})

	t.Run("successful payout for multiple sellers", func(t *testing.T) {
		items := item.Items{item1, item2, item3, item4}
		txl := payout.NewTransactionLimit(limit)

		// asuming 1:1
		mockExchange.On("Convert", mock.Anything, valueobject.Currency_GBP).Return(money, nil).Times(4)

		total, err := service.CalculateTotalAmountBySeller(items, valueobject.Currency_GBP)
		assert.NoError(t, err)
		po, err := service.GeneratePayouts(items, txl, valueobject.Currency_GBP, total)

		assert.NoError(t, err)
		// one seller gets 1 payout batch
		assert.Len(t, po, 3)
		mockExchange.AssertExpectations(t)

		for _, payout := range po {
			if seller1 == payout.SellerRef() {
				// limit is 50 so seller1 should get 4 payouts sold 200 worth
				assert.Len(t, payout.GetPayouts(), 4)
				for _, payment := range payout.GetPayouts() {
					assert.Equal(t, payment.TotalAmount().Amount(), txl.Limit().Amount())
				}
			}
		}
	})

	t.Run("errored payout", func(t *testing.T) {
		items := item.Items{item1}
		txl := payout.NewTransactionLimit(limit)

		mockExchange.On("Convert", mock.Anything, valueobject.Currency_GBP).Return(nil, assert.AnError).Once()

		total, _ := service.CalculateTotalAmountBySeller(items, valueobject.Currency_GBP)
		po, err := service.GeneratePayouts(items, txl, valueobject.Currency_GBP, total)

		assert.Error(t, err)
		assert.Nil(t, po)
		mockExchange.AssertExpectations(t)
	})

	t.Run("successfull conversation", func(t *testing.T) {
		limit, _ := valueobject.NewMoney(100, valueobject.Currency_GBP)
		convertedLimit, _ := valueobject.NewMoney(100, valueobject.Currency_USD)

		mockExchange.On("Convert", limit, valueobject.Currency_USD).Return(convertedLimit, nil).Once()

		txLimit, err := service.ConvertTransactionLimit(limit, valueobject.Currency_USD)

		assert.NoError(t, err)
		assert.NotNil(t, txLimit)
		assert.Equal(t, convertedLimit.Amount(), txLimit.Limit().Amount())
		mockExchange.AssertExpectations(t)
	})

	t.Run("failed conversation", func(t *testing.T) {
		limit, _ := valueobject.NewMoney(100, valueobject.Currency_GBP)
		mockExchange.On("Convert", limit, valueobject.Currency_USD).Return(nil, assert.AnError).Once()

		txLimit, err := service.ConvertTransactionLimit(limit, valueobject.Currency_USD)

		assert.Error(t, err)
		assert.Nil(t, txLimit)
		mockExchange.AssertExpectations(t)
	})

	t.Run("successfull calculation per user", func(t *testing.T) {
		limit, _ := valueobject.NewMoney(100, valueobject.Currency_GBP)

		items := item.Items{item1, item2, item3}
		mockExchange.On("Convert", limit, valueobject.Currency_GBP).Return(money, nil).Times(3)

		calculations, err := service.CalculateTotalAmountBySeller(items, valueobject.Currency_GBP)

		assert.NoError(t, err)
		assert.NotNil(t, calculations)

		assert.Equal(t, calculations[seller1].Amount(), float32(200))
		assert.Equal(t, calculations[seller2].Amount(), float32(100))
		mockExchange.AssertExpectations(t)
	})

	t.Run("failed calculation", func(t *testing.T) {
		limit, _ := valueobject.NewMoney(100, valueobject.Currency_GBP)

		items := item.Items{item1}
		mockExchange.On("Convert", limit, valueobject.Currency_GBP).Return(nil, assert.AnError).Once()

		calculations, err := service.CalculateTotalAmountBySeller(items, valueobject.Currency_GBP)

		assert.Error(t, err)
		assert.Nil(t, calculations)
		mockExchange.AssertExpectations(t)
	})
}
