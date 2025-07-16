package event

import (
	"context"
	"encoding/json"
	"fmt"
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
	fmt.Println("Event Source:", os.Getenv(EVENT_SOURCE_NAME))
	fmt.Println("Event Bus:", os.Getenv(EVENT_BUS_NAME))

	entry := types.PutEventsRequestEntry{
		EventBusName: aws.String(os.Getenv(EVENT_BUS_NAME)),
		Source:       aws.String(os.Getenv(EVENT_SOURCE_NAME)),
		DetailType:   aws.String(detailType),
		Detail:       aws.String(string(message)),
	}

	// Optional: print exactly what will be sent to EventBridge
	debug := map[string]interface{}{
		"source":       *entry.Source,
		"detail-type":  *entry.DetailType,
		"eventBusName": *entry.EventBusName,
		"detail":       json.RawMessage(*entry.Detail), // keep JSON format
	}
	logBytes, _ := json.MarshalIndent(debug, "", "  ")
	fmt.Println("Event to be sent to EventBridge:\n", string(logBytes))

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
