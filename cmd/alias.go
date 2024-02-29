package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
> sheet get @range:myrangealias

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
	default:
		log.Println("Unknown alias command", args[0])
		cmd.Help()
	}
}

func doAliasGet(cmd *cobra.Command, args []string) {

	all := viper.GetStringMap("aliases")

	if len(args) == 1 { // Print all aliases on 'sheet alias get'
		if all != nil {
			wb, ok := all["workbook"]
			if ok {
				for k, v := range wb.(map[string]interface{}) {
					fmt.Printf("%v -> [workbook] %v\n", k, v)
				}
			}
			ws, ok := all["worksheet"]
			if ok {
				for k, v := range ws.(map[string]interface{}) {
					fmt.Printf("%v -> [worksheet] %v\n", k, v)
				}
			}
			r, ok := all["range"]
			if ok {
				for k, v := range r.(map[string]interface{}) {
					fmt.Printf("%v -> [range] (", k)
					for t, u := range v.(map[string]interface{}) {
						fmt.Printf("%v:%v ", t, u)
					}
					fmt.Println(")")
				}
			}
		}
	}
	alias := viper.Get("aliases." + args[0])
	if alias != nil {
		for k, v := range alias.(map[string]interface{}) {
			fmt.Printf("%s: %v\n", k, v)
		}
	}
}

func doAliasSet(cmd *cobra.Command, args []string) {
	if len(args) < 3 || len(args) > 4 {
		fmt.Println("alias set requires 2, 3 or 4 arguments")
		cmd.Help()
		return
	}
	if len(args) == 3 {
		viper.Set("aliases.workbook."+args[1], args[2])
		fmt.Println("Setting alias", args[1], "to", args[2])
		fmt.Println(viper.ConfigFileUsed())
		doAliasGet(cmd, []string{})
		viper.WriteConfig()
		return
	}
	if len(args) == 4 {
		if strings.Contains(args[3], "!") {
			fragments := strings.Split(args[3], "!")
			viper.Set("aliases.range."+args[1]+".workbook", args[2])
			viper.Set("aliases.range."+args[1]+".worksheet", fragments[0])
			viper.Set("aliases.range."+args[1]+".range", fragments[1])
		} else {
			viper.Set("aliases.worksheet."+args[1]+".workbook", args[2])
			viper.Set("aliases.worksheet."+args[1]+".worksheet", args[3])
		}
		fmt.Println("Setting alias", args[1], "to", args[2], args[3])
		viper.WriteConfig()
		return
	}
}
