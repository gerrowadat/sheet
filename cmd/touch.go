package cmd

import (
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/sheets/v4"
)

// touchCmd represents the touch command
var (
	defaultWorkbookTitle string
	touchCmd             = &cobra.Command{
		Use:   "touch",
		Short: "touch a worksheet (creating it if it doesn't exist)",
		Long: `This command simply touches a worksheet, creating it if it doesn't exist.

# Create a new workbook (will be in the root of your drive)
sheet touch workbook "TPS Reports"
sheet touch workbook # Will use the default workbook title from the config or --default-workbook-title

# Create a new worksheet in an existing workbook
sheet touch worksheet MyWoRkBoOk mynewsheet

# Create a new worksheet in a workbook alias
sheet touch worksheet @myworkbook mynewsheet

# Create a new sheet that's referred to by an existing alias.
# (Only works for worksheets, not workbooks)
sheet touch worksheet @mynewsheet
`,
		Run: func(cmd *cobra.Command, args []string) {
			doTouch(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(touchCmd)
	touchCmd.PersistentFlags().StringVar(&defaultWorkbookTitle, "default-workbook-title", "", "The default title of touched workbooks")
	viper.BindPFlag("default-workbook-title", touchCmd.PersistentFlags().Lookup("default-workbook-title"))
}

func doTouch(_ *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatalf("touch command requires a subcommand: workbook or worksheet")
	}

	srv, err := sheet.GetService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	switch args[0] {
	case "workbook":
		doTouchWorkbook(srv, args[1:])
	case "worksheet":
		doTouchWorksheet(srv, args[1:])
	default:
		log.Fatalf("Unknown touch command: %v", args[0])
	}
}

func doTouchWorkbook(srv *sheets.Service, args []string) {
	// Only argument is the workbook title.
	if len(args) > 1 {
		log.Fatalf("touch workbook requires 0 or 1 arguments")
	}
	var workbookTitle string

	if len(args) == 1 {
		workbookTitle = args[0]
	} else {
		// Flag, then config
		if defaultWorkbookTitle != "" {
			workbookTitle = defaultWorkbookTitle
		} else {
			// If this is also "", then the gsheets default "Untitled Spreadsheet" will be used.
			workbookTitle = viper.GetString("default-workbook-title")
		}
	}

	resp, err := srv.Spreadsheets.Create(&sheets.Spreadsheet{Properties: &sheets.SpreadsheetProperties{Title: workbookTitle}}).Do()

	if err != nil {
		log.Fatalf("Unable to create workbook: %v", err)
	}
	// Simply print the new spreadsheet ID, for doing terrifying scripts.
	fmt.Println(resp.SpreadsheetId)
}

func doTouchWorksheet(srv *sheets.Service, args []string) {
	dataspec, err := sheet.ExpandArgsToDataSpec(args)
	if err != nil {
		log.Fatalf("Unable to expand data spec: %v", err)
	}
	if !dataspec.IsWorksheet() {
		log.Fatalf("touch worksheet requires a worksheet spec")
	}

	// Get the existing worksheets
	resp, err := srv.Spreadsheets.Get(dataspec.Workbook).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve workbook: %v", err)
	}

	// Check if the worksheet exists
	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == dataspec.Worksheet {
			// The worksheet exists, nothing to do
			return
		}
	}

	// The worksheet doesn't exist, create it
	_, err = srv.Spreadsheets.BatchUpdate(dataspec.Workbook,
		&sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					AddSheet: &sheets.AddSheetRequest{
						Properties: &sheets.SheetProperties{
							Title: dataspec.Worksheet,
						},
					},
				},
			},
		}).Do()

	if err != nil {
		log.Fatalf("Unable to create worksheet: %v", err)
	}
}
