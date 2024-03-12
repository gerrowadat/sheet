package sheet

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/api/sheets/v4"
)

func TestDataFormat_Set(t *testing.T) {
	tests := []struct {
		name    string
		f       *DataFormat
		ftype   string
		wantErr bool
	}{
		{
			name:    "Simple",
			ftype:   "csv",
			wantErr: false,
		},
		{
			name:    "UnknownFType",
			ftype:   "blah",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f DataFormat = "csv"
			if err := f.Set(tt.ftype); (err != nil) != tt.wantErr {
				t.Errorf("DataFormat.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataFormat_Separator(t *testing.T) {
	tests := []struct {
		name string
		f    DataFormat
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Csv",
			f:    CsvFormat,
			want: ",",
		},
		{
			name: "Tsv",
			f:    TsvFormat,
			want: "\t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Separator(); got != tt.want {
				t.Errorf("DataFormat.Separator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatValues(t *testing.T) {
	type args struct {
		v *sheets.ValueRange
		f DataFormat
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "SimpleCsv",
			args: args{v: &sheets.ValueRange{Values: [][]interface{}{{"a", "b"}, {"c", "d"}}}, f: CsvFormat},
			want: "a,b\nc,d\n",
		},
		{
			name: "SimpleTsv",
			args: args{v: &sheets.ValueRange{Values: [][]interface{}{{"a", "b"}, {"c", "d"}}}, f: TsvFormat},
			want: "a\tb\nc\td\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatValues(tt.args.v, tt.args.f); got != tt.want {
				t.Errorf("FormatValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanValues(t *testing.T) {
	type args struct {
		r *bufio.Reader
		f DataFormat
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name:    "SimpleCsv",
			args:    args{r: bufio.NewReader(strings.NewReader("a,b\nc,d\n")), f: CsvFormat},
			want:    [][]string{{"a", "b"}, {"c", "d"}},
			wantErr: false,
		},
		{
			name:    "SimpleTsv",
			args:    args{r: bufio.NewReader(strings.NewReader("a\tb\nc\td\n")), f: TsvFormat},
			want:    [][]string{{"a", "b"}, {"c", "d"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanValues(tt.args.r, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
