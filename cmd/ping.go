package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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

type pingResults map[string]pingResult

type pingResult struct {
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
	Failed  bool          `json:"failed,omitempty"`
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

	pingAsync := func(address string, wg *sync.WaitGroup, resCh chan<- pingResult) {
		defer wg.Done()
		var res pingResult
		ping, err := connectivity.Ping(address)
		if err != nil {
			res = pingResult{address, time.Second * 0, true}
		} else {
			res = pingResult{address, ping, false}
		}
		resCh <- res
	}

	resultCh := make(chan pingResult)

	wg := &sync.WaitGroup{}

	for _, addr := range addrs {
		wg.Add(1)
		go pingAsync(addr, wg, resultCh)
	}

	results := make(pingResults)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		results[res.Address] = res
	}

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

func (r pingResult) String() string {
	if r.Failed {
		return fmt.Sprintf("%s: failed", r.Address)
	}
	return fmt.Sprintf("%s: %v", r.Address, r.RTT)
}
