package sheet

import (
	"os"
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
	localconfig, err := os.Open("testdata/dataspec.yaml")
	if err != nil {
		t.Errorf("Error opening testdata/dataspec.yaml")
	}
	viper.SetConfigFile("testdata/dataspec.yaml")
	err = viper.ReadConfig(localconfig)
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
