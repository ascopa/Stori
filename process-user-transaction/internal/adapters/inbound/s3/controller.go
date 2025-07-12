package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gocarina/gocsv"
	"log"
	"process-user-transaction/internal/core/domain"
	"process-user-transaction/internal/core/service"
)

type Controller struct {
	s  service.IService
	s3 IS3CustomClient
}

func NewController(service service.IService, s3Client IS3CustomClient) Controller {
	return Controller{
		s:  service,
		s3: s3Client,
	}
}

func (c *Controller) Handle(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3Record := record.S3
		bucket := s3Record.Bucket.Name
		key := s3Record.Object.Key

		log.Printf("New file uploaded: bucket=%s, key=%s\n", bucket, key)

		result, err := c.s3.GetObject(ctx, bucket, key)

		if err != nil {
			return fmt.Errorf("failed to get S3 object: %w", err)
		}

		var transactions []domain.Transaction
		err = gocsv.Unmarshal(result, &transactions)
		if err != nil {
			return fmt.Errorf("failed to parse CSV: %w", err)
		}

		userTransactionInfo, err := c.s.ProcessUserTransactions(transactions)
		if err != nil {
			return fmt.Errorf("error processing file: %w", err)
		}

		fmt.Printf("File processed: s3://%s/%s\n", bucket, key)
		fmt.Println("Transaction user info: ", userTransactionInfo)
	}

	return nil
}
