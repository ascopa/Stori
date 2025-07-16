package s3

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mockService.On("ProcessUserTransactions", mock.Anything, mock.Anything, mock.Anything).Return(domain.UserTransactionInfo{Balance: decimal.NewFromFloat(15.55)}, nil)
	mockS3Client := new(MockS3Client)
	mockS3Client.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(file), nil)
	mockEBClient := new(MockEBClient)
	mockEBClient.On("PutEvents", mock.Anything, mock.Anything).Return(nil)
	controller := NewController(mockService, mockS3Client, mockEBClient)

	err = controller.Handle(context.Background(), events.S3Event{})
	assert.NoError(t, err)
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
	reader := args.Get(0)
	if reader == nil {
		return nil, args.Error(1)
	}
	return reader.(io.ReadCloser), args.Error(1)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) ProcessUserTransactions(transactions []domain.Transaction) (domain.UserTransactionInfo, error) {
	args := m.Called(transactions)
	return args.Get(0).(domain.UserTransactionInfo), args.Error(1)
}
