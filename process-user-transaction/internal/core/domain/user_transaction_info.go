package domain

import "github.com/shopspring/decimal"

type UserTransactionInfo struct {
	MonthlyCreditAverages     map[int]decimal.Decimal
	MonthlyDebitAverages      map[int]decimal.Decimal
	MonthlyTransactionsAmount map[int]decimal.Decimal
	MonthlyTransactions       map[int]int64
	Balance                   decimal.Decimal
	AccountId                 string
}
