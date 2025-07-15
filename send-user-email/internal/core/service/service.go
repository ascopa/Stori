package service

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"process-user-transaction/internal/adapters/inbound/repository"
	"process-user-transaction/internal/adapters/outbound/ses"
	"process-user-transaction/internal/core/domain"
)

type Service struct {
	r repository.IUsersRepository
	s ses.ISesCustomClient
}

type IService interface {
	SendUserEmail(ctx context.Context, message domain.Message) error
}

func NewService(repository repository.IUsersRepository, sesClient ses.ISesCustomClient) *Service {
	return &Service{
		r: repository,
		s: sesClient,
	}
}

type EmailData struct {
	Name string
}

func (s *Service) SendUserEmail(ctx context.Context, message domain.Message) error {

	user, err := s.r.GetUserByAccountId(ctx, message.AccountId)
	if err != nil {
		return fmt.Errorf("failed to retrieve username: %w", err)
	}

	template, err := buildTemplate(err, user)
	if err != nil {
		return fmt.Errorf("failed to build template: %w", err)
	}

	err = s.s.SendEmail(user.Email, "trx notification", template)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func buildTemplate(err error, user *domain.User) (bytes.Buffer, error) {
	var emailTemplate embed.FS

	tmplBytes, err := emailTemplate.ReadFile("template/email.html")
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to read template: %w", err)
	}

	tmpl, err := template.New("email").Parse(string(tmplBytes))
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to parse template: %w", err)
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, EmailData{Name: user.Name})
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("failed to execute template: %w", err)
	}
	return rendered, nil
}
