/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// catCmd represents the cat command
var catCmd = &cobra.Command{
	Use:   "cat [data spec]",
	Short: "Output the contents of a worksheet",
	Long: `Data spec must specify a worksheet, i.e.:
> sheet cat SpreAdSheeTiD myworksheet
> sheet cat @myworkbook myworksheet
> sheet cat @myworksheet `,
	Run: func(cmd *cobra.Command, args []string) {
		doCat(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(catCmd)
}

func doCat(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	client := sheet.GetClient(clientSecretFile, authTokenFile)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
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

	start := 1
	// --read-chunksize
	end := chunkSize
	chunkspec := fmt.Sprintf("%v!%v:%v", dataspec.Worksheet, start, end)

	resp, err := srv.Spreadsheets.Values.Get(dataspec.Workbook, chunkspec).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for {
		sheet.PrintValues(resp)

		if len(resp.Values) < chunkSize {
			break
		}

		start = 1 + end
		end = start + (chunkSize - 1)
		chunkspec = fmt.Sprintf("%v!%v:%v", dataspec.Worksheet, start, end)
		resp, err = srv.Spreadsheets.Values.Get(dataspec.Workbook, chunkspec).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet: %v", err)
		}
	}
}
