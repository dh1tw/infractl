package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dh1tw/infractl/mf823"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// status4gCmd represents the status4g command
var status4gCmd = &cobra.Command{
	Use:   "status4g",
	Short: "Request the status from a ZTE MF823 4G USB Modem",
	Long: `Request the status from a ZTE MF823 4G USB Modem

The status of a MF823 can be queried through a REST interface. The list
of possible parameters is pretty long. The example config file provided
with the source code (https://github.com/dh1tw/infractl/.infractl.toml)
should be complete.
However a list of parameters has to be supplied when calling this command.

The result can be optionally written to stdio in JSON.
`,
	Run: status4g,
}

func init() {
	rootCmd.AddCommand(status4gCmd)
	status4gCmd.Flags().String("address", "192.168.3.1", "address of the ZTE MF823 4G Stick")
	status4gCmd.Flags().StringSlice("parameters", []string{"network_type", "network_provider", "signalbar"}, "list of status parameters")
	status4gCmd.Flags().Bool("json", false, "outputs the result as json")
}

func status4g(cmd *cobra.Command, args []string) {

	// Try to read config file
	configFileMsg := ""

	if err := viper.ReadInConfig(); err == nil {
		configFileMsg = fmt.Sprintf("Using config file: %s\n", viper.ConfigFileUsed())
	} else {
		if strings.Contains(err.Error(), "Not Found in") {
			configFileMsg = fmt.Sprintf("no config file found\n")
		} else {
			fmt.Println("Error parsing config file", viper.ConfigFileUsed())
			fmt.Println(err)
			os.Exit(1)
		}
	}

	viper.BindPFlag("mf823.address", cmd.Flags().Lookup("address"))
	viper.BindPFlag("mf823.parameters", cmd.Flags().Lookup("parameters"))
	viper.BindPFlag("mf823.json", cmd.Flags().Lookup("json"))

	address := viper.GetString("mf823.address")
	params := viper.GetStringSlice("mf823.parameters")
	outputJSON := viper.GetBool("mf823.json")

	if !outputJSON {
		fmt.Println(configFileMsg)
	}

	res, err := mf823.Status(address, params...)

	if outputJSON {
		if err != nil {
			// if there is a problem, return an empty json object
			res = make(map[string]interface{})
		}
		j, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(j)) // to stdio
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status MF823 (%s):\n", address)
	for k, v := range res {
		fmt.Printf("%s: %v\n", k, v)
	}
}
