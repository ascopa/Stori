package domain

type MonthSummary struct {
	Name             string
	TransactionCount int
	CreditAverage    string
	DebitAverage     string
}

type EmailData struct {
	Name    string
	Summary []MonthSummary
	Balance string
}
