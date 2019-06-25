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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.infractl.yaml)")
	rootCmd.Flags().StringP("microtik-address", "a", "192.168.0.1", "address of your microtik router")
	rootCmd.Flags().IntP("microtik-port", "p", 8728, "API port of your microtik router")
	rootCmd.Flags().StringP("microtik-username", "U", "admin", "username for your microtik router")
	rootCmd.Flags().StringP("microtik-password", "P", "admin", "password for your microtik router")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".remoteRotator") // name of config file (without extension)
		viper.AddConfigPath("$HOME")          // adding home directory as first search path
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match
}
