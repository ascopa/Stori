package service

import (
	"fmt"
	"process-user-transaction/internal/adapters/outbound/repository"
	"strconv"
	"strings"
)

const ()

type Service struct {
	r repository.ITransactionRepository
}

type IService interface {
	ProcessUserTransactions([]repository.Transaction) error
}

func NewService(repository repository.ITransactionRepository) *Service {
	return &Service{
		r: repository,
	}
}

type MonthAverage struct {
	averages    map[int]float64
	trxPerMonth map[int]int
}

func (s *Service) ProcessUserTransactions(transactions []repository.Transaction) error {
	var balance float64
	var monthAverage MonthAverage

	monthAverage.averages = make(map[int]float64)
	monthAverage.trxPerMonth = make(map[int]int)

	for _, transaction := range transactions {
		//err := s.r.PutUser(context.Background(), transaction)
		trxMonth := strings.Split(transaction.CreatedDate, "/")[0]
		month, err := strconv.Atoi(trxMonth)
		if err != nil {
			return err
		}
		monthAverage.trxPerMonth[month] += 1
		if monthAverage.averages[month] != 0.0 {
			monthAverage.averages[month] = (monthAverage.averages[month] + transaction.Amount) / 2
		} else {
			monthAverage.averages[month] = transaction.Amount
		}
		balance += transaction.Amount

		// Print progress
		fmt.Printf("Processed transaction: %+v\n", transaction)
		fmt.Printf("Month: %d, Transaction count: %d, Average amount: %.2f\n",
			month, monthAverage.trxPerMonth[month], monthAverage.averages[month])
		fmt.Printf("Running balance: %.2f\n", balance)
		fmt.Println("------")
	}

	// 🔚 Final Summary
	fmt.Println("========== Final Summary ==========")
	fmt.Printf("Total balance: %.2f\n", balance)
	fmt.Println("Transactions per month:")
	for month, count := range monthAverage.trxPerMonth {
		fmt.Printf("  Month %02d: %d transactions\n", month, count)
	}
	fmt.Println("Average amount per month:")
	for month, avg := range monthAverage.averages {
		fmt.Printf("  Month %02d: %.2f average amount\n", month, avg)
	}
	fmt.Println("===================================")

	return nil
}
