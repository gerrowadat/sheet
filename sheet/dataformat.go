package sheet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

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

func ScanValues(r *bufio.Reader, f DataFormat) ([][]string, error) {
	ret := [][]string{}

	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	// Split the buffer into lines. We don't hold truck with any multi-line data format fuckery, for now.
	lines := strings.Split(buf.String(), "\n")

	for l := range lines {
		if lines[l] != "" {
			ret = append(ret, strings.Split(lines[l], f.Separator()))
		}
	}

	return ret, nil
}

// Implement an enum-a-like for the [input|output]-format flag
type DataFormatValue interface {
	String() string
	Set(string) error
	Type() string
	Separator() string
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
func (f *DataFormat) Separator() string {
	switch *f {
	case "csv":
		return ","
	case "tsv":
		return "\t"
	default:
		return ","
	}
}
