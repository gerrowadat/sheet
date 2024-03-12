package cmd

import (
	"testing"

	"github.com/gerrowadat/sheet/sheet"
)

func Test_mayDelete(t *testing.T) {
	type args struct {
		spec *sheet.DataSpec
	}
	tests := []struct {
		name      string
		args      args
		config    string
		forceflag bool
		want      bool
	}{
		{
			name:   "rmworkbook",
			args:   args{spec: &sheet.DataSpec{Workbook: "mywb"}},
			config: "rm_protect_none",
			want:   false,
		},
		{
			name:      "rmworkbookwithforce",
			args:      args{spec: &sheet.DataSpec{Workbook: "mywb"}},
			config:    "rm_protect_all",
			forceflag: true,
			want:      false,
		},
		{
			name:   "rmunprotectedworksheet",
			args:   args{spec: &sheet.DataSpec{Workbook: "mywb", Worksheet: "myws"}},
			config: "rm_protect_none",
			want:   true,
		},
		{
			name:   "rmprotectedworksheet",
			args:   args{spec: &sheet.DataSpec{Workbook: "mywb", Worksheet: "myws"}},
			config: "rm_protect_all",
			want:   false,
		},
		{
			name:      "rmprotectedworksheetwithforce",
			args:      args{spec: &sheet.DataSpec{Workbook: "mywb", Worksheet: "myws"}},
			config:    "rm_protect_all",
			forceflag: true,
			want:      true,
		},
	}
	for _, tt := range tests {
		sheet.SetupTempConfig(t, tt.config)
		forceDelete = tt.forceflag
		t.Run(tt.name, func(t *testing.T) {
			if got := mayDelete(tt.args.spec); got != tt.want {
				t.Errorf("mayDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}
