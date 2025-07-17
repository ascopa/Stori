package factory

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"process-user-transaction/internal/adapters/inbound/s3"
	"process-user-transaction/internal/adapters/inbound/sqs"
	"process-user-transaction/internal/adapters/outbound/event"
	"process-user-transaction/internal/adapters/outbound/repository"
	"process-user-transaction/internal/core/service"
)

type Factory struct {
}

func (f *Factory) Start(ctx context.Context, sqsEvent events.SQSEvent) error {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("couldn't load default configuration: %w", err)
	}

	repo := repository.NewTransactionsRepository(sdkConfig)
	s3Client := s3.NewS3CustomClient(sdkConfig)
	eventsClient := event.NewEventBridgeCustomClient(sdkConfig)
	s := service.NewService(repo)
	c := sqs.NewController(s, s3Client, eventsClient)

	return c.Handle(ctx, sqsEvent)
}
