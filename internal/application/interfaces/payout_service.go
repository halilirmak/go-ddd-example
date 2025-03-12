package interfaces

import (
	"context"

	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
)

type PayoutService interface {
	// CreatePayouts creates payouts for multiple sellers based on given command
	CreatePayouts(ctx context.Context, command *command.CreatePayoutCommand) (*command.CreatePayoutCommandResult, error)
}
