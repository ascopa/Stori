package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"process-user-transaction/internal/core/domain"
	"testing"
)

func TestHandleS3EventWithLocalCSV(t *testing.T) {
	file, err := os.Open("trx.csv")
	assert.NoError(t, err)
	defer file.Close()

	mockService := new(MockService)
	mockS3Client := new(MockS3Client)
	mockEBClient := new(MockEBClient)
	mockS3Client.On("GetObject", mock.Anything, "stori-user-transactions", "transactions.csv").Return(io.NopCloser(file), nil)
	mockService.On("ProcessUserTransactions", mock.Anything, mock.Anything, mock.Anything).Return(domain.UserTransactionInfo{}, nil)
	mockEBClient.On("PutEvents", mock.Anything, mock.Anything).Return(nil)
	controller := NewController(mockService, mockS3Client, mockEBClient)

	err = controller.Handle(context.Background(), buildSQSEvent())
	assert.NoError(t, err)
	mockS3Client.AssertExpectations(t)
	mockService.AssertExpectations(t)
	mockEBClient.AssertExpectations(t)
}

func TestHandle_GetObjectError(t *testing.T) {
	mockService := new(MockService)
	mockS3Client := new(MockS3Client)
	mockEBClient := new(MockEBClient)
	mockS3Client.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(nil), fmt.Errorf("failed to get S3 object"))
	controller := NewController(mockService, mockS3Client, mockEBClient)

	err := controller.Handle(context.Background(), buildSQSEvent())
	require.ErrorContains(t, err, "failed to get S3 object")
}

func TestHandle_ProcessUserTransactionsError(t *testing.T) {
	file, err := os.Open("trx.csv")
	assert.NoError(t, err)
	defer file.Close()

	mockService := new(MockService)
	mockS3Client := new(MockS3Client)
	mockEBClient := new(MockEBClient)
	mockS3Client.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(file), nil)
	mockService.On("ProcessUserTransactions", mock.Anything, mock.Anything, mock.Anything).Return(domain.UserTransactionInfo{}, fmt.Errorf("failed to process trx"))
	controller := NewController(mockService, mockS3Client, mockEBClient)

	err = controller.Handle(context.Background(), buildSQSEvent())
	require.ErrorContains(t, err, "failed to process trx")
}
func TestHandle_PutEventsError(t *testing.T) {
	file, err := os.Open("trx.csv")
	assert.NoError(t, err)
	defer file.Close()

	mockService := new(MockService)
	mockS3Client := new(MockS3Client)
	mockEBClient := new(MockEBClient)
	mockS3Client.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(file), nil)
	mockService.On("ProcessUserTransactions", mock.Anything, mock.Anything, mock.Anything).Return(domain.UserTransactionInfo{}, nil)
	mockEBClient.On("PutEvents", mock.Anything, mock.Anything).Return(fmt.Errorf("failed to put event"))
	controller := NewController(mockService, mockS3Client, mockEBClient)

	err = controller.Handle(context.Background(), buildSQSEvent())
	require.ErrorContains(t, err, "failed to put event")
}

type MockEBClient struct {
	mock.Mock
}

func (m *MockEBClient) PutEvents(message []byte, detailType string) error {
	args := m.Called(message, detailType)
	return args.Error(0)
}

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(ctx context.Context, bucket string, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) ProcessUserTransactions(transactions []domain.Transaction) (domain.UserTransactionInfo, error) {
	args := m.Called(transactions)
	return args.Get(0).(domain.UserTransactionInfo), args.Error(1)
}

func buildSQSEvent() events.SQSEvent {
	bodyStruct := struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	}{
		Bucket: "stori-user-transactions",
		Key:    "transactions.csv",
	}

	bodyBytes, _ := json.Marshal(bodyStruct)

	return events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(bodyBytes),
			},
		},
	}
}
