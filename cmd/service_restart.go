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

// restartCmd represents the service command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart one or more systemd services",
	Long: `restart one or more systemd services

This command requires elevated privileges`,
	Run: restart,
}

func init() {
	serviceCmd.AddCommand(restartCmd)
	restartCmd.Flags().StringSliceP("service", "s", []string{"myservice"}, "list of services to be restarted")
}

func restart(cmd *cobra.Command, args []string) {

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

	viper.BindPFlag("service.service", cmd.Flags().Lookup("service"))

	ss := viper.GetStringSlice("service.service")

	for _, s := range ss {
		err := services.Restart(s)
		if err != nil {
			log.Printf("%s: %v", s, err)
		}
	}
}
