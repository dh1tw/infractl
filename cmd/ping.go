package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dh1tw/infractl/connectivity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping one or more hosts",
	Long: `Ping will send ICMP pings to one or more hosts and return the
the average round trip time.

Under Linux, sending pings require elevated privileges. You can either:
1. run this command with root privileges (sudo)
2. enable unprivileged pings with:
	'sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"'

The result can be optionally written to stdio in JSON. In this case the
ping's round trip time will be returned in nano seconds.`,

	Run: checkPing,
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringSliceP("address", "a", []string{"google.com", "8.8.8.8"}, "one or more address/url to ping")
	pingCmd.Flags().Bool("json", false, "outputs the result as json")
}

func checkPing(cmd *cobra.Command, args []string) {

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

	viper.BindPFlag("ping.address", cmd.Flags().Lookup("address"))
	viper.BindPFlag("ping.json", cmd.Flags().Lookup("json"))

	addrs := viper.GetStringSlice("ping.address")
	outputJSON := viper.GetBool("ping.json")

	if !outputJSON {
		fmt.Println(configFileMsg)
	}

	results := connectivity.PingHosts(addrs)

	if outputJSON {
		j, err := json.Marshal(results)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(j))
		return
	}

	for _, r := range results {
		fmt.Printf("%+v\n", r)
	}
}
