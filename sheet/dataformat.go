package sheet

import (
	"errors"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

func PrintValues(v *sheets.ValueRange) {
	for _, row := range v.Values {
		for i := range row {
			fmt.Print(row[i])
			if i != len(row)-1 {
				fmt.Print(",")
			}
		}
		fmt.Print("\n")
	}

}

// Implement an enum-a-like for the [input|output]-format flag
type DataFormatValue interface {
	String() string
	Set(string) error
	Type() string
}

type DataFormat string

const (
	CsvFormat DataFormat = "csv"
	TsvFormat DataFormat = "tsv"
)

func (f *DataFormat) String() string { return string(*f) }
func (f *DataFormat) Type() string   { return "DataFormat" }
func (f *DataFormat) Set(v string) error {
	switch v {
	case "csv", "tsv":
		*f = DataFormat(v)
		return nil
	default:
		return errors.New("invalid DataFormat. Allowed [csv|tsv]")
	}
}
