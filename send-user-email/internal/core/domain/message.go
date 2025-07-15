package domain

import "github.com/shopspring/decimal"

type Message struct {
	MonthlyAverages           map[int]decimal.Decimal
	MonthlyTransactionsAmount map[int]decimal.Decimal
	MonthlyTransactions       map[int]int64
	Balance                   string
	AccountId                 string
}
