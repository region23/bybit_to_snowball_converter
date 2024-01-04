package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/tealeg/xlsx/v3"
)

var writer *csv.Writer

func main() {
	// read filename from params
	if len(os.Args) < 2 {
		fmt.Println("Please specify filename")
		return
	}
	filename := os.Args[1]

	// Открыть XLS-файл
	xlsFile, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Remove the extension from the filename
	filenameWithoutExt := filename[0 : len(filename)-len(filepath.Ext(filename))]

	// Add the new extension
	newFilename := "snowball_" + filenameWithoutExt + ".csv"

	// Создать CSV-файл
	csvFile, err := os.Create(newFilename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer csvFile.Close()

	writer = csv.NewWriter(csvFile)
	defer writer.Flush()

	// Здесь вы должны определить структуру данных в соответствии с вашим CSV-шаблоном
	headers := []string{
		"Event",
		"Date",
		"Symbol",
		"Price",
		"Quantity",
		"Currency",
		"FeeTax",
		"Exchange",
		"NKD",
		"FeeCurrency",
		"DoNotAdjustCash",
		"Note",
	}

	// Запись заголовков
	if err := writer.Write(headers); err != nil {
		fmt.Println("Error writing header to csv:", err)
		return
	}

	// Обработка данных из XLS-файла и запись в CSV
	for _, sheet := range xlsFile.Sheets {
		fmt.Println("Max row is", sheet.MaxRow)
		err = sheet.ForEachRow(rowProcessing)
		fmt.Println("Err=", err)
	}
}

func rowProcessing(r *xlsx.Row) error {
	row := r.GetCoordinate()
	if row == 0 {
		return nil
	}

	event, err := r.GetCell(3).FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	}
	event = cases.Title(language.English, cases.Compact).String(event)

	date, err := r.GetCell(0).FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	}

	parsedDate, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		fmt.Println(err.Error())
	}

	formattedDate := parsedDate.Format("2006-01-02")

	symbol, err := r.GetCell(1).FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	}
	symbol = strings.ReplaceAll(symbol, " ", "")

	price, err := r.GetCell(5).FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	}
	clearedPrice := strings.Replace(price, "USDT", "", -1)
	clearedPrice = strings.ReplaceAll(clearedPrice, " ", "")

	quantity, err := r.GetCell(6).FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	}
	quantity = strings.ReplaceAll(quantity, " ", "")

	clearedQuantity := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return -1
		}
		return r
	}, quantity)

	re := regexp.MustCompile("[^a-zA-Z]+")
	currency := re.ReplaceAllString(quantity, "")

	fee, err := r.GetCell(7).FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	}
	fee = strings.ReplaceAll(fee, " ", "")

	clearedFee := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return -1
		}
		return r
	}, fee)

	record := []string{event, formattedDate, symbol, clearedPrice, clearedQuantity, currency, clearedFee, "ByBit", "", "", "", ""}
	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}
