/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/lib"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		// Exactly 1 arg.
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Use:   "ls <spreadsheet ID/alias>",
	Short: "List worksheets in the sheet.",
	Run: func(cmd *cobra.Command, args []string) {
		doLs(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func doLs(_ *cobra.Command, args []string) {
	srv, err := sheet.GetService()

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	dataspec, err := sheet.ExpandArgsToDataSpec(args)

	if err != nil {
		log.Fatalf("Unable to expand data spec: %v", err)
	}

	if !dataspec.IsWorkbook() {
		log.Fatalf("data spec must specify a workbook: %v", args)
	}

	resp, err := srv.Spreadsheets.Get(dataspec.Workbook).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve sheet Id %v: %v", args[0], err)
	}

	for _, sheet := range resp.Sheets {
		fmt.Println(sheet.Properties.Title)
	}
}
