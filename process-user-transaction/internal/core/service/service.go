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

	monthlyDebitTransactions := make(map[int]map[string]decimal.Decimal)
	monthlyCreditTransactions := make(map[int]map[string]decimal.Decimal)

	debitTrxInfoChan := make(chan TransactionInfo)
	creditTrxInfoChan := make(chan TransactionInfo)
	errorChan := make(chan error)
	quitChan := make(chan int)

	go func() {
		for _, transaction := range transactions {
			err := s.r.PutTransaction(context.Background(), transaction)
			if err != nil {
				return
			}
			extractTransactionInfo(transaction, debitTrxInfoChan, creditTrxInfoChan, errorChan)
			fmt.Printf("Processed transaction: %+v\n", transaction)
		}
		close(debitTrxInfoChan)
		close(creditTrxInfoChan)
		quitChan <- 0
	}()

	var done bool

	for !done {
		select {
		case value, ok := <-debitTrxInfoChan:
			if ok {
				if _, exists := monthlyDebitTransactions[value.Month]; !exists {
					monthlyDebitTransactions[value.Month] = make(map[string]decimal.Decimal)
				}
				monthlyDebitTransactions[value.Month][value.TransactionId] = value.Amount
			}
		case value, ok := <-creditTrxInfoChan:
			if ok {
				if _, exists := monthlyCreditTransactions[value.Month]; !exists {
					monthlyCreditTransactions[value.Month] = make(map[string]decimal.Decimal)
				}
				monthlyCreditTransactions[value.Month][value.TransactionId] = value.Amount
			}
		case err := <-errorChan:
			return domain.UserTransactionInfo{}, err
		case <-quitChan:
			done = true
		}
	}

	var balance decimal.Decimal

	for i, monthlyAverage := range monthlyDebitTransactions {
		var monthlyDebit decimal.Decimal
		for _, amount := range monthlyAverage {
			monthlyDebit = monthlyDebit.Add(amount)
			userTransactionInfo.MonthlyTransactions[i] += 1
		}
		userTransactionInfo.MonthlyDebitAverages[i] = monthlyDebit.DivRound(decimal.NewFromInt(int64(len(monthlyAverage))), 2)
		balance = balance.Add(monthlyDebit)
	}

	for i, monthlyAverage := range monthlyCreditTransactions {
		var monthlyCredit decimal.Decimal
		for _, amount := range monthlyAverage {
			monthlyCredit = monthlyCredit.Add(amount)
			userTransactionInfo.MonthlyTransactions[i] += 1
		}
		userTransactionInfo.MonthlyCreditAverages[i] = monthlyCredit.DivRound(decimal.NewFromInt(int64(len(monthlyAverage))), 2)
		balance = balance.Add(monthlyCredit)
	}

	userTransactionInfo.Balance = balance

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

type TransactionInfo struct {
	TransactionId string
	Amount        decimal.Decimal
	Month         int
}

func extractTransactionInfo(transaction domain.Transaction, debitChannel chan TransactionInfo, creditChannel chan TransactionInfo, errChan chan error) {
	trxMonth := strings.Split(transaction.CreatedDate, "/")[0]
	month, err := strconv.Atoi(trxMonth)
	if err != nil {
		fmt.Printf("error retrieving transaction month for transaction with transactionId %s and err %s", transaction.TransactionId, err)
		errChan <- err
	}

	transactionInfo := TransactionInfo{
		TransactionId: transaction.TransactionId,
		Amount:        transaction.Amount,
		Month:         month,
	}

	if transaction.Amount.GreaterThan(decimal.NewFromInt(0)) {
		creditChannel <- transactionInfo
	} else {
		debitChannel <- transactionInfo
	}
}
