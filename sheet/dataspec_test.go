package sheet

import (
	"reflect"
	"testing"
)

func TestDataSpec_GetInSheetDataSpec(t *testing.T) {
	type fields struct {
		Workbook  string
		Worksheet string
		Range     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "BareWorksheet",
			fields: fields{Worksheet: "mysheet"},
			want:   "mysheet",
		},
		{
			name:   "BareRange",
			fields: fields{Range: "A1:B10"},
			want:   "A1:B10",
		},
		{
			name:   "Combined",
			fields: fields{Worksheet: "mysheet", Range: "A1:B10"},
			want:   "mysheet!A1:B10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataSpec{
				Workbook:  tt.fields.Workbook,
				Worksheet: tt.fields.Worksheet,
				Range:     tt.fields.Range,
			}
			if got := d.GetInSheetDataSpec(); got != tt.want {
				t.Errorf("DataSpec.GetInSheetDataSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataSpec_FromString(t *testing.T) {
	type fields struct {
		Workbook  string
		Worksheet string
		Range     string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *DataSpec
	}{
		{
			name: "Blank",
			args: args{s: ""},
			want: &DataSpec{},
		},
		{
			name: "JustWorksheet",
			args: args{s: "mysheet"},
			want: &DataSpec{Worksheet: "mysheet"},
		},
		{
			name: "WithRange",
			args: args{s: "mysheet!A1:B100"},
			want: &DataSpec{Worksheet: "mysheet", Range: "A1:B100"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataSpec{
				Workbook:  tt.fields.Workbook,
				Worksheet: tt.fields.Worksheet,
				Range:     tt.fields.Range,
			}
			if got := d.FromString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataSpec.FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeDataSpecs(t *testing.T) {
	type args struct {
		specs []*DataSpec
	}
	tests := []struct {
		name    string
		args    args
		want    *DataSpec
		wantErr bool
	}{
		{
			name:    "AllBlank",
			args:    args{specs: []*DataSpec{{}, {}}},
			want:    &DataSpec{},
			wantErr: false,
		},
		{
			name:    "FullNoClashes",
			args:    args{specs: []*DataSpec{{Workbook: "mybook"}, {Worksheet: "mysheet"}, {Range: "myrange"}}},
			want:    &DataSpec{Workbook: "mybook", Worksheet: "mysheet", Range: "myrange"},
			wantErr: false,
		},
		{
			name:    "PartialNoClashes",
			args:    args{specs: []*DataSpec{{}, {Worksheet: "mysheet"}, {Range: "myrange"}}},
			want:    &DataSpec{Worksheet: "mysheet", Range: "myrange"},
			wantErr: false,
		},
		{
			name:    "SimpleWorkbookclash",
			args:    args{specs: []*DataSpec{{}, {Workbook: "mybook"}, {Workbook: "myotherbook"}}},
			wantErr: true,
		},
		{
			name:    "SimpleWorksheetClash",
			args:    args{specs: []*DataSpec{{}, {Worksheet: "mysheet"}, {Worksheet: "myothersheet"}}},
			wantErr: true,
		},
		{
			name:    "SimpleRangeClash",
			args:    args{specs: []*DataSpec{{}, {Range: "myrange"}, {Range: "myotherrange"}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeDataSpecs(tt.args.specs)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeDataSpecs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeDataSpecs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandArgsToDataSpec(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    *DataSpec
		wantErr bool
	}{
		{
			name:    "NoArgs",
			args:    args{args: []string{}},
			want:    &DataSpec{},
			wantErr: false,
		},
		{
			name:    "TooManyArgs",
			args:    args{args: []string{"this", "is", "too", "many"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "BareWorkbook",
			args:    args{args: []string{"myworkbook"}},
			want:    &DataSpec{Workbook: "myworkbook"},
			wantErr: false,
		},
		{
			name:    "BareWorkbookAndSheet",
			args:    args{args: []string{"myworkbook", "myworksheet"}},
			want:    &DataSpec{Workbook: "myworkbook", Worksheet: "myworksheet"},
			wantErr: false,
		},
		{
			name:    "BareWorkbookAndSheetWithRange",
			args:    args{args: []string{"myworkbook", "myworksheet!A1:B100"}},
			want:    &DataSpec{Workbook: "myworkbook", Worksheet: "myworksheet", Range: "A1:B100"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandArgsToDataSpec(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandArgsToDataSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpandArgsToDataSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}
