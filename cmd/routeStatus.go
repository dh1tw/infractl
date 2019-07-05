// Copyright © 2019 Tobias Wellnitz, DH1TW <Tobias.Wellnitz@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dh1tw/infractl/microtik"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// routeStatusCmd represents the routeStatus command
var routeStatusCmd = &cobra.Command{
	Use:   "route-status",
	Short: "Get the status of one or more routes of a microtik router",
	Long: `Get the status of one or more routes of a microtik router

This command will connect to a microtik router and retrieve the status if
the requested route (ip/route) is disabled and/or active. Since the routes
don't have static labels, the most reliable way to find the route is by
it's comment field. Therefore a shorthand identified / name (string) and the
corresponding comment (string) has to be set in the config file when calling
this command.

See the example config file for more details:
https://github.com/dh1tw/infractl/blob/master/.infractl.toml

`,

	Run: routeStatus,
}

func init() {
	rootCmd.AddCommand(routeStatusCmd)
	routeStatusCmd.Flags().StringP("address", "a", "192.168.0.1", "address of your microtik router")
	routeStatusCmd.Flags().IntP("port", "p", 8728, "API port of your microtik router")
	routeStatusCmd.Flags().StringP("username", "U", "admin", "username for your microtik router")
	routeStatusCmd.Flags().StringP("password", "P", "admin", "password for your microtik router")
	routeStatusCmd.Flags().Bool("json", false, "outputs the result as json")
}

func routeStatus(cmd *cobra.Command, args []string) {

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

	viper.BindPFlag("microtik.address", cmd.Flags().Lookup("address"))
	viper.BindPFlag("microtik.port", cmd.Flags().Lookup("port"))
	viper.BindPFlag("microtik.username", cmd.Flags().Lookup("username"))
	viper.BindPFlag("microtik.password", cmd.Flags().Lookup("password"))
	viper.BindPFlag("microtik.routes.json", cmd.Flags().Lookup("json"))

	outputJSON := viper.GetBool("microtik.routes.json")

	if !outputJSON {
		fmt.Println(configFileMsg)
	}

	mConfig := microtik.Config{
		Address:  viper.GetString("microtik.address"),
		Port:     viper.GetInt("microtik.port"),
		Username: viper.GetString("microtik.username"),
		Password: viper.GetString("microtik.password"),
	}

	if !viper.IsSet("microtik.routes.routes") {
		log.Fatal("key microtik.routes not found in config file")
	}

	routes := viper.GetStringSlice("microtik.routes.routes")
	if len(routes) == 0 {
		log.Fatal("key routes in microtik.routes is missing or empty")
	}

	opts := []microtik.Option{}
	routeNames := []string{}

	for _, r := range routes {
		rMap := viper.GetStringMapString(r)
		if len(rMap) == 0 {
			log.Fatalf("hashmap for route %s not found in config file", r)
		}
		name, ok := rMap["name"]
		if !ok {
			log.Fatalf("hashmap for route %s missing parameter 'name'", r)
		}
		comment, ok := rMap["comment"]
		if !ok {
			log.Fatalf("hashmap for route %s missing parameter 'comment'", r)
		}

		opt := microtik.RouteID(name, comment)
		opts = append(opts, opt)
		routeNames = append(routeNames, name)
	}

	mt := microtik.New(mConfig, opts...)

	results := make(routeStatusResults)

	for _, r := range routeNames {
		res, err := mt.RouteStatus(r)
		if err != nil {
			log.Fatal(err)
		}
		results[r] = res
	}

	if outputJSON {
		j, err := json.Marshal(results)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(j))
		return
	}

	fmt.Printf("%+v\n", results)
}

type routeStatusResults map[string]microtik.RouteResult

func (r routeStatusResults) String() string {
	s := ""
	for name, res := range r {
		s = s + fmt.Sprintf("%s:\n", name)
		for k, v := range res {
			s = s + fmt.Sprintf(" %s: %v\n", k, v)
		}
	}
	return s
}
