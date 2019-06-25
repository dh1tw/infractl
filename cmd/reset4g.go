package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dh1tw/infractl/microtik"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reset4gCmd represents the reset4g command
var reset4gCmd = &cobra.Command{
	Use:   "reset4g",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: reset4g,
}

func init() {
	rootCmd.AddCommand(reset4gCmd)
}

func reset4g(cmd *cobra.Command, args []string) {
	// Try to read config file
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		if strings.Contains(err.Error(), "Not Found in") {
			fmt.Println("no config file found")
		} else {
			fmt.Println("Error parsing config file", viper.ConfigFileUsed())
			fmt.Println(err)
			os.Exit(1)
		}
	}

	viper.BindPFlag("microtik.address", rootCmd.Flags().Lookup("microtik-address"))
	viper.BindPFlag("microtik.port", rootCmd.Flags().Lookup("microtik-port"))
	viper.BindPFlag("microtik.username", rootCmd.Flags().Lookup("microtik-username"))
	viper.BindPFlag("microtik.password", rootCmd.Flags().Lookup("microtik-password"))

	mConfig := microtik.Config{
		Address:  viper.GetString("microtik.address"),
		Port:     viper.GetInt("microtik.port"),
		Username: viper.GetString("microtik.username"),
		Password: viper.GetString("microtik.password"),
	}

	mt := microtik.New(mConfig)
	if err := mt.Reset4G(); err != nil {
		log.Fatal(err)
	}
	log.Println("4G reset sucessfully initiated")
}
