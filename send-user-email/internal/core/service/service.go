package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"send-user-email/internal/adapters/inbound/repository"
	"send-user-email/internal/adapters/inbound/s3"
	"send-user-email/internal/adapters/outbound/ses"
	"send-user-email/internal/core/domain"
	"time"
)

const (
	BUCKET_NAME   = "BUCKET_NAME"
	BUCKET_KEY    = "BUCKET_KEY"
	EMAIL_SUBJECT = "EMAIL_SUBJECT"
)

type Service struct {
	r   repository.IUsersRepository
	ses ses.ISesCustomClient
	s3  s3.IS3CustomClient
}

type IService interface {
	SendUserEmail(ctx context.Context, message domain.Message) error
}

func NewService(repository repository.IUsersRepository, sesClient ses.ISesCustomClient, s3 s3.IS3CustomClient) *Service {
	return &Service{
		r:   repository,
		ses: sesClient,
		s3:  s3,
	}
}

func (s *Service) SendUserEmail(ctx context.Context, message domain.Message) error {
	user, err := s.r.GetUserByAccountId(ctx, message.Detail.AccountId)
	if err != nil {
		return fmt.Errorf("failed to retrieve username: %w", err)
	}

	builtTemplate, err := s.buildTemplate(ctx, user, message)
	if err != nil {
		return fmt.Errorf("failed to build template: %w", err)
	}

	err = s.ses.SendEmail(user.Email, os.Getenv(EMAIL_SUBJECT), builtTemplate)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Service) buildTemplate(ctx context.Context, user *domain.User, msg domain.Message) (bytes.Buffer, error) {
	resp, err := s.s3.GetObject(ctx, os.Getenv(BUCKET_NAME), os.Getenv(BUCKET_KEY))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to retrieve s3 template: %w", err)
	}

	bodyBytes, err := io.ReadAll(resp)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to read template content: %w", err)
	}

	tmpl, err := template.New("email").Parse(string(bodyBytes))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to parse template: %w", err)
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, getTemplateInfo(user, msg.Detail))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute template: %w", err)
	}

	return rendered, nil
}

func getTemplateInfo(user *domain.User, detail domain.Detail) domain.EmailData {
	var summary []domain.MonthSummary

	for monthNum := 1; monthNum <= 12; monthNum++ {
		count, hasTx := detail.MonthlyTransactions[monthNum]
		if !hasTx {
			continue
		}
		summary = append(summary, domain.MonthSummary{
			Name:             time.Month(monthNum).String(),
			TransactionCount: count,
			CreditAverage:    detail.MonthlyCreditAverages[monthNum],
			DebitAverage:     detail.MonthlyDebitAverages[monthNum],
		})
	}

	emailData := domain.EmailData{
		Name:    user.Name,
		Summary: summary,
		Balance: detail.Balance,
	}
	return emailData
}
