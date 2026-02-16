// Package main demonstrates how to use the sheet library in your own Go programs.
//
// Before running this example, make sure you have:
// 1. OAuth credentials configured (see main README.md)
// 2. Set up viper configuration with clientsecretfile and authtokenfile
//
// This file is for documentation purposes and may not run as-is without proper setup.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gerrowadat/sheet/lib"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	// Example 1: Basic reading from a worksheet
	readExample()

	// Example 2: Writing data to a worksheet
	writeExample()

	// Example 3: Working with ranges
	rangeExample()

	// Example 4: Data format conversion (works without Google credentials)
	formatExample()

	// Example 5: Direct authentication without viper
	directAuthExample()

	// Example 6: Working with DataSpec
	dataSpecExample()
}

func readExample() {
	fmt.Println("=== Reading Example ===")

	// Get an authenticated Google Sheets service
	srv, err := sheet.GetService()
	if err != nil {
		log.Fatal(err)
	}

	// Create a DataSpec to reference a worksheet
	spec := &sheet.DataSpec{
		Workbook:  "your-spreadsheet-id-here",
		Worksheet: "Sheet1",
	}

	// Read data from the worksheet
	resp, err := srv.Spreadsheets.Values.Get(
		spec.Workbook,
		spec.GetInSheetDataSpec(),
	).Do()
	if err != nil {
		log.Fatal(err)
	}

	// Print the data as CSV
	fmt.Print(sheet.FormatValues(resp, sheet.CsvFormat))
}

func writeExample() {
	fmt.Println("\n=== Writing Example ===")

	srv, err := sheet.GetService()
	if err != nil {
		log.Fatal(err)
	}

	spec := &sheet.DataSpec{
		Workbook:  "your-spreadsheet-id-here",
		Worksheet: "Sheet1",
	}

	// Prepare data to write
	data := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "30", "NYC"},
		{"Bob", "25", "SF"},
		{"Carol", "35", "LA"},
	}

	// Write data to the worksheet
	// Parameters: service, spec, data, protect, force
	err = sheet.WriteDataToWorksheet(srv, spec, data, false, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data written successfully!")
}

func rangeExample() {
	fmt.Println("\n=== Range Example ===")

	srv, err := sheet.GetService()
	if err != nil {
		log.Fatal(err)
	}

	// Work with a specific range
	spec := &sheet.DataSpec{
		Workbook:  "your-spreadsheet-id-here",
		Worksheet: "Sheet1",
		Range:     sheet.RangeFromString("A1:C3"),
	}

	// Data to write (must fit in the specified range)
	data := [][]string{
		{"X", "Y", "Z"},
		{"1", "2", "3"},
		{"4", "5", "6"},
	}

	err = sheet.WriteDataToRange(srv, spec, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Range updated successfully!")
}

func formatExample() {
	fmt.Println("\n=== Format Conversion Example ===")

	// Parse CSV data
	csvData := "a,b,c\nd,e,f\ng,h,i"
	reader := bufio.NewReader(strings.NewReader(csvData))
	data, err := sheet.ScanValues(reader, sheet.CsvFormat)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsed CSV data:")
	for _, row := range data {
		fmt.Printf("  %v\n", row)
	}

	// Convert to ValueRange format for Google Sheets
	valueRange := &sheets.ValueRange{
		Values: make([][]interface{}, len(data)),
	}
	for i, row := range data {
		valueRange.Values[i] = make([]interface{}, len(row))
		for j, cell := range row {
			valueRange.Values[i][j] = cell
		}
	}

	// Format as TSV
	tsvOutput := sheet.FormatValues(valueRange, sheet.TsvFormat)
	fmt.Println("\nFormatted as TSV:")
	fmt.Print(tsvOutput)
}

func directAuthExample() {
	fmt.Println("\n=== Direct Auth Example ===")

	// Authenticate without viper by providing credential paths directly
	client := sheet.GetClient("/path/to/client_secret.json", "/path/to/token.json")

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatal(err)
	}

	spec := &sheet.DataSpec{
		Workbook:  "your-spreadsheet-id-here",
		Worksheet: "Sheet1",
	}

	resp, err := srv.Spreadsheets.Values.Get(spec.Workbook, spec.GetInSheetDataSpec()).Do()
	if err != nil {
		log.Fatal(err)
	}

	sheet.PrintValues(resp, sheet.CsvFormat)
}

func dataSpecExample() {
	fmt.Println("\n=== DataSpec Example ===")

	// A DataSpec can represent a workbook, worksheet, or range
	wbSpec := &sheet.DataSpec{Workbook: "my-spreadsheet-id"}
	fmt.Printf("IsWorkbook: %v\n", wbSpec.IsWorkbook()) // true

	wsSpec := &sheet.DataSpec{Workbook: "my-spreadsheet-id", Worksheet: "Sheet1"}
	fmt.Printf("IsWorksheet: %v\n", wsSpec.IsWorksheet()) // true

	rngSpec := &sheet.DataSpec{
		Workbook:  "my-spreadsheet-id",
		Worksheet: "Sheet1",
		Range:     sheet.RangeFromString("A1:C10"),
	}
	fmt.Printf("IsRange: %v\n", rngSpec.IsRange())                   // true
	fmt.Printf("InSheetRef: %v\n", rngSpec.GetInSheetDataSpec())     // Sheet1!A1:C10
	fmt.Printf("String: %v\n", rngSpec.String())                     // Workbook: my-spreadsheet-id, Worksheet: Sheet1, Range: A1:C10

	// DataRange utilities
	r := sheet.RangeFromString("A1:D10")
	cols, rows := r.SizeXY()
	fmt.Printf("Range %v: %d cols x %d rows, fixed=%v\n", r.String(), cols, rows, r.IsFixedSize())
}
