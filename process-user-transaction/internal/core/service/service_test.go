package service

import (
	"context"
	_ "context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	_ "path/filepath"
	"process-user-transaction/internal/adapters/outbound/repository"
	"testing"
)

func TestHandleS3EventWithLocalCSV(t *testing.T) {
	var SampleTransactions = []repository.Transaction{
		{TransactionId: "0", CreatedDate: "7/15", Amount: 60.5},
		{TransactionId: "1", CreatedDate: "7/28", Amount: -10.3},
		{TransactionId: "2", CreatedDate: "8/2", Amount: -20.46},
		{TransactionId: "3", CreatedDate: "8/13", Amount: 10},
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	dbClient := dynamodb.NewFromConfig(cfg)
	repo := repository.NewTransactionsRepository(dbClient, "pepe")
	service := NewService(repo)
	err = service.ProcessUserTransactions(SampleTransactions)
	assert.NoError(t, err)
}
