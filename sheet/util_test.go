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
		// TODO: Add test cases.
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
