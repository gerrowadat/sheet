package cmd

import (
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/sheets/v4"
)

// rmCmd represents the rm command
var (
	forceDelete bool
	rmCmd       = &cobra.Command{
		Use:   "rm",
		Short: "Delete a worksheet or range",
		Long: `Delete a worksheet or range. Expands aliases.

Note: This command will delete data. No, really. It works just like rm.

This command cannot delete workbooks.
Why can it create workbooks (with 'sheet touch') but not delete them?
Because...that's how the Sheets API works :-/

If you delete a worksheet, it will be gone.
If you delete a range, that range will be filled with empty cells.

If you're worried that your amazing scripting skillz will cause you to delete something important,
you can specify --protect-worksheets on your command line (or protect-worksheets in the config file),
This will protect all sheets from deletion.

If you're sure, you can also use the --force-delete flag to override the protection.
`,
		Run: func(cmd *cobra.Command, args []string) {
			doRm(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.PersistentFlags().BoolVar(&forceDelete, "force-delete", false, "Override protect-workbooks and protect-worksheets")
}

func doRm(_ *cobra.Command, args []string) {
	srv, err := sheet.GetService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spec, err := sheet.ExpandArgsToDataSpec(args)

	if err != nil {
		log.Fatalf("Unable to expand data spec: %v", err)
	}

	if spec.IsWorkbook() {
		log.Fatalf("You can't delete a workbook with this command.")
	}

	if mayDelete(spec) {
		err = DeleteSpecified(srv, spec)
		if err != nil {
			log.Fatalf("Unable to delete (%v): %v", spec.String(), err)
		}
	} else {
		log.Fatalf("Protection prevents deletion of: (%v). Use --force-delete to force", spec.String())
	}

}

func mayDelete(spec *sheet.DataSpec) bool {

	if spec.IsWorkbook() {
		return false
	}

	if spec.IsWorksheet() {
		if (viper.GetBool("protect-worksheets") || protectWorksheets) && !forceDelete {
			return false
		}
	}

	return true
}

func DeleteSpecified(srv *sheets.Service, spec *sheet.DataSpec) error {
	fmt.Printf("Deleting: %v\n", spec.String())

	if spec.IsWorkbook() {
		return fmt.Errorf("you can't delete a workbook with this command")
	}

	wb, err := sheets.NewSpreadsheetsService(srv).Get(spec.Workbook).Do()

	if err != nil {
		return fmt.Errorf("unable to retrieve workbook: %v", err)
	}

	var rm_id int64
	for _, sheet := range wb.Sheets {
		if sheet.Properties.Title == spec.Worksheet {
			rm_id = sheet.Properties.SheetId
		}
	}

	if rm_id == 0 {
		return fmt.Errorf("unable to find worksheet: %v", spec.Worksheet)
	}

	if spec.IsWorksheet() {
		srv.Spreadsheets.BatchUpdate(spec.Workbook,
			&sheets.BatchUpdateSpreadsheetRequest{
				Requests: []*sheets.Request{
					{
						DeleteSheet: &sheets.DeleteSheetRequest{
							SheetId: rm_id},
					},
				},
			}).Do()
	}

	if spec.IsRange() {
		_, err := srv.Spreadsheets.Values.Clear(spec.Workbook, spec.Worksheet+"!"+spec.Range.String(), &sheets.ClearValuesRequest{}).Do()
		if err != nil {
			return fmt.Errorf("unable to clear range (%v): %v", spec, err)
		}
	}

	return nil
}
