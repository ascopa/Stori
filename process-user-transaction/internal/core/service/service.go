package service

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"process-user-transaction/internal/adapters/outbound/repository"
	"process-user-transaction/internal/core/domain"
	"strconv"
	"strings"
	"sync"
)

const (
	NUM_WORKERS  = "NUM_WORKERS"
	ENABLE_DEBUG = "ENABLE_DEBUG"
)

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

	numWorkers, err := strconv.Atoi(os.Getenv(NUM_WORKERS))
	if err != nil {
		return domain.UserTransactionInfo{}, fmt.Errorf("failed to convert NUM_WORKERS to int: %w", err)
	}

	resultChan := make(chan PartialResult)
	errChan := make(chan error, numWorkers)

	var wgDB sync.WaitGroup
	var wgProcessTrx sync.WaitGroup
	wgDB.Add(numWorkers)
	wgProcessTrx.Add(numWorkers)

	chunkSize := len(transactions) / numWorkers
	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(transactions)
		}

		chunk := transactions[start:end]

		go func(chunk []domain.Transaction) {
			defer wgProcessTrx.Done()
			processTransactionsChunk(chunk, resultChan, errChan)
		}(chunk)
		go func(chunk []domain.Transaction) {
			defer wgDB.Done()
			s.putTransactions(chunk, errChan)
		}(chunk)
	}

	go func() {
		wgProcessTrx.Wait()
		close(resultChan)
	}()

	go func() {
		wgDB.Wait()
		close(errChan)
	}()

	final := PartialResult{
		SumMonthlyCredit: make(map[int]decimal.Decimal),
		SumMonthlyDebit:  make(map[int]decimal.Decimal),
		CountCredit:      make(map[int]int),
		CountDebit:       make(map[int]int),
	}

	for partial := range resultChan {
		for month, sum := range partial.SumMonthlyCredit {
			final.SumMonthlyCredit[month] = final.SumMonthlyCredit[month].Add(sum)
			final.CountCredit[month] += partial.CountCredit[month]
		}
		for month, sum := range partial.SumMonthlyDebit {
			final.SumMonthlyDebit[month] = final.SumMonthlyDebit[month].Add(sum)
			final.CountDebit[month] += partial.CountDebit[month]
		}
	}

	for err = range errChan {
		fmt.Printf("error processing trx %s", err)
	}

	for month, total := range final.SumMonthlyCredit {
		avg := total.DivRound(decimal.NewFromInt(int64(final.CountCredit[month])), 2)
		userTransactionInfo.MonthlyCreditAverages[month] = avg
		userTransactionInfo.Balance = userTransactionInfo.Balance.Add(total)
		userTransactionInfo.MonthlyTransactions[month] += final.CountCredit[month]
	}

	for month, total := range final.SumMonthlyDebit {
		avg := total.DivRound(decimal.NewFromInt(int64(final.CountDebit[month])), 2)
		userTransactionInfo.MonthlyDebitAverages[month] = avg
		userTransactionInfo.Balance = userTransactionInfo.Balance.Add(total)
		userTransactionInfo.MonthlyTransactions[month] += final.CountDebit[month]
	}

	if os.Getenv(ENABLE_DEBUG) == "true" {
		log(userTransactionInfo)
	}

	return userTransactionInfo, nil
}

type PartialResult struct {
	SumMonthlyCredit map[int]decimal.Decimal
	SumMonthlyDebit  map[int]decimal.Decimal
	CountCredit      map[int]int
	CountDebit       map[int]int
}

func (s *Service) putTransactions(transactions []domain.Transaction, errChan chan error) {
	for _, trx := range transactions {
		err := s.r.PutTransaction(context.Background(), trx)
		if err != nil {
			err = fmt.Errorf("error saving transaction on dynamoDB for trx with transactionId %s and err %s", trx.TransactionId, err)
			errChan <- err
			continue
		}
	}
}

func processTransactionsChunk(transactions []domain.Transaction, partialResultChan chan PartialResult, errChan chan error) {
	result := PartialResult{
		SumMonthlyCredit: make(map[int]decimal.Decimal),
		SumMonthlyDebit:  make(map[int]decimal.Decimal),
		CountCredit:      make(map[int]int),
		CountDebit:       make(map[int]int),
	}

	for _, trx := range transactions {
		month, err := extractMonth(trx)
		if err != nil {
			err = fmt.Errorf("error extracting month for trx with transactionId %s and err %s", trx.TransactionId, err)
			errChan <- err
			continue
		}

		if trx.Amount.GreaterThan(decimal.Zero) {
			result.SumMonthlyCredit[month] = result.SumMonthlyCredit[month].Add(trx.Amount)
			result.CountCredit[month]++
		} else {
			result.SumMonthlyDebit[month] = result.SumMonthlyDebit[month].Add(trx.Amount)
			result.CountDebit[month]++
		}
	}

	partialResultChan <- result
}

func extractMonth(transaction domain.Transaction) (int, error) {
	trxMonth := strings.Split(transaction.CreatedDate, "/")[0]
	month, err := strconv.Atoi(trxMonth)
	if err != nil {
		return 0, fmt.Errorf("error retrieving transaction month for transaction with transactionId %s and err %w", transaction.TransactionId, err)
	}
	return month, nil
}

func log(userTransactionInfo domain.UserTransactionInfo) {
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
}
