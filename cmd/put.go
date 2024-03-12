package cmd

import (
	"log"

	"github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
)

// putCmd represents the put command
var (
	forcePut bool
	putCmd   = &cobra.Command{
		Use:   "put",
		Short: "Write data to gsheets",
		Long: `Write data from stdin to a range or worksheet.

e.g.:

# Write the contents of a file to a worksheet
> sheet put @myworkbook myworksheet < mydata.csv

# Write the contents of a file to a range
> sheet put @myworkbook 'myworksheet!A1:B2' < mydata.csv

This subcommand respects the --protect-worksheets flag and config item.

Note: When writing to worksheet, the worksheet will be cleared first.
 - If you want to append data, use the append subcommand.

Note: When writing to a range, the range size must match the size of the data being written. 
 - If the range is larger, the extra cells will be cleared. If the range is smaller, the write will fail.

`,
		Run: func(cmd *cobra.Command, args []string) {
			doPut(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.PersistentFlags().BoolVar(&forcePut, "force-put", false, "Override protect-worksheets and put data")
}

func doPut(_ *cobra.Command, args []string) {
	srv, err := sheet.GetService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spec, err := sheet.ExpandArgsToDataSpec(args)

	if err != nil {
		log.Fatalf("Unable to expand data spec: %v", err)
	}

	if spec.IsWorkbook() {
		log.Fatalf("Workbooks cannot be....putten to.")
	}

	if spec.IsWorksheet() {
		err = sheet.ClearWorksheet(srv, spec, protectWorksheets, forcePut)

		if err != nil {
			log.Fatalf("Unable to clear worksheet: %v", err)
		}
	}
}
