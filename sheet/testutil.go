package sheet

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func SetupTempConfig(t *testing.T, cfname string) {
	tempdir := t.TempDir()
	tempconfigfile := tempdir + "/" + cfname + ".yaml"
	fmt.Printf("Creating temp file %v\n", tempconfigfile)
	os.Create(tempconfigfile)

	localconfig, err := os.Open("testdata/" + cfname + ".yaml")

	if err != nil {
		t.Errorf("Error opening testdata/" + cfname + ".yaml ")
	}

	tempconfig, err := os.OpenFile(tempconfigfile, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		t.Errorf("Error opening %v", tempconfigfile)
	}
	// Setup viper config for each run
	// Copy our actual test file to a tempdir, since we edit it.
	_, err = io.Copy(tempconfig, localconfig)
	if err != nil {
		t.Errorf("Error copying testdata/%v.yaml to tempdir: %v", cfname, err)
	}

	tempconfig.Close()

	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetConfigName(cfname)
	viper.AddConfigPath(tempdir)
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("Error reading test config %v", err)
	}
}

func RangeFromString(s string) DataRange {
	ret := DataRange{}
	_, err := ret.FromString(s)
	if err != nil {
		log.Fatalf("error parsing range %v :  %v", s, err)
	}
	return ret
}
