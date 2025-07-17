package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	const filename = "transactions3000.csv"
	const numTransactions = 3000

	file, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("failed to create file: %v", err))
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"TransactionId", "Date", "Amount", "AccountId"})

	accountIds := []string{"axz"}
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numTransactions; i++ {
		id := strconv.Itoa(i)

		// Random month and day
		month := rand.Intn(12) + 1
		day := rand.Intn(28) + 1
		date := fmt.Sprintf("%d/%d", month, day)

		// Random amount between -100 and +100
		amount := fmt.Sprintf("%.2f", (rand.Float64()*200 - 100))

		accountId := accountIds[rand.Intn(len(accountIds))]

		writer.Write([]string{id, date, amount, accountId})
	}

	fmt.Printf("File '%s' with %d transactions created successfully\n", filename, numTransactions)
}
