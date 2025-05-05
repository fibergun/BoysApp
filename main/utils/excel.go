package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func SaveSentences(sentences []string) error {
	file := excelize.NewFile()
	sheetName := "Sentences"
	file.NewSheet(sheetName)

	for i, sentence := range sentences {
		cell := fmt.Sprintf("A%d", i+1)
		file.SetCellValue(sheetName, cell, sentence)
	}

	if err := file.SaveAs("sentences.xlsx"); err != nil {
		return err
	}

	return nil
}
