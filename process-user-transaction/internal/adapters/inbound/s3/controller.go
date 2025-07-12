package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gocarina/gocsv"
	"log"
	"process-user-transaction/internal/adapters/outbound/repository"
	"process-user-transaction/internal/core/service"
)

type Controller struct {
	s        service.IService
	s3Client *s3.Client
}

func NewController(service service.IService, s3Client *s3.Client) Controller {
	return Controller{
		s:        service,
		s3Client: s3Client,
	}
}

func (c *Controller) Handle(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		s3Record := record.S3
		bucket := s3Record.Bucket.Name
		key := s3Record.Object.Key

		log.Printf("New file uploaded: bucket=%s, key=%s\n", bucket, key)

		result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		if err != nil {
			return fmt.Errorf("failed to get S3 object: %w", err)
		}

		defer result.Body.Close()

		var transactions []repository.Transaction
		err = gocsv.Unmarshal(result.Body, &transactions)
		if err != nil {
			return fmt.Errorf("failed to parse CSV: %w", err)
		}

		err = c.s.ProcessUserTransactions(transactions)
		if err != nil {
			return fmt.Errorf("error processing file: %w", err)
		}

		fmt.Printf("File processed: s3Client://%s/%s\n", bucket, key)
	}

	return nil
}
