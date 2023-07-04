package pkg

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

type PersonData struct {
	FirstName    string
	LastName     string
	Organization string
	Email        string
	Domain       string
	Title        string
}

// SaveToExcel saves the collected person data into an Excel file
func SaveToExcel(personData []PersonData, organization string) error {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Employee Info")

	for _, data := range personData {
		row := sheet.AddRow()
		row.AddCell().Value = data.FirstName
		row.AddCell().Value = data.LastName
		row.AddCell().Value = data.Organization
		row.AddCell().Value = data.Email
		row.AddCell().Value = data.Domain
		row.AddCell().Value = data.Title
	}

	fileName := fmt.Sprintf("apollonator_%s.xlsx", organization)
	err := file.Save(fileName)
	if err != nil {
		return err
	}

	return nil
}
