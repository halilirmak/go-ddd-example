package services

import (
	"context"

	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/application/interfaces"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger"
)

type PayoutServiceWithLogger struct {
	service interfaces.PayoutService
	logger  logger.Logger
}

func NewPayoutServiceWithLogger(service interfaces.PayoutService, logger logger.Logger) interfaces.PayoutService {
	return &PayoutServiceWithLogger{
		service: service,
		logger:  logger,
	}
}

func (p *PayoutServiceWithLogger) CreatePayouts(ctx context.Context, command *command.CreatePayoutCommand) (*command.CreatePayoutCommandResult, error) {
	cmd, err := p.service.CreatePayouts(ctx, command)
	if err != nil {
		if err, ok := err.(errors.DetailedError); ok {
			p.logger.Errorf("context: %s | errType: %s | errMsg: %s | unwrapped: %s\n", err.Context(), err.ErrorType(), err.Error(), err.UnWrap())
			return nil, err
		}
		p.logger.Errorf("unknown err: %s\n", err)
		return nil, err
	}

	return cmd, nil
}
