package sheet

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func setupTempConfig(t *testing.T) {
	tempconfigfile := t.TempDir() + "/alias.yaml"
	os.Create(tempconfigfile)

	localconfig, err := os.Open("testdata/alias.yaml")

	if err != nil {
		t.Errorf("Error opening testdata/alias.yaml ")
	}

	tempconfig, err := os.OpenFile(tempconfigfile, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		t.Errorf("Error opening %v", tempconfigfile)
	}
	// Setup viper config for each run
	// Copy our actual test file to a tempdir, since we edit it.
	_, err = io.Copy(tempconfig, localconfig)
	if err != nil {
		t.Errorf("Error copying testdata/alias.yaml to tempdir: %v", err)
	}

	tempconfig.Close()

	tempconfig, err = os.Open(tempconfigfile)
	if err != nil {
		t.Errorf("Error opening %v", tempconfigfile)
	}

	viper.SetConfigFile(tempconfigfile)
	err = viper.ReadConfig(tempconfig)
	if err != nil {
		t.Errorf("Error reading test config %v", err)
	}
}
func TestGetAlias(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *DataSpec
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "NotFound",
			args:    args{name: "notfound"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "JustWorkbook",
			args:    args{name: "myworkbook"},
			want:    &DataSpec{Workbook: "mywb"},
			wantErr: false,
		},
		{
			name:    "AllFields",
			args:    args{name: "myrange"},
			want:    &DataSpec{Workbook: "mywb", Worksheet: "myws", Range: "myr"},
			wantErr: false,
		},
	}
	setupTempConfig(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAlias(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetAlias(t *testing.T) {
	type args struct {
		name string
		spec *DataSpec
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantAfter *DataSpec
	}{
		// TODO: Add test cases.
		{
			name:      "ReplaceExistingSameType",
			args:      args{name: "myworkbook", spec: &DataSpec{Workbook: "myotherwb"}},
			wantErr:   false,
			wantAfter: &DataSpec{Workbook: "myotherwb"},
		},
		{
			name:      "ReplaceExistingDifferentType",
			args:      args{name: "myrange", spec: &DataSpec{Workbook: "mythirdwb"}},
			wantErr:   false,
			wantAfter: &DataSpec{Workbook: "mythirdwb"},
		},
		{
			name:      "ReplaceAllFields",
			args:      args{name: "myrange", spec: &DataSpec{Workbook: "a", Worksheet: "b", Range: "c"}},
			wantErr:   false,
			wantAfter: &DataSpec{Workbook: "a", Worksheet: "b", Range: "c"},
		},
	}
	setupTempConfig(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetAlias(tt.args.name, tt.args.spec); (err != nil) != tt.wantErr {
				t.Errorf("SetAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := GetAlias(tt.args.name)
			if err != nil {
				t.Errorf("GetAlias() (after SetAlias) error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.wantAfter) {
				t.Errorf("GetAlias() = %v, want %v", got, tt.wantAfter)
			}
		})
	}
}
