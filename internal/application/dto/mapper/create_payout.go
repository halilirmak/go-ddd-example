package application_mapper

import (
	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/domain/shared/valueobject"
)

func NewCurrencyFromCommand(command *command.CreatePayoutCommand) (valueobject.Currency, error) {
	return valueobject.NewCurrency(valueobject.Currency(command.Currency))
}
