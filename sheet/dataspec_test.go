package sheet

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
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
			name:    "AliasedWorkbook",
			args:    args{args: []string{"@myworkbook"}},
			want:    &DataSpec{Workbook: "mywb"},
			wantErr: false,
		},
		{
			name:    "AliasedWorkbookBadAlias",
			args:    args{args: []string{"@mything"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "BareWorkbookAndSheet",
			args:    args{args: []string{"myworkbook", "myworksheet"}},
			want:    &DataSpec{Workbook: "myworkbook", Worksheet: "myworksheet"},
			wantErr: false,
		},
		{
			name:    "AliasedWorksheet",
			args:    args{args: []string{"@myworksheet"}},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myws"},
			wantErr: false,
		},
		{
			name:    "AliasedWorksheetWithRange",
			args:    args{args: []string{"@myworksheet!A3:F6"}},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myws", Range: "A3:F6"},
			wantErr: false,
		},
		{
			name:    "BareWorkbookAndSheetWithRange",
			args:    args{args: []string{"myworkbook", "myworksheet!A1:B100"}},
			want:    &DataSpec{Workbook: "myworkbook", Worksheet: "myworksheet", Range: "A1:B100"},
			wantErr: false,
		},
		{
			name:    "AliasedWorkbookAndSheet",
			args:    args{args: []string{"@myworkbook", "myworksheet"}},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myworksheet"},
			wantErr: false,
		},
		{
			name:    "AliasedWorkbookAndSheetWithRange",
			args:    args{args: []string{"@myworkbook", "myworksheet!A1:B100"}},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myworksheet", Range: "A1:B100"},
			wantErr: false,
		},
		{
			name:    "AliasedWorkbookAndSheetWithRangeBadAlias",
			args:    args{args: []string{"@mything", "myworksheet!A1:B100"}},
			want:    nil,
			wantErr: true,
		},
	}

	// Setup viper config for each run
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigName("dataspec")
	viper.AddConfigPath("testdata")
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("Error reading test config %v", err)
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

func Test_dataSpecFromAlias(t *testing.T) {
	type args struct {
		aliasname string
	}
	tests := []struct {
		name    string
		args    args
		want    *DataSpec
		wantErr bool
	}{
		{
			name:    "NoSuchAlias",
			args:    args{aliasname: "fiddledeedee"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "WorkbookAlias",
			args:    args{aliasname: "myworkbook"},
			want:    &DataSpec{Workbook: "mywb"},
			wantErr: false,
		},
		{
			name:    "WorksheetAlias",
			args:    args{aliasname: "myworksheet"},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myws"},
			wantErr: false,
		},
		{
			name:    "RangeAlias",
			args:    args{aliasname: "myrange"},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myws", Range: "myr"},
			wantErr: false,
		},
		{
			name:    "BangNotationWithWorksheetAlias",
			args:    args{aliasname: "myworksheet!A1:C5"},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myws", Range: "A1:C5"},
			wantErr: false,
		},
		{
			name:    "BangNotationWithWorkbookAlias",
			args:    args{aliasname: "myworkbook!A1:C5"},
			want:    nil,
			wantErr: true,
		},
	}
	// Setup viper config for each run
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigName("dataspec")
	viper.AddConfigPath("testdata")
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("Error reading test config %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dataSpecFromAlias(tt.args.aliasname)
			if (err != nil) != tt.wantErr {
				t.Errorf("dataSpecFromAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSpecFromAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRange_String(t *testing.T) {
	type fields struct {
		StartRow int
		StartCol int
		EndRow   int
		EndCol   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Simple",
			fields: fields{StartRow: 1, StartCol: 1, EndRow: 10, EndCol: 10},
			want:   "A1:J10",
		},
		{
			name:   "WholeRows",
			fields: fields{StartRow: 1, StartCol: 0, EndRow: 10, EndCol: 0},
			want:   "1:10",
		},
		{
			name:   "WholeCols",
			fields: fields{StartRow: 0, StartCol: 2, EndRow: 0, EndCol: 30},
			want:   "B:AD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataRange{
				StartRow: tt.fields.StartRow,
				StartCol: tt.fields.StartCol,
				EndRow:   tt.fields.EndRow,
				EndCol:   tt.fields.EndCol,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("DataRange.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRange_FromString(t *testing.T) {
	type fields struct {
		StartRow int
		StartCol int
		EndRow   int
		EndCol   int
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DataRange
		wantErr bool
	}{
		{
			name:    "Simple",
			fields:  fields{},
			args:    args{s: "A1:J10"},
			want:    &DataRange{StartRow: 1, StartCol: 1, EndRow: 10, EndCol: 10},
			wantErr: false,
		},
		{
			name:    "SimpleError",
			fields:  fields{},
			args:    args{s: "doot"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "TooManyColons",
			fields:  fields{},
			args:    args{s: "doot:doot:doot"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NonAlphaNum",
			fields:  fields{},
			args:    args{s: "A@1:B2"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "AlsoNonAlphaNum",
			fields:  fields{},
			args:    args{s: "A1:@B2"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataRange{
				StartRow: tt.fields.StartRow,
				StartCol: tt.fields.StartCol,
				EndRow:   tt.fields.EndRow,
				EndCol:   tt.fields.EndCol,
			}
			got, err := d.FromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataRange.FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataRange.FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_colToLetter(t *testing.T) {
	type args struct {
		col int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "A", args: args{col: 1}, want: "A"},
		{name: "Z", args: args{col: 26}, want: "Z"},
		{name: "AA", args: args{col: 27}, want: "AA"},
		{name: "AZ", args: args{col: 52}, want: "AZ"},
		{name: "BA", args: args{col: 53}, want: "BA"},
		{name: "ALM", args: args{col: 1001}, want: "ALM"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colToLetter(tt.args.col); got != tt.want {
				t.Errorf("colToLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_letterToCol(t *testing.T) {
	type args struct {
		letter string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "A",
			args: args{letter: "A"},
			want: 1,
		},
		{
			name: "AA",
			args: args{letter: "AA"},
			want: 27,
		},
		{
			name: "ALM",
			args: args{letter: "ALM"},
			want: 1001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := letterToCol(tt.args.letter); got != tt.want {
				t.Errorf("letterToCol() = %v, want %v", got, tt.want)
			}
		})
	}
}
