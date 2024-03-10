package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or manipulate the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		doConfig(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func doConfig(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("config called without arguments")
		return
	}
	if args[0] == "get" {
		if len(args) == 1 {
			// Print entire config
			fmt.Printf("# For aliases, use 'sheet alias get' instead\n")
			for k, v := range viper.AllSettings() {
				if k != "aliases" {
					fmt.Printf("%s: %v\n", k, v)
				}
			}
		}
		if len(args) == 2 {
			// Print a specific config item
			fmt.Printf("%s: %v\n", args[1], viper.Get(args[1]))
		}
	}
	if args[0] == "set" {
		current := viper.Get(args[1])
		if current == nil {
			fmt.Printf("No such config item: %s\n", args[1])
			return
		}
		if len(args) == 3 {
			viper.Set(args[1], args[2])
			viper.WriteConfig()
		}
	}
}
