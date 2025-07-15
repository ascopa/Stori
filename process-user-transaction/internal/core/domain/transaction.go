package domain

import "github.com/shopspring/decimal"

type Transaction struct {
	TransactionId string          `json:"transactionId" csv:"TransactionId"`
	CreatedDate   string          `json:"createdDate"   csv:"Date"`
	UpdatedDate   string          `json:"updatedDate"   csv:"UpdatedDate"`
	Amount        decimal.Decimal `json:"amount"        csv:"Amount"`
	AccountId     string          `json:"accountId"     csv:"AccountId"`
}

type TransactionItem struct {
	TransactionId string `DynamoDB:"TransactionId"`
	CreatedDate   string `DynamoDB:"CreatedDate"`
	UpdatedDate   string `DynamoDB:"UpdatedDate"`
	Amount        string `DynamoDB:"Amount"`
	AccountId     string `DynamoDB:"AccountId"`
}

func ToTransactionItem(t Transaction) TransactionItem {
	return TransactionItem{
		TransactionId: t.TransactionId,
		CreatedDate:   t.CreatedDate,
		UpdatedDate:   t.UpdatedDate,
		Amount:        t.Amount.String(),
		AccountId:     t.AccountId,
	}
}
