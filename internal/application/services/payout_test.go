package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/application/services"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/item"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func TestPayoutService(t *testing.T) {
	ctx := context.Background()
	mItempRepo := new(MockItemRepository)
	mPayoutRepo := new(MockPayoutRepository)
	mPayoutUsecase := new(MockPayoutUsecase)

	service := services.NewPayoutService(mItempRepo, mPayoutRepo, mPayoutUsecase)

	t.Run("should return error if currency is unexpected value", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "TRY"}
		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "invalid currency", err.Error())
	})

	t.Run("should return error if GetNotSoldItemsByID call fails", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "USD"}

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mItempRepo.AssertExpectations(t)

		mItempRepo.Reset()
	})

	t.Run("should fail if provided len(ids) not matching to GetNotSoldItemsByID call", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		cmd := &command.CreatePayoutCommand{Currency: "USD", Items: command.Items{id1, id2}}
		money, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(item.Items{
			item.NewItem(id1, item.NewProductName("something"), money, valueobject.NewSellerReference("test")),
		}, nil).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "some items already paid or invalid", err.Error())

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mItempRepo.AssertExpectations(t)

		mItempRepo.Reset()
	})

	t.Run("should return error if GetTransactionLimit call fails", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "USD"}

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutRepo.On("GetTransactionLimit", mock.Anything).Return(nil, assert.AnError).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mPayoutRepo.AssertNumberOfCalls(t, "GetTransactionLimit", 1)

		mPayoutRepo.AssertExpectations(t)
		mItempRepo.AssertExpectations(t)

		mPayoutRepo.Reset()
		mItempRepo.Reset()
	})

	t.Run("should return error if ConvertTransactionLimit call fails", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "USD"}

		money, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)
		txl := payout.NewTransactionLimit(money)

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutRepo.On("GetTransactionLimit", mock.Anything).Return(txl, nil).Once()
		mPayoutUsecase.On("ConvertTransactionLimit", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mPayoutRepo.AssertNumberOfCalls(t, "GetTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "ConvertTransactionLimit", 1)

		mPayoutRepo.AssertExpectations(t)
		mItempRepo.AssertExpectations(t)
		mPayoutUsecase.AssertExpectations(t)
		mItempRepo.Reset()
		mPayoutRepo.Reset()
		mPayoutUsecase.Reset()
	})

	t.Run("should return error if CalculateTotalAmountBySeller call fails", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "USD"}

		money, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)
		txl := payout.NewTransactionLimit(money)

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutRepo.On("GetTransactionLimit", mock.Anything).Return(txl, nil).Once()
		mPayoutUsecase.On("ConvertTransactionLimit", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutUsecase.On("CalculateTotalAmountBySeller", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mPayoutRepo.AssertNumberOfCalls(t, "GetTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "ConvertTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "CalculateTotalAmountBySeller", 1)

		mPayoutRepo.AssertExpectations(t)
		mItempRepo.AssertExpectations(t)
		mPayoutUsecase.AssertExpectations(t)
		mItempRepo.Reset()
		mPayoutRepo.Reset()
		mPayoutUsecase.Reset()
	})

	t.Run("should return error if GeneratePayouts call fails", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "USD"}

		money, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)
		txl := payout.NewTransactionLimit(money)

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutRepo.On("GetTransactionLimit", mock.Anything).Return(txl, nil).Once()
		mPayoutUsecase.On("ConvertTransactionLimit", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutUsecase.On("CalculateTotalAmountBySeller", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutUsecase.On("GeneratePayouts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mPayoutRepo.AssertNumberOfCalls(t, "GetTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "ConvertTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "CalculateTotalAmountBySeller", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "GeneratePayouts", 1)

		mPayoutRepo.AssertExpectations(t)
		mItempRepo.AssertExpectations(t)
		mPayoutUsecase.AssertExpectations(t)
		mItempRepo.Reset()
		mPayoutRepo.Reset()
		mPayoutUsecase.Reset()
	})

	t.Run("should return error if TxCreatePayouts call fails", func(t *testing.T) {
		cmd := &command.CreatePayoutCommand{Currency: "USD"}

		money, err := valueobject.NewMoney(100, valueobject.Currency_EUR)
		assert.NoError(t, err)
		txl := payout.NewTransactionLimit(money)

		mItempRepo.On("GetNotSoldItemsByID", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutRepo.On("GetTransactionLimit", mock.Anything).Return(txl, nil).Once()
		mPayoutUsecase.On("ConvertTransactionLimit", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutUsecase.On("CalculateTotalAmountBySeller", mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutUsecase.On("GeneratePayouts", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.Anything, nil).Once()
		mPayoutRepo.On("TxCreatePayouts", mock.Anything, mock.Anything).Return(assert.AnError).Once()

		result, err := service.CreatePayouts(ctx, cmd)
		assert.Nil(t, result)
		assert.Error(t, err)

		mItempRepo.AssertNumberOfCalls(t, "GetNotSoldItemsByID", 1)
		mPayoutRepo.AssertNumberOfCalls(t, "GetTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "ConvertTransactionLimit", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "CalculateTotalAmountBySeller", 1)
		mPayoutUsecase.AssertNumberOfCalls(t, "GeneratePayouts", 1)
		mPayoutRepo.AssertNumberOfCalls(t, "TxCreatePayouts", 1)

		mPayoutRepo.AssertExpectations(t)
		mItempRepo.AssertExpectations(t)
		mPayoutUsecase.AssertExpectations(t)
		mItempRepo.Reset()
		mPayoutRepo.Reset()
		mPayoutUsecase.Reset()
	})
}

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetNotSoldItemsByID(ctx context.Context, ids []uuid.UUID) (item.Items, error) {
	args := m.Called(ctx, ids)
	return mocks[item.Items](args)
}

func (m *MockItemRepository) Reset() {
	m.Calls = []mock.Call{}
}

type MockPayoutRepository struct {
	mock.Mock
}

func (m *MockPayoutRepository) GetTransactionLimit(ctx context.Context) (*payout.TransactionLimit, error) {
	args := m.Called(ctx)
	return mocks[*payout.TransactionLimit](args)
}

func (m *MockPayoutRepository) TxCreatePayouts(ctx context.Context, payouts payout.BatchPayouts) error {
	args := m.Called(ctx, payouts)
	return args.Error(0)
}

func (m *MockPayoutRepository) Reset() {
	m.Calls = []mock.Call{}
}

type MockPayoutUsecase struct {
	mock.Mock
}

func (m *MockPayoutUsecase) GeneratePayouts(items item.Items, txl *payout.TransactionLimit, cur valueobject.Currency, amounts payout.TotalAmountsBySeller) (payout.BatchPayouts, error) {
	args := m.Called(items, txl, cur, amounts)
	return mocks[payout.BatchPayouts](args)
}

func (m *MockPayoutUsecase) ConvertTransactionLimit(limit *valueobject.Money, cur valueobject.Currency) (*payout.TransactionLimit, error) {
	args := m.Called(limit, cur)
	return mocks[*payout.TransactionLimit](args)
}

func (m *MockPayoutUsecase) CalculateTotalAmountBySeller(items item.Items, cur valueobject.Currency) (payout.TotalAmountsBySeller, error) {
	args := m.Called(items, cur)
	return mocks[payout.TotalAmountsBySeller](args)
}

func (m *MockPayoutUsecase) Reset() {
	m.Calls = []mock.Call{}
}

func mocks[T any](args mock.Arguments) (T, error) {
	var zero T
	result, ok := args.Get(0).(T)
	if !ok {
		return zero, args.Error(1)
	}
	return result, args.Error(1)
}
