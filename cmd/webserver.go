package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/dh1tw/infractl/app"
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

	addr := app.Address(viper.GetString("web.address"))
	port := app.Port(viper.GetInt("web.port"))
	mt := app.Microtik(microtik.New(mConfig))

	webserver := app.New(addr, port, mt)

	errorCh := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go webserver.ListenHTTP(errorCh)

	select {
	case <-errorCh:
	case <-c:
	}

}
