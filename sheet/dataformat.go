package sheet

import (
	"errors"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

func PrintValues(v *sheets.ValueRange, f DataFormat) {
	fmt.Print(FormatValues(v, f))
}

func FormatValues(v *sheets.ValueRange, f DataFormat) string {
	ret := ""

	sep := ","
	if f == "tsv" {
		sep = "\t"
	}

	for _, row := range v.Values {
		for i := range row {
			ret += row[i].(string)
			if i != len(row)-1 {
				ret += sep
			}
		}
		ret += "\n"
	}
	return ret
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
