package sqs

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"send-user-email/internal/core/domain"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) SendUserEmail(ctx context.Context, message domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func TestHandle_Success(t *testing.T) {
	mockSvc := new(MockService)
	controller := NewController(mockSvc)

	msg := domain.Message{
		Detail: domain.Detail{
			AccountId: "abc123",
			Balance:   "100.00",
			MonthlyTransactions: map[int]int{
				7: 2,
			},
			MonthlyCreditAverages: map[int]string{
				7: "50.00",
			},
			MonthlyDebitAverages: map[int]string{
				7: "-10.00",
			},
		},
	}

	bodyBytes, _ := json.Marshal(msg)
	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(bodyBytes),
			},
		},
	}

	mockSvc.On("SendUserEmail", mock.Anything, msg).Return(nil)

	err := controller.Handle(context.Background(), event)

	assert.NoError(t, err)
	mockSvc.AssertExpectations(t)
}
