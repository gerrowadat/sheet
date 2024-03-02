package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias [get|set] [alias] [workbook] [worksheet][!range] ",
	Short: "Manipulate aliases",
	Long: `Get, set or delete sheet, worksheet and range aliases.
Aliases are used to refer to workbooks, worksheets and ranges by a short name.
You may specify just a workbook, or just a worksheet for a workbook or worksheet-level alias
You may specify a range by appending the range name to the worksheet name with a !.

You may then specify aliases to regular commands using @aliasname:

i.e.:
# Set an alias to a range, then get the range
> sheet alias set myrangealias myworkbook myworksheet!myrange
> sheet get @myrangealias

# Set an alias to a workbook, then get a range in a worksheet in that workbook
> sheet alias set mywbalias myworkbook
> sheet get @mywbalias worksheet!range

# Set an alias to a worksheet, then tail that worksheet
> sheet alias set mydata myworkbook myworksheet
> sheet tail @mydata
		`,
	Run: func(cmd *cobra.Command, args []string) {
		doAlias(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}

func doAlias(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}
	switch args[0] {
	case "get":
		doAliasGet(cmd, args)
	case "set":
		doAliasSet(cmd, args)
	case "rm":
		doAliasRm(cmd, args)
	default:
		log.Println("Unknown alias command", args[0])
		cmd.Help()
	}
}

func doAliasAll(cmd *cobra.Command) {
	aliases := sheet.GetAllAliases()
	for k := range aliases {
		spec, err := sheet.GetAlias(k)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v => %v\n", k, spec.String())
	}
}

func doAliasGet(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		doAliasAll(cmd)
		return
	}
	dataspec, err := sheet.GetAlias(args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v => (%v)\n", args[1], dataspec.String())
}

func doAliasSet(cmd *cobra.Command, args []string) {
	if len(args) < 3 || len(args) > 4 {
		fmt.Println("alias set requires 2, 3 or 4 arguments")
		cmd.Help()
		return
	}
	spec := &sheet.DataSpec{}
	if len(args) == 3 {
		// Workbook-level alias
		spec.Workbook = args[2]
	}
	if len(args) == 4 {
		if strings.Contains(args[3], "!") {
			fragments := strings.Split(args[3], "!")
			spec.Workbook = args[2]
			spec.Worksheet = fragments[0]
			spec.Range = fragments[1]
		} else {
			spec.Workbook = args[2]
			spec.Worksheet = args[3]
		}
		fmt.Println("Setting alias", args[1], "to", spec.String())
		sheet.SetAlias(args[1], spec)
		return
	}
}

func doAliasRm(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		fmt.Println("alias rm requires 1 argument")
		cmd.Help()
		return
	}
	sheet.DeleteAlias(args[1])
	fmt.Println("Deleting alias", args[1])
}
