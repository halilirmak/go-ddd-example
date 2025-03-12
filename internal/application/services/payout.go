package services

import (
	"context"
	"fmt"

	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	application_mapper "github.com/cryptoPickle/go-ddd-example/internal/application/dto/mapper"
	"github.com/cryptoPickle/go-ddd-example/internal/application/interfaces"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/item"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/payout"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/usecase"
)

type PayoutService struct {
	itemsRepository  item.ItemRepository
	payoutRepository payout.PayoutRepository
	payoutusecase    usecase.PayoutCalculator
}

func NewPayoutService(itemrepository item.ItemRepository, payoutRepository payout.PayoutRepository, payoutusecase usecase.PayoutCalculator) interfaces.PayoutService {
	return &PayoutService{
		itemsRepository:  itemrepository,
		payoutRepository: payoutRepository,
		payoutusecase:    payoutusecase,
	}
}

var ErrInvalidInput = fmt.Errorf("some items already paid or invalid")

func (p *PayoutService) CreatePayouts(ctx context.Context, command *command.CreatePayoutCommand) (*command.CreatePayoutCommandResult, error) {
	currency, err := application_mapper.NewCurrencyFromCommand(command)
	if err != nil {
		return nil, err
	}

	items, err := p.itemsRepository.GetNotSoldItemsByID(ctx, command.Items)
	if err != nil {
		return nil, err
	}

	if !command.Items.IsEqual(len(items)) {
		return nil, errors.NewInvalidInputError("some items already paid or invalid", "payout-application")
	}

	txl, err := p.payoutRepository.GetTransactionLimit(ctx)
	if err != nil {
		return nil, err
	}

	txlConverted, err := p.payoutusecase.ConvertTransactionLimit(txl.Limit(), currency)
	if err != nil {
		return nil, err
	}

	totalPayoutBySeller, err := p.payoutusecase.CalculateTotalAmountBySeller(items, currency)
	if err != nil {
		return nil, err
	}

	payouts, err := p.payoutusecase.GeneratePayouts(items, txlConverted, currency, totalPayoutBySeller)
	if err != nil {
		return nil, err
	}

	if err := p.payoutRepository.TxCreatePayouts(ctx, payouts); err != nil {
		return nil, err
	}

	return application_mapper.NewPayoutResultFromEntity(payouts, currency), nil
}
