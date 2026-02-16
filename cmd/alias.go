package cmd

import (
	"fmt"
	"log"

	"github.com/gerrowadat/sheet/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aliasCmd represents the alias command
var (
	aliasSpecPrefix string
	aliasCmd        = &cobra.Command{
		Use:   "alias [get|set] [alias] [workbook] [worksheet][!range] ",
		Short: "Manipulate aliases",
		Long: `Get, set or delete sheet, worksheet and range aliases.
	Aliases are used to refer to workbooks, worksheets and ranges by a short name.
	You may specify just a workbook, or just a worksheet for a workbook or worksheet-level alias
	You may specify a range by appending the range name to the worksheet name with a !.

	You may then specify aliases to regular commands using @aliasname:

	i.e.:

	# Print all aliases
	> sheet alias get

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
)

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.PersistentFlags().StringVar(&aliasSpecPrefix, "alias-spec-prefix", "@", "The default prefix used to specify aliases")
	viper.BindPFlag("alias-spec-prefix", aliasCmd.PersistentFlags().Lookup("alias-spec-prefix"))
}

func doAlias(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		doAliasAll()
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

func doAliasAll() {
	aliases := sheet.GetAllAliases()
	for k := range aliases {
		spec, err := sheet.GetAlias(k)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v => %v\n", k, spec.String())
	}
}

func doAliasGet(_ *cobra.Command, args []string) {
	if len(args) == 1 {
		doAliasAll()
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
	spec, err := sheet.ExpandArgsToDataSpec(args[2:])
	if err != nil {
		log.Fatal(err)
	}
	sheet.SetAlias(args[1], spec)
	fmt.Printf("%v => (%v)\n", args[1], spec.String())
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
