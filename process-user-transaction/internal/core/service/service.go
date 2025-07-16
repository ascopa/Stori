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

	userTransactionInfo.AccountId = transactions[0].AccountId

	userTransactionInfo.MonthlyDebitAverages = make(map[int]decimal.Decimal)
	userTransactionInfo.MonthlyCreditAverages = make(map[int]decimal.Decimal)
	userTransactionInfo.MonthlyTransactions = make(map[int]int)

	monthlyCreditTransactions := make(map[int]int)
	monthlyDebitTransactions := make(map[int]int)

	for _, transaction := range transactions {
		err := s.r.PutTransaction(context.Background(), transaction)
		if err != nil {
			return domain.UserTransactionInfo{}, fmt.Errorf("error saving transaction with transactionId %s and err %w", transaction.TransactionId, err)
		}
		trxMonth := strings.Split(transaction.CreatedDate, "/")[0]
		month, err := strconv.Atoi(trxMonth)
		if err != nil {
			return domain.UserTransactionInfo{}, fmt.Errorf("error retrieving transaction month for transaction with transactionId %s and err %w", transaction.TransactionId, err)
		}

		userTransactionInfo.MonthlyTransactions[month] += 1
		if transaction.Amount.GreaterThan(decimal.NewFromInt(0)) {
			userTransactionInfo.MonthlyCreditAverages[month] = userTransactionInfo.MonthlyCreditAverages[month].Add(transaction.Amount)
			monthlyCreditTransactions[month] += 1
		} else {
			userTransactionInfo.MonthlyDebitAverages[month] = userTransactionInfo.MonthlyDebitAverages[month].Add(transaction.Amount)
			monthlyDebitTransactions[month] += 1
		}

		userTransactionInfo.Balance = userTransactionInfo.Balance.Add(transaction.Amount)

		fmt.Printf("Processed transaction: %+v\n", transaction)
		fmt.Printf("Month: %d, Transaction count: %d, ", month, userTransactionInfo.MonthlyTransactions[month])
		fmt.Println("Average credit amount: ", userTransactionInfo.MonthlyCreditAverages[month])
		fmt.Println("Average debit amount: ", userTransactionInfo.MonthlyDebitAverages[month])
		fmt.Println("Running balance: ", userTransactionInfo.Balance)
		fmt.Println("------")
	}

	for i, monthlyAverage := range userTransactionInfo.MonthlyDebitAverages {
		userTransactionInfo.MonthlyDebitAverages[i] = monthlyAverage.DivRound(decimal.NewFromInt(int64(monthlyDebitTransactions[i])), 2)
	}

	for i, monthlyAverage := range userTransactionInfo.MonthlyCreditAverages {
		userTransactionInfo.MonthlyCreditAverages[i] = monthlyAverage.DivRound(decimal.NewFromInt(int64(monthlyCreditTransactions[i])), 2)
	}

	fmt.Println("========== Final Summary ==========")
	fmt.Println("Total balance: ", userTransactionInfo.Balance)
	fmt.Println("Transactions per month:")
	for month, count := range userTransactionInfo.MonthlyTransactions {
		fmt.Printf("Month %02d: %d transactions\n", month, count)
	}
	fmt.Println("Average credit amount per month:")
	for month, avg := range userTransactionInfo.MonthlyCreditAverages {
		fmt.Printf("Month %02d: ,", month)
		fmt.Println("average amount: ", avg)
	}
	fmt.Println("Average debit amount per month:")
	for month, avg := range userTransactionInfo.MonthlyDebitAverages {
		fmt.Printf("Month %02d: ,", month)
		fmt.Println("average amount: ", avg)
	}
	fmt.Println("===================================")

	return userTransactionInfo, nil
}
