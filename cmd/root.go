package cmd

import (
	"log"
	"os"

	"github.com/gerrowadat/sheet/sheet"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	configFormat     = "yaml"
	outputFormat     = sheet.CsvFormat
	inputFormat      = sheet.CsvFormat
	clientSecretFile string
	authTokenFile    string
	readChunkSize    int
	writeChunkSize   int

	rootCmd = &cobra.Command{
		Use:   "sheet",
		Short: "Manipulate google sheet data",
		Long: `A utility to send and recieve data to/from a google
sheet from the command line in various forms.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&clientSecretFile, "clientsecretfile", "", "Client secret file")
	viper.BindPFlag("clientsecretfile", rootCmd.PersistentFlags().Lookup("clientsecretfile"))
	rootCmd.PersistentFlags().StringVar(&authTokenFile, "authtokenfile", "", "where to store our oauth token")
	viper.BindPFlag("authtokenfile", rootCmd.PersistentFlags().Lookup("authtokenfile"))
	// This is passed directly to viper.SetConfigType
	rootCmd.PersistentFlags().StringVar(&configFormat, "configformat", "yaml", "config file format")

	rootCmd.PersistentFlags().IntVar(&readChunkSize, "read-chunksize", 500, "How many rows at a time to read while fetching data")
	viper.BindPFlag("read-chunksize", rootCmd.PersistentFlags().Lookup("read-chunksize"))
	rootCmd.PersistentFlags().IntVar(&writeChunkSize, "write-chunksize", 500, "How many rows at a time to write at a time while updating data")
	viper.BindPFlag("write-chunksize", rootCmd.PersistentFlags().Lookup("write-chunksize"))

	rootCmd.PersistentFlags().Var(&outputFormat, "output-format", "Output format ([csv|tsv])")
	viper.BindPFlag("output-format", rootCmd.PersistentFlags().Lookup("output-format"))
	rootCmd.PersistentFlags().Var(&inputFormat, "input-format", "Input format ([csv|tsv])")
	viper.BindPFlag("input-format", rootCmd.PersistentFlags().Lookup("input-format"))
}

func initializeConfig(_ *cobra.Command) error {
	// With thanks to https://github.com/carolynvs/stingoftheviper

	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelTrace)

	viper.SetConfigName("sheet")
	viper.SetConfigType("yaml")

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not determine home directory")
	}

	configdir := homedir + "/.config/sheet"

	viper.AddConfigPath(configdir)

	err = os.MkdirAll(configdir, os.ModeDir)

	if err != nil {
		log.Fatalf("could not create config directory %v", configdir)
	}

	if err := viper.ReadInConfig(); err != nil {
		// No configs found anywhere, so create a default one
		homepath := os.Getenv("HOME") + "/.config/sheet"
		filename := homepath + "/sheet." + configFormat
		_, err := os.Create(filename)
		return err
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SafeWriteConfig()

	return nil
}
