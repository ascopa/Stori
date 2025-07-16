package event

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
	"os"
)

const EVENT_BUS_NAME = "EVENT_BUS_NAME"
const EVENT_SOURCE_NAME = "EVENT_SOURCE_NAME"

type IEventBridgeCustomClient interface {
	PutEvents(message []byte, detailType string) error
}

type EventBridgeCustomClient struct {
	client *cloudwatchevents.Client
}

func (e *EventBridgeCustomClient) PutEvents(message []byte, detailType string) error {
	entry := types.PutEventsRequestEntry{
		EventBusName: aws.String(os.Getenv(EVENT_BUS_NAME)),
		Source:       aws.String(os.Getenv(EVENT_SOURCE_NAME)),
		DetailType:   aws.String(detailType),
		Detail:       aws.String(string(message)),
	}

	_, err := e.client.PutEvents(context.Background(), &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{entry},
	})
	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}
	return nil
}

func NewEventBridgeCustomClient(cfg aws.Config) *EventBridgeCustomClient {
	client := cloudwatchevents.NewFromConfig(cfg)

	return &EventBridgeCustomClient{client: client}
}
