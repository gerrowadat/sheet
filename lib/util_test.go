package sheet

import (
	"testing"

	"google.golang.org/api/sheets/v4"
)

func TestClearWorksheet(t *testing.T) {
	type args struct {
		srv     *sheets.Service
		spec    *DataSpec
		protect bool
		force   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestProtection",
			args:    args{srv: &sheets.Service{}, spec: &DataSpec{Worksheet: "doot"}, protect: true, force: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ClearWorksheet(tt.args.srv, tt.args.spec, tt.args.protect, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("ClearWorksheet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_valueRangeFromStrings(t *testing.T) {
	tests := []struct {
		name string
		data [][]string
		want int // expected number of rows
	}{
		{
			name: "SimpleData",
			data: [][]string{{"a", "b"}, {"c", "d"}},
			want: 2,
		},
		{
			name: "SingleRow",
			data: [][]string{{"x", "y", "z"}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := valueRangeFromStrings(tt.data)
			if len(got.Values) != tt.want {
				t.Errorf("valueRangeFromStrings() rows = %v, want %v", len(got.Values), tt.want)
			}
			// Verify the data round-trips correctly
			for i, row := range tt.data {
				for j, cell := range row {
					if got.Values[i][j].(string) != cell {
						t.Errorf("valueRangeFromStrings()[%d][%d] = %v, want %v", i, j, got.Values[i][j], cell)
					}
				}
			}
		})
	}
}

func TestDataRange_IsFixedSize(t *testing.T) {
	tests := []struct {
		name string
		rng  DataRange
		want bool
	}{
		{
			name: "FixedRange",
			rng:  RangeFromString("A1:B2"),
			want: true,
		},
		{
			name: "WholeColumns",
			rng:  RangeFromString("A:B"),
			want: false,
		},
		{
			name: "SingleCell",
			rng:  RangeFromString("A1:A1"),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rng.IsFixedSize(); got != tt.want {
				t.Errorf("DataRange.IsFixedSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkDataFitsInRange(t *testing.T) {
	type args struct {
		spec *DataSpec
		data [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestDataOverflowCols",
			args:    args{spec: &DataSpec{Range: RangeFromString("A1:B2")}, data: [][]string{{"1", "2", "3"}, {"4", "5", "6"}}},
			wantErr: true,
		},
		{
			name:    "TestDataOverflowRows",
			args:    args{spec: &DataSpec{Range: RangeFromString("A1:D1")}, data: [][]string{{"1", "2", "3"}, {"4", "5", "6"}}},
			wantErr: true,
		},
		{
			name:    "TestExactFit",
			args:    args{spec: &DataSpec{Range: RangeFromString("A1:C2")}, data: [][]string{{"1", "2", "3"}, {"4", "5", "6"}}},
			wantErr: false,
		},
		{
			name:    "TestFewerCols",
			args:    args{spec: &DataSpec{Range: RangeFromString("A1:E2")}, data: [][]string{{"1", "2", "3"}, {"4", "5", "6"}}},
			wantErr: false,
		},
		{
			name:    "TestFewerRows",
			args:    args{spec: &DataSpec{Range: RangeFromString("A1:C3")}, data: [][]string{{"1", "2", "3"}, {"4", "5", "6"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkDataFitsInRange(tt.args.spec, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("checkDataFitsInRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
