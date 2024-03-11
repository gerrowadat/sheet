package cmd

import (
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/sheets/v4"
)

// rmCmd represents the rm command
var (
	protectWorkbooks  bool
	protectWorksheets bool
	forceDelete       bool
	rmCmd             = &cobra.Command{
		Use:   "rm",
		Short: "Delete a workbook, worksheet or range",
		Long: `Delete a workbook, worksheet or range. Expands aliases.

Note: This command will delete data. No, really. It works just like rm.

If you delete a workbook, it will be gone.
If you delete a worksheet, it will be gone.
If you delete a range, that range will be filled with empty cells.

If you're worried that your amazing scripting skillz will cause you to delete something important,
you can specify --protect-workbooks and --protect-worksheets on your command line
(or protect-workbooks and protect-worksheets in the config file
you know you're not generally going to be deleting workbooks or worksheets). This will cause 'sheet rm' to
never actually delete worksheets or workbooks. It will still delete ranges, though. It will return non-zero if
it fails to delete a workbook or worksheet in this way.

If you're sure, you can also use the --force-delete flag to override the protection.
`,
		Run: func(cmd *cobra.Command, args []string) {
			doRm(cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)

	rmCmd.PersistentFlags().BoolVar(&protectWorkbooks, "protect-workbooks", false, "Never delete any workbooks")
	viper.BindPFlag("protect-workbooks", rmCmd.PersistentFlags().Lookup("protect-workbooks"))

	rmCmd.PersistentFlags().BoolVar(&protectWorksheets, "protect-worksheets", false, "Never delete any worksheets")
	viper.BindPFlag("protect-worksheets", rmCmd.PersistentFlags().Lookup("protect-worksheets"))

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
		if (viper.GetBool("protect-workbooks") || protectWorkbooks) && !forceDelete {
			return false
		}
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
	return nil
}
