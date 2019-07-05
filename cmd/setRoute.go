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

// setRouteCmd represents the setRoute command
var setRouteCmd = &cobra.Command{
	Use:   "set-route",
	Short: "set a parameter for a route of a microtik router",
	Long: `set a parameter for a route of a microtik router

This command will connect to a microtik router and execute the set command
on a particular route (ip/route). Since the routes don't have static labels,
the routes are identified by their comment fields. Therefore a shorthand / name
(string) and the corresponding comment (string) has to be set in the config file
when calling this command.

See the example config file for more details:
https://github.com/dh1tw/infractl/blob/master/.infractl.toml

WARNING:
Executing the wrong commands (e.g. disable=yes) might lock you out!

Example:
./infractl set-route --config=.myconfig.toml -r adsl -c "disabled=false"
	`,
	Run: setRoute,
}

func init() {
	rootCmd.AddCommand(setRouteCmd)
	setRouteCmd.Flags().StringP("address", "a", "192.168.0.1", "address of your microtik router")
	setRouteCmd.Flags().IntP("port", "p", 8728, "API port of your microtik router")
	setRouteCmd.Flags().StringP("username", "U", "admin", "username for your microtik router")
	setRouteCmd.Flags().StringP("password", "P", "admin", "password for your microtik router")
	setRouteCmd.Flags().StringP("route", "r", "adsl", "route name (route must be in config file")
	setRouteCmd.Flags().StringP("command", "c", "disabled=false", "command")
}

func setRoute(cmd *cobra.Command, args []string) {
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

	fmt.Println(configFileMsg)

	route, err := cmd.Flags().GetString("route")
	if err != nil {
		log.Fatal(err)
	}

	command, err := cmd.Flags().GetString("command")
	if err != nil {
		log.Fatal(err)
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

	err = mt.SetRoute(route, command)
	if err != nil {
		log.Fatal(err)
	}
}
