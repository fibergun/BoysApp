package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func SaveSentences(sentences []string) error {
	file := excelize.NewFile()
	sheetName := "Sheet1" // Default sheet name in Excelize
	file.SetSheetName("Sheet1", sheetName)

	for i, sentence := range sentences {
		cell := fmt.Sprintf("A%d", i+1) // Write sentences in column A
		if err := file.SetCellValue(sheetName, cell, sentence); err != nil {
			return err
		}
	}
	err := file.SaveAs("sentences.xlsx")
	// Save the file to the current directory

	if err != nil {
		fmt.Println("File saved as sentences.xlsx")
		return err
	}

	return nil
}

func LoadSentences() ([]string, error) {
	file, err := excelize.OpenFile("sentences.xlsx")
	if err != nil {
		return nil, err
	}

	sheetName := "Sheet1"
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var sentences []string
	for _, row := range rows {
		if len(row) > 0 {
			sentences = append(sentences, row[0]) // Read sentences from column A
		}
	}

	return sentences, nil
}

func SaveSentencesWithPlayers(sentences []string, players []string) error {
	file := excelize.NewFile()
	sheetName := "Sheet1"
	file.SetSheetName("Sheet1", sheetName)

	for i, sentence := range sentences {
		cellSentence := fmt.Sprintf("A%d", i+1)
		file.SetCellValue(sheetName, cellSentence, sentence)

		if players != nil && i < len(players) {
			cellPlayer := fmt.Sprintf("B%d", i+1)
			file.SetCellValue(sheetName, cellPlayer, players[i])
		} else {
			cellPlayer := fmt.Sprintf("B%d", i+1)
			file.SetCellValue(sheetName, cellPlayer, "") // Clear the player name
		}
	}

	// Debug: Print the sentences and players being saved
	fmt.Println("Saving sentences:", sentences)
	fmt.Println("Saving players:", players)

	if err := file.SaveAs("sentences.xlsx"); err != nil {
		fmt.Println("Error saving file:", err)
		return err
	}

	fmt.Println("File saved successfully.")
	return nil
}
