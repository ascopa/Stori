package domain

type Message struct {
	Detail Detail `json:"detail"`
}

type Detail struct {
	MonthlyCreditAverages map[int]string
	MonthlyDebitAverages  map[int]string
	MonthlyTransactions   map[int]int
	Balance               string
	AccountId             string
}
