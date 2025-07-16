package service

import (
	"bytes"
	"context"
	"github.com/gocarina/gocsv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"path"
	"process-user-transaction/internal/core/domain"
	"testing"
)

func TestHandleS3EventWithLocalCSV(t *testing.T) {
	os.Setenv("NUM_WORKERS", "8")
	os.Setenv("ENABLE_DEBUG", "false")
	var transactions []domain.Transaction
	file, err := os.ReadFile(path.Join("trx.csv"))
	csv := io.NopCloser(bytes.NewReader(file))
	err = gocsv.Unmarshal(csv, &transactions)
	assert.NoError(t, err)

	repo := new(MockRepository)
	repo.On("PutTransaction", mock.Anything, mock.Anything).Return(nil)

	service := NewService(repo)

	userTransactionInfo, err := service.ProcessUserTransactions(transactions)

	assert.Equal(t, 3, userTransactionInfo.MonthlyTransactions[7])
	assert.Equal(t, 5, userTransactionInfo.MonthlyTransactions[8])
	assert.True(t, userTransactionInfo.Balance.Equal(decimal.NewFromFloat(69.5)))
	assert.True(t, userTransactionInfo.MonthlyCreditAverages[7].Equal(decimal.NewFromFloat(60.5)))
	assert.True(t, userTransactionInfo.MonthlyCreditAverages[8].Equal(decimal.NewFromFloat(12.5)))
	assert.True(t, userTransactionInfo.MonthlyDebitAverages[7].Equal(decimal.NewFromFloat(-15.5)))
	assert.True(t, userTransactionInfo.MonthlyDebitAverages[8].Equal(decimal.NewFromFloat(-10)))

	assert.NoError(t, err)
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) PutTransaction(ctx context.Context, transaction domain.Transaction) error {
	args := m.Called(ctx, transaction)

	return args.Error(0)
}
