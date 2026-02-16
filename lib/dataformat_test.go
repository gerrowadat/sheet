package sheet

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/api/sheets/v4"
)

func TestDataFormat_String(t *testing.T) {
	tests := []struct {
		name string
		f    DataFormat
		want string
	}{
		{name: "Csv", f: CsvFormat, want: "csv"},
		{name: "Tsv", f: TsvFormat, want: "tsv"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("DataFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataFormat_Type(t *testing.T) {
	var f DataFormat = CsvFormat
	if got := f.Type(); got != "DataFormat" {
		t.Errorf("DataFormat.Type() = %v, want DataFormat", got)
	}
}

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
			name:    "Tsv",
			ftype:   "tsv",
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
		{
			name: "SingleCell",
			args: args{v: &sheets.ValueRange{Values: [][]interface{}{{"hello"}}}, f: CsvFormat},
			want: "hello\n",
		},
		{
			name: "EmptyValues",
			args: args{v: &sheets.ValueRange{Values: [][]interface{}{}}, f: CsvFormat},
			want: "",
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
		{
			name:    "EmptyInput",
			args:    args{r: bufio.NewReader(strings.NewReader("")), f: CsvFormat},
			want:    [][]string{},
			wantErr: false,
		},
		{
			name:    "SingleValue",
			args:    args{r: bufio.NewReader(strings.NewReader("hello\n")), f: CsvFormat},
			want:    [][]string{{"hello"}},
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
