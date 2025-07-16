package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"send-user-email/internal/adapters/inbound/repository"
	"send-user-email/internal/adapters/inbound/s3"
	"send-user-email/internal/adapters/outbound/ses"
	"send-user-email/internal/core/domain"
	"time"
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

type MonthSummary struct {
	Name             string
	TransactionCount int
	CreditAverage    string
	DebitAverage     string
}

type EmailData struct {
	Name    string
	Summary []MonthSummary
	Balance string
}

func (s *Service) SendUserEmail(ctx context.Context, message domain.Message) error {

	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	user, err := s.r.GetUserByAccountId(ctx, message.Detail.AccountId)
	if err != nil {
		return fmt.Errorf("failed to retrieve username: %w", err)
	}

	template, err := s.buildTemplate(ctx, user, message)
	if err != nil {
		return fmt.Errorf("failed to build template: %w", err)
	}

	err = s.ses.SendEmail(user.Email, "trx notification", template)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *Service) buildTemplate(ctx context.Context, user *domain.User, msg domain.Message) (bytes.Buffer, error) {

	resp, err := s.s3.GetObject(ctx, "stori-user-transactions-email-templates", "template.html")
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

	fmt.Println("HTML:\n", string(bodyBytes))

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, GetTemplateInfo(user, msg.Detail))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Println("HTML with content:\n", string(rendered.Bytes()))

	return rendered, nil
}

func GetTemplateInfo(user *domain.User, detail domain.Detail) EmailData {
	var summary []MonthSummary

	for monthNum := 1; monthNum <= 12; monthNum++ {
		count, hasTx := detail.MonthlyTransactions[monthNum]
		if !hasTx {
			continue
		}
		summary = append(summary, MonthSummary{
			Name:             time.Month(monthNum).String(),
			TransactionCount: count,
			CreditAverage:    detail.MonthlyCreditAverages[monthNum],
			DebitAverage:     detail.MonthlyDebitAverages[monthNum],
		})
	}

	emailData := EmailData{
		Name:    user.Name,
		Summary: summary,
		Balance: detail.Balance,
	}
	return emailData
}
