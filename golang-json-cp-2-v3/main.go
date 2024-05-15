package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type LoanData struct {
	StartBalance int
	Data         []Loan
	Employees    []Employee
}

type Loan struct {
	Date        string
	EmployeeIDs []string
}

type Employee struct {
	ID       string
	Name     string
	Position string
}

// json structure
type LoanRecord struct {
	MonthDate    string     `json:"month_date"`
	StartBalance int        `json:"start_balance"`
	EndBalance   int        `json:"end_balance"`
	Borrowers    []Borrower `json:"borrowers"`
}

type Borrower struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Total int    `json:"total_loan"`
}

func FindEmployee(id string, employees []Employee) (Employee, bool) {
	for _, employee := range employees {
		if employee.ID == id {
			return employee, true
		}
	}

	return Employee{}, false
}

func GetEndBalanceAndBorrowers(data LoanData) (int, []Borrower) {
	if len(data.Data) != 0 {

		// recap the data in map
		var tempMap = map[string]int{}

		for _, v := range data.Data {
			for _, id := range v.EmployeeIDs {
				if data.StartBalance >= 50000 {
					tempMap[id] += 50000
					data.StartBalance -= 50000
				}
			}
		}

		// sort temp
		var keys []string

		for k := range tempMap {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		var borrowers []Borrower

		// set data to struct
		for _, k := range keys {
			if employee, ok := FindEmployee(k, data.Employees); ok {
				borrowers = append(borrowers, Borrower{
					ID:    employee.ID,
					Name:  employee.Name,
					Total: tempMap[k],
				})
			}
		}

		// set the start balance and end balance
		return data.StartBalance, borrowers
	} else {
		return 0, []Borrower{}
	}
}

func LoanReport(data LoanData) LoanRecord {
	if len(data.Data) != 0 {
		var loanRecord = LoanRecord{
			StartBalance: data.StartBalance,
		}

		endBalance, borrowers := GetEndBalanceAndBorrowers(data)

		date := strings.Split(data.Data[0].Date, "-")
		loanRecord.MonthDate = fmt.Sprintf("%s %s", date[1], date[2])

		loanRecord.EndBalance = endBalance
		loanRecord.Borrowers = borrowers

		return loanRecord
	} else {
		return LoanRecord{}
	}
}

func RecordJSON(record LoanRecord, path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, _ = os.Create(path)
		defer file.Close()
		return nil
	}

	byteJSON, _ := json.Marshal(record)
	_ = ioutil.WriteFile(path, byteJSON, 0644)

	return nil // TODO: replace this
}

// gunakan untuk debug
func main() {
	records := LoanReport(LoanData{
		StartBalance: 500000,
		Data: []Loan{
			{"01-January-2021", []string{"1", "2"}},
			{"02-January-2021", []string{"1", "2", "3"}},
			{"03-January-2021", []string{"2", "3"}},
			{"04-January-2021", []string{"1", "3"}},
		},
		Employees: []Employee{
			{"1", "Employee A", "Manager"},
			{"2", "Employee B", "Staff"},
			{"3", "Employee C", "Staff"},
		},
	})

	err := RecordJSON(records, "loan-records.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(records)
}
