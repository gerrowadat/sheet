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

func doConfig(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("config called without arguments")
		return
	}
	if args[0] == "get" {
		if len(args) == 1 {
			// Print entire config
			for k, v := range viper.AllSettings() {
				fmt.Printf("%s: %v\n", k, v)
			}
		}
	}
}
