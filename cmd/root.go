package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "infractl",
	Short: "infractl is a collection of tools for controlling network infrastructure",
	Long: `infractl is a collection of tools for controlling network infrastructure.
To get a list of supported commands, please check --help`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.infractl.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".infractl") // name of config file (without extension)
		viper.AddConfigPath("$HOME")     // adding home directory as first search path
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match
}
