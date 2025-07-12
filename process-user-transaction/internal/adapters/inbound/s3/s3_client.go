package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type IS3CustomClient interface {
	GetObject(ctx context.Context, bucket string, key string) (io.ReadCloser, error)
}

type S3CustomClient struct {
	client *s3.Client
}

func (r *S3CustomClient) GetObject(ctx context.Context, bucket string, key string) (io.ReadCloser, error) {
	out, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func NewS3CustomClient(cfg aws.Config) *S3CustomClient {
	s3Client := s3.NewFromConfig(cfg)

	return &S3CustomClient{client: s3Client}
}
