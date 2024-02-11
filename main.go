package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Transaction struct {
	TxID        string
	Fee         int
	Weight      int
	ParentTxIDs []string
	IsSelected  bool
	Children    []*Transaction
	TotalFee    int
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Skip the header if it exists
	if len(records) > 0 {
		records = records[1:]
	}

	return records, nil
}

func createTransactions(records [][]string) map[string]*Transaction {
	transactions := make(map[string]*Transaction)

	for _, record := range records {
		txID := record[0]
		fee, _ := strconv.Atoi(record[1])
		weight, _ := strconv.Atoi(record[2])
		parentTxIDs := strings.Split(record[3], ";")

		transaction := &Transaction{
			TxID:        txID,
			Fee:         fee,
			Weight:      weight,
			ParentTxIDs: parentTxIDs,
			TotalFee:    fee,
			Children:    []*Transaction{},
		}

		transactions[txID] = transaction
	}

	return transactions
}

func sortAndPrintTransactions(transactions map[string]*Transaction) {
	var sortedTransactions []*Transaction
	for _, transaction := range transactions {
		if !transaction.IsSelected {
			sortTransactions(transaction, transactions, &sortedTransactions)
		}
	}

	var selectedTransactions []string
	for _, transaction := range sortedTransactions {
		selectTransactions(transaction, transactions, &selectedTransactions)
	}

	// Sort selected transactions based on TotalFee
	sort.Slice(selectedTransactions, func(i, j int) bool {
		return transactions[selectedTransactions[i]].TotalFee > transactions[selectedTransactions[j]].TotalFee
	})

	// Print selected transactions
	for _, txID := range selectedTransactions {
		fmt.Println(txID)
	}
}

func sortTransactions(transaction *Transaction, transactions map[string]*Transaction, sortedTransactions *[]*Transaction) {
	for _, child := range transaction.Children {
		sortTransactions(child, transactions, sortedTransactions)
	}

	*sortedTransactions = append(*sortedTransactions, transaction)
}

func selectTransactions(transaction *Transaction, transactions map[string]*Transaction, selectedTransactions *[]string) {
	for _, child := range transaction.Children {
		selectTransactions(child, transactions, selectedTransactions)
	}

	for _, parentID := range transaction.ParentTxIDs {
		parent, exists := transactions[parentID]
		if exists {
			transaction.TotalFee += parent.TotalFee
		}
	}

	transaction.TotalFee += transaction.Fee

	if transaction.TotalFee > 0 {
		*selectedTransactions = append(*selectedTransactions, transaction.TxID)
		transaction.IsSelected = true
	}
}

func main() {
	records, err := readCSV("mempool.csv")

	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	transactions := createTransactions(records)

	for _, transaction := range transactions {
		for _, parentID := range transaction.ParentTxIDs {
			parent := transactions[parentID]
			if parent != nil {
				parent.Children = append(parent.Children, transaction)
			}
		}
	}

	sortAndPrintTransactions(transactions)
}
