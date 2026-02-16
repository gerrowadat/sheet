package cmd

import (
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/lib"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

// tailCmd represents the tail command
var (
	tailLines int
	tailCmd   = &cobra.Command{
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Use:   "tail <spreadsheet ID> <worksheet name> <number of rows>",
		Short: "Show the last few lines of a worksheet",
		Long: `Show the last few non-blank lines in a worksheet.
	e.g.:
	# Show the last 10 lines of the 'myworksheet' worksheet. 
	> sheet tail SpReAdShEetId myworksheet --lines=10
	> sheet tail @mysheet --lines=50`,
		Run: func(cmd *cobra.Command, args []string) {
			doTail(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(tailCmd)
	tailCmd.PersistentFlags().IntVar(&tailLines, "lines", 10, "Lines to output")
}

func doTail(_ *cobra.Command, args []string) {
	srv, err := sheet.GetService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	dataspec, err := sheet.ExpandArgsToDataSpec(args)

	if err != nil {
		log.Fatalf("Unable to expand data spec: %v", err)
	}

	if !dataspec.IsWorksheet() {
		log.Fatalf("data spec must specify a worksheet: %v", args)
	}

	resp, err := srv.Spreadsheets.Get(dataspec.Workbook).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve sheet Id %v: %v", args[0], err)
	}

	if len(resp.Sheets) == 0 {
		log.Fatalf("Sheet ID %v has no sheets", args[0])
	}

	for _, sh := range resp.Sheets {
		if sh.Properties.Title == dataspec.Worksheet {
			// Properties.GridProperties.RowCount gives the grid size, not the amunt of data.
			// This seems to be 1000 for new sheets, so expensively poll through it.
			last_datarow := findLastDataRow(srv, dataspec, sh.Properties.GridProperties.RowCount)
			// We get the last line by default
			tailLines--
			chunkspec := fmt.Sprintf("%v!%v:%v", dataspec.Worksheet, max(1, last_datarow-int64(tailLines)), last_datarow)
			resp, err := srv.Spreadsheets.Values.Get(dataspec.Workbook, chunkspec).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve data from sheet at %v: %v", chunkspec, err)
			}
			sheet.PrintValues(resp, outputFormat)
			return
		}
	}
	// If we get here, we didn't find our worksheet.
	log.Fatalf("unable to find worksheet %v in spreadsheet %v", dataspec.Worksheet, args[0])
}

func findLastDataRow(srv *sheets.Service, dataspec *sheet.DataSpec, chunk_end int64) int64 {
	if chunk_end == 1 {
		return 0
	}

	chunk_start := max(1, chunk_end-int64(readChunkSize))

	// worksheet!chunk_start:chunk_end
	chunkspec := fmt.Sprintf("%v!%v:%v", dataspec.Worksheet, chunk_start, chunk_end)

	resp, err := srv.Spreadsheets.Values.Get(dataspec.Workbook, chunkspec).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// We get no values back if there is no data, so the first non-zero length chunk we see
	// scanning backwards from eof is the end of our data.
	if len(resp.Values) > 0 {
		return chunk_start + int64(len(resp.Values)-1)
	} else {
		return findLastDataRow(srv, dataspec, chunk_start)
	}
}
