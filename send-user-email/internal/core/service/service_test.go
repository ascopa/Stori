package service

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"send-user-email/internal/core/domain"
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
	mockRepo := new(MockUserRepo)
	mockSes := new(MockSesClient)
	mockS3 := new(MockS3Client)

	user := &domain.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	msg := domain.Message{
		Detail: domain.Detail{
			AccountId: "abc123",
			Balance:   "69.74",
			MonthlyTransactions: map[int]int{
				7: 2,
				9: 3,
			},
			MonthlyCreditAverages: map[int]string{
				7: "30.25",
				9: "11.11",
			},
			MonthlyDebitAverages: map[int]string{
				7: "-30.25",
				9: "-3.00",
			},
		},
	}

	htmlTemplate, _ := os.ReadFile("template.html")

	mockRepo.On("GetUserByAccountId", mock.Anything, "abc123").Return(user, nil)
	mockS3.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader(string(htmlTemplate))), nil)
	mockSes.On("SendEmail", user.Email, "trx notification", mock.Anything).Return(nil)

	svc := NewService(mockRepo, mockSes, mockS3)
	err := svc.SendUserEmail(context.Background(), msg)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockS3.AssertExpectations(t)
	mockSes.AssertExpectations(t)
}

func TestBuildTemplate_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockSes := new(MockSesClient)
	mockS3 := new(MockS3Client)

	user := &domain.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	msg := domain.Message{
		Detail: domain.Detail{
			AccountId: "abc123",
			Balance:   "69.74",
			MonthlyTransactions: map[int]int{
				7: 2,
				9: 3,
			},
			MonthlyCreditAverages: map[int]string{
				7: "30.25",
				9: "11.11",
			},
			MonthlyDebitAverages: map[int]string{
				7: "-30.25",
				9: "-3.00",
			},
		},
	}

	htmlTemplate, _ := os.ReadFile(path.Join("template", "template.html"))
	expectedTemplate, _ := os.ReadFile(path.Join("template", "expected_template.html"))

	mockRepo.On("GetUserByAccountId", mock.Anything, "abc123").Return(user, nil)
	mockS3.On("GetObject", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader(string(htmlTemplate))), nil)
	mockSes.On("SendEmail", user.Email, "trx notification", mock.Anything).Return(nil)

	svc := NewService(mockRepo, mockSes, mockS3)
	renderedTemplate, err := svc.buildTemplate(context.TODO(), user, msg)

	assert.Equal(t, string(expectedTemplate), renderedTemplate.String())

	assert.NoError(t, err)
	mockS3.AssertExpectations(t)
}
