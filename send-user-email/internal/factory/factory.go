package factory

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"send-user-email/internal/adapters/inbound/repository"
	"send-user-email/internal/adapters/inbound/s3"
	"send-user-email/internal/adapters/inbound/sqs"
	"send-user-email/internal/adapters/outbound/ses"
	"send-user-email/internal/core/service"
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
	s3CustomClient := s3.NewS3CustomClient(sdkConfig)
	s := service.NewService(repo, sesCustomClient, s3CustomClient)
	c := sqs.NewController(s)

	return c.Handle(ctx, sqsEvent)
}
