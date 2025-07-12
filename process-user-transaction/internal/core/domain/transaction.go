package domain

import "github.com/shopspring/decimal"

type Transaction struct {
	TransactionId string          `json:"transactionId" csv:"TransactionId"`
	CreatedDate   string          `json:"createdDate"   csv:"Date"`
	UpdatedDate   string          `json:"updatedDate"   csv:"UpdatedDate"`
	Amount        decimal.Decimal `json:"amount"        csv:"Amount"`
	AccountId     string          `json:"accountId"     csv:"AccountId"`
}
