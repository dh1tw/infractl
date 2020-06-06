package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dh1tw/infractl/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the service command
var listCmd = &cobra.Command{
	Use:   "list service1 service2 ...",
	Short: "get the status of one or more systemd services",
	Long: `get the status of one or more systemd services

This command requires elevated privileges`,
	Run: listService,
}

func init() {
	serviceCmd.AddCommand(listCmd)
}

func listService(cmd *cobra.Command, args []string) {

	// Try to read config file
	configFileMsg := ""

	if err := viper.ReadInConfig(); err == nil {
		configFileMsg = fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		if strings.Contains(err.Error(), "Not Found in") {
			configFileMsg = fmt.Sprintf("no config file found")
		} else {
			fmt.Println("Error parsing config file", viper.ConfigFileUsed())
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println(configFileMsg)

	ss := args

	if len(args) == 0 {
		ss = viper.GetStringSlice("service.service")
	}

	res, err := services.Status(ss...)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range res {
		fmt.Println(s)
	}
}
