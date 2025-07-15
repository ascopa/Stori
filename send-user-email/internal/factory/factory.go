package factory

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"process-user-transaction/internal/adapters/inbound/repository"
	"process-user-transaction/internal/adapters/inbound/sqs"
	"process-user-transaction/internal/adapters/outbound/ses"
	"process-user-transaction/internal/core/service"
)

type Factory struct {
}

func (f *Factory) Start(ctx context.Context, sqsEvent events.SQSEvent) error {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("couldn't load default configuration: %w", err)
	}

	repo := repository.NewUsersRepository(sdkConfig)
	sesCustomClient := ses.NewSesCustomClient(sdkConfig)
	s := service.NewService(repo, sesCustomClient)
	c := sqs.NewController(s)

	return c.Handle(ctx, sqsEvent)
}
