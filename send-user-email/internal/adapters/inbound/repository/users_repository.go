package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"process-user-transaction/internal/core/domain"
)

const (
	TABLE_NAME = "Users"
)

type UsersRepository struct {
	db        *dynamodb.Client
	tableName string
}

type IUsersRepository interface {
	GetUserByAccountId(ctx context.Context, accountId string) (*domain.User, error)
}

func NewUsersRepository(cfg aws.Config) *UsersRepository {
	client := dynamodb.NewFromConfig(cfg)

	return &UsersRepository{
		db:        client,
		tableName: TABLE_NAME,
	}
}

func (r *UsersRepository) GetUserByAccountId(ctx context.Context, accountId string) (*domain.User, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"AccountId": accountId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	out, err := r.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user from DynamoDB: %w", err)
	}

	if out.Item == nil || len(out.Item) == 0 {
		return nil, fmt.Errorf("user not found with accountId: %s", accountId)
	}

	var user domain.User
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}
