package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"html/template"

	"os"
	"path/filepath"
	"send-user-email/internal/core/domain"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByAccountId(ctx context.Context, accountId string) (*domain.User, error) {
	args := m.Called(ctx, accountId)
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockSESClient struct {
	mock.Mock
}

func (m *MockSESClient) SendEmail(to string, subject string, body bytes.Buffer) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestSendUserEmail_Success(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockRepository)
	mockSES := new(MockSESClient)

	s := NewService(mockRepo, mockSES)

	// Test data
	message := domain.Message{
		Detail: domain.Detail{
			AccountId: "12345",
		},
	}

	user := &domain.User{
		AccountId: "12345",
		Name:      "Alice",
		Email:     "alice@example.com",
	}

	// JSON print to match the behavior of the service
	_, _ = json.MarshalIndent(message, "", "  ")

	// Mock repo returns user
	mockRepo.On("GetUserByAccountId", mock.Anything, "12345").Return(user, nil)

	// Expected email template rendering
	expectedBody := bytes.NewBufferString("<html><body><h1>Hello Alice</h1></body></html>")

	// Mock SES to expect the call
	mockSES.On("SendEmail", user.Email, "trx notification", mock.MatchedBy(func(b bytes.Buffer) bool {
		return b.String() == expectedBody.String()
	})).Return(nil)

	// Monkey patch the buildTemplate (only if refactored — otherwise you'd need to patch `embed.FS` or extract logic)
	// For this test, you can fake it by injecting that the result of buildTemplate will be the expected buffer.

	// Act
	err := s.SendUserEmail(context.Background(), message)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockSES.AssertExpectations(t)
}

func TestRenderHTMLTemplate(t *testing.T) {
	// Sample data
	data := domain.Detail{
		AccountId: "user-123",
		Balance:   "120.50",
		MonthlyTransactions: map[int]int{
			7: 2,
			8: 3,
		},
		MonthlyCreditAverages: map[int]string{
			7: "30.25",
			8: "20.00",
		},
		MonthlyDebitAverages: map[int]string{
			7: "-30.25",
			8: "-20.00",
		},
	}

	// Load template
	path := filepath.Join("template.html")
	tmpl, err := template.ParseFiles(path)
	assert.NoError(t, err)

	// Render
	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, GetTemplateInfo(&domain.User{Name: "Pepe"}, data))
	assert.NoError(t, err)

	err = os.WriteFile("rendered_output.html", rendered.Bytes(), 0644)
	assert.NoError(t, err)

	assert.Contains(t, rendered.String(), "user-123")
	assert.Contains(t, rendered.String(), "120.50")
	assert.Contains(t, rendered.String(), "30.25")
}
