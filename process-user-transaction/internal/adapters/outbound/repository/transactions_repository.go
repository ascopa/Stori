package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Transaction struct {
	TransactionId string  `json:"transactionId"`
	CreatedDate   string  `json:"createdDate"`
	UpdatedDate   string  `json:"updatedDate"`
	Amount        float64 `json:"amount"`
}

type TransactionsRepository struct {
	db        *dynamodb.Client
	tableName string
}

type ITransactionRepository interface {
	PutUser(ctx context.Context, transaction Transaction) error
}

func NewTransactionsRepository(db *dynamodb.Client, tableName string) *TransactionsRepository {
	return &TransactionsRepository{
		db:        db,
		tableName: tableName,
	}
}

func (r *TransactionsRepository) PutUser(ctx context.Context, transaction Transaction) error {
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
