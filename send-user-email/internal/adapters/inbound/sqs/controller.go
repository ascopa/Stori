package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"process-user-transaction/internal/core/domain"
	"process-user-transaction/internal/core/service"
)

type Controller struct {
	s service.IService
}

func NewController(service service.IService) Controller {
	return Controller{
		s: service,
	}
}

func (c *Controller) Handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		fmt.Printf("Record: %+v\n", record)
		var msg domain.Message
		err := json.Unmarshal([]byte(record.Body), &msg)
		if err != nil {
			return fmt.Errorf("failed to parse SQS message: %w", err)
		}

		err = c.s.SendUserEmail(ctx, msg)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}
