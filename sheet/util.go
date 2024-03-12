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

	_, err := srv.Spreadsheets.Values.Clear(spec.Workbook, spec.GetInSheetDataSpec(), &sheets.ClearValuesRequest{}).Do()

	if err != nil {
		return fmt.Errorf("unable to clear worksheet (%v): %v", spec, err)
	}

	return nil
}

func ClearRange(srv *sheets.Service, spec *DataSpec) error {
	if !spec.IsRange() {
		return fmt.Errorf("not a range: %v", spec.String())
	}
	_, err := srv.Spreadsheets.Values.Clear(spec.Workbook, spec.GetInSheetDataSpec(), &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return fmt.Errorf("unable to clear range: %v", err)
	}
	return nil
}

func checkDataFitsInRange(spec *DataSpec, data [][]string) error {
	rcols, rrows := spec.Range.SizeXY()

	if len(data) > rrows {
		return fmt.Errorf("data overflow: %d rows in range, %d in data", rrows, len(data))
	}

	if len(data[0]) > rcols {
		return fmt.Errorf("data overflow: %d columns in range, %d in data", rcols, len(data[0]))
	}

	return nil
}

func valueRangeFromStrings(data [][]string) *sheets.ValueRange {
	sheet_values := make([][]interface{}, len(data))
	for i, row := range data {
		sheet_values[i] = make([]interface{}, len(row))
		for j, cell := range row {
			sheet_values[i][j] = cell
		}
	}
	return &sheets.ValueRange{Values: sheet_values}
}

func WriteDataToWorksheet(srv *sheets.Service, spec *DataSpec, data [][]string, protect bool, force bool) error {
	err := ClearWorksheet(srv, spec, protect, force)

	if err != nil {
		return err
	}

	_, err = srv.Spreadsheets.Values.Update(spec.Workbook, spec.GetInSheetDataSpec(), valueRangeFromStrings(data)).ValueInputOption("USER_ENTERED").Do()

	return err
}

func WriteDataToRange(srv *sheets.Service, spec *DataSpec, data [][]string) error {

	err := checkDataFitsInRange(spec, data)

	if err != nil {
		return err
	}

	err = ClearRange(srv, spec)

	if err != nil {
		return err
	}

	_, err = srv.Spreadsheets.Values.Update(spec.Workbook, spec.GetInSheetDataSpec(), valueRangeFromStrings(data)).ValueInputOption("USER_ENTERED").Do()

	return err
}
