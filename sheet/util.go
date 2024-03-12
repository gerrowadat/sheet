package sheet

import (
	"fmt"

	"github.com/spf13/viper"
	"google.golang.org/api/sheets/v4"
)

func ClearWorksheet(srv *sheets.Service, spec *DataSpec, protect bool, force bool) error {
	if (protect || viper.GetBool("protect-worksheets")) && !force {
		return fmt.Errorf("protection prevents clearing of: (%v)", spec.String())
	}

	fmt.Printf("Clearing: %v\n", spec.String())

	return nil
}
