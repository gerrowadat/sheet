package cmd

import (
	"fmt"
	"log"

	sheet "github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var (
	getCmd = &cobra.Command{
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("get requires a data spec : %v", args)
			}
			return nil
		},
		Use:   "get <data spec>",
		Short: "get a range of data from a sheet",
		Long: `Get data given a spreadsheet ID and a range specifier.
For example:
	  > sheet get SprEaDsHeeTiD 'rawdata'
	  > sheet get SprEaDsHeeTiD 'rawdata!A3:G5'

Or, get data based on aliases ('help alias'):

	> sheet get @mysheet worksheet!A1:B100
	> sheet get @myfavouriterange

`,
		Run: func(cmd *cobra.Command, args []string) {
			doGet(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
}

func doGet(_ *cobra.Command, args []string) {
	srv, err := sheet.GetService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	dataspec, err := sheet.ExpandArgsToDataSpec(args)

	if err != nil {
		log.Fatalf("Unable to expand data spec: %v", err)
	}

	if dataspec.IsWorkbook() {
		log.Fatalf("get command requires a data spec that is a worksheet or range, not a workbook")
	}

	resp, err := srv.Spreadsheets.Values.Get(dataspec.Workbook, dataspec.GetInSheetDataSpec()).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			for i, val := range row {
				fmt.Printf("%v", val)
				if i < len(row)-1 {
					fmt.Printf(",")
				} else {
					fmt.Println()
				}
			}
		}
	}
}
