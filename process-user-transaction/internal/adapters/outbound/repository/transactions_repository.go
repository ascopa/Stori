package repository

import (
	"context"
	"fmt"
	"process-user-transaction/internal/core/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const TABLE_NAME = "Transactions"

type TransactionsRepository struct {
	db        *dynamodb.Client
	tableName string
}

type ITransactionRepository interface {
	PutTransaction(ctx context.Context, transaction domain.Transaction) error
}

func NewTransactionsRepository(cfg aws.Config) *TransactionsRepository {
	client := dynamodb.NewFromConfig(cfg)

	return &TransactionsRepository{
		db:        client,
		tableName: TABLE_NAME,
	}
}

func (r *TransactionsRepository) PutTransaction(ctx context.Context, transaction domain.Transaction) error {
	item, err := attributevalue.MarshalMap(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = r.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}
