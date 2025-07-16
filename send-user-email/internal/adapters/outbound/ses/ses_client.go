package ses

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type ISesCustomClient interface {
	SendEmail(email, subject string, body bytes.Buffer) error
}

type SesCustomClient struct {
	sesClient *ses.Client
}

func (s *SesCustomClient) SendEmail(email, subject string, body bytes.Buffer) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String(body.String()),
				},
			},
		},
		Source: aws.String("scopaalejandro+send@gmail.com"),
	}

	_, err := s.sesClient.SendEmail(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func NewSesCustomClient(cfg aws.Config) *SesCustomClient {
	sesClient := ses.NewFromConfig(cfg)
	return &SesCustomClient{sesClient: sesClient}
}
