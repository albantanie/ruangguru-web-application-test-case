package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Report struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Date     string  `json:"date"`
	Semester int     `json:"semester"`
	Studies  []Study `json:"studies"`
	// TODO: answer here
}

type Study struct {
	Name   string `json:"study_name"`
	Credit int    `json:"study_credit"`
	Grade  string `json:"grade"`
}

// gunakan fungsi ini untuk mengambil data dari file json
// kembalian berupa struct 'Report' dan error
func ReadJSON(filename string) (Report, error) {
	byteReport, err := ioutil.ReadFile(filename)
	if err != nil {
		return Report{}, err
	}

	report := Report{}
	err = json.Unmarshal(byteReport, &report)
	if err != nil {
		return Report{}, err
	}
	return report, nil

	// TODO: answer here
}

func GradePoint(report Report) float64 {
	if report.Studies == nil || len(report.Studies) == 0 {
		return 0.0

	}
	accumIP := 0.0
	accumCredit := 0.0
	for _, study := range report.Studies {
		gradeMap := map[string]float64{
			"A":  4.0,
			"AB": 3.5,
			"B":  3.0,
			"BC": 2.5,
			"C":  2.0,
			"CD": 1.5,
			"D":  1.0,
			"DE": 0.5,
			"E":  0.0,
		}
		accumIP += float64(study.Credit) * gradeMap[study.Grade]
		accumCredit += float64(study.Credit)
		// TODO: answer here
	}

	return accumIP / accumCredit // TODO: replace this
}

func main() {
	// bisa digunakan untuk menguji test case
	report, err := ReadJSON("report.json")
	if err != nil {
		panic(err)
	}

	gradePoint := GradePoint(report)
	fmt.Println(gradePoint)
}
