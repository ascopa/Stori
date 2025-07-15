package event

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
	"os"
)

const EVENT_BUS_NAME = "EVENT_BUS_NAME"
const EVENT_SOURCE_NAME = "EVENTBUS_SOURCE_NAME"

type IEventBridgeCustomClient interface {
	PutEvents(message []byte, detailType string) error
}

type EventBridgeCustomClient struct {
	client *cloudwatchevents.Client
}

func (e *EventBridgeCustomClient) PutEvents(message []byte, detailType string) error {
	_, err := e.client.PutEvents(context.Background(),
		&cloudwatchevents.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{
				{
					EventBusName: aws.String(os.Getenv(EVENT_BUS_NAME)),
					Source:       aws.String(os.Getenv(EVENT_SOURCE_NAME)),

					DetailType: aws.String(detailType),
					Detail:     aws.String(string(message)),
				},
			},
		})
	if err != nil {
		return err
	}

	return nil
}

func NewEventBridgeCustomClient(cfg aws.Config) *EventBridgeCustomClient {
	client := cloudwatchevents.NewFromConfig(cfg)

	return &EventBridgeCustomClient{client: client}
}
