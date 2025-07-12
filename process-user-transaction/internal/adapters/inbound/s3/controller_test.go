package s3

import (
	"context"
	_ "context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	_ "path/filepath"
	"process-user-transaction/internal/adapters/outbound/repository"
	"process-user-transaction/internal/core/service"
	"testing"
)

func TestHandleS3EventWithLocalCSV(t *testing.T) {
	filePath := "trx.csv"
	_, err := os.Stat(filePath)
	assert.NoError(t, err, "CSV file should exist")

	// Simulate an S3 Event for a file upload
	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventSource: "aws:s3",
				S3: events.S3Entity{
					Bucket: events.S3Bucket{
						Name: "mock-bucket",
					},
					Object: events.S3Object{
						Key: "trx.csv", // mock S3 key
					},
				},
			},
		},
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	dbClient := dynamodb.NewFromConfig(cfg)
	repo := repository.NewTransactionsRepository(dbClient, "pepe")
	s := service.NewService(repo)

	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	controller := NewController(s, s3Client)

	file, err := os.Open("trx.csv")
	assert.NoError(t, err)
	defer file.Close()

	mockS3 := new(MockS3Client)
	mockS3.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(file), nil)

	err = controller.Handle(ctx, s3Event)
	if err != nil {
		return
	}

	assert.NoError(t, err)
}

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)

	// simulate reading file content
	reader := args.Get(0)
	if reader == nil {
		return nil, args.Error(1)
	}

	return reader.(io.ReadCloser), args.Error(1)
}
