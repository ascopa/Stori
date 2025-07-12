package service

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"process-user-transaction/internal/adapters/outbound/repository"
	"process-user-transaction/internal/core/domain"
	"strconv"
	"strings"
)

const ()

type Service struct {
	r repository.ITransactionRepository
}

type IService interface {
	ProcessUserTransactions([]domain.Transaction) (domain.UserTransactionInfo, error)
}

func NewService(repository repository.ITransactionRepository) *Service {
	return &Service{
		r: repository,
	}
}

func (s *Service) ProcessUserTransactions(transactions []domain.Transaction) (domain.UserTransactionInfo, error) {
	var userTransactionInfo domain.UserTransactionInfo

	userTransactionInfo.MonthlyAverages = make(map[int]decimal.Decimal)
	userTransactionInfo.MonthlyTransactionsAmount = make(map[int]decimal.Decimal)
	userTransactionInfo.MonthlyTransactions = make(map[int]int64)

	for _, transaction := range transactions {
		err := s.r.PutTransaction(context.Background(), transaction)
		trxMonth := strings.Split(transaction.CreatedDate, "/")[0]
		month, err := strconv.Atoi(trxMonth)
		if err != nil {
			return domain.UserTransactionInfo{}, err
		}
		userTransactionInfo.MonthlyTransactions[month] += 1

		userTransactionInfo.MonthlyAverages[month] = userTransactionInfo.MonthlyAverages[month].Add(transaction.Amount)

		userTransactionInfo.Balance = userTransactionInfo.Balance.Add(transaction.Amount)

		fmt.Printf("Processed transaction: %+v\n", transaction)
		fmt.Printf("Month: %d, Transaction count: %d, ", month, userTransactionInfo.MonthlyTransactions[month])
		fmt.Println("Average amount: ", userTransactionInfo.MonthlyAverages[month])
		fmt.Println("Running balance: ", userTransactionInfo.Balance)
		fmt.Println("------")
	}

	for i, monthlyAverage := range userTransactionInfo.MonthlyAverages {
		userTransactionInfo.MonthlyAverages[i] = monthlyAverage.DivRound(decimal.NewFromInt(userTransactionInfo.MonthlyTransactions[i]), 2)
	}

	fmt.Println("========== Final Summary ==========")
	fmt.Println("Total balance: ", userTransactionInfo.Balance)
	fmt.Println("Transactions per month:")
	for month, count := range userTransactionInfo.MonthlyTransactions {
		fmt.Printf("Month %02d: %d transactions\n", month, count)
	}
	fmt.Println("Average amount per month:")
	for month, avg := range userTransactionInfo.MonthlyAverages {
		fmt.Printf("Month %02d: ,", month)
		fmt.Println("average amount: ", avg)
	}
	fmt.Println("===================================")

	return userTransactionInfo, nil
}
