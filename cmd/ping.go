package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dh1tw/infractl/connectivity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping host1, host2, ...",
	Short: "Ping one or more hosts",
	Long: `Ping will send ICMP pings to one or more hosts and return the
the average round trip time.

Under Linux, sending pings require elevated privileges. You can either:
1. run this command with root privileges (sudo)
2. bind this executable to a raw socket (adopt the exact location):
	'setcap cap_net_raw=+ep /usr/local/bin/infractl'

The result can be optionally written to stdio in JSON. In this case the
ping's round trip time will be returned in nano seconds.

The hosts can be either specified in the config file or provided as
an argument. Example:
$ infractl ping google.com redhat.com dh1tw.de

`,

	Run: checkPing,
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().Bool("json", false, "outputs the result as json")
	pingCmd.Flags().IntP("samples", "s", 1, "amount of pings set per host")
	pingCmd.Flags().DurationP("timeout", "t", time.Second*2, "timeout for this query")
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

	viper.BindPFlag("ping.timeout", cmd.Flags().Lookup("timeout"))
	viper.BindPFlag("ping.samples", cmd.Flags().Lookup("samples"))
	viper.BindPFlag("ping.json", cmd.Flags().Lookup("json"))

	addrs := viper.GetStringSlice("ping.address")
	outputJSON := viper.GetBool("ping.json")
	timeout := viper.GetDuration("ping.timeout")
	samples := viper.GetInt("ping.samples")

	if len(args) != 0 {
		addrs = args
	}

	if len(args) == 0 {
		fmt.Println("no urls / ip addresses provided")
	}

	if !outputJSON {
		fmt.Println(configFileMsg)
	}

	results := connectivity.PingHosts(addrs, timeout, samples)

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
