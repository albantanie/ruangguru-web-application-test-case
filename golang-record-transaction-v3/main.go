package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Transaction struct {
	Date   string
	Type   string
	Amount int
}

func RecordTransactions(path string, transactions []Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date < transactions[j].Date
	})

	lastDate := transactions[0].Date
	money := 0
	output := make([]string, 0)

	for _, t := range transactions {
		if t.Date == lastDate {
			if t.Type == "income" {
				money += t.Amount
			} else {
				money -= t.Amount
			}
		} else {
			record := fmt.Sprintf("%s;%s;%d", lastDate, getType(money), abs(money))
			output = append(output, record)
			money = 0
			if t.Type == "income" {
				money += t.Amount
			} else {
				money -= t.Amount
			}
			lastDate = t.Date
		}
	}

	record := fmt.Sprintf("%s;%s;%d", lastDate, getType(money), abs(money))
	output = append(output, record)

	// Write the output to file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(output, "\n"))
	if err != nil {
		return err
	}

	return nil
}

func getType(money int) string {
	if money < 0 {
		return "expense"
	}
	return "income"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	transactions := []Transaction{
		{"01/01/2021", "income", 100000},
		{"01/01/2021", "expense", 50000},
		{"01/01/2021", "expense", 30000},
		{"01/01/2021", "income", 20000},
	}

	fmt.Println("Before", transactions)

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date < transactions[j].Date
	})

	fmt.Println("After", transactions)

	err := RecordTransactions("transactions.txt", transactions)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success")
}
