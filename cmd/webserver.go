package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	webserver "github.com/dh1tw/infractl/app"
	"github.com/dh1tw/infractl/microtik"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// webServerCmd represents the web command
var webServerCmd = &cobra.Command{
	Use:   "web",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: webServer,
}

func init() {
	rootCmd.AddCommand(webServerCmd)
	webServerCmd.Flags().StringP("address", "w", "127.0.0.1", "address of the webserver (use '0.0.0.0' to listen on all network adapters)")
	webServerCmd.Flags().IntP("port", "k", 6556, "webserver http port")
}

func webServer(cmd *cobra.Command, args []string) {

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

	viper.BindPFlag("web.address", cmd.Flags().Lookup("address"))
	viper.BindPFlag("web.port", cmd.Flags().Lookup("port"))
	errorCh := make(chan struct{})

	addr := webserver.Address(viper.GetString("web.address"))
	port := webserver.Port(viper.GetInt("web.port"))

	opts := []webserver.Option{addr, port, webserver.ErrorCh(errorCh)}

	if viper.IsSet("microtik.address") &&
		viper.IsSet("microtik.port") &&
		viper.IsSet("microtik.username") &&
		viper.IsSet("microtik.password") {

		mtConfig := microtik.Config{
			Address:  viper.GetString("microtik.address"),
			Port:     viper.GetInt("microtik.port"),
			Username: viper.GetString("microtik.username"),
			Password: viper.GetString("microtik.password"),
		}

		mtOpts := []microtik.Option{}

		if viper.IsSet("microtik.routes.routes") {
			routes := viper.GetStringSlice("microtik.routes.routes")
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

				mtOpt := microtik.RouteID(name, comment)
				mtOpts = append(mtOpts, mtOpt)
				opts = append(opts, webserver.Route(name))
			}
		}

		mt := webserver.Microtik(microtik.New(mtConfig, mtOpts...))
		opts = append(opts, mt)
	}

	if viper.IsSet("mf823.address") &&
		viper.IsSet("mf823.parameters") {
		mf832Addr := webserver.Mf823Address(viper.GetString("mf823.address"))
		mf832Params := webserver.Mf823Parameters(viper.GetStringSlice("mf823.parameters"))
		opts = append(opts, mf832Addr, mf832Params)
	}

	if viper.IsSet("ping.enabled") &&
		viper.IsSet("ping.interval") {
		pingEnabled := webserver.PingEnabled(viper.GetBool("ping.enabled"))
		pingInterval := webserver.PingInterval(viper.GetDuration("ping.interval"))
		opts = append(opts, pingEnabled, pingInterval)
	}

	services := viper.GetStringSlice("systemd.services")
	for _, s := range services {
		service := webserver.Service(s)
		opts = append(opts, service)
	}

	opts = append(opts, webserver.PingAddress(viper.GetStringSlice("ping.address")))

	webserver := webserver.New(opts...)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go webserver.Serve()

	select {
	case <-errorCh:
		log.Fatal("something failed")
	case <-c:
	}

}
