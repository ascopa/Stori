package service_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"send-user-email/internal/core/domain"
	"send-user-email/internal/core/service"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserByAccountId(ctx context.Context, accountId string) (*domain.User, error) {
	args := m.Called(ctx, accountId)
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockSesClient struct {
	mock.Mock
}

func (m *MockSesClient) SendEmail(email, subject string, body bytes.Buffer) error {
	args := m.Called(email, subject, body.String())
	return args.Error(0)
}

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestSendUserEmail_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockUserRepo)
	mockSes := new(MockSesClient)
	mockS3 := new(MockS3Client)

	user := &domain.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	htmlTemplate := `
		<html><body>
			<h1>Hello {{.Name}}</h1>
			<p>Your balance is {{.Balance}}</p>
		</body></html>`

	mockRepo.On("GetUserByAccountId", mock.Anything, "abc123").Return(user, nil)
	mockS3.On("GetObject", mock.Anything, "stori-user-transactions-email-templates", "template.html").
		Return(io.NopCloser(strings.NewReader(htmlTemplate)), nil)
	mockSes.On("SendEmail", user.Email, "trx notification", mock.MatchedBy(func(body string) bool {
		return strings.Contains(body, "John Doe") && strings.Contains(body, "69.74")
	})).Return(nil)

	msg := domain.Message{
		Detail: domain.Detail{
			AccountId: "abc123",
			Balance:   "69.74",
			MonthlyTransactions: map[int]int{
				7: 2,
			},
			MonthlyCreditAverages: map[int]string{
				7: "20.00",
			},
			MonthlyDebitAverages: map[int]string{
				7: "-5.00",
			},
		},
	}

	svc := service.NewService(mockRepo, mockSes, mockS3)
	err := svc.SendUserEmail(ctx, msg)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
	mockSes.AssertExpectations(t)
}
