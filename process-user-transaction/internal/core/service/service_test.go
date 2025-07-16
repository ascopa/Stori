package service

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"process-user-transaction/internal/core/domain"
	"testing"
)

func TestHandleS3EventWithLocalCSV(t *testing.T) {
	var SampleTransactions = []domain.Transaction{
		{TransactionId: "0", CreatedDate: "7/15", Amount: decimal.NewFromFloat(60.5)},
		{TransactionId: "1", CreatedDate: "7/28", Amount: decimal.NewFromFloat(-10.3)},
		{TransactionId: "2", CreatedDate: "8/2", Amount: decimal.NewFromFloat(-20.46)},
		{TransactionId: "3", CreatedDate: "8/13", Amount: decimal.NewFromFloat(10)},
		{TransactionId: "4", CreatedDate: "8/11", Amount: decimal.NewFromFloat(10)},
		{TransactionId: "5", CreatedDate: "8/12", Amount: decimal.NewFromFloat(-10)},
		{TransactionId: "6", CreatedDate: "8/14", Amount: decimal.NewFromFloat(15)},
		{TransactionId: "7", CreatedDate: "7/11", Amount: decimal.NewFromFloat(10)},
	}

	repo := new(MockRepository)
	repo.On("PutTransaction", mock.Anything, mock.Anything).Return(nil)

	service := NewService(repo)

	userTransactionInfo, err := service.ProcessUserTransactions(SampleTransactions)

	assert.Equal(t, userTransactionInfo.MonthlyTransactions[7], int64(2))
	assert.Equal(t, userTransactionInfo.MonthlyTransactions[8], int64(3))
	assert.True(t, userTransactionInfo.Balance.Equal(decimal.NewFromFloat(49.74)))
	//assert.True(t, userTransactionInfo.MonthlyAverages[7].Equal(decimal.NewFromFloat(25.1)))
	//assert.True(t, userTransactionInfo.MonthlyAverages[8].Equal(decimal.NewFromFloat(-0.15)))

	assert.NoError(t, err)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) PutTransaction(ctx context.Context, transaction domain.Transaction) error {
	args := m.Called(ctx, transaction)

	return args.Error(0)
}
